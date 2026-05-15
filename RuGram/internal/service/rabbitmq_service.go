package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"rugram-api/internal/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQService struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    mu      sync.RWMutex
    ctx     context.Context
    cancel  context.CancelFunc
}

var (
    rabbitMQInstance *RabbitMQService
    rabbitMQOnce     sync.Once
)

const (
    ExchangeEvents   = "app.events"
    ExchangeDLX      = "app.dlx"
    QueueUserReg     = "wp.auth.user.registered"
    QueueUserRegDLQ  = "wp.auth.user.registered.dlq"
    RoutingKeyReg    = "user.registered"
)

func NewRabbitMQService() (*RabbitMQService, error) {
    var initErr error
    rabbitMQOnce.Do(func() {
        host := os.Getenv("RABBITMQ_HOST")
        port := os.Getenv("RABBITMQ_PORT")
        user := os.Getenv("RABBITMQ_USER")
        pass := os.Getenv("RABBITMQ_PASS")

        if host == "" {
            host = "localhost"
        }
        if port == "" {
            port = "5672"
        }
        if user == "" {
            user = "guest"
        }
        if pass == "" {
            pass = "guest"
        }

        connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)
        conn, err := amqp.Dial(connStr)
        if err != nil {
            initErr = fmt.Errorf("failed to connect to RabbitMQ: %w", err)
            return
        }

        ch, err := conn.Channel()
        if err != nil {
            initErr = fmt.Errorf("failed to open channel: %w", err)
            return
        }

        ctx, cancel := context.WithCancel(context.Background())

        rabbitMQInstance = &RabbitMQService{
            conn:    conn,
            channel: ch,
            ctx:     ctx,
            cancel:  cancel,
        }

        if err := rabbitMQInstance.setupQueues(); err != nil {
            initErr = fmt.Errorf("failed to setup queues: %w", err)
            return
        }

        log.Println("RabbitMQ connected successfully")
    })

    if initErr != nil {
        return nil, initErr
    }
    return rabbitMQInstance, nil
}

func (r *RabbitMQService) setupQueues() error {
    if err := r.channel.ExchangeDeclare(ExchangeEvents, "direct", true, false, false, false, nil); err != nil {
        return err
    }

    if err := r.channel.ExchangeDeclare(ExchangeDLX, "direct", true, false, false, false, nil); err != nil {
        return err
    }

    _, err := r.channel.QueueDeclare(QueueUserReg, true, false, false, false, amqp.Table{
        "x-dead-letter-exchange":    ExchangeDLX,
        "x-dead-letter-routing-key": RoutingKeyReg,
    })
    if err != nil {
        return err
    }

    if err := r.channel.QueueBind(QueueUserReg, RoutingKeyReg, ExchangeEvents, false, nil); err != nil {
        return err
    }

    _, err = r.channel.QueueDeclare(QueueUserRegDLQ, true, false, false, false, nil)
    if err != nil {
        return err
    }

    if err := r.channel.QueueBind(QueueUserRegDLQ, RoutingKeyReg, ExchangeDLX, false, nil); err != nil {
        return err
    }

    return nil
}

func (r *RabbitMQService) PublishUserRegisteredEvent(userID, email, displayName string) error {
    event := models.Event{
        EventID:   fmt.Sprintf("%d", time.Now().UnixNano()),
        EventType: models.EventUserRegistered,
        Timestamp: time.Now(),
        Payload: models.UserRegisteredPayload{
            UserID:      userID,
            Email:       email,
            DisplayName: displayName,
        },
        Metadata: models.EventMetadata{
            Attempt:       0,
            SourceService: "auth-service",
        },
    }

    body, err := json.Marshal(event)
    if err != nil {
        return err
    }

    ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
    defer cancel()

    err = r.channel.PublishWithContext(ctx, ExchangeEvents, RoutingKeyReg, true, false, amqp.Publishing{
        ContentType:  "application/json",
        Body:         body,
        DeliveryMode: amqp.Persistent,
        Timestamp:    time.Now(),
        MessageId:    event.EventID,
    })

    if err != nil {
        return fmt.Errorf("failed to publish event: %w", err)
    }

    log.Printf("Published user.registered event for user: %s", userID)
    return nil
}

func (r *RabbitMQService) ConsumeUserRegistered(handler func(event models.Event) error) error {
    msgs, err := r.channel.Consume(QueueUserReg, "", false, false, false, false, nil)
    if err != nil {
        return err
    }

    go func() {
        for msg := range msgs {
            var event models.Event
            if err := json.Unmarshal(msg.Body, &event); err != nil {
                log.Printf("Failed to unmarshal message: %v", err)
                msg.Nack(false, false)
                continue
            }

            attempt := event.Metadata.Attempt
            maxAttempts := 3

            if err := handler(event); err != nil {
                log.Printf("Failed to handle event (attempt %d/%d): %v", attempt+1, maxAttempts, err)

                if attempt+1 >= maxAttempts {
                    log.Printf("Max attempts reached, sending to DLQ for event: %s", event.EventID)
                    msg.Nack(false, false)
                } else {
                    event.Metadata.Attempt = attempt + 1
                    updatedBody, _ := json.Marshal(event)
                    msg.Nack(true, false)

                    r.channel.Publish(ExchangeEvents, RoutingKeyReg, false, false, amqp.Publishing{
                        ContentType:  "application/json",
                        Body:         updatedBody,
                        DeliveryMode: amqp.Persistent,
                    })
                }
            } else {
                msg.Ack(false)
                log.Printf("Successfully processed event: %s", event.EventID)
            }
        }
    }()

    log.Println("Started consuming user.registered events")
    return nil
}

func (r *RabbitMQService) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.cancel != nil {
        r.cancel()
    }
    if r.channel != nil {
        r.channel.Close()
    }
    if r.conn != nil {
        return r.conn.Close()
    }
    return nil
}

func (r *RabbitMQService) IsConnected() bool {
    r.mu.RLock()
    defer r.mu.RUnlock()
    return r.conn != nil && !r.conn.IsClosed()
}
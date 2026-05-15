package queue

import (
	"log"

	"rugram-api/internal/handlers"
	"rugram-api/internal/service"
)

type Consumer struct {
    rabbitMQSvc  *service.RabbitMQService
    eventHandler *handlers.EventHandler
}

func NewConsumer(rabbitMQSvc *service.RabbitMQService, eventHandler *handlers.EventHandler) *Consumer {
    return &Consumer{
        rabbitMQSvc:  rabbitMQSvc,
        eventHandler: eventHandler,
    }
}

func (c *Consumer) Start() error {
    if c.rabbitMQSvc == nil || !c.rabbitMQSvc.IsConnected() {
        log.Println("RabbitMQ not connected, consumer not started")
        return nil
    }

    return c.rabbitMQSvc.ConsumeUserRegistered(c.eventHandler.HandleUserRegistered)
}
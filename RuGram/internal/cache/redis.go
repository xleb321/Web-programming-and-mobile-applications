package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	client *redis.Client
	ctx    context.Context
	ttl    time.Duration
}

var (
	AppPrefix = "rugram"
)

func NewCacheService() (*CacheService, error) {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")

	if host == "" {
		host = "localhost"
	}
	
	if port == "" {
		port = "6379"
	}

	ttlSeconds := 300
	if ttlStr := os.Getenv("CACHE_TTL_DEFAULT"); ttlStr != "" {
		if ttl, err := strconv.Atoi(ttlStr); err == nil {
			ttlSeconds = ttl
		}
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	})

	ctx := context.Background()

	// Проверка подключения
	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis connection failed: %v. Cache will be disabled.", err)
		return &CacheService{
			client: nil,
			ctx:    ctx,
			ttl:    time.Duration(ttlSeconds) * time.Second,
		}, nil
	}

	log.Println("Redis connected successfully")
	return &CacheService{
		client: client,
		ctx:    ctx,
		ttl:    time.Duration(ttlSeconds) * time.Second,
	}, nil
}

// IsEnabled возвращает true если Redis доступен
func (c *CacheService) IsEnabled() bool {
	return c.client != nil
}

// BuildKey строит ключ с префиксом приложения
func (c *CacheService) BuildKey(parts ...string) string {
	allParts := append([]string{AppPrefix}, parts...)
	return strings.Join(allParts, ":")
}

// Get получает значение из кеша и десериализует его
func (c *CacheService) Get(key string, dest interface{}) error {
	if !c.IsEnabled() {
		return nil
	}

	data, err := c.client.Get(c.ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, dest)
}

// Set сохраняет значение в кеш с TTL
func (c *CacheService) Set(key string, value interface{}, ttl ...time.Duration) error {
	if !c.IsEnabled() {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	expiration := c.ttl
	if len(ttl) > 0 && ttl[0] > 0 {
		expiration = ttl[0]
	}

	return c.client.Set(c.ctx, key, data, expiration).Err()
}

// Del удаляет ключ из кеша
func (c *CacheService) Del(key string) error {
	if !c.IsEnabled() {
		return nil
	}
	return c.client.Del(c.ctx, key).Err()
}

// DelByPattern удаляет все ключи, соответствующие паттерну
func (c *CacheService) DelByPattern(pattern string) error {
	if !c.IsEnabled() {
		return nil
	}

	iter := c.client.Scan(c.ctx, 0, pattern, 0).Iterator()
	for iter.Next(c.ctx) {
		if err := c.client.Del(c.ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// DelByPrefix удаляет все ключи с заданным префиксом
func (c *CacheService) DelByPrefix(prefix string) error {
	fullPrefix := c.BuildKey(prefix)
	pattern := fullPrefix + "*"
	return c.DelByPattern(pattern)
}

// SetWithDefaultTTL сохраняет значение с TTL по умолчанию
func (c *CacheService) SetWithDefaultTTL(key string, value interface{}) error {
	return c.Set(key, value, c.ttl)
}

// GetString получает строковое значение
func (c *CacheService) GetString(key string) (string, error) {
	if !c.IsEnabled() {
		return "", nil
	}
	return c.client.Get(c.ctx, key).Result()
}

// SetString сохраняет строковое значение
func (c *CacheService) SetString(key string, value string, ttl ...time.Duration) error {
	if !c.IsEnabled() {
		return nil
	}
	expiration := c.ttl
	if len(ttl) > 0 && ttl[0] > 0 {
		expiration = ttl[0]
	}
	return c.client.Set(c.ctx, key, value, expiration).Err()
}

// Exists проверяет существование ключа
func (c *CacheService) Exists(key string) (bool, error) {
	if !c.IsEnabled() {
		return false, nil
	}
	result, err := c.client.Exists(c.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// Close закрывает соединение с Redis
func (c *CacheService) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
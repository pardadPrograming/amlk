package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache[T any] struct {
	client *redis.Client
	prefix string
}

func NewRedisCache[T any](client *redis.Client, prefix string) *RedisCache[T] {
	return &RedisCache[T]{
		client: client,
		prefix: prefix,
	}
}

func (c *RedisCache[T]) Get(ctx context.Context, key string) (T, bool) {
	var zero T
	if c.client == nil {
		return zero, false
	}
	value, err := c.client.Get(ctx, c.key(key)).Bytes()
	if err != nil {
		return zero, false
	}
	var out T
	if err := json.Unmarshal(value, &out); err != nil {
		_ = c.client.Del(ctx, c.key(key)).Err()
		return zero, false
	}
	return out, true
}

func (c *RedisCache[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) {
	if c.client == nil {
		return
	}
	body, err := json.Marshal(value)
	if err != nil {
		return
	}
	_ = c.client.Set(ctx, c.key(key), body, ttl).Err()
}

func (c *RedisCache[T]) Delete(ctx context.Context, key string) {
	if c.client == nil {
		return
	}
	_ = c.client.Del(ctx, c.key(key)).Err()
}

func (c *RedisCache[T]) key(key string) string {
	if c.prefix == "" {
		return key
	}
	return c.prefix + ":" + key
}

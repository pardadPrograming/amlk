package cache

import (
	"context"
	"sync"
	"time"
)

type Cache[T any] interface {
	Get(ctx context.Context, key string) (T, bool)
	Set(ctx context.Context, key string, value T, ttl time.Duration)
	Delete(ctx context.Context, key string)
}

type item[T any] struct {
	value     T
	expiresAt time.Time
}

type TTLCache[T any] struct {
	mu    sync.RWMutex
	items map[string]item[T]
	now   func() time.Time
}

func NewTTLCache[T any]() *TTLCache[T] {
	return &TTLCache[T]{
		items: map[string]item[T]{},
		now:   time.Now,
	}
}

func (c *TTLCache[T]) Get(ctx context.Context, key string) (T, bool) {
	c.mu.RLock()
	entry, ok := c.items[key]
	c.mu.RUnlock()
	var zero T
	if !ok {
		return zero, false
	}
	if !entry.expiresAt.IsZero() && c.now().After(entry.expiresAt) {
		c.Delete(ctx, key)
		return zero, false
	}
	return entry.value, true
}

func (c *TTLCache[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) {
	expiresAt := time.Time{}
	if ttl > 0 {
		expiresAt = c.now().Add(ttl)
	}
	c.mu.Lock()
	c.items[key] = item[T]{value: value, expiresAt: expiresAt}
	c.mu.Unlock()
}

func (c *TTLCache[T]) Delete(ctx context.Context, key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

func (c *TTLCache[T]) SweepExpired() int {
	now := c.now()
	removed := 0
	c.mu.Lock()
	for key, entry := range c.items {
		if !entry.expiresAt.IsZero() && now.After(entry.expiresAt) {
			delete(c.items, key)
			removed++
		}
	}
	c.mu.Unlock()
	return removed
}

func (c *TTLCache[T]) StartJanitor(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = time.Minute
	}
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.SweepExpired()
			}
		}
	}()
}

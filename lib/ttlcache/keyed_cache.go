package ttlcache

import (
	"sync"
	"time"
)

type CachedItem[T any] struct {
	Value T
	SetAt time.Time
}

type KeyedCache[T any] struct {
	ttl time.Duration

	items map[string]CachedItem[T]
	lock  sync.RWMutex
}

func NewKeyed[T any](ttl time.Duration) *KeyedCache[T] {
	return &KeyedCache[T]{
		ttl: ttl,

		items: map[string]CachedItem[T]{},
	}
}

func (c *KeyedCache[T]) Get(key string) (item CachedItem[T], ok bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	item, ok = c.items[key]
	return
}

func (c *KeyedCache[T]) Set(key string, value T) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.items[key] = CachedItem[T]{
		Value: value,
		SetAt: time.Now(),
	}
}

func (c *KeyedCache[T]) GetOrDo(key string, fn func() T) T {
	if item, ok := c.Get(key); ok {
		if time.Since(item.SetAt) < c.ttl {
			return item.Value
		}
	}

	value := fn()

	c.Set(key, value)

	return value
}

func (c *KeyedCache[T]) GetOrDoE(key string, fn func() (T, error)) (T, error) {
	if item, ok := c.Get(key); ok {
		if time.Since(item.SetAt) < c.ttl {
			return item.Value, nil
		}
	}

	value, err := fn()
	if err != nil {
		return value, err
	}

	c.Set(key, value)

	return value, nil
}

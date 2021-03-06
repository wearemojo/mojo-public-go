package ttlcache

import (
	"sync"
	"time"
)

type CachedItem[T any] struct {
	Value T
	SetAt time.Time
}

type KeyedCache[TKey comparable, TVal any] struct {
	ttl time.Duration

	items map[TKey]CachedItem[TVal]
	lock  sync.RWMutex
}

func NewKeyed[TKey comparable, TVal any](ttl time.Duration) *KeyedCache[TKey, TVal] {
	return &KeyedCache[TKey, TVal]{
		ttl: ttl,

		items: map[TKey]CachedItem[TVal]{},
	}
}

func (c *KeyedCache[TKey, TVal]) Get(key TKey) (item CachedItem[TVal], ok bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	item, ok = c.items[key]
	return
}

func (c *KeyedCache[TKey, TVal]) Set(key TKey, value TVal) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.items[key] = CachedItem[TVal]{
		Value: value,
		SetAt: time.Now(),
	}
}

func (c *KeyedCache[TKey, TVal]) GetOrDo(key TKey, fn func() TVal) TVal {
	if item, ok := c.Get(key); ok {
		if time.Since(item.SetAt) < c.ttl {
			return item.Value
		}
	}

	value := fn()

	c.Set(key, value)

	return value
}

func (c *KeyedCache[TKey, TVal]) GetOrDoE(key TKey, fn func() (TVal, error)) (TVal, error) {
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

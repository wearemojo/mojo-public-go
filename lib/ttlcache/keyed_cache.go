package ttlcache

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

const TTLForever = -1

type CachedItem[T any] struct {
	Value T
	SetAt time.Time
}

type KeyedCache[TKey ~string, TVal any] struct {
	ttl time.Duration

	sf    singleflight.Group
	items map[TKey]CachedItem[TVal]
	lock  sync.RWMutex
}

// a TTL of -1 means that items never expire
func NewKeyed[TKey ~string, TVal any](ttl time.Duration) *KeyedCache[TKey, TVal] {
	return &KeyedCache[TKey, TVal]{
		ttl: ttl,

		items: map[TKey]CachedItem[TVal]{},
	}
}

func (c *KeyedCache[TKey, TVal]) TTL() time.Duration {
	return c.ttl
}

func (c *KeyedCache[TKey, TVal]) Get(key TKey) (item CachedItem[TVal], ok bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	item, ok = c.items[key]
	return item, ok
}

func (c *KeyedCache[TKey, TVal]) Set(key TKey, value TVal) {
	now := time.Now()

	c.lock.Lock()
	defer c.lock.Unlock()

	c.items[key] = CachedItem[TVal]{
		Value: value,
		SetAt: now,
	}
}

func (c *KeyedCache[TKey, TVal]) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.items = map[TKey]CachedItem[TVal]{}
}

func (c *KeyedCache[TKey, TVal]) GetOrDo(key TKey, fn func() TVal) TVal {
	value, _ := c.GetOrDoE(key, func() (TVal, error) {
		return fn(), nil
	})
	return value
}

func (c *KeyedCache[TKey, TVal]) GetOrDoE(key TKey, fn func() (TVal, error)) (TVal, error) {
	if item, ok := c.Get(key); ok && (c.ttl == TTLForever || time.Since(item.SetAt) < c.ttl) {
		return item.Value, nil
	}

	var value TVal
	var ok bool

	valueRaw, err, _ := c.sf.Do(string(key), func() (any, error) {
		return fn()
	})
	if err != nil {
		return value, err
	}
	value, ok = valueRaw.(TVal)
	if !ok {
		panic(fmt.Sprintf("expected value of type %T, got %T", value, valueRaw))
	}

	c.Set(key, value)

	return value, nil
}

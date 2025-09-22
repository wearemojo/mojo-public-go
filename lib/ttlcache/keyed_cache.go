package ttlcache

import (
	"fmt"
	"sync"
	"time"

	"github.com/wearemojo/mojo-public-go/lib/slicefn"
)

const TTLForever = -1

type CachedItem[T any] struct {
	Value T
	SetAt time.Time
}

type KeyedCache[TKey comparable, TVal any] struct {
	ttl time.Duration

	items map[TKey]CachedItem[TVal]
	lock  sync.RWMutex
}

// a TTL of -1 means that items never expire
func NewKeyed[TKey comparable, TVal any](ttl time.Duration) *KeyedCache[TKey, TVal] {
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

func (c *KeyedCache[TKey, TVal]) GetMany(keys []TKey) map[TKey]CachedItem[TVal] {
	res := make(map[TKey]CachedItem[TVal], len(keys))

	c.lock.RLock()
	defer c.lock.RUnlock()

	for _, key := range keys {
		if item, ok := c.items[key]; ok {
			res[key] = item
		}
	}

	return res
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

func (c *KeyedCache[TKey, TVal]) SetMany(items map[TKey]TVal) {
	now := time.Now()

	c.lock.Lock()
	defer c.lock.Unlock()

	for key, item := range items {
		c.items[key] = CachedItem[TVal]{
			Value: item,
			SetAt: now,
		}
	}
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

	value, err := fn()
	if err != nil {
		return value, err
	}

	c.Set(key, value)

	return value, nil
}

func (c *KeyedCache[TKey, TVal]) GetOrDoMany(keys []TKey, fn func([]TKey) []TVal) []TVal {
	res, _ := c.GetOrDoManyE(keys, func(keys []TKey) ([]TVal, error) {
		return fn(keys), nil
	})
	return res
}

func (c *KeyedCache[TKey, TVal]) GetOrDoManyE(keys []TKey, fn func([]TKey) ([]TVal, error)) ([]TVal, error) {
	// although `GetMany` and `SetMany` could be used in this function, we can
	// reduce the number of loops by doing everything together in a single loop

	res := make([]TVal, len(keys))
	missingKeyIndices := make([]int, 0, len(keys))

	(func() {
		c.lock.RLock()
		defer c.lock.RUnlock()

		now := time.Now()

		for idx, key := range keys {
			if item, ok := c.items[key]; ok && (c.ttl == TTLForever || now.Sub(item.SetAt) < c.ttl) {
				res[idx] = item.Value
			} else {
				missingKeyIndices = append(missingKeyIndices, idx)
			}
		}
	})()

	if len(missingKeyIndices) == 0 {
		return res, nil
	}

	missingKeys := slicefn.Map(missingKeyIndices, func(idx int) TKey { return keys[idx] })

	newItems, err := fn(missingKeys)
	if err != nil {
		return res, err
	} else if len(newItems) != len(missingKeys) {
		panic(fmt.Sprintf("fn returned %d items, expected %d", len(newItems), len(missingKeys)))
	}

	now := time.Now()

	c.lock.Lock()
	defer c.lock.Unlock()

	for idx, key := range missingKeys {
		newItem := newItems[idx]

		c.items[key] = CachedItem[TVal]{
			Value: newItem,
			SetAt: now,
		}

		res[missingKeyIndices[idx]] = newItem
	}

	return res, nil
}

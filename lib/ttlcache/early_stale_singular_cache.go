package ttlcache

import (
	"context"
	"fmt"
	"time"

	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/mlog"
)

type EarlyStaleSingularCache[T any] struct {
	cache *SingularCache[T]

	earlyStaleness time.Duration
}

// a TTL of -1 means that items never expire
func NewEarlyStaleSingular[T any](ttl, earlyStaleness time.Duration) *EarlyStaleSingularCache[T] {
	if ttl == TTLForever {
		panic("ttl cannot be Forever for EarlyStaleSingularCache")
	}
	return &EarlyStaleSingularCache[T]{
		cache:          NewSingular[T](ttl),
		earlyStaleness: earlyStaleness,
	}
}

func (c *EarlyStaleSingularCache[T]) TTL() time.Duration {
	return c.cache.TTL()
}

func (c *EarlyStaleSingularCache[T]) Get() (item CachedItem[T], ok bool) {
	return c.cache.Get()
}

func (c *EarlyStaleSingularCache[T]) Set(value T) {
	c.cache.Set(value)
}

func (c *EarlyStaleSingularCache[T]) Clear() {
	c.cache.Clear()
}

func (c *EarlyStaleSingularCache[T]) GetOrDo(ctx context.Context, fn func() T) T {
	val, _ := c.GetOrDoE(ctx, func() (T, error) {
		return fn(), nil
	})
	return val
}

func (c *EarlyStaleSingularCache[T]) GetOrDoE(ctx context.Context, fn func() (T, error)) (T, error) {
	return c.get(ctx, fn)
}

func (c *EarlyStaleSingularCache[T]) get(ctx context.Context, fn func() (T, error)) (T, error) {
	item, ok := c.cache.Get()
	if !ok || time.Since(item.SetAt) > c.cache.TTL() {
		// If the item is not in the cache or is already older than the TTL, then refresh synchronously and return the new value.
		return c.cache.GetOrDoE(fn)
	}

	if time.Since(item.SetAt) > c.earlyStaleness {
		// If the item is older than the early staleness threshold, refresh in the background and return the stale value.
		go func() {
			valueRaw, err, _ := c.cache.cache.sf.Do(string(singularCacheKey), func() (any, error) {
				return fn()
			})
			if err != nil {
				mlog.Warn(ctx, merr.New(ctx, "failed_to_refresh_early_stale_singular_cache_in_background", nil, err))
				return
			}
			value, ok := valueRaw.(T)
			if !ok {
				panic(fmt.Sprintf("expected value of type %T, got %T", value, valueRaw))
			}
			c.cache.Set(value)
		}()
	}
	return item.Value, nil
}

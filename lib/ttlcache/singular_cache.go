package ttlcache

import (
	"time"
)

var singularCacheKey = struct{}{}

type SingularCache[T any] struct {
	cache *KeyedCache[struct{}, T]
}

func NewSingular[T any](ttl time.Duration) *SingularCache[T] {
	return &SingularCache[T]{
		cache: NewKeyed[struct{}, T](ttl),
	}
}

func (c *SingularCache[T]) Get() (item CachedItem[T], ok bool) {
	return c.cache.Get(singularCacheKey)
}

func (c *SingularCache[T]) Set(value T) {
	c.cache.Set(singularCacheKey, value)
}

func (c *SingularCache[T]) GetOrDo(fn func() T) T {
	return c.cache.GetOrDo(singularCacheKey, fn)
}

func (c *SingularCache[T]) GetOrDoE(fn func() (T, error)) (T, error) {
	return c.cache.GetOrDoE(singularCacheKey, fn)
}

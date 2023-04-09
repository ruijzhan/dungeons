package cache

import (
	"time"

	"github.com/bluele/gcache"
)

type Cache interface {
	Get(key any) (any, error)
	Set(key any, value any) error
}

var _ Cache = (*ttlCache)(nil)

type ttlCache struct {
	cache Cache
}

func NewTTL() Cache {
	gc := gcache.New(100).LRU().Expiration(time.Minute).Build()
	return newTTL(gc)
}

func newTTL(cache Cache) Cache {
	return &ttlCache{cache: cache}
}

func (t *ttlCache) Get(key any) (any, error) {
	return t.cache.Get(key)
}

func (t *ttlCache) Set(key any, value any) error {
	return t.cache.Set(key, value)
}

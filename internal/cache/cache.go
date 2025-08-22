package cache

import (
	"fmt"

	"github.com/dgraph-io/ristretto/v2"
)

type Cache struct {
	cache *ristretto.Cache[string, any]
}

func (c *Cache) Set(key string, value any) {
	c.cache.Set(key, value, 0)
}

func (c *Cache) Get(key string) (any, bool) {
	return c.cache.Get(key)
}

func New() (*Cache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, any]{
		// using default options from documentation
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}

	return &Cache{cache}, nil
}

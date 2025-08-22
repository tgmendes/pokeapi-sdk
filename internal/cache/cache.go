package cache

import (
	"fmt"

	"github.com/dgraph-io/ristretto/v2"
)

type Cache struct {
	cache *ristretto.Cache[string, []byte]
}

func New() (*Cache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, []byte]{
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

func (c *Cache) Set(key string, value []byte) {
	c.cache.Set(key, value, 0)
	// wait for all writes to finish
	c.cache.Wait()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	return c.cache.Get(key)
}

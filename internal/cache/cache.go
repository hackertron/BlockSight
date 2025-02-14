package cache

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type BlockCache struct {
	mu     sync.RWMutex
	cache  map[common.Hash]time.Time
	maxAge time.Duration
}

func NewBlockCache(maxAge time.Duration) *BlockCache {
	cache := &BlockCache{
		cache:  make(map[common.Hash]time.Time),
		maxAge: maxAge,
	}

	// Start cleanup routine
	go cache.cleanup()

	return cache
}

func (c *BlockCache) Add(hash common.Hash) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[hash] = time.Now()
}

func (c *BlockCache) Exists(hash common.Hash) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.cache[hash]
	return exists
}

func (c *BlockCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for hash, timestamp := range c.cache {
			if now.Sub(timestamp) > c.maxAge {
				delete(c.cache, hash)
			}
		}
		c.mu.Unlock()
	}
}

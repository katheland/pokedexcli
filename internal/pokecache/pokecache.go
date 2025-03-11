package pokecache

import (
	"sync"
	"time"
)

// the structure of the cache
type Cache struct {
	entries map[string]cacheEntry
	mu sync.Mutex
	interval time.Duration
}
type cacheEntry struct {
	createdAt time.Time
	val []byte
}

// create a new cache
func NewCache(i time.Duration) Cache {
	c := Cache{
		entries: make(map[string]cacheEntry),
		interval: i,
	}
	go c.reapLoop()
	return c
}

// add an entry to the cache
func (c *Cache) Add(k string, v []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := cacheEntry{
		createdAt: time.Now(),
		val: v,
	}
	c.entries[k] = entry
}

// get an entry from the cache
func (c *Cache) Get(k string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, exists := c.entries[k]
	if exists {
		return entry.val, exists
	} else {
		return []byte{}, exists
	}
}

// every interval amount of time, clean out of the cache anything that's older than that interval
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for k, v := range c.entries {
				if time.Now().Sub(v.createdAt) > c.interval {
					delete(c.entries, k)
				}
			}
			c.mu.Unlock()
			ticker.Reset(c.interval)
		}
	}
}
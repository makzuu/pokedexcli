package pokecache

import (
	"sync"
	"time"
)

type cache struct {
	en map[string]cacheEntry
	mu sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) cache {
	c := cache{en: make(map[string]cacheEntry)}
	go c.reapLoop(interval)
	return c
}

func (c cache) Add(key string, val []byte) {
	entry := cacheEntry{createdAt: time.Now(), val: val}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.en[key] = entry
}

func (c cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.en[key]
	return entry.val, ok
}

func (c cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for t := range ticker.C {
		c.mu.Lock()
		for key, entry := range c.en {
			if t.Sub(entry.createdAt) > interval {
				delete(c.en, key)
			}
		}
		c.mu.Unlock()
	}
}

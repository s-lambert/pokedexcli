package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	internalCache map[string]cacheEntry
	mut           sync.Mutex
}

func (c *Cache) reapLoop() {
	go func() {
		const interval = 5 * time.Second
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				now := time.Now().Format(time.RFC3339)
				c.mut.Lock()
				for k, e := range c.internalCache {
					if e.createdAt.Add(interval).Format(time.RFC3339) < now {
						delete(c.internalCache, k)
					}
				}
				c.mut.Unlock()
			}
		}
	}()
}

func NewCache() *Cache {
	newCache := &Cache{
		internalCache: make(map[string]cacheEntry),
		mut:           sync.Mutex{},
	}
	newCache.reapLoop()
	return newCache
}

func (c *Cache) Add(key string, val []byte) {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.internalCache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mut.Lock()
	defer c.mut.Unlock()
	e, ok := c.internalCache[key]
	if !ok {
		return nil, ok
	}
	return e.val, ok
}

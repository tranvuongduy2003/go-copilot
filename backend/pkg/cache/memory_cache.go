package cache

import (
	"context"
	"path/filepath"
	"sync"
	"time"
)

type cacheEntry struct {
	value     []byte
	expiresAt time.Time
}

func (e *cacheEntry) isExpired() bool {
	if e.expiresAt.IsZero() {
		return false
	}
	return time.Now().After(e.expiresAt)
}

type MemoryCache struct {
	data          map[string]*cacheEntry
	mutex         sync.RWMutex
	cleanupTicker *time.Ticker
	stopCleanup   chan struct{}
}

func NewMemoryCache(cleanupInterval time.Duration) *MemoryCache {
	cache := &MemoryCache{
		data:        make(map[string]*cacheEntry),
		stopCleanup: make(chan struct{}),
	}

	if cleanupInterval > 0 {
		cache.cleanupTicker = time.NewTicker(cleanupInterval)
		go cache.cleanup()
	}

	return cache
}

func (c *MemoryCache) cleanup() {
	for {
		select {
		case <-c.stopCleanup:
			return
		case <-c.cleanupTicker.C:
			c.deleteExpired()
		}
	}
}

func (c *MemoryCache) deleteExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for key, entry := range c.data {
		if entry.isExpired() {
			delete(c.data, key)
		}
	}
}

func (c *MemoryCache) Close() {
	if c.cleanupTicker != nil {
		c.cleanupTicker.Stop()
	}
	close(c.stopCleanup)
}

func (c *MemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	c.mutex.RLock()
	entry, exists := c.data[key]
	c.mutex.RUnlock()

	if !exists {
		return nil, ErrCacheMiss
	}

	if entry.isExpired() {
		c.mutex.Lock()
		delete(c.data, key)
		c.mutex.Unlock()
		return nil, ErrCacheExpired
	}

	return entry.value, nil
}

func (c *MemoryCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry := &cacheEntry{
		value: make([]byte, len(value)),
	}
	copy(entry.value, value)

	if ttl > 0 {
		entry.expiresAt = time.Now().Add(ttl)
	}

	c.data[key] = entry
	return nil
}

func (c *MemoryCache) Delete(ctx context.Context, keys ...string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, key := range keys {
		delete(c.data, key)
	}
	return nil
}

func (c *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	c.mutex.RLock()
	entry, exists := c.data[key]
	c.mutex.RUnlock()

	if !exists {
		return false, nil
	}

	if entry.isExpired() {
		c.mutex.Lock()
		delete(c.data, key)
		c.mutex.Unlock()
		return false, nil
	}

	return true, nil
}

func (c *MemoryCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	c.mutex.RLock()
	entry, exists := c.data[key]
	c.mutex.RUnlock()

	if !exists {
		return -2, nil
	}

	if entry.expiresAt.IsZero() {
		return -1, nil
	}

	remaining := time.Until(entry.expiresAt)
	if remaining < 0 {
		return -2, nil
	}

	return remaining, nil
}

func (c *MemoryCache) Clear(ctx context.Context, pattern string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if pattern == "*" || pattern == "" {
		c.data = make(map[string]*cacheEntry)
		return nil
	}

	for key := range c.data {
		matched, err := filepath.Match(pattern, key)
		if err != nil {
			continue
		}
		if matched {
			delete(c.data, key)
		}
	}

	return nil
}

func (c *MemoryCache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.data)
}

// Package cache provides a generic LRU cache implementation with TTL support and centralized cache manager.
package cache

import (
	"container/list"
	"sync"
	"time"
)

// CacheEntry stores a single value with its expiration.
type CacheEntry[V any] struct {
	Value      V
	Expiry     time.Time
	AccessTime time.Time
	HitCount   int64
}

// LRUCache is a thread-safe generic LRU cache.
type LRUCache[K comparable, V any] struct {
	maxEntries    int
	defaultTTL    time.Duration
	mutex         sync.RWMutex
	cache         map[K]*list.Element
	lruList       *list.List
	cleanupTicker *time.Ticker
	cleanupStopCh chan struct{}
}

type entry[K comparable, V any] struct {
	Key   K
	Entry CacheEntry[V]
}

// NewLRUCache creates a new generic LRU cache with optional background cleanup.
func NewLRUCache[K comparable, V any](maxEntries int, defaultTTL time.Duration) *LRUCache[K, V] {
	if maxEntries <= 0 {
		panic("maxEntries must be > 0")
	}
	c := &LRUCache[K, V]{
		maxEntries:    maxEntries,
		defaultTTL:    defaultTTL,
		cache:         make(map[K]*list.Element),
		lruList:       list.New(),
		cleanupTicker: time.NewTicker(5 * time.Minute),
		cleanupStopCh: make(chan struct{}),
	}
	go c.cleanupRoutine()
	return c
}

// StopCleanup stops the background cleanup routine.
func (c *LRUCache[K, V]) StopCleanup() {
	close(c.cleanupStopCh)
	c.cleanupTicker.Stop()
}

// cleanupRoutine periodically removes expired entries.
func (c *LRUCache[K, V]) cleanupRoutine() {
	for {
		select {
		case <-c.cleanupTicker.C:
			c.cleanupExpired()
		case <-c.cleanupStopCh:
			return
		}
	}
}

// cleanupExpired removes all expired entries.
func (c *LRUCache[K, V]) cleanupExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, elem := range c.cache {
		ent := elem.Value.(*entry[K, V])
		if time.Now().After(ent.Entry.Expiry) {
			c.removeElement(elem)
		}
	}
}

// Get retrieves a value from the cache.
func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	var zero V
	c.mutex.Lock()
	defer c.mutex.Unlock()

	elem, ok := c.cache[key]
	if !ok {
		return zero, false
	}

	ent := elem.Value.(*entry[K, V])
	if time.Now().After(ent.Entry.Expiry) {
		c.removeElement(elem)
		return zero, false
	}

	ent.Entry.HitCount++
	ent.Entry.AccessTime = time.Now()
	c.lruList.MoveToFront(elem)
	return ent.Entry.Value, true
}

// BatchGet retrieves multiple values from the cache.
// It returns a map of keys to values that were found and not expired.
func (c *LRUCache[K, V]) BatchGet(keys []K) map[K]V {
	result := make(map[K]V)
	now := time.Now()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, key := range keys {
		elem, ok := c.cache[key]
		if !ok {
			continue
		}
		ent := elem.Value.(*entry[K, V])
		if now.After(ent.Entry.Expiry) {
			c.removeElement(elem)
			continue
		}
		ent.Entry.HitCount++
		ent.Entry.AccessTime = now
		c.lruList.MoveToFront(elem)
		result[key] = ent.Entry.Value
	}
	return result
}

// Set adds or updates a value in the cache.
func (c *LRUCache[K, V]) Set(key K, value V, ttl ...time.Duration) {
	expiry := time.Now().Add(c.defaultTTL)
	if len(ttl) > 0 && ttl[0] > 0 {
		expiry = time.Now().Add(ttl[0])
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, exists := c.cache[key]; exists {
		ent := elem.Value.(*entry[K, V])
		ent.Entry.Value = value
		ent.Entry.Expiry = expiry
		ent.Entry.AccessTime = time.Now()
		c.lruList.MoveToFront(elem)
		return
	}

	if c.lruList.Len() >= c.maxEntries {
		c.removeOldest()
	}

	newEntry := &entry[K, V]{
		Key: key,
		Entry: CacheEntry[V]{
			Value:      value,
			Expiry:     expiry,
			AccessTime: time.Now(),
			HitCount:   0,
		},
	}
	elem := c.lruList.PushFront(newEntry)
	c.cache[key] = elem
}

// BatchSet sets multiple key-value pairs in the cache.
// If TTL is specified, it will be applied to all entries.
func (c *LRUCache[K, V]) BatchSet(entries map[K]V, ttl ...time.Duration) {
	expiry := time.Now().Add(c.defaultTTL)
	if len(ttl) > 0 && ttl[0] > 0 {
		expiry = time.Now().Add(ttl[0])
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for key, value := range entries {
		if elem, exists := c.cache[key]; exists {
			ent := elem.Value.(*entry[K, V])
			ent.Entry.Value = value
			ent.Entry.Expiry = expiry
			ent.Entry.AccessTime = time.Now()
			c.lruList.MoveToFront(elem)
			continue
		}

		if c.lruList.Len() >= c.maxEntries {
			c.removeOldest()
		}

		newEntry := &entry[K, V]{
			Key: key,
			Entry: CacheEntry[V]{
				Value:      value,
				Expiry:     expiry,
				AccessTime: time.Now(),
				HitCount:   0,
			},
		}
		elem := c.lruList.PushFront(newEntry)
		c.cache[key] = elem
	}
}

// Delete removes a value from the cache.
func (c *LRUCache[K, V]) Delete(key K) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if elem, exists := c.cache[key]; exists {
		c.removeElement(elem)
	}
}

// Clear removes all entries.
func (c *LRUCache[K, V]) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache = make(map[K]*list.Element)
	c.lruList.Init()
}

// Len returns the number of entries in the cache.
func (c *LRUCache[K, V]) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.lruList.Len()
}

// removeElement deletes an element from the list and map.
func (c *LRUCache[K, V]) removeElement(elem *list.Element) {
	ent := elem.Value.(*entry[K, V])
	delete(c.cache, ent.Key)
	c.lruList.Remove(elem)
}

// removeOldest evicts the least recently used item.
func (c *LRUCache[K, V]) removeOldest() {
	if elem := c.lruList.Back(); elem != nil {
		c.removeElement(elem)
	}
}

// CacheManager manages multiple named caches.
type CacheManager struct {
	mutex  sync.RWMutex
	caches map[string]any // map of named caches
}

// NewCacheManager creates a new CacheManager.
func NewCacheManager() *CacheManager {
	return &CacheManager{
		caches: make(map[string]any),
	}
}

// Register registers a named cache instance.
func (m *CacheManager) Register(name string, cache any) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.caches[name] = cache
}

// Get returns a registered cache instance.
func (m *CacheManager) Get(name string) (any, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	cache, ok := m.caches[name]
	return cache, ok
}

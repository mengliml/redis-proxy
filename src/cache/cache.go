package cache

import (
    "time"
    "sync"
    "github.com/golang/groupcache/lru"
)

type CacheEntry struct {
    Value        string
    ExpiredTime  int64
}

type LRUCache struct {
    Cache        *lru.Cache
    GlobalExpiry  int64 
    mutex        *sync.Mutex
}

func NewCache(capacity int, globalExpiry int) *LRUCache {
    return &LRUCache {
        Cache:        lru.New(capacity),
        GlobalExpiry: int64(globalExpiry),
        mutex:        &sync.Mutex{},
    }
}

// Get looks up a key's value from the cache 
// if the key has not expired. Otherwise
// remove from the cache.
func (c *LRUCache) Get(key string) (string, bool) {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    if ce, ok := c.Cache.Get(key); ok {
        cacheEntry := ce.(*CacheEntry)
        if cacheEntry.ExpiredTime < int64(time.Now().UnixNano()) {
            c.Remove(key)
            return "", false
        }

        return cacheEntry.Value, true
    }

    return "", false
}

// Set adds a given entry to the cache and sets the expired time
func (c *LRUCache) Set(key string, value string) {
    now := time.Now()
    duration := time.Millisecond * time.Duration(c.GlobalExpiry)
    expiredTime := int64(now.Add(duration).UnixNano())

    ce := &CacheEntry {
        Value:  value,
        ExpiredTime: expiredTime,
    }

    c.mutex.Lock()
    c.Cache.Add(key, ce)
    c.mutex.Unlock()
}

// Remove cached key and value pair
func (c *LRUCache) Remove(key string) {
    c.Cache.Remove(key) 
}
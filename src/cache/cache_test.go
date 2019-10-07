package cache_test

import (
    "fmt"
    "time"
    "testing"
    "redis-proxy/src/cache"
    "github.com/stretchr/testify/assert"
)

var (
    capacity = 1000
    globalExpiry = 1000
)

func TestNewCache(t *testing.T) {
    c := cache.NewCache(capacity, globalExpiry)

    assert.Equal(t, capacity, c.Cache.MaxEntries, "incorrect capacity")
    assert.Equal(t, globalExpiry, int(c.GlobalExpiry), "incorrect globalExpiry")
}

func TestGetNonExistingKey(t *testing.T) {
    c := cache.NewCache(capacity, globalExpiry)

    value, ok := c.Get("key")
    assert.Equal(t, false, ok, "The cache didn't return failure status code for non-existing key")
    assert.Equal(t, "", value, "The cache didn't return empty value for non-existing key")
}

func TestGetFromCache(t *testing.T) {
    c := cache.NewCache(capacity, globalExpiry)

    c.Set("key", "value")

    value, ok := c.Get("key")
    assert.Equal(t, true, ok, "The cache didn't return value for existing key")
    assert.Equal(t, "value", value, "The cache didn't set the value correctly")
}

func TestRemoveFromCache(t *testing.T) {
    c := cache.NewCache(capacity, globalExpiry)

    c.Set("key", "value")
    c.Remove("key")

    value, ok := c.Get("key")
    assert.Equal(t, false, ok, "The cache returned found for removed key")
    assert.Equal(t, "", value, "The cache didn't return empty value for removed key")
}

func TestGlobalExpiry(t *testing.T) {
    c := cache.NewCache(capacity, globalExpiry)
    c.Set("expiredKey", "value")

    duration := time.Duration(2*globalExpiry)
    time.Sleep(duration * time.Millisecond)
    value, ok := c.Get("expiredKey")

    assert.Equal(t, false, ok, "The cache returned found for expired key")
    assert.Equal(t, "", value, "The cache didn't return empty value for expired key")
}

func TestLRUCacheEviction(t *testing.T) {
    c := cache.NewCache(capacity, globalExpiry)

    c.Set("evictedKey", "value")
    for i := 0; i < capacity; i++ {
        key := fmt.Sprintf("key%d", i)
        c.Set(key, "value")
    }

    value, ok := c.Get("evictedKey")
    assert.Equal(t, false, ok, "The cache returned found for evicted key")
    assert.Equal(t, "", value, "The cache didn't return empty value for evicted key")
}
package main

import (
    "fmt"
    "testing"
    "time"
    "redis-proxy/src/proxy"
    "redis-proxy/src/util"
    "github.com/go-redis/redis"
    "github.com/stretchr/testify/assert"
)

var (
    testRedisAddress = "redis:6379"
    testGlobalExpiry = 100
    testCapacity = 10
    testPort = 8080
    testMaxClients = 10
)

// Test whether the proxy can correctly get value from redis client
func TestGetFromRedis(t *testing.T) {
    p, _ := proxy.NewProxy(testRedisAddress, testGlobalExpiry, testCapacity, testMaxClients)
    s := proxy.NewServer(testPort, p)

    // Add (key, value) pair to the redis client
    redis := createRedisClient(t)
    redis.Set("key", "value", 0)

    // Get value from the proxy
    body, _ := util.GetValueFromServer("/GetValue/key", s, t)

    assert.Equal(t, "value", string(body), "The proxy can't get correct value from redis client")
}

// Test whether the cache can correctly cache data 
func TestGetFromCache(t *testing.T) {
    p, _ := proxy.NewProxy(testRedisAddress, testGlobalExpiry, testCapacity, testMaxClients)
    s := proxy.NewServer(testPort, p)

    // Add (key, value) pair to the redis client
    redis := createRedisClient(t)
    redis.Set("key", "value", 0)

    // Get value from the proxy
    util.GetValueFromServer("/GetValue/key", s, t)

    // Update the value in the redis client
    redis.Set("key", "newValue", 0)

    body, _ := util.GetValueFromServer("/GetValue/key", s, t)
    assert.Equal(t, "value", string(body), "The proxy can't get correct value from the cache")

    _, statusCode := util.GetValueFromServer("/GET/nonexistingkey", s, t)
    assert.Equal(t, 404, statusCode, "The server didn't return 404 error for non-existing key")
}

// Test whether the cache can expire the entry
func TestCacheGlobalExpiry(t *testing.T) {
    p, _ := proxy.NewProxy(testRedisAddress, testGlobalExpiry, testCapacity, testMaxClients)
    s := proxy.NewServer(testPort, p)

    // Add (key, value) pair to the redis client
    redis := createRedisClient(t)
    redis.Set("key", "value", 0)

    // Get value from the proxy
    util.GetValueFromServer("/GetValue/key", s, t)

    // Update the value in the redis client
    redis.Set("key", "newValue", 0)

    // Sleep for 2 testGlobalExpiry milliseconds and check if the entry is expired
    duration := time.Duration(2 * testGlobalExpiry)
    time.Sleep(duration * time.Millisecond)
    body, _ := util.GetValueFromServer("/GetValue/key", s, t)

    assert.NotEqual(t, "value", string(body), "The cache didn't remove the expired entry")
    assert.Equal(t, "newValue", string(body), "The proxy didn't retrieve the expired entry from redis")
}

// Test whether the cache can evict the least recently used entry
// Concurrent request to the server
func TestLRUCacheEvictionWithConcurrentClients(t *testing.T) {
    p, _ := proxy.NewProxy(testRedisAddress, testGlobalExpiry, testCapacity, testMaxClients)
    s := proxy.NewServer(testPort, p)

    // Add (testCapacity + 1) pairs of (key, value) to the redis client
    redis := createRedisClient(t)
    for i := 0; i < testCapacity + 1; i++ {
        key := fmt.Sprintf("key%d", i)
        value := fmt.Sprintf("value%d", i)
        redis.Set(key, value, 0)
    }

    // Cache the date entry when get requested 
    for i := 0; i < testCapacity + 1; i++ {
        key_path := fmt.Sprintf("/GetValue/key%d", i)
        util.GetValueFromServer(key_path, s, t)
    }

    // Add new value to the redis
    key := fmt.Sprintf("key%d", 0)
    value := fmt.Sprintf("newValue%d", 0)
    redis.Set(key, value, 0)

    // Get the value from the proxy
    evicted_key_path := fmt.Sprintf("/GetValue/key%d", 0)
    body, _ := util.GetValueFromServer(evicted_key_path, s, t)

    assert.NotEqual(t, "value0", string(body), "The cache didn't remove LRU entry upon reaching testCapacity")
    assert.Equal(t, "newValue0", string(body), "The proxy didn't properly retrieve evicted entry from redis")
}

func createRedisClient(t *testing.T) *redis.Client {
    redis := redis.NewClient(&redis.Options{
        Addr: testRedisAddress,
    })

    _, err := redis.Ping().Result()
    if err != nil {
        t.Fatal(err)
    }

    return redis
}
package redis_test

import (
    "testing"
    "redis-proxy/src/redis"
    "github.com/stretchr/testify/assert"
)

var (
    addr = "redis:6379"
)

func TestNewClient(t *testing.T) {
    _, err := redis.NewClient(addr)
    if err != nil {
        t.Fail()
    }
}

func TestGetKey(t *testing.T) {
    r, _ := redis.NewClient(addr)

    r.Client.Set("key1", "value1", 0)

    value, err := r.Get("key1")
    if err != nil {
        t.Fail()
    }
    assert.Equal(t, "value1", value, "return incorrect value from redis")
}

func TestGetNonExistingKey(t *testing.T) {
    r, _ := redis.NewClient(addr)

    r.Client.Set("key1", "value1", 0)

    _, err := r.Get("key2")
    if _, ok := err.(*redis.KeyNotFoundError); ok {
        t.Fail()
    }
}


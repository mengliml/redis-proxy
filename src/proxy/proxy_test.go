package proxy_test

import (
    "net/http"
    "testing"
    "redis-proxy/src/proxy"
    "redis-proxy/src/util"
    "github.com/stretchr/testify/assert"
)

var (
    redisAddress = "redis:6379"
    globalExpiry = 1000
    capacity = 10
    port = 8080
    maxClients = 10
)

func TestNewProxy(t *testing.T) {
    _, err := proxy.NewProxy(redisAddress, globalExpiry, capacity, maxClients)
    if err != nil {
        t.Fatal(err)
    }
}

func TestServerConnection(t *testing.T) {
    p, _ := proxy.NewProxy(redisAddress, globalExpiry, capacity, maxClients)
    s := proxy.NewServer(port, p)

    _, statusCode := util.GetValueFromServer("/GetValue/key", s, t)

    assert.Equal(t, http.StatusOK, statusCode, "Server failed to start")
}
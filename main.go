package main 

import (
    "fmt"
    "flag"
    "redis-proxy/src/proxy"
)

var redisAddress = flag.String("redis-address", "", "Redis address")
var globalExpiry = flag.Int("global-expiry", 60 * 1000, "Cache expiration in milliseconds")
var capacity = flag.Int("capacity", 1000, "Cache capacity (number of keys)")
var port = flag.Int("port", 8080, "TCP/IP port number the proxy lestens on")
var maxClients = flag.Int("max-clients", 10, "maximum number of concurrent clients can connect")

func main() {
    flag.Parse()

    // initialize proxy
    p, err := proxy.NewProxy(*redisAddress, *globalExpiry, *capacity, *maxClients)
    if err != nil {
        return
    }

    // initialize and start proxy server
    server := proxy.NewServer(*port, p)
    server.ListenAndServe()
    fmt.Println("Server is running")
}   

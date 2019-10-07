package proxy 

import (
    "fmt"
    "strconv"
    "net/http"
    "redis-proxy/src/redis"
    "redis-proxy/src/cache"
    "github.com/gorilla/mux"
)

type Proxy struct {
    redis       *redis.RedisClient
    cache       *cache.LRUCache
    maxClients  int
    semaphore   chan struct{}
}

func NewProxy(redisAddress string, globalExpiry int, capacity int, maxClients int) (*Proxy, error) {
    r, err := redis.NewClient(redisAddress)

    if err != nil {
        return nil, err
    }

    c := cache.NewCache(capacity, globalExpiry)

    return &Proxy{
        redis: r,
        cache: c,
        maxClients: maxClients,
        semaphore: make(chan struct{}, maxClients),
    }, nil
}

func NewServer(port int, p *Proxy) *http.Server {
    router := mux.NewRouter()
    router.HandleFunc("/", p.IndexHandler).Methods("GET")
    router.HandleFunc("/GetValue/{key}", p.CachedGetHandler).Methods("GET")

    server := &http.Server {
        Addr: ":" + strconv.Itoa(port),
        Handler: router,
    }

    return server
}

func (p *Proxy) IndexHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "The redis-proxy server is running.")
}

func (p *Proxy) CachedGetHandler(w http.ResponseWriter, r *http.Request) {
    p.semaphore <- struct{}{}
    defer func() { <-p.semaphore }()
    
    defer r.Body.Close()

    vars := mux.Vars(r)
    key, _ := vars["key"]

    // return the value from the cache, if the cache contains the key
    if value, ok := p.cache.Get(key); ok {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(value))
        return
    }

    // otherwise fetch the value from the redis
    value, err := p.redis.Get(key)
    if err != nil {
        if keyNotFoundError, ok := err.(*redis.KeyNotFoundError); ok {
            w.WriteHeader(http.StatusNotFound)
            w.Write([]byte("The key is not found: " + keyNotFoundError.Key))
            return 
        }
        
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Internal server error"))
        return   
    }

    // store it in the cache
    p.cache.Set(key, value)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(value))
    return    
}

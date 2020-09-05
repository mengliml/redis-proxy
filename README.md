# redis-proxy

## High-level Architecture Overview

**1. Redis**

The single backing Redis instance is instantiated for each instacne of the proxy service. It can be configurable at proxy startup with address of backing redis like `redis:6379`.

**2. LRU Cache**

The cache is implemented as an LRU cache from [`groupcache/lru`](https://github.com/golang/groupcache). It will cache the data requested by clients. Two parameters can be configurable at the proxy startup:

- `GLOBAL_EXPIRY`: Entries added to the proxy cache are expired after being in the cache for a time duration (in milliseconds) that is globally configured (per instance). After an entry is expired, a GET request will act as if the value associated with the key was never stored in the cache. 

- `CAPACITY`: The maximum number of keys it retains. Once the cache fills to capacity, the least recently used (i.e. read) key is evicted each time a new key needs to be added to the cache.

**3. Proxy server**

A GET request like `HOST/{key}`, directed at the proxy server, returns the value of the specified key from the proxy’s local cache if the local cache contains a value for that key. If the local cache does not contain a value for the specified key, it fetches the value from the backing Redis instance, using the Redis `GET` command, and stores it in the local cache, associated with the specified key.

- `GET '/'` returns `'The redis-proxy server is running.'`, if the server is up
- `GET 'GetValue/{key}'` returns the value associated with the secified `key`. Returns from cache if available, otherwise retrieves from the redis instance.

**4. LRU eviction**

Once the cache fills to capacity, the least recently used (i.e. read) key is evicted each time a new key needs to be added to the cache.

**5. Parallel concurrent access** 

In order to provide multiple clients concurrent access to the proxy, the simplest way to achieve that would be to use the `semaphore` with a `buffered channel`. We can easily create a fixed size (maximum number of clients) of buffered channel. Each client request will send a `struct{}` to the channel when the call starts and remove it from the channel once the call returns. And put a `mutex` in front of cache access function ensure that only one goroutine could modify it at a time. However, one limitation to this, if the buffered channel reaches its limit and another client request to connect to the server, it would be blocked until one slot opens up, making it a bottleneck. To eliminate this problem, we can use the cache sharding to reduce the blocking. But it will take more time to implement the solution.  

For the sequential concurrent processing, we can just remove the the `semaphore` and `mutex` from the codebase. It will process the multiple requests in a sequential order. 


## Code Overview

The code for redis proxy is split into four main packages:

**`src/redis/redis.go`**

This instantiate a Redis client using the Redis protocol and provides an in-memory storage for the data. The Redis already has in-built concurrency control.

**`src/cache/cache.go`**

The cache is implemented as an LRU cache from [`groupcache/lru`](https://github.com/golang/groupcache). It is not safe for concurrent access, so all accesses must require a mutex. The `LRUCache` struct stores the pointer to the LRU cache, global expiry and the mutex. 
 
**`src/proxy/proxy.go`**

This file provides two methods to instantiate a redis proxy service. 

- `NewProxy(redisAddress string, globalExpiry int, capacity int, maxClients int) (*Proxy, error)`
Returns a new instance of `Proxy`, with an instance of a redis client and cache. 

- `NewServer(port int, p *Proxy) *http.Server`
Returns an web server.

**`src/util/util.go`**

Provides a helper method which makes a request to the server for testing.

**`main.go`** 

This creates a web server and set up the configurable options passed through the command line.

The following files are for the unit tests. 

**`src/cache/cache_test.go`**

**`src/redis/redis_test.go`**

**`src/proxy/proxy_test.go`**

**`main_test.go`** 

## Algorithmic Complexity of Cache Operations

If the requested `key` is stored in the local cache and it's not expired, we are able to retrieve the value in `O(1)` time from the cache since the `groupcache/lru` library which provides `O(1)` amortized lookup and set value. 

However, if the local cache does not contain a value for the specified key, it fetches the value from the backing Redis instance, using the Redis `GET` command. Redis provides `O(1)` amortized lookup time if all the data fits in memory, or `O(1+n/k)` where n is the number of items and k the number of buckets. 
And storing it associated with the specified key in the local cache cost `O(1)`. 


## How To Use

To clone and run this application, you'll need `make`, `docker`, `docker-compose` and `bash` installed on your computer. 

The following parameters are configurable in the .env file

```
REDIS_ADDRESS=redis:6379
GLOBAL_EXPIRY=600000
CAPACITY=10000
PORT=8080
MAX_CLIENTS=100
```

From your command lin: 

```bash
# Clone this repository
git clone https://github.com/Mengsuper/redis-proxy.git

# Go into the repository
cd redis-proxy

# Build and run unit tests
make test

# Run the app in docker container
make run

# Stop container
make stop

# Help info
make help
```

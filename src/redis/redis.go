package redis

import (
    "fmt"
    "github.com/go-redis/redis"
)

type RedisClient struct {
    Client *redis.Client
}

func NewClient(addr string) (*RedisClient, error) {
    client := redis.NewClient(&redis.Options{
        Addr: addr,
    })

    _, err := client.Ping().Result()

    if err != nil {
        return nil, err
    }

    return &RedisClient{
        Client: client,
    }, nil
}

func (r *RedisClient) Get(key string) (string, error) {
    value, err := r.Client.Get(key).Result()

    if err != nil {
        if err == redis.Nil {
            return "", &KeyNotFoundError{key}
        }
        return "", err
    }

    return value, nil
}

type KeyNotFoundError struct {
    Key string
}

func (e *KeyNotFoundError) Error() string {
    return fmt.Sprintf("key not found %s", e.Key)
}
package app

import (
	. "nova/cache"
	"sync"
)

type Cache struct {
	userCache UserCache
}

type UserCache struct {
	userSet []User
	mutex   sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{}
}

type RedisCache struct {
	redisCache *RedisClient
}

func NewRedisCache() (*RedisCache, error) {
	// set redis configure
	cfg := &RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
	// create redis cache
	client, err := NewRedisClient(cfg)
	if err != nil {
		return nil, err
	}
	return &RedisCache{client}, nil
}

func (cache *RedisCache) Close() error {
	// close redis client
	if err := cache.redisCache.Close(); err != nil {
		return err
	}
	return nil
}

package app

import (
	. "nova/cache"
)

type Cache struct {
	redisCache *RedisClient
}

func NewCache() (*Cache, error) {
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
	return &Cache{client}, nil
}

func (cache *Cache) Close() error {
	// close redis client
	if err := cache.redisCache.Close(); err != nil {
		return err
	}
	return nil
}

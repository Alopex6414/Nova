package app

import (
	"context"
	. "nova/cache"
	"sync"
)

type Cache struct {
	userCache     UserCache
	questionCache QuestionCache
}

type UserCache struct {
	userSet []User
	mutex   sync.RWMutex
}

type QuestionCache struct {
	singleChoiceSet   []QuestionSingleChoice
	multipleChoiceSet []QuestionMultipleChoice
	judgementSet      []QuestionJudgement
	essaySet          []QuestionEssay
	mutex             sync.RWMutex
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

func (cache *RedisCache) CreateUser(user *User) error {
	// create user in redis cache
	if err := cache.redisCache.HSet(context.Background(), user.UserId, user); err != nil {
		return err
	}
	return nil
}

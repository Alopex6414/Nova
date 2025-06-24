package app

import (
	"context"
	. "nova/cache"
	"sync"
)

type Cache struct {
	userCache      UserCache
	questionsCache QuestionsCache
}

type UserCache struct {
	userSet []User
	mutex   sync.RWMutex
}

type QuestionsCache struct {
	singleChoiceCache   QuestionSingleChoiceCache
	multipleChoiceCache QuestionMultipleChoiceCache
	judgementCache      QuestionJudgementCache
	essayCache          QuestionEssayCache
}

type QuestionSingleChoiceCache struct {
	singleChoiceSet []QuestionSingleChoice
	mutex           sync.RWMutex
}

type QuestionMultipleChoiceCache struct {
	multipleChoiceSet []QuestionMultipleChoice
	mutex             sync.RWMutex
}

type QuestionJudgementCache struct {
	judgementSet []QuestionJudgement
	mutex        sync.RWMutex
}

type QuestionEssayCache struct {
	essaySet []QuestionEssay
	mutex    sync.RWMutex
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

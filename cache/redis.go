package cache

import (
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

// RedisClient Redis Client
type RedisClient struct {
	client *redis.Client
	config *Config
	logger log.Logger
}

// Config Redis Configure Options
type Config struct {
	Addr           string
	Password       string
	DB             int
	PoolSize       int
	MinIdleConns   int
	IdleTimeout    time.Duration
	DialTimeout    time.Duration
	CommandTimeout time.Duration
}

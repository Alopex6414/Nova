package cache

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient struct {
	client *redis.Client
	config *RedisConfig
}

type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	MaxRetries   int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewRedisClient(cfg *RedisConfig) (*RedisClient, error) {
	// get configure address and max retry times
	if cfg.Addr == "" {
		cfg.Addr = "localhost:6379"
	}
	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = 3
	}
	// create redis client
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     100,
		MinIdleConns: 10,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		MaxRetries:   cfg.MaxRetries,
	})
	// context without timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %v", err)
	}
	// return redis client
	return &RedisClient{
		client: client,
		config: cfg,
	}, nil
}

func (rc *RedisClient) Close() error {
	if rc.client != nil {
		return rc.client.Close()
	}
	return nil
}

// Set key pair
func (rc *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rc.client.Set(ctx, key, value, expiration).Err()
}

// Get key pair
func (rc *RedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := rc.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("key %s does not exist", key)
	}
	return val, err
}

// Exists key pair
func (rc *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	exist, err := rc.client.Exists(ctx, key).Result()
	return exist == 1, err
}

// Expire key pair
func (rc *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return rc.client.Expire(ctx, key, expiration).Err()
}

// Delete key pair
func (rc *RedisClient) Delete(ctx context.Context, keys ...string) (int64, error) {
	return rc.client.Del(ctx, keys...).Result()
}

// Pipeline pipeline operation
func (rc *RedisClient) Pipeline(ctx context.Context, fn func(pipe redis.Pipeliner) error) error {
	_, err := rc.client.Pipelined(ctx, fn)
	return err
}

// HSet set hash
func (rc *RedisClient) HSet(ctx context.Context, key string, fieldValues ...interface{}) error {
	return rc.client.HSet(ctx, key, fieldValues...).Err()
}

// HGet get hash
func (rc *RedisClient) HGet(ctx context.Context, key, field string) (string, error) {
	val, err := rc.client.HGet(ctx, key, field).Result()
	if errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("field %s not found in hash %s", field, key)
	}
	return val, err
}

// HGetAll get all hash
func (rc *RedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return rc.client.HGetAll(ctx, key).Result()
}

// LPush left push list
func (rc *RedisClient) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return rc.client.LPush(ctx, key, values...).Result()
}

// BRPop block right pop list
func (rc *RedisClient) BRPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	return rc.client.BRPop(ctx, timeout, keys...).Result()
}

// Lock get distribute lock (auto generate random token)
func (rc *RedisClient) Lock(ctx context.Context, key string, ttl time.Duration) (string, error) {
	token := func() string {
		b := make([]byte, 16)
		rand.Read(b)
		return base64.StdEncoding.EncodeToString(b)
	}()
	ok, err := rc.client.SetNX(ctx, key, token, ttl).Result()
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New("acquire lock failed")
	}
	return token, nil
}

// Unlock release distribute lock (atom operate)
func (rc *RedisClient) Unlock(ctx context.Context, key, token string) error {
	script := redis.NewScript(`
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
	`)
	return script.Run(ctx, rc.client, []string{key}, token).Err()
}

// Transaction execute transaction (auto retry)
func (rc *RedisClient) Transaction(ctx context.Context, fn func(tx *redis.Tx) error, keys ...string) error {
	return rc.client.Watch(ctx, fn, keys...)
}

// TxIncr increase transaction
func (rc *RedisClient) TxIncr(ctx context.Context, key string) (int64, error) {
	var result int64
	err := rc.Transaction(ctx, func(tx *redis.Tx) error {
		// get current value
		n, err := tx.Get(ctx, key).Int64()
		if err != nil && !errors.Is(err, redis.Nil) {
			return err
		}
		// business logic
		n++
		// commit transaction
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, key, n, 0)
			return nil
		})
		return err
	}, key)
	if err == nil {
		result = 1
	}
	return result, err
}

// IncrBy increase atom
func (rc *RedisClient) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return rc.client.IncrBy(ctx, key, value).Result()
}

// DecrBy decrease atom
func (rc *RedisClient) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	return rc.client.DecrBy(ctx, key, value).Result()
}

// PoolStats get pool status
func (rc *RedisClient) PoolStats() *redis.PoolStats {
	return rc.client.PoolStats()
}

// Subscribe event
func (rc *RedisClient) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return rc.client.Subscribe(ctx, channels...)
}

// Publish event
func (rc *RedisClient) Publish(ctx context.Context, channel string, message interface{}) error {
	return rc.client.Publish(ctx, channel, message).Err()
}

// HealthCheck health check
func (rc *RedisClient) HealthCheck(ctx context.Context) error {
	_, err := rc.client.Ping(ctx).Result()
	return err
}

// ScanKeys scan keys
func (rc *RedisClient) ScanKeys(ctx context.Context, pattern string, batchSize int) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		var cursor uint64
		for {
			keys, nextCursor, err := rc.client.Scan(ctx, cursor, pattern, int64(batchSize)).Result()
			if err != nil {
				return
			}
			for _, key := range keys {
				ch <- key
			}
			if nextCursor == 0 {
				break
			}
			cursor = nextCursor
		}
	}()
	return ch
}

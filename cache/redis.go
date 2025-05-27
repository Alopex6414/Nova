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

// RedisClient
type RedisClient struct {
	client *redis.Client
	config *RedisConfig
}

// RedisConfig
type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	MaxRetries   int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
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
	token := generateToken()
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

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// NewRedisClient
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

// use distribute lock
func useDistributedLock(client *RedisClient) {
	ctx := context.Background()
	lockKey := "resource_lock"
	// get lock
	token, err := client.Lock(ctx, lockKey, 30*time.Second)
	if err != nil {
		fmt.Println("get lock failed:", err)
		return
	}
	defer client.Unlock(ctx, lockKey, token)
	// perform operate with protection
	fmt.Println("execute...")
	time.Sleep(10 * time.Second)
}

// use transaction
func useTransaction(client *RedisClient) {
	ctx := context.Background()
	key := "transaction_counter"

	result, err := client.TxIncr(ctx, key)
	if err != nil {
		fmt.Println("transaction execute failed:", err)
		return
	}
	fmt.Printf("transaction increase result: %d\n", result)
}

// user hash operate
func useHash(client *RedisClient) {
	ctx := context.Background()
	key := "user:1001"
	// set hash segment
	err := client.HSet(ctx, key, "name", "Alice", "age", 30)
	if err != nil {
		panic(err)
	}
	// get single segment
	name, err := client.HGet(ctx, key, "name")
	if err != nil {
		panic(err)
	}
	fmt.Println("User name:", name)
	// get all segments
	allFields, err := client.HGetAll(ctx, key)
	if err != nil {
		panic(err)
	}
	fmt.Println("All fields:", allFields)
}

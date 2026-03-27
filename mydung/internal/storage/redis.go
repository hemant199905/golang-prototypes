package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mydung/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(RedisConfig models.RedisConfig) (*RedisStore, error) {
	addr := fmt.Sprintf("%s:%d", RedisConfig.Host, RedisConfig.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: RedisConfig.Password,
		DB:       RedisConfig.DB,
	})

	var lastErr error
	for i := 0; i < RedisConfig.MaxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err := rdb.Ping(ctx).Err()
		cancel()

		if err == nil {
			log.Printf("✅ Successfully connected to Redis at %s!", addr)
			return &RedisStore{client: rdb}, nil
		}

		lastErr = err
		log.Printf("⚠️ Attempt %d/%d: Redis unreachable. Retrying in %ds...", i+1, RedisConfig.MaxRetries, RedisConfig.RetryDelay)
		
		time.Sleep(time.Duration(RedisConfig.RetryDelay) * time.Second)
	}

	return nil, fmt.Errorf("redis connection failed after %d attempts: %v", RedisConfig.MaxRetries, lastErr)
}

// --- Persistence Methods (Key-Value) ---

func (r *RedisStore) SaveJob(ctx context.Context, job *models.Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, job.ID, data, 0).Err()
}

func (r *RedisStore) GetJob(ctx context.Context, jobID string) (*models.Job, bool) {
	val, err := r.client.Get(ctx, jobID).Result()
	if err != nil {
		return nil, false
	}

	var job models.Job
	if err := json.Unmarshal([]byte(val), &job); err != nil {
		return nil, false
	}
	return &job, true
}

// --- Queue Methods (Producer/Consumer) ---

// Enqueue: Job ko line ke piche (Left) dalta hai
func (r *RedisStore) Enqueue(ctx context.Context, queueName string, job *models.Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	// LPUSH job into the list 📥
	return r.client.LPush(ctx, queueName, data).Err()
}

// Dequeue: Blocking tareeke se aage (Right) se job nikalta hai
func (r *RedisStore) Dequeue(ctx context.Context, queueName string) (*models.Job, error) {
	// BRPop: Jab tak data na aaye, wait karo (blocking) ⏳
	// Index 0: queue name, Index 1: actual job data
	results, err := r.client.BRPop(ctx, 0, queueName).Result()
	if err != nil {
		return nil, err
	}

	var job models.Job
	if err := json.Unmarshal([]byte(results), &job); err != nil {
		return nil, err
	}

	return &job, nil
}
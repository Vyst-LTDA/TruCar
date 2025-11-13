package repositories

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
}

type redisCacheRepository struct {
	client *redis.Client
}

func NewRedisCacheRepository(client *redis.Client) CacheRepository {
	return &redisCacheRepository{client: client}
}

func (r *redisCacheRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, bytes, expiration).Err()
}

func (r *redisCacheRepository) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (r *redisCacheRepository) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

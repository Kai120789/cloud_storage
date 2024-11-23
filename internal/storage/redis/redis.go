package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStorage(addr string, password string, db int) *RedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()

	return &RedisStorage{
		client: client,
		ctx:    ctx,
	}
}

func (r *RedisStorage) WriteRefreshToken(userId uint, refreshTokenValue string) error {
	key := fmt.Sprintf("user:session:%d", userId)

	expiration := 31 * 24 * time.Hour

	err := r.client.Set(r.ctx, key, refreshTokenValue, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to write refresh token to Redis: %w", err)
	}

	return nil
}

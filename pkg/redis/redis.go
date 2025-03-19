package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"domashka-backend/config"
)

type Redis struct {
	client *redis.Client
	ctx    context.Context
}

func New(config *config.RedisConfig) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.DB,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к Redis: %v", err)
	}

	return &Redis{
		client: client,
		ctx:    ctx,
	}, nil
}

func (r *Redis) Set(key string, value string, ttl time.Duration) error {
	return r.client.Set(r.ctx, key, value, ttl).Err()
}

func (r *Redis) Get(key string) (string, error) {
	result, err := r.client.Get(r.ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return result, err
}

func (r *Redis) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *Redis) Close() error {
	return r.client.Close()
}

func (r *Redis) Ping() error {
	_, err := r.client.Ping(r.ctx).Result()
	return err
}

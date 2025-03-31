package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr string) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisClient{client: client}, nil
}

func (r *RedisClient) SetChannels(ctx context.Context, channels interface{}) error {
	data, err := json.Marshal(channels)
	if err != nil {
		return fmt.Errorf("failed to marshal channels: %v", err)
	}

	return r.client.Set(ctx, "channels", data, 24*time.Hour).Err()
}

func (r *RedisClient) GetChannels(ctx context.Context) ([]byte, error) {
	data, err := r.client.Get(ctx, "channels").Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get channels from Redis: %v", err)
	}

	return data, nil
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

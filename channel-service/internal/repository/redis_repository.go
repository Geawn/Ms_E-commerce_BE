package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yourusername/channel-service/internal/model"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(addr string) (*RedisRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisRepository{client: client}, nil
}

func (r *RedisRepository) GetChannels(ctx context.Context) ([]*model.Channel, error) {
	data, err := r.client.Get(ctx, "channels").Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get channels from Redis: %v", err)
	}

	var channels []*model.Channel
	if err := json.Unmarshal(data, &channels); err != nil {
		return nil, fmt.Errorf("failed to unmarshal channels: %v", err)
	}

	return channels, nil
}

func (r *RedisRepository) SetChannels(ctx context.Context, channels []*model.Channel) error {
	data, err := json.Marshal(channels)
	if err != nil {
		return fmt.Errorf("failed to marshal channels: %v", err)
	}

	return r.client.Set(ctx, "channels", data, 24*time.Hour).Err()
}

func (r *RedisRepository) Close() error {
	return r.client.Close()
}

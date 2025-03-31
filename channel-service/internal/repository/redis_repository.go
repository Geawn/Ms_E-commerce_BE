package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Geawn/Ms_E-commerce_BE/channel-service/internal/model"
	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(addr string) (*RedisRepository, error) {
	log.Printf("Connecting to Redis at %s", addr)
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}
	log.Printf("Successfully connected to Redis")

	return &RedisRepository{client: client}, nil
}

func (r *RedisRepository) GetChannels(ctx context.Context) ([]*model.Channel, error) {
	log.Printf("Getting channels from Redis")
	data, err := r.client.Get(ctx, "channels").Bytes()
	if err == redis.Nil {
		log.Printf("No channels found in Redis")
		return nil, nil
	}
	if err != nil {
		log.Printf("Error getting channels from Redis: %v", err)
		return nil, fmt.Errorf("failed to get channels from Redis: %v", err)
	}

	var channels []*model.Channel
	if err := json.Unmarshal(data, &channels); err != nil {
		log.Printf("Error unmarshaling channels: %v", err)
		return nil, fmt.Errorf("failed to unmarshal channels: %v", err)
	}

	log.Printf("Successfully retrieved %d channels from Redis", len(channels))
	return channels, nil
}

func (r *RedisRepository) SetChannels(ctx context.Context, channels []*model.Channel) error {
	log.Printf("Setting %d channels in Redis", len(channels))
	data, err := json.Marshal(channels)
	if err != nil {
		log.Printf("Error marshaling channels: %v", err)
		return fmt.Errorf("failed to marshal channels: %v", err)
	}

	if err := r.client.Set(ctx, "channels", data, 24*time.Hour).Err(); err != nil {
		log.Printf("Error setting channels in Redis: %v", err)
		return fmt.Errorf("failed to set channels in Redis: %v", err)
	}

	log.Printf("Successfully set channels in Redis")
	return nil
}

func (r *RedisRepository) Close() error {
	log.Printf("Closing Redis connection")
	return r.client.Close()
}

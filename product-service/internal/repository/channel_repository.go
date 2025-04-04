package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/Geawn/Ms_E-commerce_BE/product-service/internal/models"
)

const (
	channelListKey = "channels:list"
	expirationTime = 24 * time.Hour
)

type ChannelRepository struct {
	redisClient *redis.Client
}

func NewChannelRepository(redisClient *redis.Client) *ChannelRepository {
	return &ChannelRepository{
		redisClient: redisClient,
	}
}

func (r *ChannelRepository) GetChannels(ctx context.Context) ([]models.Channel, error) {
	// Try to get from Redis first
	val, err := r.redisClient.Get(ctx, channelListKey).Result()
	if err == nil {
		var channels []models.Channel
		if err := json.Unmarshal([]byte(val), &channels); err != nil {
			return nil, fmt.Errorf("error unmarshaling channels: %v", err)
		}
		return channels, nil
	}

	// If not in Redis, get from database (implement this part)
	// For now, return mock data
	channels := []models.Channel{
		{
			ID:           "1",
			Name:         "Default Channel",
			Slug:         "default",
			IsActive:     true,
			CurrencyCode: "USD",
			Countries: []models.Country{
				{Country: "United States", Code: "US"},
				{Country: "Canada", Code: "CA"},
			},
		},
	}

	// Cache the result in Redis
	channelsJSON, err := json.Marshal(channels)
	if err != nil {
		return nil, fmt.Errorf("error marshaling channels: %v", err)
	}

	err = r.redisClient.Set(ctx, channelListKey, channelsJSON, expirationTime).Err()
	if err != nil {
		return nil, fmt.Errorf("error caching channels: %v", err)
	}

	return channels, nil
}

func (r *ChannelRepository) InvalidateCache(ctx context.Context) error {
	return r.redisClient.Del(ctx, channelListKey).Err()
}

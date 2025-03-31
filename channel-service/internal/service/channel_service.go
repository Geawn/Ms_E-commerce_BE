package service

import (
	"context"
	"fmt"

	"github.com/Geawn/Ms_E-commerce_BE/channel-service/internal/model"
	"github.com/Geawn/Ms_E-commerce_BE/channel-service/internal/repository"
)

type ChannelService struct {
	repo repository.ChannelRepository
}

func NewChannelService(repo repository.ChannelRepository) *ChannelService {
	return &ChannelService{
		repo: repo,
	}
}

func (s *ChannelService) GetChannels(ctx context.Context) ([]*model.Channel, error) {
	// Try to get from repository first
	channels, err := s.repo.GetChannels(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get channels from repository: %v", err)
	}

	if channels != nil {
		return channels, nil
	}

	// If not in repository, fetch from your data source
	// TODO: Implement your data fetching logic here
	channels = []*model.Channel{
		{
			ID:           "1",
			Name:         "Default Channel",
			Slug:         "default-channel",
			IsActive:     true,
			CurrencyCode: "USD",
			Countries: []model.Country{
				{
					Country: "United States",
					Code:    "US",
				},
			},
		},
	}

	// Cache the result in repository
	if err := s.repo.SetChannels(ctx, channels); err != nil {
		// Log the error but don't fail the request
		fmt.Printf("Failed to cache channels: %v\n", err)
	}

	return channels, nil
}

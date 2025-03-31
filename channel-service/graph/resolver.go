package graph

import (
	"context"
	"fmt"

	"github.com/yourusername/channel-service/internal/repository"
	"github.com/yourusername/channel-service/internal/service"
)

type Resolver struct {
	channelService *service.ChannelService
}

func NewResolver(redisAddr string) (*Resolver, error) {
	repo, err := repository.NewRedisRepository(redisAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %v", err)
	}

	channelService := service.NewChannelService(repo)

	return &Resolver{
		channelService: channelService,
	}, nil
}

func (r *Resolver) Channels(ctx context.Context) ([]*Channel, error) {
	channels, err := r.channelService.GetChannels(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get channels: %v", err)
	}

	// Convert model.Channel to graph.Channel
	var result []*Channel
	for _, ch := range channels {
		countries := make([]*Country, len(ch.Countries))
		for i, c := range ch.Countries {
			countries[i] = &Country{
				Country: c.Country,
				Code:    c.Code,
			}
		}

		result = append(result, &Channel{
			ID:           ch.ID,
			Name:         ch.Name,
			Slug:         ch.Slug,
			IsActive:     ch.IsActive,
			CurrencyCode: ch.CurrencyCode,
			Countries:    countries,
		})
	}

	return result, nil
}

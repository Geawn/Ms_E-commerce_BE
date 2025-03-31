package graphql

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"github.com/Geawn/Ms_E-commerce_BE/channel-service/internal/repository"
	"github.com/Geawn/Ms_E-commerce_BE/channel-service/internal/service"
)

type Resolver struct {
	channelService *service.ChannelService
}

func NewResolver(redisAddr string) (*Resolver, error) {
	repo, err := repository.NewRedisRepository(redisAddr)
	if err != nil {
		return nil, err
	}

	channelService := service.NewChannelService(repo)

	return &Resolver{
		channelService: channelService,
	}, nil
}

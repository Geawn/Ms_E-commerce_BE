package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.70

import (
	"context"
	"log"

	"github.com/Geawn/Ms_E-commerce_BE/channel-service/internal/model"
)

// Channels is the resolver for the channels field.
func (r *queryResolver) Channels(ctx context.Context) ([]*model.Channel, error) {
	log.Printf("Resolver: Getting channels")
	channels, err := r.Resolver.channelService.GetChannels(ctx)
	if err != nil {
		log.Printf("Resolver: Error getting channels: %v", err)
		return nil, err
	}
	log.Printf("Resolver: Successfully retrieved %d channels", len(channels))
	return channels, nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

package database

import (
	"context"
	"fmt"

	"github.com/Geawn/Ms_E-commerce_BE/channel-service/internal/model"
	"github.com/Geawn/Ms_E-commerce_BE/channel-service/internal/repository"
)

func MigrateChannels(ctx context.Context, repo repository.ChannelRepository) error {
	channels := []*model.Channel{
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
				{
					Country: "Canada",
					Code:    "CA",
				},
			},
		},
		{
			ID:           "2",
			Name:         "European Channel",
			Slug:         "european-channel",
			IsActive:     true,
			CurrencyCode: "EUR",
			Countries: []model.Country{
				{
					Country: "Germany",
					Code:    "DE",
				},
				{
					Country: "France",
					Code:    "FR",
				},
				{
					Country: "Italy",
					Code:    "IT",
				},
			},
		},
		{
			ID:           "3",
			Name:         "Asian Channel",
			Slug:         "asian-channel",
			IsActive:     true,
			CurrencyCode: "JPY",
			Countries: []model.Country{
				{
					Country: "Japan",
					Code:    "JP",
				},
				{
					Country: "South Korea",
					Code:    "KR",
				},
				{
					Country: "Singapore",
					Code:    "SG",
				},
			},
		},
	}

	if err := repo.SetChannels(ctx, channels); err != nil {
		return fmt.Errorf("failed to migrate channels: %v", err)
	}

	return nil
}

package repository

import (
	"context"

	"github.com/Geawn/Ms_E-commerce_BE/channel-service/internal/model"
)

type ChannelRepository interface {
	GetChannels(ctx context.Context) ([]*model.Channel, error)
	SetChannels(ctx context.Context, channels []*model.Channel) error
}

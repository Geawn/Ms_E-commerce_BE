package repository

import (
	"context"

	"github.com/yourusername/channel-service/internal/model"
)

type ChannelRepository interface {
	GetChannels(ctx context.Context) ([]*model.Channel, error)
	SetChannels(ctx context.Context, channels []*model.Channel) error
}

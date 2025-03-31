package generated

import (
	"context"

	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/models"
)

type Resolver interface {
	Query() QueryResolver
	Menu() MenuResolver
	Page() PageResolver
	MenuItem() MenuItemResolver
}

type QueryResolver interface {
	Menu(ctx context.Context, slug string, channel string) (*models.Menu, error)
	Page(ctx context.Context, slug string) (*models.Page, error)
}

type MenuResolver interface {
	Items(ctx context.Context, obj *models.Menu) ([]*models.MenuItem, error)
}

type PageResolver interface {
	Content(ctx context.Context, obj *models.Page) (string, error)
}

type MenuItemResolver interface {
	Category(ctx context.Context, obj *models.MenuItem) (*models.Category, error)
	Collection(ctx context.Context, obj *models.MenuItem) (*models.Collection, error)
	Page(ctx context.Context, obj *models.MenuItem) (*models.Page, error)
	Children(ctx context.Context, obj *models.MenuItem) ([]*models.MenuItem, error)
} 
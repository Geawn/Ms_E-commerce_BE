package graphql

import (
	"context"

	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/models"
	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/service"
)

type Resolver struct {
	pageService *service.PageService
	menuService *service.MenuService
}

func NewResolver(pageService *service.PageService, menuService *service.MenuService) *Resolver {
	return &Resolver{
		pageService: pageService,
		menuService: menuService,
	}
}

func (r *Resolver) Page() *pageResolver {
	return &pageResolver{r}
}

func (r *Resolver) Menu() *menuResolver {
	return &menuResolver{r}
}

type pageResolver struct{ *Resolver }

func (r *pageResolver) Page(ctx context.Context, slug string) (*models.Page, error) {
	return r.pageService.GetBySlug(ctx, slug)
}

type menuResolver struct{ *Resolver }

func (r *menuResolver) Menu(ctx context.Context, slug string, channel string) (*models.Menu, error) {
	return r.menuService.GetBySlugAndChannel(ctx, slug, channel)
} 
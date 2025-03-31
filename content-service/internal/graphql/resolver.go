package graphql

import (
	"context"
	"fmt"

	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/models"
	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

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

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Page returns PageResolver implementation.
func (r *Resolver) Page() PageResolver { return &pageResolver{r} }

// Category returns CategoryResolver implementation.
func (r *Resolver) Category() CategoryResolver { return &categoryResolver{r} }

// Collection returns CollectionResolver implementation.
func (r *Resolver) Collection() CollectionResolver { return &collectionResolver{r} }

// MenuItem returns MenuItemResolver implementation.
func (r *Resolver) MenuItem() MenuItemResolver { return &menuItemResolver{r} }

type queryResolver struct{ *Resolver }

func (r *queryResolver) Page(ctx context.Context, slug string) (*models.Page, error) {
	return r.pageService.GetBySlug(ctx, slug)
}

func (r *queryResolver) Menu(ctx context.Context, slug string, channel string) (*models.Menu, error) {
	return r.menuService.GetBySlugAndChannel(ctx, slug, channel)
}

type pageResolver struct{ *Resolver }

func (r *pageResolver) ID(ctx context.Context, obj *models.Page) (string, error) {
	return fmt.Sprintf("%d", obj.ID), nil
}

func (r *pageResolver) Slug(ctx context.Context, obj *models.Page) (string, error) {
	return obj.Slug, nil
}

func (r *pageResolver) Title(ctx context.Context, obj *models.Page) (string, error) {
	return obj.Title, nil
}

func (r *pageResolver) SeoTitle(ctx context.Context, obj *models.Page) (*string, error) {
	return &obj.SeoTitle, nil
}

func (r *pageResolver) SeoDescription(ctx context.Context, obj *models.Page) (*string, error) {
	return &obj.SeoDescription, nil
}

func (r *pageResolver) Content(ctx context.Context, obj *models.Page) (string, error) {
	return obj.Content, nil
}

type categoryResolver struct{ *Resolver }

func (r *categoryResolver) ID(ctx context.Context, obj *models.Category) (string, error) {
	return fmt.Sprintf("%d", obj.ID), nil
}

type collectionResolver struct{ *Resolver }

func (r *collectionResolver) ID(ctx context.Context, obj *models.Collection) (string, error) {
	return fmt.Sprintf("%d", obj.ID), nil
}

type menuItemResolver struct{ *Resolver }

func (r *menuItemResolver) ID(ctx context.Context, obj *models.MenuItem) (string, error) {
	return fmt.Sprintf("%d", obj.ID), nil
}

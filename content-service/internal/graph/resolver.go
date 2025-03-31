package graph

import (
	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/repository"
	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/service"
	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/graph/generated"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB           *gorm.DB
	Redis        *redis.Client
	menuService  *service.MenuService
	pageService  *service.PageService
}

func NewResolver(db *gorm.DB, redis *redis.Client) *Resolver {
	menuRepo := repository.NewMenuRepository(db)
	pageRepo := repository.NewPageRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	collectionRepo := repository.NewCollectionRepository(db)

	return &Resolver{
		DB:           db,
		Redis:        redis,
		menuService:  service.NewMenuService(menuRepo, pageRepo, categoryRepo, collectionRepo, redis),
		pageService:  service.NewPageService(pageRepo, redis),
	}
}

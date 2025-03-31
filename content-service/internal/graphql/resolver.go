package graphql

import (
	"github.com/Geawn/Ms_E-commerce_BE/content-service/internal/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PageService *service.PageService
	MenuService *service.MenuService
}

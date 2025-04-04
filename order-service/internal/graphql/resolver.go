package graphql

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"github.com/Geawn/Ms_E-commerce_BE/order-service/internal/service"
)

type Resolver struct {
	UserService    *service.UserService
	ProductService *service.ProductService
	OrderService   *service.OrderService
}

func NewResolver(userService *service.UserService, productService *service.ProductService, orderService *service.OrderService) *Resolver {
	return &Resolver{
		UserService:    userService,
		ProductService: productService,
		OrderService:   orderService,
	}
}

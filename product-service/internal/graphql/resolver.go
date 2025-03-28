package graphql

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"github.com/Geawn/Ms_E-commerce_BE/product-service/internal/service"
)

type Resolver struct {
	productService *service.ProductService
}

func NewResolver(productService *service.ProductService) *Resolver {
	return &Resolver{
		productService: productService,
	}
}

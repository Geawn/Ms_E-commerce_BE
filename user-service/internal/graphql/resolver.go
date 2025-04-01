package graphql

import (
	"context"

	"github.com/Geawn/Ms_E-commerce_BE/user-service/internal/models"
	"github.com/Geawn/Ms_E-commerce_BE/user-service/internal/service"
)

type Resolver struct {
	userService service.UserService
}

func NewResolver(userService service.UserService) *Resolver {
	return &Resolver{
		userService: userService,
	}
}

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
	// In a real implementation, you would get the user ID from the context
	// This is typically set by your authentication middleware
	userID := "current-user-id" // Replace with actual user ID from context
	return r.userService.GetCurrentUser(ctx, userID)
}

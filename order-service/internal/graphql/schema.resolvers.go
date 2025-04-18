package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.70

import (
	"context"
	"fmt"

	"github.com/Geawn/Ms_E-commerce_BE/order-service/internal/model"
)

// Amount is the resolver for the amount field.
func (r *moneyResolver) Amount(ctx context.Context, obj *model.Money) (float64, error) {
	panic(fmt.Errorf("not implemented: Amount - amount"))
}

// Created is the resolver for the created field.
func (r *orderResolver) Created(ctx context.Context, obj *model.Order) (*model.DateTime, error) {
	panic(fmt.Errorf("not implemented: Created - created"))
}

// Total is the resolver for the total field.
func (r *orderResolver) Total(ctx context.Context, obj *model.Order) (*model.Money, error) {
	panic(fmt.Errorf("not implemented: Total - total"))
}

// Edges is the resolver for the edges field.
func (r *orderConnectionResolver) Edges(ctx context.Context, obj *model.OrderConnection) ([]*OrderEdge, error) {
	panic(fmt.Errorf("not implemented: Edges - edges"))
}

// Variant is the resolver for the variant field.
func (r *orderLineResolver) Variant(ctx context.Context, obj *model.OrderLine) (*model.ProductVariant, error) {
	panic(fmt.Errorf("not implemented: Variant - variant"))
}

// Price is the resolver for the price field.
func (r *pricingResolver) Price(ctx context.Context, obj *model.Pricing) (*model.Money, error) {
	panic(fmt.Errorf("not implemented: Price - price"))
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	token := ctx.Value("token").(string)
	resp, err := r.UserService.GetCurrentUser(ctx, token)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:        resp.Id,
		Email:     resp.Email,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
		Avatar: &model.Image{
			URL: resp.Avatar.Url,
			Alt: resp.Avatar.Alt,
		},
	}, nil
}

// Orders is the resolver for the orders field.
func (r *userResolver) Orders(ctx context.Context, obj *model.User, first *int) (*model.OrderConnection, error) {
	page := 1
	perPage := 10
	if first != nil {
		perPage = *first
	}

	orders, _, err := r.OrderService.GetUserOrders(ctx, obj.ID, page, perPage)
	if err != nil {
		return nil, err
	}

	edges := make([]*model.OrderEdge, len(orders))
	for i, order := range orders {
		edges[i] = &model.OrderEdge{
			Node: order,
		}
	}

	return &model.OrderConnection{
		Edges: edges,
	}, nil
}

// Money returns MoneyResolver implementation.
func (r *Resolver) Money() MoneyResolver { return &moneyResolver{r} }

// Order returns OrderResolver implementation.
func (r *Resolver) Order() OrderResolver { return &orderResolver{r} }

// OrderConnection returns OrderConnectionResolver implementation.
func (r *Resolver) OrderConnection() OrderConnectionResolver { return &orderConnectionResolver{r} }

// OrderLine returns OrderLineResolver implementation.
func (r *Resolver) OrderLine() OrderLineResolver { return &orderLineResolver{r} }

// Pricing returns PricingResolver implementation.
func (r *Resolver) Pricing() PricingResolver { return &pricingResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

type moneyResolver struct{ *Resolver }
type orderResolver struct{ *Resolver }
type orderConnectionResolver struct{ *Resolver }
type orderLineResolver struct{ *Resolver }
type pricingResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }

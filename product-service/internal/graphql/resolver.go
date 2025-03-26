package graphql

import (
	"context"

	"github.com/yourusername/product-service/internal/models"
	"github.com/yourusername/product-service/internal/service"
)

type Resolver struct {
	productService *service.ProductService
}

func NewResolver(productService *service.ProductService) *Resolver {
	return &Resolver{
		productService: productService,
	}
}

func (r *Resolver) Product() ProductResolver { return &productResolver{r} }

type productResolver struct{ *Resolver }

func (r *productResolver) ID(ctx context.Context, obj *models.Product) (string, error) {
	return obj.ID.String(), nil
}

func (r *productResolver) Name(ctx context.Context, obj *models.Product) (string, error) {
	return obj.Name, nil
}

func (r *productResolver) Slug(ctx context.Context, obj *models.Product) (string, error) {
	return obj.Slug, nil
}

func (r *productResolver) Description(ctx context.Context, obj *models.Product) (string, error) {
	return obj.Description, nil
}

func (r *productResolver) SeoTitle(ctx context.Context, obj *models.Product) (string, error) {
	return obj.SeoTitle, nil
}

func (r *productResolver) SeoDescription(ctx context.Context, obj *models.Product) (string, error) {
	return obj.SeoDesc, nil
}

func (r *productResolver) Thumbnail(ctx context.Context, obj *models.Product) (*models.Image, error) {
	return &models.Image{
		URL: obj.Thumbnail,
		Alt: obj.Name,
	}, nil
}

func (r *productResolver) Category(ctx context.Context, obj *models.Product) (*models.Category, error) {
	return &obj.Category, nil
}

func (r *productResolver) Variants(ctx context.Context, obj *models.Product) ([]*models.Variant, error) {
	var variants []*models.Variant
	for i := range obj.Variants {
		variants = append(variants, &obj.Variants[i])
	}
	return variants, nil
}

func (r *productResolver) Pricing(ctx context.Context, obj *models.Product) (*models.ProductPricing, error) {
	return &models.ProductPricing{
		PriceRange: models.PriceRange{
			Start: models.Price{
				Amount:   obj.Price,
				Currency: obj.Currency,
			},
			Stop: models.Price{
				Amount:   obj.Price,
				Currency: obj.Currency,
			},
		},
	}, nil
}

type queryResolver struct{ *Resolver }

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

func (r *queryResolver) Product(ctx context.Context, slug string, channel string) (*models.Product, error) {
	return r.productService.GetProductBySlug(ctx, slug)
}

func (r *queryResolver) Products(ctx context.Context, first *int, after *string) (*models.ProductConnection, error) {
	limit := 10
	if first != nil {
		limit = *first
	}

	offset := 0
	if after != nil {
		// TODO: Parse cursor and get offset
	}

	products, err := r.productService.ListProducts(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return &models.ProductConnection{
		Edges: func() []*models.ProductEdge {
			var edges []*models.ProductEdge
			for _, product := range products {
				edges = append(edges, &models.ProductEdge{
					Node: &product,
				})
			}
			return edges
		}(),
	}, nil
}

func (r *queryResolver) ProductsByCategory(ctx context.Context, categorySlug string, first *int, after *string) (*models.ProductConnection, error) {
	limit := 10
	if first != nil {
		limit = *first
	}

	offset := 0
	if after != nil {
		// TODO: Parse cursor and get offset
	}

	products, err := r.productService.ListProductsByCategory(ctx, categorySlug, limit, offset)
	if err != nil {
		return nil, err
	}

	return &models.ProductConnection{
		Edges: func() []*models.ProductEdge {
			var edges []*models.ProductEdge
			for _, product := range products {
				edges = append(edges, &models.ProductEdge{
					Node: &product,
				})
			}
			return edges
		}(),
	}, nil
}

func (r *queryResolver) SearchProducts(ctx context.Context, query string, first *int, after *string) (*models.ProductConnection, error) {
	limit := 10
	if first != nil {
		limit = *first
	}

	offset := 0
	if after != nil {
		// TODO: Parse cursor and get offset
	}

	products, err := r.productService.SearchProducts(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return &models.ProductConnection{
		Edges: func() []*models.ProductEdge {
			var edges []*models.ProductEdge
			for _, product := range products {
				edges = append(edges, &models.ProductEdge{
					Node: &product,
				})
			}
			return edges
		}(),
	}, nil
}

package service

import (
	"context"
	"fmt"

	"github.com/Geawn/Ms_E-commerce_BE/order-service/internal/model"
	pb "github.com/Geawn/Ms_E-commerce_BE/order-service/proto"
	"google.golang.org/grpc"
)

type ProductService struct {
	client pb.ProductServiceClient
}

func NewProductService(conn *grpc.ClientConn) *ProductService {
	return &ProductService{
		client: pb.NewProductServiceClient(conn),
	}
}

func (s *ProductService) GetProductDetails(ctx context.Context, slug, channel string) (*model.Product, error) {
	resp, err := s.client.GetProductDetails(ctx, &pb.GetProductDetailsRequest{
		Slug:    slug,
		Channel: channel,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get product details: %w", err)
	}

	// Convert variants
	variants := make([]*model.ProductVariant, len(resp.Variants))
	for i, v := range resp.Variants {
		variants[i] = &model.ProductVariant{
			ID:   v.Id,
			Name: v.Name,
			Pricing: &model.Pricing{
				Price: &model.Price{
					Gross: &model.Money{
						Amount:   v.Pricing.Price.Gross.Amount,
						Currency: v.Pricing.Price.Gross.Currency,
					},
				},
			},
		}
	}

	return &model.Product{
		ID:          resp.Id,
		Name:        resp.Name,
		Description: resp.Description,
		Slug:        resp.Slug,
		Thumbnail: &model.Image{
			URL: resp.Thumbnail.Url,
			Alt: resp.Thumbnail.Alt,
		},
		Category: &model.Category{
			ID:   resp.Category.Id,
			Name: resp.Category.Name,
		},
		Variants: variants,
	}, nil
}

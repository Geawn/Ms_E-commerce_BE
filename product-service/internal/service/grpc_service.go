package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Geawn/Ms_E-commerce_BE/product-service/internal/models"
	pb "github.com/Geawn/Ms_E-commerce_BE/product-service/proto"
	"gorm.io/gorm"
)

type GRPCProductService struct {
	pb.UnimplementedProductServiceServer
	ProductService *ProductService
}

func NewGRPCProductService(productService *ProductService) *GRPCProductService {
	return &GRPCProductService{
		ProductService: productService,
	}
}

func (s *GRPCProductService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	product, err := s.ProductService.GetProductBySlug(ctx, req.Slug)
	if err != nil {
		return nil, err
	}
	return convertToProtoProduct(product), nil
}

func (s *GRPCProductService) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	offset := 0
	if req.After != "" {
		// Convert cursor to offset
		cursorID, err := strconv.ParseUint(req.After, 10, 32)
		if err != nil {
			return nil, err
		}
		// Get the position of the cursor in the list
		position, err := s.ProductService.GetProductPosition(ctx, uint(cursorID))
		if err != nil {
			return nil, err
		}
		offset = position + 1
	}

	products, err := s.ProductService.ListProducts(ctx, int(req.Limit), offset)
	if err != nil {
		return nil, err
	}

	protoProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = convertToProtoProduct(product)
	}

	// Get total count for pagination
	totalCount, err := s.ProductService.GetTotalProducts(ctx)
	if err != nil {
		return nil, err
	}

	hasNextPage := offset+len(products) < totalCount

	return &pb.ListProductsResponse{
		Products: protoProducts,
		Total:    int32(totalCount),
		PageInfo: &pb.PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: offset > 0,
			StartCursor:     fmt.Sprintf("%d", products[0].ID),
			EndCursor:       fmt.Sprintf("%d", products[len(products)-1].ID),
		},
	}, nil
}

func (s *GRPCProductService) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.ListProductsResponse, error) {
	offset := 0
	if req.After != "" {
		// Convert cursor to offset
		cursorID, err := strconv.ParseUint(req.After, 10, 32)
		if err != nil {
			return nil, err
		}
		// Get the position of the cursor in the list
		position, err := s.ProductService.GetProductPosition(ctx, uint(cursorID))
		if err != nil {
			return nil, err
		}
		offset = position + 1
	}

	products, err := s.ProductService.SearchProducts(ctx, req.Query, int(req.Limit), offset, req.SortBy.String(), req.SortDirection.String())
	if err != nil {
		return nil, err
	}

	protoProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = convertToProtoProduct(product)
	}

	// Get total count for pagination
	totalCount, err := s.ProductService.GetTotalSearchResults(ctx, req.Query)
	if err != nil {
		return nil, err
	}

	hasNextPage := offset+len(products) < totalCount

	return &pb.ListProductsResponse{
		Products: protoProducts,
		Total:    int32(totalCount),
		PageInfo: &pb.PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: offset > 0,
			StartCursor:     fmt.Sprintf("%d", products[0].ID),
			EndCursor:       fmt.Sprintf("%d", products[len(products)-1].ID),
		},
	}, nil
}

func (s *GRPCProductService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	product := &models.Product{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		CategoryID:  uint(req.CategoryId),
		Rating:      req.Rating,
	}

	if req.ThumbnailUrl != "" {
		product.Thumbnail = &models.Image{
			URL:    req.ThumbnailUrl,
			Alt:    req.ThumbnailAlt,
			Size:   int(req.ThumbnailSize),
			Format: req.ThumbnailFormat,
		}
	}

	if err := s.ProductService.CreateProduct(ctx, product); err != nil {
		return nil, err
	}

	return convertToProtoProduct(product), nil
}

func (s *GRPCProductService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {
	product := &models.Product{
		Model:       gorm.Model{ID: uint(req.Id)},
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		CategoryID:  uint(req.CategoryId),
		Rating:      req.Rating,
	}

	if req.ThumbnailUrl != "" {
		product.Thumbnail = &models.Image{
			URL:    req.ThumbnailUrl,
			Alt:    req.ThumbnailAlt,
			Size:   int(req.ThumbnailSize),
			Format: req.ThumbnailFormat,
		}
	}

	if err := s.ProductService.UpdateProduct(ctx, product); err != nil {
		return nil, err
	}

	return convertToProtoProduct(product), nil
}

func (s *GRPCProductService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	if err := s.ProductService.DeleteProduct(ctx, uint(req.Id)); err != nil {
		return nil, err
	}

	return &pb.DeleteProductResponse{
		Success: true,
	}, nil
}

func (s *GRPCProductService) mustEmbedUnimplementedProductServiceServer() {}

func convertToProtoProduct(product *models.Product) *pb.Product {
	protoProduct := &pb.Product{
		Id:             uint32(product.ID),
		Name:           product.Name,
		Slug:           product.Slug,
		Description:    product.Description,
		SeoTitle:       product.SeoTitle,
		SeoDescription: product.SeoDescription,
		CategoryId:     uint32(product.CategoryID),
		CreatedAt:      product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if product.Category != nil && product.Category.ID != 0 {
		protoProduct.Category = &pb.Category{
			Id:             uint32(product.Category.ID),
			Name:           product.Category.Name,
			Slug:           product.Category.Slug,
			Description:    product.Category.Description,
			SeoTitle:       product.Category.SeoTitle,
			SeoDescription: product.Category.SeoDescription,
		}
	}

	if product.Thumbnail != nil {
		protoProduct.Thumbnail = &pb.Image{
			Url:    product.Thumbnail.URL,
			Alt:    product.Thumbnail.Alt,
			Size:   int32(product.Thumbnail.Size),
			Format: product.Thumbnail.Format,
		}
	}

	protoProduct.Variants = make([]*pb.ProductVariant, len(product.Variants))
	for i, variant := range product.Variants {
		protoProduct.Variants[i] = &pb.ProductVariant{
			Id:                uint32(variant.ID),
			Name:              variant.Name,
			QuantityAvailable: int32(variant.QuantityAvailable),
			Pricing: &pb.ProductPricing{
				PriceRange: &pb.PriceRange{
					Start: &pb.Price{
						Gross: &pb.Money{
							Amount:   variant.Pricing.PriceRange.Start.Amount,
							Currency: variant.Pricing.PriceRange.Start.Currency,
						},
					},
					Stop: &pb.Price{
						Gross: &pb.Money{
							Amount:   variant.Pricing.PriceRange.Stop.Amount,
							Currency: variant.Pricing.PriceRange.Stop.Currency,
						},
					},
				},
			},
		}
	}

	protoProduct.Pricing = &pb.ProductPricing{
		PriceRange: &pb.PriceRange{
			Start: &pb.Price{
				Gross: &pb.Money{
					Amount:   product.Pricing.PriceRange.Start.Amount,
					Currency: product.Pricing.PriceRange.Start.Currency,
				},
			},
			Stop: &pb.Price{
				Gross: &pb.Money{
					Amount:   product.Pricing.PriceRange.Stop.Amount,
					Currency: product.Pricing.PriceRange.Stop.Currency,
				},
			},
		},
	}

	return protoProduct
}

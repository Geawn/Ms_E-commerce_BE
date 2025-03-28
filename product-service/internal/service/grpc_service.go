package service

import (
	"context"

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
	products, err := s.ProductService.ListProducts(ctx, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	protoProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = convertToProtoProduct(product)
	}

	return &pb.ListProductsResponse{
		Products: protoProducts,
		Total:    int32(len(products)),
	}, nil
}

func (s *GRPCProductService) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.ListProductsResponse, error) {
	products, err := s.ProductService.SearchProducts(ctx, req.Query, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	protoProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = convertToProtoProduct(product)
	}

	return &pb.ListProductsResponse{
		Products: protoProducts,
		Total:    int32(len(products)),
	}, nil
}

func (s *GRPCProductService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	product := &models.Product{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		CategoryID:  uint(req.CategoryId),
		Rating:      req.Rating,
		Thumbnail: &models.Image{
			URL: req.ThumbnailUrl,
			Alt: req.ThumbnailAlt,
		},
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
		Thumbnail: &models.Image{
			URL: req.ThumbnailUrl,
			Alt: req.ThumbnailAlt,
		},
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
		Id:           uint32(product.ID),
		Name:         product.Name,
		Slug:         product.Slug,
		Description:  product.Description,
		CategoryId:   uint32(product.CategoryID),
		Rating:       product.Rating,
		ThumbnailUrl: product.Thumbnail.URL,
		ThumbnailAlt: product.Thumbnail.Alt,
		CreatedAt:    product.CreatedAt.String(),
		UpdatedAt:    product.UpdatedAt.String(),
	}

	if product.Category != nil && product.Category.ID != 0 {
		protoProduct.Category = &pb.Category{
			Id:          uint32(product.Category.ID),
			Name:        product.Category.Name,
			Slug:        product.Category.Slug,
			Description: product.Category.Description,
		}
	}

	protoProduct.Variants = make([]*pb.Variant, len(product.Variants))
	for i, variant := range product.Variants {
		protoProduct.Variants[i] = &pb.Variant{
			Id:    uint32(variant.ID),
			Name:  variant.Name,
			Stock: int32(variant.Stock),
			Price: variant.Pricing.PriceRange.Start.Amount,
		}

		protoProduct.Variants[i].Attributes = make([]*pb.VariantAttribute, len(variant.Attributes))
		for j, attr := range variant.Attributes {
			protoProduct.Variants[i].Attributes[j] = &pb.VariantAttribute{
				Id:    uint32(attr.ID),
				Name:  attr.Name,
				Value: attr.Value,
			}
		}
	}

	protoProduct.Attributes = make([]*pb.Attribute, len(product.Attributes))
	for i, attr := range product.Attributes {
		protoProduct.Attributes[i] = &pb.Attribute{
			Id:     uint32(attr.ID),
			Name:   attr.Name,
			Values: attr.Values,
		}
	}

	protoProduct.Collections = make([]*pb.Collection, len(product.Collections))
	for i, collection := range product.Collections {
		protoProduct.Collections[i] = &pb.Collection{
			Id:          uint32(collection.ID),
			Name:        collection.Name,
			Slug:        collection.Slug,
			Description: collection.Description,
		}
	}

	protoProduct.Reviews = make([]*pb.Review, len(product.Reviews))
	for i, review := range product.Reviews {
		protoProduct.Reviews[i] = &pb.Review{
			Id:        uint32(review.ID),
			Rating:    review.Rating,
			Comment:   review.Comment,
			CreatedAt: review.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			User: &pb.User{
				Id:    uint32(review.User.ID),
				Name:  review.User.Name,
				Email: review.User.Email,
			},
		}
	}

	return protoProduct
}

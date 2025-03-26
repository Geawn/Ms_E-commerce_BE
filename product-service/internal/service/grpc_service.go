package service

import (
	"context"

	"github.com/yourusername/product-service/internal/models"
	pb "github.com/yourusername/product-service/proto"
	"gorm.io/gorm"
)

type GRPCProductService struct {
	pb.UnimplementedProductServiceServer
	ProductService *ProductService
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
		protoProducts[i] = convertToProtoProduct(&product)
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
		protoProducts[i] = convertToProtoProduct(&product)
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
		Price:       req.Price,
		Stock:       int(req.Stock),
		CategoryID:  uint(req.CategoryId),
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
		Price:       req.Price,
		Stock:       int(req.Stock),
		CategoryID:  uint(req.CategoryId),
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

	return &pb.DeleteProductResponse{Success: true}, nil
}

func convertToProtoProduct(product *models.Product) *pb.Product {
	protoProduct := &pb.Product{
		Id:          uint32(product.ID),
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		Price:       product.Price,
		Stock:       int32(product.Stock),
		CategoryId:  uint32(product.CategoryID),
		CreatedAt:   product.CreatedAt.String(),
		UpdatedAt:   product.UpdatedAt.String(),
	}

	if product.Category.ID != 0 {
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
			Price: variant.Price,
			Stock: int32(variant.Stock),
		}
	}

	return protoProduct
}

package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/youngprinnce/product-microservice/internal/service"
	"github.com/youngprinnce/product-microservice/internal/service/product"
	pb "github.com/youngprinnce/product-microservice/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProductHandler implements the ProductService gRPC interface
type ProductHandler struct {
	pb.UnimplementedProductServiceServer
	productService product.ProductBC
}

// NewProductHandler creates a new product gRPC handler
func NewProductHandler(productService product.ProductBC) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct creates a new product
func (h *ProductHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	// Convert protobuf request to domain request
	createReq := product.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Type:        convertFromProtobufProductType(req.Type),
	}

	// Set type-specific fields
	switch req.Type {
	case pb.ProductType_DIGITAL:
		if req.DigitalProduct != nil {
			createReq.DigitalProduct = &product.DigitalProductInfo{
				FileSize:     req.DigitalProduct.FileSize,
				DownloadLink: req.DigitalProduct.DownloadLink,
			}
		}
	case pb.ProductType_PHYSICAL:
		if req.PhysicalProduct != nil {
			createReq.PhysicalProduct = &product.PhysicalProductInfo{
				Weight:     req.PhysicalProduct.Weight,
				Dimensions: req.PhysicalProduct.Dimensions,
			}
		}
	case pb.ProductType_SUBSCRIPTION:
		if req.SubscriptionProduct != nil {
			createReq.SubscriptionProduct = &product.SubscriptionProductInfo{
				SubscriptionPeriod: req.SubscriptionProduct.SubscriptionPeriod,
				RenewalPrice:       req.SubscriptionProduct.RenewalPrice,
			}
		}
	}

	prod, err := h.productService.CreateProduct(ctx, createReq)
	if err != nil {
		return nil, convertToGRPCError(err)
	}

	return &pb.CreateProductResponse{
		Product: convertToProtobufProduct(prod),
	}, nil
}

// GetProduct retrieves a product by ID
func (h *ProductHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}

	prod, err := h.productService.GetProduct(ctx, id)
	if err != nil {
		return nil, convertToGRPCError(err)
	}

	return &pb.GetProductResponse{
		Product: convertToProtobufProduct(prod),
	}, nil
}

// UpdateProduct updates a product
func (h *ProductHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}

	updateReq := product.UpdateProductRequest{
		Name:        req.Name,
		Description: req.Description,
	}

	if req.Price > 0 {
		updateReq.Price = &req.Price
	}

	// Set type-specific fields
	if req.DigitalProduct != nil {
		updateReq.DigitalProduct = &product.DigitalProductInfo{
			FileSize:     req.DigitalProduct.FileSize,
			DownloadLink: req.DigitalProduct.DownloadLink,
		}
	}
	if req.PhysicalProduct != nil {
		updateReq.PhysicalProduct = &product.PhysicalProductInfo{
			Weight:     req.PhysicalProduct.Weight,
			Dimensions: req.PhysicalProduct.Dimensions,
		}
	}
	if req.SubscriptionProduct != nil {
		updateReq.SubscriptionProduct = &product.SubscriptionProductInfo{
			SubscriptionPeriod: req.SubscriptionProduct.SubscriptionPeriod,
			RenewalPrice:       req.SubscriptionProduct.RenewalPrice,
		}
	}

	prod, err := h.productService.UpdateProduct(ctx, id, updateReq)
	if err != nil {
		return nil, convertToGRPCError(err)
	}

	return &pb.UpdateProductResponse{
		Product: convertToProtobufProduct(prod),
	}, nil
}

// DeleteProduct deletes a product
func (h *ProductHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}

	err = h.productService.DeleteProduct(ctx, id)
	if err != nil {
		return nil, convertToGRPCError(err)
	}

	return &pb.DeleteProductResponse{
		Success: true,
	}, nil
}

// ListProducts lists products with optional filtering and pagination
func (h *ProductHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	var typeFilter *product.ProductType
	// In protobuf, if the type field is set to any value, we apply the filter
	// If client wants to filter by DIGITAL (which is 0), they still need to explicitly set it
	// For now, we'll assume no filtering if type is DIGITAL (0) - this could be improved with a separate "has_filter" field
	if req.Type == pb.ProductType_PHYSICAL || req.Type == pb.ProductType_SUBSCRIPTION {
		prodType := convertFromProtobufProductType(req.Type)
		typeFilter = &prodType
	}

	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	products, total, err := h.productService.ListProducts(ctx, typeFilter, page, pageSize)
	if err != nil {
		return nil, convertToGRPCError(err)
	}

	var pbProducts []*pb.Product
	for _, prod := range products {
		pbProducts = append(pbProducts, convertToProtobufProduct(prod))
	}

	return &pb.ListProductsResponse{
		Products: pbProducts,
		Total:    total,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}

// Helper functions for conversion
func convertToProtobufProduct(prod *product.Product) *pb.Product {
	pbProd := &pb.Product{
		Id:          prod.ID.String(),
		Name:        prod.Name,
		Description: prod.Description,
		Price:       prod.Price,
		Type:        convertToProtobufProductType(prod.Type),
		CreatedAt:   timestamppb.New(prod.CreatedAt),
		UpdatedAt:   timestamppb.New(prod.UpdatedAt),
	}

	// Set type-specific fields
	if prod.DigitalProductInfo != nil {
		pbProd.DigitalProduct = &pb.DigitalProduct{
			FileSize:     prod.DigitalProductInfo.FileSize,
			DownloadLink: prod.DigitalProductInfo.DownloadLink,
		}
	}
	if prod.PhysicalProductInfo != nil {
		pbProd.PhysicalProduct = &pb.PhysicalProduct{
			Weight:     prod.PhysicalProductInfo.Weight,
			Dimensions: prod.PhysicalProductInfo.Dimensions,
		}
	}
	if prod.SubscriptionProductInfo != nil {
		pbProd.SubscriptionProduct = &pb.SubscriptionProduct{
			SubscriptionPeriod: prod.SubscriptionProductInfo.SubscriptionPeriod,
			RenewalPrice:       prod.SubscriptionProductInfo.RenewalPrice,
		}
	}

	return pbProd
}

func convertToProtobufProductType(prodType product.ProductType) pb.ProductType {
	switch prodType {
	case product.DigitalProduct:
		return pb.ProductType_DIGITAL
	case product.PhysicalProduct:
		return pb.ProductType_PHYSICAL
	case product.SubscriptionProduct:
		return pb.ProductType_SUBSCRIPTION
	default:
		return pb.ProductType_DIGITAL
	}
}

func convertFromProtobufProductType(pbType pb.ProductType) product.ProductType {
	switch pbType {
	case pb.ProductType_DIGITAL:
		return product.DigitalProduct
	case pb.ProductType_PHYSICAL:
		return product.PhysicalProduct
	case pb.ProductType_SUBSCRIPTION:
		return product.SubscriptionProduct
	default:
		return product.DigitalProduct
	}
}

func convertToGRPCError(err error) error {
	switch err.(type) {
	case service.BadRequest:
		return status.Error(codes.InvalidArgument, err.Error())
	case service.NotFound:
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}

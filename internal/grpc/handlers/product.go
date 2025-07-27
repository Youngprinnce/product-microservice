package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/youngprinnce/product-microservice/internal/service"
	"github.com/youngprinnce/product-microservice/internal/service/product"
	"github.com/youngprinnce/product-microservice/internal/validation"
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
	// Basic input validation
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "product name is required")
	}
	if len(req.Name) > 255 {
		return nil, status.Error(codes.InvalidArgument, "product name must be at most 255 characters")
	}
	if len(req.Description) > 1000 {
		return nil, status.Error(codes.InvalidArgument, "product description must be at most 1000 characters")
	}
	if req.Price < 0 {
		return nil, status.Error(codes.InvalidArgument, "product price cannot be negative")
	}

	// Sanitize input
	req.Name = validation.SanitizeString(req.Name)
	req.Description = validation.SanitizeString(req.Description)

	// Validate type-specific fields at handler level
	if err := h.validateTypeSpecificFields(req.Type, req.DigitalProduct, req.PhysicalProduct, req.SubscriptionProduct); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

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
	// Input validation and sanitization
	if err := h.validateAndSanitizeUpdateProductRequest(req); err != nil {
		return nil, err
	}

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
	// With optional protobuf field, we can now properly detect if type filter was provided
	if req.Type != nil {
		prodType := convertFromProtobufProductType(*req.Type)
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

func (h *ProductHandler) validateAndSanitizeUpdateProductRequest(req *pb.UpdateProductRequest) error {
	// Required field validation
	if req.Id == "" {
		return status.Error(codes.InvalidArgument, "id is required")
	}

	// UUID validation for id
	if _, err := uuid.Parse(req.Id); err != nil {
		return status.Error(codes.InvalidArgument, "invalid id format")
	}

	// Sanitize text inputs if provided
	if req.Name != "" {
		req.Name = validation.SanitizeString(req.Name)
		if len(req.Name) < 2 {
			return status.Error(codes.InvalidArgument, "name must be at least 2 characters")
		}
		if len(req.Name) > 255 {
			return status.Error(codes.InvalidArgument, "name must be at most 255 characters")
		}
	}

	if req.Description != "" {
		req.Description = validation.SanitizeString(req.Description)
		if len(req.Description) > 1000 {
			return status.Error(codes.InvalidArgument, "description must be at most 1000 characters")
		}
	}

	// Business rule validation for optional fields
	if req.Price != 0 {
		if req.Price < 0 {
			return status.Error(codes.InvalidArgument, "price cannot be negative")
		}
		if req.Price > 1000000 {
			return status.Error(codes.InvalidArgument, "price cannot exceed 1,000,000")
		}
	}

	// Validate type-specific fields if provided
	if req.DigitalProduct != nil {
		if req.DigitalProduct.DownloadLink != "" {
			sanitizedURL := validation.SanitizeURL(req.DigitalProduct.DownloadLink)
			if sanitizedURL == "" {
				return status.Error(codes.InvalidArgument, "invalid download_link format - must be a valid URL")
			}
			req.DigitalProduct.DownloadLink = sanitizedURL
		}
		if req.DigitalProduct.FileSize < 0 {
			return status.Error(codes.InvalidArgument, "file_size cannot be negative")
		}
	}

	if req.PhysicalProduct != nil {
		if req.PhysicalProduct.Weight < 0 {
			return status.Error(codes.InvalidArgument, "weight cannot be negative")
		}
		if req.PhysicalProduct.Dimensions != "" && len(req.PhysicalProduct.Dimensions) > 50 {
			return status.Error(codes.InvalidArgument, "dimensions too long")
		}
	}

	if req.SubscriptionProduct != nil {
		if req.SubscriptionProduct.SubscriptionPeriod != "" {
			validPeriods := []string{"daily", "weekly", "monthly", "quarterly", "yearly"}
			isValidPeriod := false
			for _, period := range validPeriods {
				if req.SubscriptionProduct.SubscriptionPeriod == period {
					isValidPeriod = true
					break
				}
			}
			if !isValidPeriod {
				return status.Error(codes.InvalidArgument, "invalid subscription_period. Must be one of: daily, weekly, monthly, quarterly, yearly")
			}
		}
		if req.SubscriptionProduct.RenewalPrice < 0 {
			return status.Error(codes.InvalidArgument, "renewal_price cannot be negative")
		}
	}

	return nil
}

func (h *ProductHandler) validateTypeSpecificFields(productType pb.ProductType, digitalProduct *pb.DigitalProduct, physicalProduct *pb.PhysicalProduct, subscriptionProduct *pb.SubscriptionProduct) error {
	switch productType {
	case pb.ProductType_DIGITAL:
		if digitalProduct == nil {
			return status.Error(codes.InvalidArgument, "digital_product is required for digital product type")
		}
		// Validate digital product fields
		if digitalProduct.DownloadLink != "" {
			// Simple URL validation
			sanitizedURL := validation.SanitizeURL(digitalProduct.DownloadLink)
			if sanitizedURL == "" {
				return status.Error(codes.InvalidArgument, "invalid download_link format - must be a valid URL")
			}
		}
		if digitalProduct.FileSize < 0 {
			return status.Error(codes.InvalidArgument, "file_size cannot be negative")
		}

	case pb.ProductType_PHYSICAL:
		if physicalProduct == nil {
			return status.Error(codes.InvalidArgument, "physical_product is required for physical product type")
		}
		// Validate physical product fields
		if physicalProduct.Weight < 0 {
			return status.Error(codes.InvalidArgument, "weight cannot be negative")
		}
		if physicalProduct.Dimensions != "" {
			// Basic validation for dimensions format
			if len(physicalProduct.Dimensions) > 50 {
				return status.Error(codes.InvalidArgument, "dimensions too long")
			}
		}

	case pb.ProductType_SUBSCRIPTION:
		if subscriptionProduct == nil {
			return status.Error(codes.InvalidArgument, "subscription_product is required for subscription product type")
		}
		// Validate subscription product fields
		if subscriptionProduct.SubscriptionPeriod == "" {
			return status.Error(codes.InvalidArgument, "subscription_period is required for subscription products")
		}
		validPeriods := []string{"daily", "weekly", "monthly", "quarterly", "yearly"}
		isValidPeriod := false
		for _, period := range validPeriods {
			if subscriptionProduct.SubscriptionPeriod == period {
				isValidPeriod = true
				break
			}
		}
		if !isValidPeriod {
			return status.Error(codes.InvalidArgument, "invalid subscription_period. Must be one of: daily, weekly, monthly, quarterly, yearly")
		}
		if subscriptionProduct.RenewalPrice < 0 {
			return status.Error(codes.InvalidArgument, "renewal_price cannot be negative")
		}
	}
	return nil
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

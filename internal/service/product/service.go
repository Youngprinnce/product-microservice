package product

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/youngprinnce/product-microservice/internal/service"
	"gorm.io/gorm"
)

// ProductBC defines the business logic interface for products
type ProductBC interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (*Product, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*Product, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req UpdateProductRequest) (*Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	ListProducts(ctx context.Context, typeFilter *ProductType, page, pageSize int) ([]*Product, int64, error)
}

// ProductService implements ProductBC
type ProductService struct {
	store ProductStore
}

// NewProductService creates a new product service
func NewProductService(store ProductStore) *ProductService {
	return &ProductService{
		store: store,
	}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, req CreateProductRequest) (*Product, error) {
	// Validate product type (business rule)
	if !req.Type.IsValid() {
		return nil, service.BadRequest{Err: errors.New("invalid product type")}
	}

	// Validate type-specific fields (business rules)
	if err := s.validateTypeSpecificFields(req.Type, req.DigitalProduct, req.PhysicalProduct, req.SubscriptionProduct); err != nil {
		return nil, service.BadRequest{Err: err}
	}

	product := &Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Type:        req.Type,
	}

	// Set type-specific fields
	switch req.Type {
	case DigitalProduct:
		product.DigitalProductInfo = req.DigitalProduct
	case PhysicalProduct:
		product.PhysicalProductInfo = req.PhysicalProduct
	case SubscriptionProduct:
		product.SubscriptionProductInfo = req.SubscriptionProduct
	}

	err := s.store.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(ctx context.Context, id uuid.UUID) (*Product, error) {
	product, err := s.store.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.NotFound{Err: errors.New("product not found")}
		}
		return nil, err
	}
	return product, nil
}

// UpdateProduct updates a product
func (s *ProductService) UpdateProduct(ctx context.Context, id uuid.UUID, req UpdateProductRequest) (*Product, error) {
	// Check if product exists
	existingProduct, err := s.store.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.NotFound{Err: errors.New("product not found")}
		}
		return nil, err
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}

	// Update type-specific fields based on existing product type
	switch existingProduct.Type {
	case DigitalProduct:
		if req.DigitalProduct != nil {
			if req.DigitalProduct.FileSize > 0 {
				updates["digital_file_size"] = req.DigitalProduct.FileSize
			}
			if req.DigitalProduct.DownloadLink != "" {
				updates["digital_download_link"] = req.DigitalProduct.DownloadLink
			}
		}
	case PhysicalProduct:
		if req.PhysicalProduct != nil {
			if req.PhysicalProduct.Weight > 0 {
				updates["physical_weight"] = req.PhysicalProduct.Weight
			}
			if req.PhysicalProduct.Dimensions != "" {
				updates["physical_dimensions"] = req.PhysicalProduct.Dimensions
			}
		}
	case SubscriptionProduct:
		if req.SubscriptionProduct != nil {
			if req.SubscriptionProduct.SubscriptionPeriod != "" {
				updates["subscription_period"] = req.SubscriptionProduct.SubscriptionPeriod
			}
			if req.SubscriptionProduct.RenewalPrice > 0 {
				updates["subscription_renewal_price"] = req.SubscriptionProduct.RenewalPrice
			}
		}
	}

	if len(updates) == 0 {
		return nil, service.BadRequest{Err: errors.New("no fields to update")}
	}

	return s.store.Update(ctx, id, updates)
}

// DeleteProduct deletes a product
func (s *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	// Check if product exists
	_, err := s.store.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.NotFound{Err: errors.New("product not found")}
		}
		return err
	}

	return s.store.Delete(ctx, id)
}

// ListProducts retrieves products with pagination and optional type filtering
func (s *ProductService) ListProducts(ctx context.Context, typeFilter *ProductType, page, pageSize int) ([]*Product, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	products, err := s.store.GetAll(ctx, typeFilter, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.store.Count(ctx, typeFilter)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// validateTypeSpecificFields validates that the correct type-specific fields are provided
func (s *ProductService) validateTypeSpecificFields(productType ProductType, digital *DigitalProductInfo, physical *PhysicalProductInfo, subscription *SubscriptionProductInfo) error {
	switch productType {
	case DigitalProduct:
		if digital == nil {
			return errors.New("digital product information is required for digital products")
		}
		// Business logic validation only
		if digital.FileSize <= 0 {
			return errors.New("file size must be greater than 0 for digital products")
		}
		if digital.DownloadLink == "" {
			return errors.New("download link is required for digital products")
		}
	case PhysicalProduct:
		if physical == nil {
			return errors.New("physical product information is required for physical products")
		}
		// Business logic validation only
		if physical.Weight <= 0 {
			return errors.New("weight must be greater than 0 for physical products")
		}
		if physical.Dimensions == "" {
			return errors.New("dimensions are required for physical products")
		}
	case SubscriptionProduct:
		if subscription == nil {
			return errors.New("subscription product information is required for subscription products")
		}
		// Business logic validation only
		if subscription.SubscriptionPeriod == "" {
			return errors.New("subscription period is required for subscription products")
		}
		if subscription.RenewalPrice <= 0 {
			return errors.New("renewal price must be greater than 0 for subscription products")
		}
	}
	return nil
}

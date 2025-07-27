package product

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductStore defines the interface for product data operations
type ProductStore interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	GetAll(ctx context.Context, typeFilter *ProductType, limit, offset int) ([]*Product, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context, typeFilter *ProductType) (int64, error)
}

// ProductRepo implements ProductStore using GORM
type ProductRepo struct {
	db *gorm.DB
}

// NewProductRepo creates a new product repository
func NewProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

// Create creates a new product
func (r *ProductRepo) Create(ctx context.Context, product *Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// GetByID retrieves a product by ID
func (r *ProductRepo) GetByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	var product Product
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetAll retrieves all products with optional type filtering and pagination
func (r *ProductRepo) GetAll(ctx context.Context, typeFilter *ProductType, limit, offset int) ([]*Product, error) {
	var products []*Product
	query := r.db.WithContext(ctx)

	if typeFilter != nil {
		query = query.Where("type = ?", *typeFilter)
	}

	err := query.Limit(limit).Offset(offset).Find(&products).Error
	return products, err
}

// Update updates a product
func (r *ProductRepo) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*Product, error) {
	var product Product
	err := r.db.WithContext(ctx).Model(&product).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	// Fetch updated product
	err = r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// Delete permanently deletes a product
func (r *ProductRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&Product{}).Error
}

// Count returns the total number of products with optional type filtering
func (r *ProductRepo) Count(ctx context.Context, typeFilter *ProductType) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&Product{})

	if typeFilter != nil {
		query = query.Where("type = ?", *typeFilter)
	}

	err := query.Count(&count).Error
	return count, err
}

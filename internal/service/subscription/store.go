package subscription

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubscriptionStore defines the interface for subscription plan data operations
type SubscriptionStore interface {
	Create(ctx context.Context, plan *SubscriptionPlan) error
	GetByID(ctx context.Context, id uuid.UUID) (*SubscriptionPlan, error)
	GetByProductID(ctx context.Context, productID uuid.UUID, limit, offset int) ([]*SubscriptionPlan, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*SubscriptionPlan, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CountByProductID(ctx context.Context, productID uuid.UUID) (int64, error)
}

// SubscriptionRepo implements SubscriptionStore using GORM
type SubscriptionRepo struct {
	db *gorm.DB
}

// NewSubscriptionRepo creates a new subscription repository
func NewSubscriptionRepo(db *gorm.DB) *SubscriptionRepo {
	return &SubscriptionRepo{db: db}
}

// Create creates a new subscription plan
func (r *SubscriptionRepo) Create(ctx context.Context, plan *SubscriptionPlan) error {
	return r.db.WithContext(ctx).Create(plan).Error
}

// GetByID retrieves a subscription plan by ID
func (r *SubscriptionRepo) GetByID(ctx context.Context, id uuid.UUID) (*SubscriptionPlan, error) {
	var plan SubscriptionPlan
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// GetByProductID retrieves subscription plans for a specific product with pagination
func (r *SubscriptionRepo) GetByProductID(ctx context.Context, productID uuid.UUID, limit, offset int) ([]*SubscriptionPlan, error) {
	var plans []*SubscriptionPlan
	err := r.db.WithContext(ctx).Where("product_id = ?", productID).Limit(limit).Offset(offset).Find(&plans).Error
	return plans, err
}

// Update updates a subscription plan
func (r *SubscriptionRepo) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*SubscriptionPlan, error) {
	var plan SubscriptionPlan
	err := r.db.WithContext(ctx).Model(&plan).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	// Fetch updated plan
	err = r.db.WithContext(ctx).Where("id = ?", id).First(&plan).Error
	if err != nil {
		return nil, err
	}

	return &plan, nil
}

// Delete permanently deletes a subscription plan
func (r *SubscriptionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&SubscriptionPlan{}).Error
}

// CountByProductID returns the total number of subscription plans for a product
func (r *SubscriptionRepo) CountByProductID(ctx context.Context, productID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&SubscriptionPlan{}).Where("product_id = ?", productID).Count(&count).Error
	return count, err
}

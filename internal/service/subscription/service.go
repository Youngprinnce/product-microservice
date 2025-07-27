package subscription

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/youngprinnce/product-microservice/internal/service"
	"gorm.io/gorm"
)

// SubscriptionBC defines the business logic interface for subscription plans
type SubscriptionBC interface {
	CreateSubscriptionPlan(ctx context.Context, req CreateSubscriptionPlanRequest) (*SubscriptionPlan, error)
	GetSubscriptionPlan(ctx context.Context, id uuid.UUID) (*SubscriptionPlan, error)
	UpdateSubscriptionPlan(ctx context.Context, id uuid.UUID, req UpdateSubscriptionPlanRequest) (*SubscriptionPlan, error)
	DeleteSubscriptionPlan(ctx context.Context, id uuid.UUID) error
	ListSubscriptionPlans(ctx context.Context, productID uuid.UUID, page, pageSize int) ([]*SubscriptionPlan, int64, error)
}

// SubscriptionService implements SubscriptionBC
type SubscriptionService struct {
	store SubscriptionStore
}

// NewSubscriptionService creates a new subscription service
func NewSubscriptionService(store SubscriptionStore) *SubscriptionService {
	return &SubscriptionService{
		store: store,
	}
}

// CreateSubscriptionPlan creates a new subscription plan
func (s *SubscriptionService) CreateSubscriptionPlan(ctx context.Context, req CreateSubscriptionPlanRequest) (*SubscriptionPlan, error) {
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return nil, service.BadRequest{Err: errors.New("invalid product ID format")}
	}

	plan := &SubscriptionPlan{
		ID:        uuid.New(),
		ProductID: productID,
		PlanName:  req.PlanName,
		Duration:  req.Duration,
		Price:     req.Price,
	}

	err = s.store.Create(ctx, plan)
	if err != nil {
		return nil, err
	}

	return plan, nil
}

// GetSubscriptionPlan retrieves a subscription plan by ID
func (s *SubscriptionService) GetSubscriptionPlan(ctx context.Context, id uuid.UUID) (*SubscriptionPlan, error) {
	plan, err := s.store.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.NotFound{Err: errors.New("subscription plan not found")}
		}
		return nil, err
	}
	return plan, nil
}

// UpdateSubscriptionPlan updates a subscription plan
func (s *SubscriptionService) UpdateSubscriptionPlan(ctx context.Context, id uuid.UUID, req UpdateSubscriptionPlanRequest) (*SubscriptionPlan, error) {
	// Check if plan exists
	_, err := s.store.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.NotFound{Err: errors.New("subscription plan not found")}
		}
		return nil, err
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.PlanName != "" {
		updates["plan_name"] = req.PlanName
	}
	if req.Duration != nil {
		updates["duration"] = *req.Duration
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}

	if len(updates) == 0 {
		return nil, service.BadRequest{Err: errors.New("no fields to update")}
	}

	return s.store.Update(ctx, id, updates)
}

// DeleteSubscriptionPlan deletes a subscription plan
func (s *SubscriptionService) DeleteSubscriptionPlan(ctx context.Context, id uuid.UUID) error {
	// Check if plan exists
	_, err := s.store.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.NotFound{Err: errors.New("subscription plan not found")}
		}
		return err
	}

	return s.store.Delete(ctx, id)
}

// ListSubscriptionPlans retrieves subscription plans for a product with pagination
func (s *SubscriptionService) ListSubscriptionPlans(ctx context.Context, productID uuid.UUID, page, pageSize int) ([]*SubscriptionPlan, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	plans, err := s.store.GetByProductID(ctx, productID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.store.CountByProductID(ctx, productID)
	if err != nil {
		return nil, 0, err
	}

	return plans, total, nil
}

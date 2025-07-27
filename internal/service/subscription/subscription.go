package subscription

import (
	"time"

	"github.com/google/uuid"
)

// SubscriptionPlan represents a subscription plan entity
type SubscriptionPlan struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid"`
	PlanName  string    `json:"plan_name"`
	Duration  int       `json:"duration"` // number of days
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateSubscriptionPlanRequest represents the request to create a subscription plan
type CreateSubscriptionPlanRequest struct {
	ProductID string  `json:"product_id"`
	PlanName  string  `json:"plan_name"`
	Duration  int     `json:"duration"` // max 10 years
	Price     float64 `json:"price"`
}

// UpdateSubscriptionPlanRequest represents the request to update a subscription plan
type UpdateSubscriptionPlanRequest struct {
	PlanName string   `json:"plan_name,omitempty"`
	Duration *int     `json:"duration,omitempty"`
	Price    *float64 `json:"price,omitempty"`
}

// ListSubscriptionPlansRequest represents the request to list subscription plans
type ListSubscriptionPlansRequest struct {
	ProductID string `json:"product_id"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
}

// TableName returns the table name for the SubscriptionPlan model
func (SubscriptionPlan) TableName() string {
	return "subscription_plans"
}

package subscription

import (
	"time"

	"github.com/google/uuid"
)

// SubscriptionPlan represents a subscription plan entity
type SubscriptionPlan struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	PlanName  string    `json:"plan_name" gorm:"not null"`
	Duration  int       `json:"duration" gorm:"not null"` // number of days
	Price     float64   `json:"price" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateSubscriptionPlanRequest represents the request to create a subscription plan
type CreateSubscriptionPlanRequest struct {
	ProductID string  `json:"product_id" binding:"required"`
	PlanName  string  `json:"plan_name" binding:"required"`
	Duration  int     `json:"duration" binding:"required,gt=0"`
	Price     float64 `json:"price" binding:"required,gt=0"`
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

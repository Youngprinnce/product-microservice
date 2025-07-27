package product

import (
	"time"

	"github.com/google/uuid"
)

// ProductType represents the type of product
type ProductType string

const (
	DigitalProduct      ProductType = "digital"
	PhysicalProduct     ProductType = "physical"
	SubscriptionProduct ProductType = "subscription"
)

// Product represents the base product entity
type Product struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Price       float64     `json:"price"`
	Type        ProductType `json:"type"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`

	// Type-specific embedded structs
	DigitalProductInfo      *DigitalProductInfo      `json:"digital_product,omitempty" gorm:"embedded"`
	PhysicalProductInfo     *PhysicalProductInfo     `json:"physical_product,omitempty" gorm:"embedded"`
	SubscriptionProductInfo *SubscriptionProductInfo `json:"subscription_product,omitempty" gorm:"embedded"`
}

// DigitalProductInfo contains digital product specific fields
type DigitalProductInfo struct {
	FileSize     int64  `json:"file_size" gorm:"column:digital_file_size"`
	DownloadLink string `json:"download_link" gorm:"column:digital_download_link"`
}

// PhysicalProductInfo contains physical product specific fields
type PhysicalProductInfo struct {
	Weight     float64 `json:"weight" gorm:"column:physical_weight"`
	Dimensions string  `json:"dimensions" gorm:"column:physical_dimensions"`
}

// SubscriptionProductInfo contains subscription product specific fields
type SubscriptionProductInfo struct {
	SubscriptionPeriod string  `json:"subscription_period" gorm:"column:subscription_period"`
	RenewalPrice       float64 `json:"renewal_price" gorm:"column:subscription_renewal_price"`
}

// CreateProductRequest represents the request to create a product
type CreateProductRequest struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Price       float64     `json:"price"`
	Type        ProductType `json:"type"`

	// Type-specific fields
	DigitalProduct      *DigitalProductInfo      `json:"digital_product,omitempty"`
	PhysicalProduct     *PhysicalProductInfo     `json:"physical_product,omitempty"`
	SubscriptionProduct *SubscriptionProductInfo `json:"subscription_product,omitempty"`
}

// UpdateProductRequest represents the request to update a product
type UpdateProductRequest struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`

	// Type-specific fields
	DigitalProduct      *DigitalProductInfo      `json:"digital_product,omitempty"`
	PhysicalProduct     *PhysicalProductInfo     `json:"physical_product,omitempty"`
	SubscriptionProduct *SubscriptionProductInfo `json:"subscription_product,omitempty"`
}

// TableName returns the table name for the Product model
func (Product) TableName() string {
	return "products"
}

// IsValid checks if the product type is valid
func (pt ProductType) IsValid() bool {
	switch pt {
	case DigitalProduct, PhysicalProduct, SubscriptionProduct:
		return true
	default:
		return false
	}
}

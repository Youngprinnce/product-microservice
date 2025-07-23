package product

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockProductStore is a mock implementation of ProductStore
type MockProductStore struct {
	mock.Mock
}

func (m *MockProductStore) Create(ctx context.Context, product *Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductStore) GetByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Product), args.Error(1)
}

func (m *MockProductStore) GetAll(ctx context.Context, typeFilter *ProductType, limit, offset int) ([]*Product, error) {
	args := m.Called(ctx, typeFilter, limit, offset)
	return args.Get(0).([]*Product), args.Error(1)
}

func (m *MockProductStore) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*Product, error) {
	args := m.Called(ctx, id, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Product), args.Error(1)
}

func (m *MockProductStore) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductStore) Count(ctx context.Context, typeFilter *ProductType) (int64, error) {
	args := m.Called(ctx, typeFilter)
	return args.Get(0).(int64), args.Error(1)
}

func TestProductService_CreateProduct(t *testing.T) {
	mockStore := new(MockProductStore)
	service := NewProductService(mockStore)

	tests := []struct {
		name    string
		request CreateProductRequest
		setup   func()
		wantErr bool
	}{
		{
			name: "successful digital product creation",
			request: CreateProductRequest{
				Name:        "Test Digital Product",
				Description: "A test digital product",
				Price:       29.99,
				Type:        DigitalProduct,
				DigitalProduct: &DigitalProductInfo{
					FileSize:     1024000,
					DownloadLink: "https://example.com/download",
				},
			},
			setup: func() {
				mockStore.On("Create", mock.Anything, mock.AnythingOfType("*product.Product")).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name: "successful physical product creation",
			request: CreateProductRequest{
				Name:        "Test Physical Product",
				Description: "A test physical product",
				Price:       49.99,
				Type:        PhysicalProduct,
				PhysicalProduct: &PhysicalProductInfo{
					Weight:     2.5,
					Dimensions: "10x5x3 inches",
				},
			},
			setup: func() {
				mockStore.On("Create", mock.Anything, mock.AnythingOfType("*product.Product")).Return(nil).Once()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			product, err := service.CreateProduct(context.Background(), tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, product)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, product)
				assert.Equal(t, tt.request.Name, product.Name)
				assert.Equal(t, tt.request.Description, product.Description)
				assert.Equal(t, tt.request.Price, product.Price)
				assert.Equal(t, tt.request.Type, product.Type)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestProductService_GetProduct(t *testing.T) {
	mockStore := new(MockProductStore)
	service := NewProductService(mockStore)

	productID := uuid.New()
	expectedProduct := &Product{
		ID:          productID,
		Name:        "Test Product",
		Description: "A test product",
		Price:       29.99,
		Type:        DigitalProduct,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tests := []struct {
		name    string
		id      uuid.UUID
		setup   func()
		want    *Product
		wantErr bool
	}{
		{
			name: "successful get product",
			id:   productID,
			setup: func() {
				mockStore.On("GetByID", mock.Anything, productID).Return(expectedProduct, nil).Once()
			},
			want:    expectedProduct,
			wantErr: false,
		},
		{
			name: "product not found",
			id:   productID,
			setup: func() {
				mockStore.On("GetByID", mock.Anything, productID).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			product, err := service.GetProduct(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, product)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, product)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestProductService_ListProducts(t *testing.T) {
	mockStore := new(MockProductStore)
	service := NewProductService(mockStore)

	expectedProducts := []*Product{
		{
			ID:          uuid.New(),
			Name:        "Product 1",
			Description: "First product",
			Price:       29.99,
			Type:        DigitalProduct,
		},
		{
			ID:          uuid.New(),
			Name:        "Product 2",
			Description: "Second product",
			Price:       49.99,
			Type:        PhysicalProduct,
		},
	}

	t.Run("successful list all products", func(t *testing.T) {
		mockStore.On("GetAll", mock.Anything, (*ProductType)(nil), 10, 0).Return(expectedProducts, nil).Once()
		mockStore.On("Count", mock.Anything, (*ProductType)(nil)).Return(int64(2), nil).Once()

		products, total, err := service.ListProducts(context.Background(), nil, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, expectedProducts, products)
		assert.Equal(t, int64(2), total)

		mockStore.AssertExpectations(t)
	})
}

func TestProductService_DeleteProduct(t *testing.T) {
	mockStore := new(MockProductStore)
	service := NewProductService(mockStore)

	productID := uuid.New()
	existingProduct := &Product{
		ID:          productID,
		Name:        "Test Product",
		Description: "A test product",
		Price:       29.99,
		Type:        DigitalProduct,
	}

	t.Run("successful delete", func(t *testing.T) {
		mockStore.On("GetByID", mock.Anything, productID).Return(existingProduct, nil).Once()
		mockStore.On("Delete", mock.Anything, productID).Return(nil).Once()

		err := service.DeleteProduct(context.Background(), productID)

		assert.NoError(t, err)
		mockStore.AssertExpectations(t)
	})

	t.Run("product not found", func(t *testing.T) {
		nonExistentID := uuid.New()
		mockStore.On("GetByID", mock.Anything, nonExistentID).Return(nil, gorm.ErrRecordNotFound).Once()

		err := service.DeleteProduct(context.Background(), nonExistentID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		mockStore.AssertExpectations(t)
	})
}

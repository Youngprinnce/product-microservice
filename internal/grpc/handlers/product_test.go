package handlers

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/youngprinnce/product-microservice/internal/service/product"
	pb "github.com/youngprinnce/product-microservice/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockProductService is a mock implementation of ProductBC
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) CreateProduct(ctx context.Context, req product.CreateProductRequest) (*product.Product, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Product), args.Error(1)
}

func (m *MockProductService) GetProduct(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Product), args.Error(1)
}

func (m *MockProductService) UpdateProduct(ctx context.Context, id uuid.UUID, req product.UpdateProductRequest) (*product.Product, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Product), args.Error(1)
}

func (m *MockProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductService) ListProducts(ctx context.Context, typeFilter *product.ProductType, page, pageSize int) ([]*product.Product, int64, error) {
	args := m.Called(ctx, typeFilter, page, pageSize)
	return args.Get(0).([]*product.Product), args.Get(1).(int64), args.Error(2)
}

func TestProductHandler_CreateProduct(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	productID := uuid.New()
	expectedProduct := &product.Product{
		ID:          productID,
		Name:        "Test Digital Product",
		Description: "A test digital product",
		Price:       29.99,
		Type:        product.DigitalProduct,
		DigitalProductInfo: &product.DigitalProductInfo{
			FileSize:     1024000,
			DownloadLink: "https://example.com/download",
		},
	}

	t.Run("successful create digital product", func(t *testing.T) {
		req := &pb.CreateProductRequest{
			Name:        "Test Digital Product",
			Description: "A test digital product",
			Price:       29.99,
			Type:        pb.ProductType_DIGITAL,
			DigitalProduct: &pb.DigitalProduct{
				FileSize:     1024000,
				DownloadLink: "https://example.com/download",
			},
		}

		mockService.On("CreateProduct", mock.Anything, mock.AnythingOfType("product.CreateProductRequest")).Return(expectedProduct, nil).Once()

		resp, err := handler.CreateProduct(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Product)
		assert.Equal(t, expectedProduct.Name, resp.Product.Name)
		assert.Equal(t, expectedProduct.Description, resp.Product.Description)
		assert.Equal(t, expectedProduct.Price, resp.Product.Price)

		mockService.AssertExpectations(t)
	})
}

func TestProductHandler_GetProduct(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	productID := uuid.New()
	expectedProduct := &product.Product{
		ID:          productID,
		Name:        "Test Product",
		Description: "A test product",
		Price:       29.99,
		Type:        product.DigitalProduct,
	}

	tests := []struct {
		name           string
		request        *pb.GetProductRequest
		setup          func()
		wantErr        bool
		expectedStatus codes.Code
	}{
		{
			name: "successful get product",
			request: &pb.GetProductRequest{
				Id: productID.String(),
			},
			setup: func() {
				mockService.On("GetProduct", mock.Anything, productID).Return(expectedProduct, nil).Once()
			},
			wantErr: false,
		},
		{
			name: "invalid product ID",
			request: &pb.GetProductRequest{
				Id: "invalid-uuid",
			},
			setup:          func() {},
			wantErr:        true,
			expectedStatus: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			resp, err := handler.GetProduct(context.Background(), tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedStatus, st.Code())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Product)
				assert.Equal(t, expectedProduct.Name, resp.Product.Name)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestProductHandler_ListProducts(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	expectedProducts := []*product.Product{
		{
			ID:          uuid.New(),
			Name:        "Product 1",
			Description: "First product",
			Price:       29.99,
			Type:        product.DigitalProduct,
		},
		{
			ID:          uuid.New(),
			Name:        "Product 2",
			Description: "Second product",
			Price:       49.99,
			Type:        product.PhysicalProduct,
		},
	}

	t.Run("successful list products", func(t *testing.T) {
		req := &pb.ListProductsRequest{
			Page:     1,
			PageSize: 10,
		}

		mockService.On("ListProducts", mock.Anything, (*product.ProductType)(nil), 1, 10).Return(expectedProducts, int64(2), nil).Once()

		resp, err := handler.ListProducts(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Products, 2)
		assert.Equal(t, int64(2), resp.Total)
		assert.Equal(t, int32(1), resp.Page)
		assert.Equal(t, int32(10), resp.PageSize)

		mockService.AssertExpectations(t)
	})
}

func TestProductHandler_DeleteProduct(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	productID := uuid.New()

	t.Run("successful delete product", func(t *testing.T) {
		req := &pb.DeleteProductRequest{
			Id: productID.String(),
		}

		mockService.On("DeleteProduct", mock.Anything, productID).Return(nil).Once()

		resp, err := handler.DeleteProduct(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.Success)

		mockService.AssertExpectations(t)
	})
}

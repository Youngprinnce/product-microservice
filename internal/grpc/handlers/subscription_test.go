package handlers

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/youngprinnce/product-microservice/internal/service/subscription"
	pb "github.com/youngprinnce/product-microservice/proto"
)

// MockSubscriptionService is a mock implementation of SubscriptionBC
type MockSubscriptionService struct {
	mock.Mock
}

func (m *MockSubscriptionService) CreateSubscriptionPlan(ctx context.Context, req subscription.CreateSubscriptionPlanRequest) (*subscription.SubscriptionPlan, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*subscription.SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionService) GetSubscriptionPlan(ctx context.Context, id uuid.UUID) (*subscription.SubscriptionPlan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*subscription.SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionService) UpdateSubscriptionPlan(ctx context.Context, id uuid.UUID, req subscription.UpdateSubscriptionPlanRequest) (*subscription.SubscriptionPlan, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*subscription.SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionService) DeleteSubscriptionPlan(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSubscriptionService) ListSubscriptionPlans(ctx context.Context, productID uuid.UUID, page, pageSize int) ([]*subscription.SubscriptionPlan, int64, error) {
	args := m.Called(ctx, productID, page, pageSize)
	return args.Get(0).([]*subscription.SubscriptionPlan), args.Get(1).(int64), args.Error(2)
}

func TestSubscriptionHandler_CreateSubscriptionPlan(t *testing.T) {
	mockService := new(MockSubscriptionService)
	handler := NewSubscriptionHandler(mockService)

	subscriptionID := uuid.New()
	productID := uuid.New()
	expectedPlan := &subscription.SubscriptionPlan{
		ID:        subscriptionID,
		ProductID: productID,
		PlanName:  "Premium Plan",
		Duration:  30,
		Price:     29.99,
	}

	t.Run("successful create subscription plan", func(t *testing.T) {
		req := &pb.CreateSubscriptionPlanRequest{
			ProductId: productID.String(),
			PlanName:  "Premium Plan",
			Duration:  30,
			Price:     29.99,
		}

		mockService.On("CreateSubscriptionPlan", mock.Anything, mock.AnythingOfType("subscription.CreateSubscriptionPlanRequest")).Return(expectedPlan, nil).Once()

		resp, err := handler.CreateSubscriptionPlan(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Plan)
		assert.Equal(t, expectedPlan.PlanName, resp.Plan.PlanName)
		assert.Equal(t, expectedPlan.Duration, int(resp.Plan.Duration))
		assert.Equal(t, expectedPlan.Price, resp.Plan.Price)

		mockService.AssertExpectations(t)
	})
}

func TestSubscriptionHandler_GetSubscriptionPlan(t *testing.T) {
	mockService := new(MockSubscriptionService)
	handler := NewSubscriptionHandler(mockService)

	subscriptionID := uuid.New()
	productID := uuid.New()
	expectedPlan := &subscription.SubscriptionPlan{
		ID:        subscriptionID,
		ProductID: productID,
		PlanName:  "Premium Plan",
		Duration:  30,
		Price:     29.99,
	}

	t.Run("successful get subscription plan", func(t *testing.T) {
		req := &pb.GetSubscriptionPlanRequest{
			Id: subscriptionID.String(),
		}

		mockService.On("GetSubscriptionPlan", mock.Anything, subscriptionID).Return(expectedPlan, nil).Once()

		resp, err := handler.GetSubscriptionPlan(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Plan)
		assert.Equal(t, expectedPlan.PlanName, resp.Plan.PlanName)
		assert.Equal(t, expectedPlan.Duration, int(resp.Plan.Duration))

		mockService.AssertExpectations(t)
	})
}

func TestSubscriptionHandler_ListSubscriptionPlans(t *testing.T) {
	mockService := new(MockSubscriptionService)
	handler := NewSubscriptionHandler(mockService)

	productID := uuid.New()
	expectedPlans := []*subscription.SubscriptionPlan{
		{
			ID:        uuid.New(),
			ProductID: productID,
			PlanName:  "Basic Plan",
			Duration:  30,
			Price:     9.99,
		},
		{
			ID:        uuid.New(),
			ProductID: productID,
			PlanName:  "Premium Plan",
			Duration:  30,
			Price:     29.99,
		},
	}

	t.Run("successful list subscription plans", func(t *testing.T) {
		req := &pb.ListSubscriptionPlansRequest{
			ProductId: productID.String(),
			Page:      1,
			PageSize:  10,
		}

		mockService.On("ListSubscriptionPlans", mock.Anything, productID, 1, 10).Return(expectedPlans, int64(2), nil).Once()

		resp, err := handler.ListSubscriptionPlans(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Plans, 2)
		assert.Equal(t, int64(2), resp.Total)
		assert.Equal(t, int32(1), resp.Page)
		assert.Equal(t, int32(10), resp.PageSize)

		mockService.AssertExpectations(t)
	})
}

func TestSubscriptionHandler_DeleteSubscriptionPlan(t *testing.T) {
	mockService := new(MockSubscriptionService)
	handler := NewSubscriptionHandler(mockService)

	subscriptionID := uuid.New()

	t.Run("successful delete subscription plan", func(t *testing.T) {
		req := &pb.DeleteSubscriptionPlanRequest{
			Id: subscriptionID.String(),
		}

		mockService.On("DeleteSubscriptionPlan", mock.Anything, subscriptionID).Return(nil).Once()

		resp, err := handler.DeleteSubscriptionPlan(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.Success)

		mockService.AssertExpectations(t)
	})
}

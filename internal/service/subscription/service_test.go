package subscription

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSubscriptionStore is a mock implementation of SubscriptionStore
type MockSubscriptionStore struct {
	mock.Mock
}

func (m *MockSubscriptionStore) Create(ctx context.Context, plan *SubscriptionPlan) error {
	args := m.Called(ctx, plan)
	return args.Error(0)
}

func (m *MockSubscriptionStore) GetByID(ctx context.Context, id uuid.UUID) (*SubscriptionPlan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionStore) GetByProductID(ctx context.Context, productID uuid.UUID, limit, offset int) ([]*SubscriptionPlan, error) {
	args := m.Called(ctx, productID, limit, offset)
	return args.Get(0).([]*SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionStore) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*SubscriptionPlan, error) {
	args := m.Called(ctx, id, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionStore) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSubscriptionStore) CountByProductID(ctx context.Context, productID uuid.UUID) (int64, error) {
	args := m.Called(ctx, productID)
	return args.Get(0).(int64), args.Error(1)
}

func TestSubscriptionService_CreateSubscriptionPlan(t *testing.T) {
	mockStore := new(MockSubscriptionStore)
	service := NewSubscriptionService(mockStore)

	productID := uuid.New()
	request := CreateSubscriptionPlanRequest{
		ProductID: productID.String(),
		PlanName:  "Monthly Plan",
		Duration:  30,
		Price:     19.99,
	}

	t.Run("successful subscription plan creation", func(t *testing.T) {
		mockStore.On("Create", mock.Anything, mock.AnythingOfType("*subscription.SubscriptionPlan")).Return(nil).Once()

		plan, err := service.CreateSubscriptionPlan(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, plan)
		assert.Equal(t, productID, plan.ProductID)
		assert.Equal(t, request.PlanName, plan.PlanName)
		assert.Equal(t, request.Duration, plan.Duration)
		assert.Equal(t, request.Price, plan.Price)

		mockStore.AssertExpectations(t)
	})
}

func TestSubscriptionService_GetSubscriptionPlan(t *testing.T) {
	mockStore := new(MockSubscriptionStore)
	service := NewSubscriptionService(mockStore)

	planID := uuid.New()
	expectedPlan := &SubscriptionPlan{
		ID:        planID,
		ProductID: uuid.New(),
		PlanName:  "Monthly Plan",
		Duration:  30,
		Price:     19.99,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("successful get subscription plan", func(t *testing.T) {
		mockStore.On("GetByID", mock.Anything, planID).Return(expectedPlan, nil).Once()

		plan, err := service.GetSubscriptionPlan(context.Background(), planID)

		assert.NoError(t, err)
		assert.Equal(t, expectedPlan, plan)

		mockStore.AssertExpectations(t)
	})
}

func TestSubscriptionService_ListSubscriptionPlans(t *testing.T) {
	mockStore := new(MockSubscriptionStore)
	service := NewSubscriptionService(mockStore)

	productID := uuid.New()
	expectedPlans := []*SubscriptionPlan{
		{
			ID:        uuid.New(),
			ProductID: productID,
			PlanName:  "Monthly Plan",
			Duration:  30,
			Price:     19.99,
		},
		{
			ID:        uuid.New(),
			ProductID: productID,
			PlanName:  "Annual Plan",
			Duration:  365,
			Price:     199.99,
		},
	}

	t.Run("successful list subscription plans", func(t *testing.T) {
		mockStore.On("GetByProductID", mock.Anything, productID, 10, 0).Return(expectedPlans, nil).Once()
		mockStore.On("CountByProductID", mock.Anything, productID).Return(int64(2), nil).Once()

		plans, total, err := service.ListSubscriptionPlans(context.Background(), productID, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, expectedPlans, plans)
		assert.Equal(t, int64(2), total)

		mockStore.AssertExpectations(t)
	})
}

func TestSubscriptionService_DeleteSubscriptionPlan(t *testing.T) {
	mockStore := new(MockSubscriptionStore)
	service := NewSubscriptionService(mockStore)

	planID := uuid.New()
	existingPlan := &SubscriptionPlan{
		ID:        planID,
		ProductID: uuid.New(),
		PlanName:  "Test Plan",
		Duration:  30,
		Price:     29.99,
	}

	t.Run("successful delete", func(t *testing.T) {
		mockStore.On("GetByID", mock.Anything, planID).Return(existingPlan, nil).Once()
		mockStore.On("Delete", mock.Anything, planID).Return(nil).Once()

		err := service.DeleteSubscriptionPlan(context.Background(), planID)

		assert.NoError(t, err)
		mockStore.AssertExpectations(t)
	})
}

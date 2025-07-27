package subscription

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock
}

func createTestSubscriptionPlan() *SubscriptionPlan {
	return &SubscriptionPlan{
		ID:        uuid.New(),
		ProductID: uuid.New(),
		PlanName:  "Test Subscription Plan",
		Duration:  30,
		Price:     19.99,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestSubscriptionRepo_Create(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewSubscriptionRepo(db)
		ctx := context.Background()

		plan := createTestSubscriptionPlan()

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "subscription_plans"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(plan.ID))
		mock.ExpectCommit()

		err := repo.Create(ctx, plan)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("create with database error", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewSubscriptionRepo(db)
		ctx := context.Background()

		plan := createTestSubscriptionPlan()

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "subscription_plans"`)).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Create(ctx, plan)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSubscriptionRepo_GetByID(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewSubscriptionRepo(db)
		ctx := context.Background()

		planID := uuid.New()
		expectedPlan := createTestSubscriptionPlan()
		expectedPlan.ID = planID

		rows := sqlmock.NewRows([]string{
			"id", "product_id", "plan_name", "duration", "price", "created_at", "updated_at",
		}).AddRow(
			expectedPlan.ID, expectedPlan.ProductID, expectedPlan.PlanName,
			expectedPlan.Duration, expectedPlan.Price, expectedPlan.CreatedAt, expectedPlan.UpdatedAt,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subscription_plans" WHERE id = $1 ORDER BY "subscription_plans"."id" LIMIT $2`)).
			WithArgs(planID, 1).
			WillReturnRows(rows)

		plan, err := repo.GetByID(ctx, planID)

		assert.NoError(t, err)
		assert.NotNil(t, plan)
		assert.Equal(t, expectedPlan.ID, plan.ID)
		assert.Equal(t, expectedPlan.PlanName, plan.PlanName)
		assert.Equal(t, expectedPlan.Price, plan.Price)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("subscription plan not found", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewSubscriptionRepo(db)
		ctx := context.Background()

		planID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subscription_plans" WHERE id = $1 ORDER BY "subscription_plans"."id" LIMIT $2`)).
			WithArgs(planID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		plan, err := repo.GetByID(ctx, planID)

		assert.Error(t, err)
		assert.Nil(t, plan)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSubscriptionRepo_GetByProductID(t *testing.T) {
	t.Run("get plans by product ID", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewSubscriptionRepo(db)
		ctx := context.Background()

		productID := uuid.New()
		rows := sqlmock.NewRows([]string{
			"id", "product_id", "plan_name", "duration", "price", "created_at", "updated_at",
		}).AddRow(
			uuid.New(), productID, "Monthly Plan", 30, 19.99, time.Now(), time.Now(),
		).AddRow(
			uuid.New(), productID, "Annual Plan", 365, 199.99, time.Now(), time.Now(),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subscription_plans" WHERE product_id = $1 LIMIT $2`)).
			WithArgs(productID, 10).
			WillReturnRows(rows)

		plans, err := repo.GetByProductID(ctx, productID, 10, 0)

		assert.NoError(t, err)
		assert.Len(t, plans, 2)
		for _, plan := range plans {
			assert.Equal(t, productID, plan.ProductID)
		}
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("get plans with pagination", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewSubscriptionRepo(db)
		ctx := context.Background()

		productID := uuid.New()
		rows := sqlmock.NewRows([]string{
			"id", "product_id", "plan_name", "duration", "price", "created_at", "updated_at",
		}).AddRow(
			uuid.New(), productID, "Premium Plan", 30, 29.99, time.Now(), time.Now(),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subscription_plans" WHERE product_id = $1 LIMIT $2`)).
			WithArgs(productID, 1).
			WillReturnRows(rows)

		plans, err := repo.GetByProductID(ctx, productID, 1, 0)

		assert.NoError(t, err)
		assert.Len(t, plans, 1)
		assert.Equal(t, productID, plans[0].ProductID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSubscriptionRepo_Update(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewSubscriptionRepo(db)
		ctx := context.Background()

		planID := uuid.New()
		updates := map[string]interface{}{
			"plan_name": "Updated Plan Name",
			"price":     29.99,
		}

		// Mock the update operation
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "subscription_plans" SET`)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		// Mock the fetch operation
		rows := sqlmock.NewRows([]string{
			"id", "product_id", "plan_name", "duration", "price", "created_at", "updated_at",
		}).AddRow(
			planID, uuid.New(), "Updated Plan Name", 30, 29.99, time.Now(), time.Now(),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subscription_plans" WHERE id = $1 ORDER BY "subscription_plans"."id" LIMIT $2`)).
			WithArgs(planID, 1).
			WillReturnRows(rows)

		plan, err := repo.Update(ctx, planID, updates)

		assert.NoError(t, err)
		assert.NotNil(t, plan)
		assert.Equal(t, "Updated Plan Name", plan.PlanName)
		assert.Equal(t, 29.99, plan.Price)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("update with database error", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewSubscriptionRepo(db)
		ctx := context.Background()

		planID := uuid.New()
		updates := map[string]interface{}{
			"plan_name": "Updated Name",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "subscription_plans" SET`)).
			WillReturnError(errors.New("update failed"))
		mock.ExpectRollback()

		plan, err := repo.Update(ctx, planID, updates)

		assert.Error(t, err)
		assert.Nil(t, plan)
		assert.Contains(t, err.Error(), "update failed")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSubscriptionRepo_Delete(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewSubscriptionRepo(db)
		ctx := context.Background()

		planID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "subscription_plans" WHERE`)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Delete(ctx, planID)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("delete with database error", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewSubscriptionRepo(db)
		ctx := context.Background()

		planID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "subscription_plans" WHERE`)).
			WillReturnError(errors.New("delete failed"))
		mock.ExpectRollback()

		err := repo.Delete(ctx, planID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSubscriptionRepo_CountByProductID(t *testing.T) {
	t.Run("count plans by product ID", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewSubscriptionRepo(db)
		ctx := context.Background()

		productID := uuid.New()
		rows := sqlmock.NewRows([]string{"count"}).AddRow(3)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "subscription_plans" WHERE product_id = $1`)).
			WithArgs(productID).
			WillReturnRows(rows)

		count, err := repo.CountByProductID(ctx, productID)

		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("count with database error", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewSubscriptionRepo(db)
		ctx := context.Background()

		productID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "subscription_plans" WHERE product_id = $1`)).
			WithArgs(productID).
			WillReturnError(errors.New("count failed"))

		count, err := repo.CountByProductID(ctx, productID)

		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
		assert.Contains(t, err.Error(), "count failed")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

package product

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

func createTestProduct() *Product {
	return &Product{
		ID:          uuid.New(),
		Name:        "Test Product",
		Description: "A test product",
		Price:       29.99,
		Type:        DigitalProduct,
		DigitalProductInfo: &DigitalProductInfo{
			FileSize:     1024000,
			DownloadLink: "https://example.com/download",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestProductRepo_Create(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewProductRepo(db)
		ctx := context.Background()

		product := createTestProduct()

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "products"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(product.ID))
		mock.ExpectCommit()

		err := repo.Create(ctx, product)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("create with database error", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewProductRepo(db)
		ctx := context.Background()

		product := createTestProduct()

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "products"`)).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Create(ctx, product)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductRepo_GetByID(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewProductRepo(db)
		ctx := context.Background()

		productID := uuid.New()
		expectedProduct := createTestProduct()
		expectedProduct.ID = productID

		rows := sqlmock.NewRows([]string{
			"id", "name", "description", "price", "type", "created_at", "updated_at",
			"digital_file_size", "digital_download_link", "physical_weight",
			"physical_dimensions", "subscription_period", "subscription_renewal_price",
		}).AddRow(
			expectedProduct.ID, expectedProduct.Name, expectedProduct.Description,
			expectedProduct.Price, expectedProduct.Type, expectedProduct.CreatedAt, expectedProduct.UpdatedAt,
			expectedProduct.DigitalProductInfo.FileSize, expectedProduct.DigitalProductInfo.DownloadLink,
			nil, nil, nil, nil,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE id = $1 ORDER BY "products"."id" LIMIT $2`)).
			WithArgs(productID, 1).
			WillReturnRows(rows)

		product, err := repo.GetByID(ctx, productID)

		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, expectedProduct.ID, product.ID)
		assert.Equal(t, expectedProduct.Name, product.Name)
		assert.Equal(t, expectedProduct.Price, product.Price)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("product not found", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewProductRepo(db)
		ctx := context.Background()

		productID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE id = $1 ORDER BY "products"."id" LIMIT $2`)).
			WithArgs(productID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		product, err := repo.GetByID(ctx, productID)

		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductRepo_GetAll(t *testing.T) {
	t.Run("get all products without filter", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewProductRepo(db)
		ctx := context.Background()

		rows := sqlmock.NewRows([]string{
			"id", "name", "description", "price", "type", "created_at", "updated_at",
			"digital_file_size", "digital_download_link", "physical_weight",
			"physical_dimensions", "subscription_period", "subscription_renewal_price",
		}).AddRow(
			uuid.New(), "Product 1", "Description 1", 19.99, DigitalProduct, time.Now(), time.Now(),
			500000, "https://example.com/1", nil, nil, nil, nil,
		).AddRow(
			uuid.New(), "Product 2", "Description 2", 29.99, PhysicalProduct, time.Now(), time.Now(),
			nil, nil, 2.5, "10x10x5", nil, nil,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" LIMIT $1`)).
			WithArgs(10).
			WillReturnRows(rows)

		products, err := repo.GetAll(ctx, nil, 10, 0)

		assert.NoError(t, err)
		assert.Len(t, products, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("get products with type filter", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewProductRepo(db)
		ctx := context.Background()

		digitalType := DigitalProduct
		rows := sqlmock.NewRows([]string{
			"id", "name", "description", "price", "type", "created_at", "updated_at",
			"digital_file_size", "digital_download_link", "physical_weight",
			"physical_dimensions", "subscription_period", "subscription_renewal_price",
		}).AddRow(
			uuid.New(), "Digital Product", "Description", 19.99, DigitalProduct, time.Now(), time.Now(),
			500000, "https://example.com/digital", nil, nil, nil, nil,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE type = $1 LIMIT $2`)).
			WithArgs(DigitalProduct, 10).
			WillReturnRows(rows)

		products, err := repo.GetAll(ctx, &digitalType, 10, 0)

		assert.NoError(t, err)
		assert.Len(t, products, 1)
		assert.Equal(t, DigitalProduct, products[0].Type)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductRepo_Update(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewProductRepo(db)
		ctx := context.Background()

		productID := uuid.New()
		updates := map[string]interface{}{
			"name":  "Updated Product Name",
			"price": 39.99,
		}

		// Mock the update operation
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "products" SET`)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		// Mock the fetch operation
		rows := sqlmock.NewRows([]string{
			"id", "name", "description", "price", "type", "created_at", "updated_at",
			"digital_file_size", "digital_download_link", "physical_weight",
			"physical_dimensions", "subscription_period", "subscription_renewal_price",
		}).AddRow(
			productID, "Updated Product Name", "Description", 39.99, DigitalProduct, time.Now(), time.Now(),
			500000, "https://example.com/download", nil, nil, nil, nil,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE id = $1 ORDER BY "products"."id" LIMIT $2`)).
			WithArgs(productID, 1).
			WillReturnRows(rows)

		product, err := repo.Update(ctx, productID, updates)

		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, "Updated Product Name", product.Name)
		assert.Equal(t, 39.99, product.Price)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("update with database error", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewProductRepo(db)
		ctx := context.Background()

		productID := uuid.New()
		updates := map[string]interface{}{
			"name": "Updated Name",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "products" SET`)).
			WillReturnError(errors.New("update failed"))
		mock.ExpectRollback()

		product, err := repo.Update(ctx, productID, updates)

		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Contains(t, err.Error(), "update failed")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductRepo_Delete(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewProductRepo(db)
		ctx := context.Background()

		productID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "products" WHERE`)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Delete(ctx, productID)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("delete with database error", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewProductRepo(db)
		ctx := context.Background()

		productID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "products" WHERE`)).
			WillReturnError(errors.New("delete failed"))
		mock.ExpectRollback()

		err := repo.Delete(ctx, productID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductRepo_Count(t *testing.T) {
	t.Run("count all products", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewProductRepo(db)
		ctx := context.Background()

		rows := sqlmock.NewRows([]string{"count"}).AddRow(5)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "products"`)).
			WillReturnRows(rows)

		count, err := repo.Count(ctx, nil)

		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("count products with type filter", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewProductRepo(db)
		ctx := context.Background()

		digitalType := DigitalProduct
		rows := sqlmock.NewRows([]string{"count"}).AddRow(3)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "products" WHERE type = $1`)).
			WithArgs(DigitalProduct).
			WillReturnRows(rows)

		count, err := repo.Count(ctx, &digitalType)

		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("count with database error", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewProductRepo(db)
		ctx := context.Background()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "products"`)).
			WillReturnError(errors.New("count failed"))

		count, err := repo.Count(ctx, nil)

		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
		assert.Contains(t, err.Error(), "count failed")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

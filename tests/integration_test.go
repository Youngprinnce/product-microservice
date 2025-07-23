package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb_product "github.com/youngprinnce/product-microservice/proto"
	pb_subscription "github.com/youngprinnce/product-microservice/proto"
)

// IntegrationTestSuite defines the integration test suite
type IntegrationTestSuite struct {
	suite.Suite
	productClient      pb_product.ProductServiceClient
	subscriptionClient pb_subscription.SubscriptionServiceClient
	conn               *grpc.ClientConn
}

// SetupSuite runs before all tests in the suite
func (suite *IntegrationTestSuite) SetupSuite() {
	// Connect to the gRPC server
	// Note: This assumes the server is running on localhost:50051
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(suite.T(), err)

	suite.conn = conn
	suite.productClient = pb_product.NewProductServiceClient(conn)
	suite.subscriptionClient = pb_subscription.NewSubscriptionServiceClient(conn)
}

// TearDownSuite runs after all tests in the suite
func (suite *IntegrationTestSuite) TearDownSuite() {
	if suite.conn != nil {
		suite.conn.Close()
	}
}

// TestProductWorkflow tests the complete product workflow
func (suite *IntegrationTestSuite) TestProductWorkflow() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test 1: Create a digital product
	suite.T().Log("Creating digital product...")
	createReq := &pb_product.CreateProductRequest{
		Name:        "Integration Test E-book",
		Description: "A test e-book for integration testing",
		Price:       19.99,
		Type:        pb_product.ProductType_DIGITAL,
		DigitalProduct: &pb_product.DigitalProduct{
			FileSize:     1024000,
			DownloadLink: "https://example.com/test-ebook.pdf",
		},
	}

	createResp, err := suite.productClient.CreateProduct(ctx, createReq)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), createResp.Product)

	productID := createResp.Product.Id
	assert.NotEmpty(suite.T(), productID)
	assert.Equal(suite.T(), "Integration Test E-book", createResp.Product.Name)
	assert.Equal(suite.T(), pb_product.ProductType_DIGITAL, createResp.Product.Type)

	// Test 2: Get the created product
	suite.T().Log("Retrieving created product...")
	getReq := &pb_product.GetProductRequest{Id: productID}
	getResp, err := suite.productClient.GetProduct(ctx, getReq)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), productID, getResp.Product.Id)
	assert.Equal(suite.T(), "Integration Test E-book", getResp.Product.Name)

	// Test 3: Update the product
	suite.T().Log("Updating product...")
	updateReq := &pb_product.UpdateProductRequest{
		Id:    productID,
		Name:  "Updated Integration Test E-book",
		Price: 24.99,
	}
	updateResp, err := suite.productClient.UpdateProduct(ctx, updateReq)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated Integration Test E-book", updateResp.Product.Name)
	assert.Equal(suite.T(), 24.99, updateResp.Product.Price)

	// Test 4: List products (should include our product)
	suite.T().Log("Listing products...")
	listReq := &pb_product.ListProductsRequest{
		Page:     1,
		PageSize: 10,
	}
	listResp, err := suite.productClient.ListProducts(ctx, listReq)
	require.NoError(suite.T(), err)
	assert.Greater(suite.T(), len(listResp.Products), 0)

	// Find our product in the list
	var foundProduct *pb_product.Product
	for _, product := range listResp.Products {
		if product.Id == productID {
			foundProduct = product
			break
		}
	}
	assert.NotNil(suite.T(), foundProduct, "Created product should be in the list")

	// Test 5: Delete the product
	suite.T().Log("Deleting product...")
	deleteReq := &pb_product.DeleteProductRequest{Id: productID}
	deleteResp, err := suite.productClient.DeleteProduct(ctx, deleteReq)
	require.NoError(suite.T(), err)
	assert.True(suite.T(), deleteResp.Success)

	// Test 6: Verify product is deleted (should return not found)
	suite.T().Log("Verifying product deletion...")
	_, err = suite.productClient.GetProduct(ctx, getReq)
	assert.Error(suite.T(), err, "Getting deleted product should return an error")
}

// TestSubscriptionWorkflow tests the complete subscription workflow
func (suite *IntegrationTestSuite) TestSubscriptionWorkflow() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// First, create a product for the subscription plan
	suite.T().Log("Creating product for subscription...")
	createProductReq := &pb_product.CreateProductRequest{
		Name:        "Subscription Product",
		Description: "A product for subscription testing",
		Price:       0.00, // Base price, subscription will have its own pricing
		Type:        pb_product.ProductType_SUBSCRIPTION,
		SubscriptionProduct: &pb_product.SubscriptionProduct{
			SubscriptionPeriod: "monthly",
			RenewalPrice:       29.99,
		},
	}

	productResp, err := suite.productClient.CreateProduct(ctx, createProductReq)
	require.NoError(suite.T(), err)
	productID := productResp.Product.Id

	// Test 1: Create subscription plan
	suite.T().Log("Creating subscription plan...")
	createSubReq := &pb_subscription.CreateSubscriptionPlanRequest{
		ProductId: productID,
		PlanName:  "Monthly Premium Plan",
		Duration:  30,
		Price:     29.99,
	}

	createSubResp, err := suite.subscriptionClient.CreateSubscriptionPlan(ctx, createSubReq)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), createSubResp.Plan)

	planID := createSubResp.Plan.Id
	assert.NotEmpty(suite.T(), planID)
	assert.Equal(suite.T(), "Monthly Premium Plan", createSubResp.Plan.PlanName)
	assert.Equal(suite.T(), productID, createSubResp.Plan.ProductId)

	// Test 2: Get the subscription plan
	suite.T().Log("Retrieving subscription plan...")
	getSubReq := &pb_subscription.GetSubscriptionPlanRequest{Id: planID}
	getSubResp, err := suite.subscriptionClient.GetSubscriptionPlan(ctx, getSubReq)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), planID, getSubResp.Plan.Id)
	assert.Equal(suite.T(), "Monthly Premium Plan", getSubResp.Plan.PlanName)

	// Test 3: Update subscription plan
	suite.T().Log("Updating subscription plan...")
	updateSubReq := &pb_subscription.UpdateSubscriptionPlanRequest{
		Id:       planID,
		PlanName: "Updated Premium Plan",
		Duration: 45, // Increase duration from 30 to 45 days
		Price:    34.99,
	}
	updateSubResp, err := suite.subscriptionClient.UpdateSubscriptionPlan(ctx, updateSubReq)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated Premium Plan", updateSubResp.Plan.PlanName)
	assert.Equal(suite.T(), 34.99, updateSubResp.Plan.Price)
	assert.Equal(suite.T(), int32(45), updateSubResp.Plan.Duration)

	// Test 4: List subscription plans
	suite.T().Log("Listing subscription plans...")
	listSubReq := &pb_subscription.ListSubscriptionPlansRequest{
		ProductId: productID,
		Page:      1,
		PageSize:  10,
	}
	listSubResp, err := suite.subscriptionClient.ListSubscriptionPlans(ctx, listSubReq)
	require.NoError(suite.T(), err)
	assert.Greater(suite.T(), len(listSubResp.Plans), 0)

	// Find our plan in the list
	var foundPlan *pb_subscription.SubscriptionPlan
	for _, plan := range listSubResp.Plans {
		if plan.Id == planID {
			foundPlan = plan
			break
		}
	}
	assert.NotNil(suite.T(), foundPlan, "Created subscription plan should be in the list")

	// Test 5: Delete subscription plan
	suite.T().Log("Deleting subscription plan...")
	deleteSubReq := &pb_subscription.DeleteSubscriptionPlanRequest{Id: planID}
	deleteSubResp, err := suite.subscriptionClient.DeleteSubscriptionPlan(ctx, deleteSubReq)
	require.NoError(suite.T(), err)
	assert.True(suite.T(), deleteSubResp.Success)

	// Test 6: Clean up - delete the product
	suite.T().Log("Cleaning up - deleting product...")
	deleteProductReq := &pb_product.DeleteProductRequest{Id: productID}
	_, err = suite.productClient.DeleteProduct(ctx, deleteProductReq)
	require.NoError(suite.T(), err)
}

// TestIntegrationSuite runs the integration test suite
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

# Product Microservice

A gRPC-based microservice for managing products and subscription plans, built with Go, PostgreSQL, and Protocol Buffers.

## Table of Contents

- [Features](#features)
- [Setup Options](#setup-options)
  - [Option 1: Docker Setup (Recommended)](#option-1-docker-setup-recommended)
  - [Option 2: Manual Setup](#option-2-manual-setup)
- [API Documentation](#api-documentation)
- [Development](#development)

## Features

### Product Management

- **Create Products**: Support for digital, physical, and subscription products
- **CRUD Operations**: Full create, read, update, delete functionality
- **Product Types**:
  - **Digital Products**: File size and download links
  - **Physical Products**: Weight and dimensions
  - **Subscription Products**: Subscription periods and renewal pricing
- **Product Listing**: Paginated listing with optional type filtering

### Subscription Plan Management

- **Plan Management**: Create and manage subscription plans
- **Product Association**: Link subscription plans to products
- **Duration & Pricing**: Flexible duration (in days) and pricing models
- **CRUD Operations**: Complete lifecycle management for subscription plans

### Security & Authentication

- **Basic Authentication**: All gRPC endpoints protected with username/password authentication
- **Multiple Users**: Supports multiple user accounts with different credentials
- **Secure Headers**: Uses standard Authorization header with Base64 encoding
- **Default Users**: Pre-configured users for testing (admin, client, test)

## Setup Options

Choose your preferred setup method:

| Feature | Docker (Option 1) ✅ | Manual (Option 2) |
|---------|---------------------|-------------------|
| **Setup Time** | ~2 minutes | ~5-10 minutes |
| **Database Setup** | Automatic | Manual configuration required |
| **Dependencies** | Only Docker | Go, PostgreSQL, protoc |
| **Isolation** | Complete | Uses system resources |
| **Recommended For** | Development, Testing | Production, Custom setups |

## Option 1: Docker Setup (Recommended)

**Prerequisites**: Docker and Docker Compose installed

Docker setup provides the easiest and most consistent development experience with zero configuration required.

### 1. Start the Application

```bash
# Clone the repository
git clone https://github.com/youngprinnce/product-microservice.git
cd product-microservice

# Start all services (PostgreSQL + Application)
make docker-up
```

The application will be available on:

- **gRPC Server**: `localhost:50051`
- **PostgreSQL**: `localhost:5434`

### 2. Test the Docker Setup

```bash
# View logs
make docker-logs

# Test API
grpcurl -plaintext localhost:50051 list

# Stop services
make docker-down
```

## Option 2: Manual Setup

**Prerequisites**: Go 1.21+, PostgreSQL 13+, Protocol Buffers compiler

Manual setup requires you to configure your own PostgreSQL database.

### 1. Database Setup

1. **Create a PostgreSQL database** using pgAdmin or command line:
   ```sql
   CREATE DATABASE product_microservice;
   ```

2. **Update configuration** in `etc/config.yaml` with your database details:
   ```yaml
   database:
     host: "localhost"
     port: 5432
     user: "your_postgresql_username"
     password: "your_postgresql_password"
     db_name: "your_database_name"
   ```

### 2. Application Setup

```bash
# Install dependencies
go mod download

# Build the application
make build

# Run tests
make test

# Start the server
make run
```

The application will be available on:

- **gRPC Server**: `localhost:50051`
- **Your PostgreSQL**: Your configured host and port

### 3. Verify Manual Setup

```bash
# List available services
### 3. Verify Manual Setup

```bash
# Check if server is running
grpcurl -plaintext localhost:50051 list

# Test creating a product
grpcurl -plaintext -d '{
  "name": "Local Test Product",
  "description": "Testing manual setup",
  "price": 19.99,
  "type": "DIGITAL",
  "digital_product": {
    "file_size": 2048,
    "download_link": "https://example.com/test"
  }
}' localhost:50051 product.ProductService/CreateProduct
```

**Note**: Tables (`products`, `subscription_plans`) will be created automatically when the server starts.

---

## API Documentation

### Authentication

**All API endpoints require basic authentication.** Include the Authorization header with your requests.

#### Default Users

| Username | Password | Description |
|----------|----------|-------------|
| `admin` | `password123` | Administrative access |
| `client` | `client456` | Client access |
| `test` | `test789` | Testing access |

#### Authentication Examples

```bash
# Using grpcurl with authentication
grpcurl -plaintext \
  -H "authorization: Basic YWRtaW46cGFzc3dvcmQxMjM=" \
  -d '{"page": 1, "page_size": 10}' \
  localhost:50051 product.ProductService.ListProducts

# Generate auth header (admin:password123)
echo -n "admin:password123" | base64
# Result: YWRtaW46cGFzc3dvcmQxMjM=

# Generate auth header (client:client456)
echo -n "client:client456" | base64
# Result: Y2xpZW50OmNsaWVudDQ1Ng==
```

**⚠️ Important**: All examples below include authentication headers. Without proper authentication, you'll receive `Unauthenticated` errors.

### Product Service

#### CreateProduct

Creates a new product with type-specific information.

```bash
# Create Digital Product
grpcurl -plaintext \
  -H "authorization: Basic YWRtaW46cGFzc3dvcmQxMjM=" \
  -d '{
  "name": "E-book: Go Programming",
  "description": "Complete guide to Go programming",
  "price": 29.99,
  "type": "DIGITAL",
  "digital_product": {
    "file_size": 5242880,
    "download_link": "https://example.com/download/go-book.pdf"
  }
}' localhost:50051 product.ProductService.CreateProduct

# Create Physical Product
grpcurl -plaintext \
  -H "authorization: Basic YWRtaW46cGFzc3dvcmQxMjM=" \
  -d '{
  "name": "Go Programming Book",
  "description": "Physical copy of Go programming guide",
  "price": 49.99,
  "type": "PHYSICAL",
  "physical_product": {
    "weight": 0.5,
    "dimensions": "20x15x3 cm"
  }
}' localhost:50051 product.ProductService.CreateProduct
```

#### ListProducts

```bash
# List all products
grpcurl -plaintext \
  -H "authorization: Basic YWRtaW46cGFzc3dvcmQxMjM=" \
  -d '{
  "page": 1,
  "page_size": 10
}' localhost:50051 product.ProductService.ListProducts

# Filter by type
grpcurl -plaintext \
  -H "authorization: Basic YWRtaW46cGFzc3dvcmQxMjM=" \
  -d '{
  "type": "DIGITAL",
  "page": 1,
  "page_size": 10
}' localhost:50051 product.ProductService.ListProducts
```

#### GetProduct

```bash
# Get product by ID
grpcurl -plaintext \
  -H "authorization: Basic YWRtaW46cGFzc3dvcmQxMjM=" \
  -d '{"id": "your-product-id"}' \
  localhost:50051 product.ProductService.GetProduct
```

#### UpdateProduct

```bash
# Update product details
grpcurl -plaintext \
  -H "authorization: Basic YWRtaW46cGFzc3dvcmQxMjM=" \
  -d '{
  "id": "your-product-id",
  "name": "Updated Product Name",
  "price": 39.99
}' localhost:50051 product.ProductService.UpdateProduct

# Update digital product with new download link
grpcurl -plaintext \
  -H "authorization: Basic YWRtaW46cGFzc3dvcmQxMjM=" \
  -d '{
  "id": "your-product-id",
  "name": "Updated E-book",
  "description": "Updated description",
  "price": 49.99,
  "digital_product": {
    "file_size": 2048000,
    "download_link": "https://example.com/new-download-link.pdf"
  }
}' localhost:50051 product.ProductService.UpdateProduct
```

#### DeleteProduct

```bash
# Delete product
grpcurl -plaintext \
  -H "authorization: Basic YWRtaW46cGFzc3dvcmQxMjM=" \
  -d '{"id": "your-product-id"}' \
  localhost:50051 product.ProductService.DeleteProduct
```

### Subscription Service

#### CreateSubscriptionPlan

```bash
# Create a subscription plan for a product
grpcurl -plaintext \
  -H "authorization: Basic YWRtaW46cGFzc3dvcmQxMjM=" \
  -d '{
  "product_id": "your-product-id",
  "plan_name": "Monthly Premium",
  "duration": 30,
  "price": 29.99
}' localhost:50051 subscription.SubscriptionService.CreateSubscriptionPlan

# Create annual subscription plan
grpcurl -plaintext \
  -H "authorization: Basic YWRtaW46cGFzc3dvcmQxMjM=" \
  -d '{
  "product_id": "your-product-id", 
  "plan_name": "Annual Premium",
  "duration": 365,
  "price": 299.99
}' localhost:50051 subscription.SubscriptionService.CreateSubscriptionPlan
```

#### GetSubscriptionPlan

```bash
# Get subscription plan by ID
grpcurl -plaintext \
  -H "authorization: Basic YWRtaW46cGFzc3dvcmQxMjM=" \
  -d '{"id": "your-plan-id"}' \
  localhost:50051 subscription.SubscriptionService.GetSubscriptionPlan
```

#### ListSubscriptionPlans

```bash
# List all subscription plans for a product
grpcurl -plaintext \
  -H "authorization: Basic YWRtaW46cGFzc3dvcmQxMjM=" \
  -d '{
  "product_id": "your-product-id",
  "page": 1,
  "page_size": 10
}' localhost:50051 subscription.SubscriptionService.ListSubscriptionPlans

# List all subscription plans (no product filter)
grpcurl -plaintext \
  -H "authorization: Basic YWRtaW46cGFzc3dvcmQxMjM=" \
  -d '{
  "page": 1,
  "page_size": 10
}' localhost:50051 subscription.SubscriptionService.ListSubscriptionPlans
```

#### UpdateSubscriptionPlan

```bash
# Update subscription plan details
grpcurl -plaintext \
  -H "authorization: Basic YWRtaW46cGFzc3dvcmQxMjM=" \
  -d '{
  "id": "your-plan-id",
  "plan_name": "Updated Premium Plan",
  "duration": 60,
  "price": 49.99
}' localhost:50051 subscription.SubscriptionService.UpdateSubscriptionPlan
```

#### DeleteSubscriptionPlan

```bash
# Delete subscription plan
grpcurl -plaintext \
  -H "authorization: Basic YWRtaW46cGFzc3dvcmQxMjM=" \
  -d '{"id": "your-plan-id"}' \
  localhost:50051 subscription.SubscriptionService.DeleteSubscriptionPlan
```

## Development

### Available Make Commands

```bash
# Docker (Recommended)
make docker-up        # Start all services with docker-compose
make docker-down      # Stop docker-compose services
make docker-logs      # View logs from all services

# Local Development (requires Go and PostgreSQL)
make build            # Build the application
make run              # Run the server
make test             # Run tests with coverage
make proto            # Generate protobuf code
```

### Architecture

The service follows **Clean Architecture** principles:

```text
├── cmd/                    # Application entry points
├── internal/              # Private application code
│   ├── grpc/handlers/     # gRPC handlers (presentation layer)
│   ├── service/           # Business logic (use case layer)
│   ├── postgres/          # Database connection
│   └── logger/            # Logging utilities
├── proto/                 # Protocol buffer definitions
├── config/                # Configuration management
└── etc/                   # Configuration files
```

**Built with ❤️ using Go, gRPC, and PostgreSQL**

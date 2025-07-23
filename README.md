# Product Microservice

A gRPC-based microservice for managing products and subscription plans, built with Go, PostgreSQL, and Protocol Buffers.

## Table of Contents

- [Features](#features)
- [Quick Start (Docker)](#quick-start-docker)
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

## Quick Start (Docker)

**Prerequisites**: Docker and Docker Compose installed

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

### 2. Test the API

```bash
# List available services
grpcurl -plaintext localhost:50051 list

# Create a product
grpcurl -plaintext -d '{
  "name": "Test Product",
  "description": "A test product",
  "price": 29.99,
  "type": "DIGITAL",
  "digital_product": {
    "file_size": 1024,
    "download_link": "https://example.com/download"
  }
}' localhost:50051 product.ProductService.CreateProduct

# List products
grpcurl -plaintext -d '{"page": 1, "page_size": 10}' localhost:50051 product.ProductService.ListProducts
```

### 3. Manage Services

```bash
# View logs
make docker-logs

# Stop services
make docker-down

# Clean up (remove containers and volumes)
make docker-clean
```

## API Documentation

### Product Service

#### CreateProduct

Creates a new product with type-specific information.

```bash
# Create Digital Product
grpcurl -plaintext -d '{
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
grpcurl -plaintext -d '{
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
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10
}' localhost:50051 product.ProductService.ListProducts

# Filter by type
grpcurl -plaintext -d '{
  "type": "DIGITAL",
  "page": 1,
  "page_size": 10
}' localhost:50051 product.ProductService.ListProducts
```

#### Other Operations

```bash
# Get product by ID
grpcurl -plaintext -d '{"id": "your-product-id"}' localhost:50051 product.ProductService.GetProduct

# Update product
grpcurl -plaintext -d '{
  "id": "your-product-id",
  "name": "Updated Product Name",
  "price": 39.99
}' localhost:50051 product.ProductService.UpdateProduct

# Delete product
grpcurl -plaintext -d '{"id": "your-product-id"}' localhost:50051 product.ProductService.DeleteProduct
```

### Subscription Service

#### CreateSubscriptionPlan

```bash
grpcurl -plaintext -d '{
  "product_id": "your-product-id",
  "plan_name": "Monthly Premium",
  "duration": 30,
  "price": 29.99
}' localhost:50051 subscription.SubscriptionService.CreateSubscriptionPlan
```

#### Other Operations

```bash
# Get subscription plan
grpcurl -plaintext -d '{"id": "your-plan-id"}' localhost:50051 subscription.SubscriptionService.GetSubscriptionPlan

# List subscription plans
grpcurl -plaintext -d '{
  "product_id": "your-product-id",
  "page": 1,
  "page_size": 10
}' localhost:50051 subscription.SubscriptionService.ListSubscriptionPlans
```

## Development

### Available Make Commands

```bash
# Docker (Recommended)
make docker-up        # Start all services with docker-compose
make docker-down      # Stop docker-compose services
make docker-logs      # View logs from all services
make docker-clean     # Stop and clean up everything

# Local Development (requires Go and PostgreSQL)
make build            # Build the application
make run              # Run the server
make test             # Run tests
make clean            # Clean build artifacts
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

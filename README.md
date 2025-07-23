# Product Microservice

A gRPC-based microservice for managing products and subscription plans, built with Go, PostgreSQL, and Protocol Buffers.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Database Setup](#database-setup)
- [Running the Service](#running-the-service)
- [API Documentation](#api-documentation)
- [Testing](#testing)
- [Docker Deployment](#docker-deployment)
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

## Architecture

The service follows **Clean Architecture** principles with clear separation of concerns:

```
├── cmd/                    # Application entry points
│   └── server/            # gRPC server configuration
├── internal/              # Private application code
│   ├── grpc/             # gRPC handlers (presentation layer)
│   ├── service/          # Business logic (use case layer)
│   ├── postgres/         # Database connection
│   └── logger/           # Logging utilities
├── proto/                # Protocol buffer definitions
├── config/               # Configuration management
└── db/migrations/        # Database migrations
```

### Key Components

- **gRPC Handlers**: HTTP/gRPC interface layer
- **Business Services**: Core business logic and validation
- **Repository Pattern**: Data access abstraction
- **PostgreSQL**: Persistent data storage with GORM ORM
- **Protocol Buffers**: API contract definition

## Prerequisites

- **Go** 1.20 or later
- **PostgreSQL** 12 or later
- **Protocol Buffers Compiler** (protoc)
- **Docker** (optional, for containerized deployment)

### Go Tools Required

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/youngprinnce/product-microservice.git
   cd product-microservice
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Generate protobuf code** (if modified):
   ```bash
   make proto
   ```

## Configuration

The service uses YAML configuration. Copy the example configuration:

```bash
cp config.example.yaml config.yaml
```

### Configuration Options

```yaml
server:
  grpc_port: 50051
  
database:
  host: localhost
  port: 5432
  username: postgres
  password: postgres
  database: product_microservice
  ssl_mode: disable
  
logger:
  level: info
  format: json
```

### Environment Variables

You can override configuration with environment variables:

- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_USERNAME`: Database username
- `DB_PASSWORD`: Database password
- `DB_DATABASE`: Database name
- `GRPC_PORT`: gRPC server port

## Database Setup

### 1. Create Database

```sql
CREATE DATABASE product_microservice;
```

### 2. Run Migrations

```bash
make migrate-up
```

### Available Migration Commands

```bash
make migrate-up      # Apply all pending migrations
make migrate-down    # Rollback the last migration
make migrate-create  # Create a new migration file
```

### Database Schema

#### Products Table
```sql
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    type VARCHAR(20) NOT NULL,
    -- Type-specific JSON fields
    digital_product_info JSONB,
    physical_product_info JSONB,
    subscription_product_info JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);
```

#### Subscription Plans Table
```sql
CREATE TABLE subscription_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id),
    plan_name VARCHAR NOT NULL,
    duration INTEGER NOT NULL, -- days
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);
```

## Running the Service

### Development Mode

```bash
# Using Go directly
go run main.go server

# Using Makefile
make run
```

### Production Mode

```bash
# Build binary
make build

# Run binary
./bin/go-boilerplate server
```

### Verify Service is Running

```bash
# Check if gRPC server is listening
lsof -i :50051

# Test with grpcurl
grpcurl -plaintext localhost:50051 list
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
}' localhost:50051 product.ProductService/CreateProduct
```

```bash
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
}' localhost:50051 product.ProductService/CreateProduct
```

#### GetProduct
Retrieves a product by ID.

```bash
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000"
}' localhost:50051 product.ProductService/GetProduct
```

#### ListProducts
Lists products with pagination and optional filtering.

```bash
# List all products
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10
}' localhost:50051 product.ProductService/ListProducts

# Filter by type
grpcurl -plaintext -d '{
  "type": "DIGITAL",
  "page": 1,
  "page_size": 10
}' localhost:50051 product.ProductService/ListProducts
```

#### UpdateProduct
Updates an existing product.

```bash
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Updated Product Name",
  "price": 39.99
}' localhost:50051 product.ProductService/UpdateProduct
```

#### DeleteProduct
Soft deletes a product.

```bash
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000"
}' localhost:50051 product.ProductService/DeleteProduct
```

### Subscription Service

#### CreateSubscriptionPlan
Creates a subscription plan linked to a product.

```bash
grpcurl -plaintext -d '{
  "product_id": "550e8400-e29b-41d4-a716-446655440000",
  "plan_name": "Monthly Premium",
  "duration": 30,
  "price": 29.99
}' localhost:50051 subscription.SubscriptionService/CreateSubscriptionPlan
```

#### GetSubscriptionPlan
Retrieves a subscription plan by ID.

```bash
grpcurl -plaintext -d '{
  "id": "660e8400-e29b-41d4-a716-446655440000"
}' localhost:50051 subscription.SubscriptionService/GetSubscriptionPlan
```

#### ListSubscriptionPlans
Lists subscription plans for a product.

```bash
grpcurl -plaintext -d '{
  "product_id": "550e8400-e29b-41d4-a716-446655440000",
  "page": 1,
  "page_size": 10
}' localhost:50051 subscription.SubscriptionService/ListSubscriptionPlans
```

#### UpdateSubscriptionPlan
Updates an existing subscription plan.

```bash
grpcurl -plaintext -d '{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "plan_name": "Updated Premium Plan",
  "price": 34.99
}' localhost:50051 subscription.SubscriptionService/UpdateSubscriptionPlan
```

#### DeleteSubscriptionPlan
Soft deletes a subscription plan.

```bash
grpcurl -plaintext -d '{
  "id": "660e8400-e29b-41d4-a716-446655440000"
}' localhost:50051 subscription.SubscriptionService/DeleteSubscriptionPlan
```

## Testing

The service includes comprehensive unit tests with mock implementations.

### Run All Tests

```bash
make test
```

### Run Tests with Coverage

```bash
make test-coverage
```

### Test Coverage Report

```bash
make test-coverage-html
```

### Test Structure

- **Service Layer Tests**: `internal/service/*/service_test.go`
  - Business logic validation
  - Mock repository implementations
  - Error handling scenarios

- **Handler Layer Tests**: `internal/grpc/handlers/*_test.go`
  - gRPC endpoint testing
  - Request/response validation
  - Mock service implementations

### Example Test Execution

```bash
# Run specific package tests
go test ./internal/service/product/... -v
go test ./internal/service/subscription/... -v
go test ./internal/grpc/handlers/... -v

# Run with race detection
go test -race ./...
```

## Docker Deployment

### Quick Start with Docker Compose

```bash
# Start all services (app + postgres)
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down

# Clean up (stop services and remove volumes)
make docker-clean
```

### Manual Docker Commands

```bash
# Build Docker image
make docker-build

# Run single container (requires external database)
make docker-run
```

### Docker Environment Variables

The application supports the following environment variables for Docker deployment:

```bash
DATABASE_HOST=postgres        # Database host
DATABASE_PORT=5432           # Database port
DATABASE_USER=postgres       # Database username
DATABASE_PASSWORD=admin      # Database password
DATABASE_NAME=product_microservice  # Database name
SERVER_PORT=50051           # gRPC server port
```

These are automatically set in the `docker-compose.yml` file.

## Development

### Makefile Commands

```bash
make help          # Show available commands
make build         # Build the application
make run           # Run the application
make test          # Run tests
make clean         # Clean build artifacts
make proto         # Generate protobuf code
make lint          # Run linters
make format        # Format code
```

### Code Generation

```bash
# Generate protobuf code
make proto

# Generate mocks (if using mockery)
make mocks
```

### Adding New Features

1. **Define Protocol Buffers**: Update `.proto` files
2. **Generate Code**: Run `make proto`
3. **Implement Service Logic**: Add business logic in `internal/service/`
4. **Add Handlers**: Implement gRPC handlers in `internal/grpc/handlers/`
5. **Write Tests**: Add comprehensive unit tests
6. **Update Documentation**: Update this README

### Database Migrations

```bash
# Create new migration
make migrate-create name=add_new_table

# Apply migrations
make migrate-up

# Rollback migrations
make migrate-down
```

## Monitoring and Observability

### Health Checks

The service exposes health check endpoints:

```bash
# gRPC health check
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check
```

### Logging

Structured logging with configurable levels:

- `debug`: Detailed debugging information
- `info`: General operational messages
- `warn`: Warning conditions
- `error`: Error conditions

### Metrics

Consider adding metrics collection with:
- Prometheus metrics
- Request duration histograms
- Error rate counters
- Database connection metrics

## Troubleshooting

### Common Issues

1. **Port Already in Use**
   ```bash
   # Find process using port 50051
   lsof -i :50051
   
   # Kill process
   kill -9 <PID>
   ```

2. **Database Connection Failed**
   ```bash
   # Check PostgreSQL is running
   pg_isready -h localhost -p 5432
   
   # Verify database exists
   psql -h localhost -U postgres -l
   ```

3. **Migration Errors**
   ```bash
   # Check migration status
   make migrate-status
   
   # Force migration version
   make migrate-force version=<VERSION>
   ```

### Debug Mode

Run with debug logging:

```bash
LOG_LEVEL=debug go run main.go server
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Standards

- Follow Go conventions and best practices
- Write comprehensive tests for new features
- Update documentation for API changes
- Use meaningful commit messages

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions:

- **Issues**: [GitHub Issues](https://github.com/youngprinnce/product-microservice/issues)
- **Discussions**: [GitHub Discussions](https://github.com/youngprinnce/product-microservice/discussions)
- **Email**: [your-email@example.com]

---

**Built with ❤️ using Go, gRPC, and PostgreSQL**
│   ├── grpc/handlers/   # gRPC handlers
│   ├── service/         # Business logic layer
│   │   ├── product/     # Product domain
│   │   └── subscription/# Subscription domain
│   ├── postgres/        # Database layer
│   └── db/migrations/   # Database migrations
├── config/              # Configuration management
└── etc/                 # Configuration files
```

## Prerequisites

- Go 1.21+
- PostgreSQL 13+
- Protocol Buffers compiler (`protoc`)
- golang-migrate (for database migrations)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd product-microservice
```

2. Install dependencies:
```bash
make deps
```

3. Set up PostgreSQL database:
```bash
createdb product_microservice
```

4. Run database migrations:
```bash
make migrate-up
```

5. Build the application:
```bash
make build
```

## Configuration

Update `etc/config.yaml` with your database and server settings:

```yaml
app:
  name: "product-microservice"
  version: "1.0.0"
  env: "development"

server:
  listen: "0.0.0.0"
  port: "50051"

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  db_name: "product_microservice"
```

## Usage

### Running the Server

```bash
make run
```

The gRPC server will start on port 50051 (or as configured).

### gRPC Services

#### ProductService

- `CreateProduct` - Create a new product
- `GetProduct` - Retrieve a product by ID
- `UpdateProduct` - Update an existing product
- `DeleteProduct` - Delete a product
- `ListProducts` - List products with pagination and filtering

#### SubscriptionService

- `CreateSubscriptionPlan` - Create a new subscription plan
- `GetSubscriptionPlan` - Retrieve a subscription plan by ID
- `UpdateSubscriptionPlan` - Update an existing subscription plan
- `DeleteSubscriptionPlan` - Delete a subscription plan
- `ListSubscriptionPlans` - List subscription plans with pagination and filtering

### Testing with grpcurl

Install grpcurl for testing:
```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

List available services:
```bash
grpcurl -plaintext localhost:50051 list
```

Create a product:
```bash
grpcurl -plaintext -d '{
  "name": "Sample Digital Product",
  "description": "A sample digital product",
  "price": 29.99,
  "type": "DIGITAL",
  "digital_product": {
    "file_size": 1024000,
    "download_link": "https://example.com/download"
  }
}' localhost:50051 ProductService/CreateProduct
```

## Product Types

### Digital Products
- `file_size`: Size of the digital file in bytes
- `download_link`: URL for downloading the product

### Physical Products
- `weight`: Weight in kg
- `dimensions`: Physical dimensions (e.g., "10x5x2 cm")

### Subscription Products
- `subscription_period`: Billing period (e.g., "monthly", "yearly")
- `renewal_price`: Price for renewal

## Database Schema

### Products Table
- Supports all three product types in a single table
- Type-specific fields are nullable
- Indexed by type, name, and creation date

### Subscription Plans Table
- Links to products via foreign key
- Supports trial periods and discount percentages
- Can be active/inactive for subscription management

## Development

### Available Make Commands

```bash
make build        # Build the server binary
make run          # Run the server
make test         # Run tests
make clean        # Clean build artifacts
make proto        # Generate protobuf code
make migrate-up   # Run database migrations up
make migrate-down # Run database migrations down
make deps         # Download dependencies
make tidy         # Tidy dependencies
make fmt          # Format code
make lint         # Run linter
make help         # Show available commands
```

### Code Generation

Regenerate protobuf code after proto file changes:
```bash
make proto
```

### Database Migrations

Create new migration:
```bash
migrate create -ext sql -dir internal/db/migrations -seq migration_name
```

Apply migrations:
```bash
make migrate-up
```

Rollback migrations:
```bash
make migrate-down
```

## Testing

Run the test suite:
```bash
make test
```

## Docker

Quick start with Docker Compose:
```bash
make docker-up    # Start all services
make docker-logs  # View logs
make docker-down  # Stop services
```

Build Docker image:
```bash
make docker-build
```

Run with Docker:
```bash
make docker-run
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and ensure they pass
5. Submit a pull request

## License

This project is licensed under the MIT License.

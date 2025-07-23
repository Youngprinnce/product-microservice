# Docker Setup Guide

## Quick Start

This guide provides a simple way to run the Product Microservice using Docker.

### Prerequisites

- Docker installed on your system
- Docker Compose installed (usually comes with Docker Desktop)

### Option 1: Using Docker Compose (Recommended)

This is the easiest way to run the complete system with the database:

```bash
# Start all services (application + PostgreSQL)
make docker-up

# View real-time logs
make docker-logs

# Stop all services
make docker-down

# Clean up (remove containers and volumes)
make docker-clean
```

The application will be available on `localhost:50051` (gRPC)

### Option 2: Using Docker directly

If you want to run just the application container:

```bash
# Build the Docker image
make docker-build

# Run the container (requires external database)
make docker-run
```

### Testing the Dockerized Application

Once the services are running, you can test the API using grpcurl:

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
}' localhost:50051 product.ProductService/CreateProduct

# List products
grpcurl -plaintext -d '{"page": 1, "page_size": 10}' localhost:50051 product.ProductService/ListProducts
```

### Environment Configuration

The Docker setup uses these environment variables (already configured in docker-compose.yml):

- `DATABASE_HOST=postgres`
- `DATABASE_PORT=5432`
- `DATABASE_USER=postgres`
- `DATABASE_PASSWORD=admin`
- `DATABASE_NAME=product_microservice`
- `SERVER_PORT=50051`

### Available Make Commands

```bash
make docker-build     # Build Docker image
make docker-run       # Run single container
make docker-up        # Start all services with docker-compose
make docker-down      # Stop docker-compose services
make docker-logs      # View logs from all services
make docker-restart   # Restart docker-compose services
make docker-clean     # Stop and clean up everything
```

### Troubleshooting

1. **Port already in use**: If port 50051 or 5432 is already in use, stop other services or modify the ports in `docker-compose.yml`

2. **Permission issues**: On Linux, you might need to run docker commands with `sudo`

3. **Database connection issues**: Make sure PostgreSQL container is healthy before the app starts (docker-compose handles this automatically)

4. **Check container status**:
   ```bash
   docker-compose ps
   ```

5. **View detailed logs**:
   ```bash
   docker-compose logs product-service
   docker-compose logs postgres
   ```

That's it! You now have a fully containerized Product Microservice running with PostgreSQL.

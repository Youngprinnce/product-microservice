.PHONY: build run run-server dev test test-coverage test-coverage-html test-integration test-all test-unit clean proto migrate-up migrate-down deps tidy fmt lint mocks docker-build docker-run help

# Build the server
build:
	go build -o bin/server ./main.go

# Run the server
run:
	go run main.go server

# Run the server with specific config
run-server:
	go run main.go server -c etc/config.yaml

# Run the server in development mode
dev:
	go run main.go server -c etc/config.yaml

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./internal/...
	go tool cover -func=coverage.out

# Generate HTML coverage report
test-coverage-html:
	go test -v -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run integration tests (requires running server)
test-integration:
	cd tests && go test -v ./...

# Run all tests (unit + integration)
test-all: test test-integration

# Run tests excluding vendor and tests directory
test-unit:
	go test -v ./internal/...

# Clean build artifacts
clean:
	rm -rf bin/

# Generate protobuf code
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/*.proto

# Database migrations (requires golang-migrate)
migrate-up:
	migrate -path internal/db/migrations -database "postgresql://postgres:admin@localhost:5432/product_microservice?sslmode=disable" up

migrate-down:
	migrate -path internal/db/migrations -database "postgresql://postgres:admin@localhost:5432/product_microservice?sslmode=disable" down

# Install dependencies
deps:
	go mod download

# Tidy dependencies
tidy:
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Generate mocks (if using mockgen)
mocks:
	mockgen -source=internal/service/product/store.go -destination=mocks/product_store_mock.go
	mockgen -source=internal/service/subscription/store.go -destination=mocks/subscription_store_mock.go

# Docker commands
docker-build:
	docker build -t product-microservice .

docker-run:
	docker run -p 50051:50051 product-microservice

# Help
help:
	@echo "Available targets:"
	@echo ""
	@echo "Build & Run:"
	@echo "  build           - Build the server binary"
	@echo "  run             - Run the server"
	@echo "  run-server      - Run the server with specific config (etc/config.yaml)"
	@echo "  dev             - Run the server in development mode"
	@echo ""
	@echo "Testing:"
	@echo "  test            - Run all tests"
	@echo "  test-unit       - Run unit tests only (internal/...)"
	@echo "  test-integration- Run integration tests"
	@echo "  test-all        - Run both unit and integration tests"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  test-coverage-html - Generate HTML coverage report"
	@echo ""
	@echo "Code Quality:"
	@echo "  fmt             - Format code with go fmt"
	@echo "  lint            - Run golangci-lint"
	@echo "  tidy            - Tidy Go modules"
	@echo ""
	@echo "Development:"
	@echo "  proto           - Generate protobuf code"
	@echo "  mocks           - Generate mock files"
	@echo "  clean           - Clean build artifacts"
	@echo "  deps            - Download dependencies"
	@echo ""
	@echo "Database:"
	@echo "  migrate-up      - Run database migrations up"
	@echo "  migrate-down    - Run database migrations down"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Run Docker container"
	@echo ""
	@echo "Usage examples:"
	@echo "  make run-server     # Start server with config"
	@echo "  make test-all       # Run all tests"
	@echo "  make test-coverage  # Check test coverage"

.DEFAULT_GOAL := help

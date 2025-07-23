.PHONY: build run test clean proto docker-build docker-run docker-up docker-down docker-logs docker-clean help

# Build the server
build:
	go build -o bin/server ./main.go

# Run the server
run:
	go run main.go server

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Generate protobuf code
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/*.proto

# Docker commands
docker-build:
	docker build -t product-microservice .

docker-run:
	docker run -p 50051:50051 product-microservice

# Docker Compose commands
docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-clean:
	docker-compose down -v
	docker system prune -f

# Help
help:
	@echo "Available commands:"
	@echo ""
	@echo "Docker (Recommended):"
	@echo "  docker-up       - Start all services with docker-compose"
	@echo "  docker-down     - Stop docker-compose services"
	@echo "  docker-logs     - View logs from all services"
	@echo "  docker-clean    - Stop and clean up everything"
	@echo ""
	@echo "Local Development:"
	@echo "  build           - Build the application"
	@echo "  run             - Run the server"
	@echo "  test            - Run tests"
	@echo "  clean           - Clean build artifacts"
	@echo "  proto           - Generate protobuf code"
	@echo ""
	@echo "Quick start:"
	@echo "  make docker-up  # Start everything with Docker"

.DEFAULT_GOAL := help

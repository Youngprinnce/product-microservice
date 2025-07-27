.PHONY: build run test proto docker-up docker-down docker-logs

# Build the server
build:
	go build -o bin/server ./main.go

# Run the server
run:
	go run main.go server

# Run tests with coverage
test:
	go test -v -cover ./...

# Generate protobuf code
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/*.proto

# Docker Compose commands
docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

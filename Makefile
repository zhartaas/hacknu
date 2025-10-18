.PHONY: build run test clean docker-up docker-down docker-build

# Build the application
build:
	go build -o bin/main .

# Run the application
run:
	go run main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Start services with Docker Compose
docker-up:
	docker-compose up -d

# Stop services
docker-down:
	docker-compose down

# Build and start services
docker-build:
	docker-compose up --build -d

# Run database migrations (if you add them later)
migrate-up:
	# Add migration commands here when you implement them

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Generate API documentation (if you add swagger)
docs:
	# Add swagger generation commands here

# Development setup
dev-setup:
	cp env.example .env
	docker-compose up -d postgres
	@echo "Waiting for database to be ready..."
	@sleep 5
	@echo "Development environment ready!"
	@echo "Run 'make run' to start the application"

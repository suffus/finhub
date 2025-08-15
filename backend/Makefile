.PHONY: build run seed clean

# Build the main application
build:
	go build -o finhub-backend main.go

# Build the seed command
build-seed:
	go build -o seed cmd/seed/main.go

# Run the main application
run:
	go run main.go

# Seed the database with comprehensive sample data
seed: build-seed
	./seed

# Enhanced seeding with all picklist types
seed-full: build-seed
	./seed

# Clean up build artifacts
clean:
	rm -f finhub-backend seed

# Install dependencies
deps:
	go mod download

# Run tests
test:
	go test ./...

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# All-in-one setup
setup: deps build-seed seed 
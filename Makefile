.PHONY: all swagger build run clean up prod down help

# Default target: generate swagger docs and build the orchestrator binary.
all: swagger build

# Generate swagger documentation.
swagger:
	swag init --parseDependency --parseInternal

# Build the orchestrator binary (after generating swagger docs).
build: swagger
	$(GOPATH) build -o bin/orchestrator

# Run the orchestrator directly (after regenerating swagger docs).
run: swagger
	$(GOPATH) run main.go

# Bring up the Docker stack in development mode (no SWAG).
up:
	docker compose up -d

# Bring up the Docker stack in production mode (includes SWAG for SSL).
prod:
	docker compose --profile production up -d

# Tear down the Docker stack.
down:
	docker compose down

# Clean build artifacts and swagger docs.
clean:
	rm -rf bin/ docs/

# Display help information.
help:
	@echo "Available commands:"
	@echo "  make all     - Generate swagger docs and build the orchestrator"
	@echo "  make swagger - Generate swagger documentation"
	@echo "  make build   - Build the orchestrator binary"
	@echo "  make run     - Run the orchestrator (with swagger docs)"
	@echo "  make up      - Bring up the Docker stack (development mode)"
	@echo "  make prod    - Bring up the Docker stack in production mode (includes SWAG for SSL)"
	@echo "  make down    - Tear down the Docker stack"
	@echo "  make clean   - Remove build artifacts and docs"

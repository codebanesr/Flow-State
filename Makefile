.PHONY: all swagger build run clean up prod down help td

# Determine which docker compose command to use
DOCKER_COMPOSE := $(shell which docker-compose 2>/dev/null || echo "docker compose")

# Default target: generate swagger docs and build the orchestrator binary.
all: swagger build

# Generate swagger documentation.
swagger:
	swag init --parseDependency --parseInternal

# Build the orchestrator binary (after generating swagger docs).
build: swagger
	$(GOPATH) build -o orchestrator

# Run the orchestrator directly (after regenerating swagger docs).
run: swagger
	$(GOPATH) run main.go

# Bring up the Docker stack in development mode (no SWAG).
up:
	$(DOCKER_COMPOSE) up -d

# Bring up the Docker stack in production mode (includes SWAG for SSL).
prod:
	$(DOCKER_COMPOSE) --profile production up -d

# Tear down the Docker stack.
down:
	$(DOCKER_COMPOSE) down

# Complete Docker teardown (removes all containers, images, volumes, and networks).
td:
	docker system prune -af --volumes

# Clean build artifacts and swagger docs.
clean:
	rm -rf bin/ docs/

logs:
	docker-compose logs -f

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
	@echo "  make td      - Complete Docker teardown (removes all containers, images, volumes, and networks)"
	@echo "  make clean   - Remove build artifacts and docs"

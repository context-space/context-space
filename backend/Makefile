.PHONY: deps init clean test lint swagger docker-db docker-migrate docker-db docker-api docker-services docker-reload-api docker-logs docker-logs-f migrate-up migrate-down build build-with-version mocks generate-mocks clean-mocks test-unit test-integration test-coverage

# Build settings
BINARY_NAME=context-space-backend
BUILD_DIR=./build
MAIN_PATH=./cmd/server/main.go
GO_ENV=CGO_ENABLED=0 GO111MODULE=on

# Version information
# Use VER if provided, otherwise auto-detect from git
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
ACTUAL_VERSION = $(if $(VER),$(VER),$(VERSION))
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_VERSION ?= $(shell go version | cut -d' ' -f3)

# Build flags
LDFLAGS=-ldflags "-s -w \
	-X 'main.Version=$(ACTUAL_VERSION)' \
	-X 'main.GitCommit=$(GIT_COMMIT)' \
	-X 'main.BuildTime=$(BUILD_TIME)' \
	-X 'main.GoVersion=$(GO_VERSION)'"

# Test settings
TEST_FLAGS=-v -race -coverprofile=coverage.out -covermode=atomic

# Docker settings
DOCKER_IMAGE=context-space-backend
DOCKER_TAG=latest
COMPOSE_PROJECT_NAME=context-space

all: clean build

# Build the application with version information
# Usage: 
#   make build                    # Use git auto-detected version
#   make build VER=v1.2.3        # Use specific version
build:
	@echo "Building ${BINARY_NAME}..."
	@if [ -n "$(VER)" ]; then \
		echo "Using specified version: $(VER)"; \
	else \
		echo "Using auto-detected version: $(VERSION)"; \
	fi
	@echo "Final version: $(ACTUAL_VERSION)"
	@echo "Git Commit: ${GIT_COMMIT}"
	@echo "Build Time: ${BUILD_TIME}"
	@echo "Go Version: ${GO_VERSION}"
	${GO_ENV} go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} ${MAIN_PATH}
	@echo "Build complete: ${BUILD_DIR}/${BINARY_NAME}"

# Print version information
version:
	@echo "Auto-detected version: ${VERSION}"
	@if [ -n "$(VER)" ]; then \
		echo "Specified version: $(VER)"; \
		echo "Final version: $(ACTUAL_VERSION)"; \
	else \
		echo "Final version: $(ACTUAL_VERSION)"; \
	fi
	@echo "Git Commit: ${GIT_COMMIT}"
	@echo "Build Time: ${BUILD_TIME}"
	@echo "Go Version: ${GO_VERSION}"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/vektra/mockery/v2@latest

# Init project (useful for new setups)
init: deps
	@echo "Initializing project..."
	@mkdir -p ${BUILD_DIR}

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf ${BUILD_DIR}
	@mkdir -p ${BUILD_DIR}

# Generate all mocks
generate-mocks:
	@echo "Generating mocks..."
	@mkdir -p internal/shared/testing/mocks
	mockery
	@echo "Mocks generated successfully"

# Clean mock files
clean-mocks:
	@echo "Cleaning mocks..."
	@rm -rf internal/shared/testing/mocks
	@echo "Mocks cleaned"

# Regenerate all mocks
mocks: clean-mocks generate-mocks

# Run unit tests (fast, no external dependencies) - Service layer with full mocks
test-unit: generate-mocks
	@echo "Running unit tests (Service layer with full mocks)..."
	go test -short ${TEST_FLAGS} ./internal/credentialmanagement/application/... ./internal/integration/application/... ./internal/identityaccess/application/... ./internal/providercore/application/... ./internal/provideradapter/application/...

# Run component tests (real components with mock dependencies) - Infrastructure layer
test-component: generate-mocks
	@echo "Running component tests (Infrastructure layer with mocked deps)..."
	go test ${TEST_FLAGS} ./internal/credentialmanagement/infrastructure/... ./internal/integration/infrastructure/... ./internal/identityaccess/infrastructure/... ./internal/providercore/infrastructure/... ./internal/provideradapter/infrastructure/...

# Run integration tests (cross-layer collaboration)
test-integration: generate-mocks
	@echo "Running integration tests (Cross-layer collaboration)..."
	go test -tags=integration ${TEST_FLAGS} ./internal/credentialmanagement/interfaces/... ./internal/integration/interfaces/... ./internal/identityaccess/interfaces/... ./internal/providercore/interfaces/... ./internal/provideradapter/interfaces/...

# Run all tests with layered approach
test: test-unit test-component test-integration
	@echo "All layered tests completed!"

# Generate comprehensive coverage report with layer breakdown
test-coverage: generate-mocks
	@echo "Running comprehensive coverage analysis..."
	@echo "🔍 Unit Tests Coverage (Service layer)..."
	go test -short -cover ./internal/credentialmanagement/application/... ./internal/integration/application/... ./internal/identityaccess/application/... ./internal/providercore/application/... ./internal/provideradapter/application/... -coverprofile=coverage_unit.out
	@echo "🔍 Component Tests Coverage (Infrastructure layer)..."
	go test -cover ./internal/credentialmanagement/infrastructure/... ./internal/integration/infrastructure/... ./internal/identityaccess/infrastructure/... ./internal/providercore/infrastructure/... ./internal/provideradapter/infrastructure/... -coverprofile=coverage_component.out
	@echo "🔍 Integration Tests Coverage (Cross-layer)..."
	go test -tags=integration -cover ./internal/credentialmanagement/interfaces/... ./internal/integration/interfaces/... ./internal/identityaccess/interfaces/... ./internal/providercore/interfaces/... ./internal/provideradapter/interfaces/... -coverprofile=coverage_integration.out
	@echo "📊 Generating combined coverage report..."
	go test ${TEST_FLAGS} ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"
	@echo "📈 Coverage summary:"
	@go tool cover -func=coverage.out | tail -1
	@echo ""
	@echo "📋 Layer-specific coverage reports:"
	@echo "   - Unit Tests: coverage_unit.out"
	@echo "   - Component Tests: coverage_component.out" 
	@echo "   - Integration Tests: coverage_integration.out"

# Quick coverage check for development
test-coverage-quick: generate-mocks
	@echo "Quick coverage check..."
	go test -cover -short ./internal/credentialmanagement/... -coverprofile=coverage_quick.out
	@echo "Coverage for credentialmanagement module:"
	@go tool cover -func=coverage_quick.out | tail -1

# Run code linting
lint:
	@echo "Linting code..."
	golangci-lint run ./...

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/server/main.go -o ./docs/swagger --parseVendor --parseDependency --exclude internal/shared --parseInternal
	@echo "Swagger documentation generated in docs/swagger"

# Start only the database
docker-db:
	@echo "Starting PostgreSQL database..."
	COMPOSE_PROJECT_NAME=${COMPOSE_PROJECT_NAME} docker compose -f docker/docker-compose.yml up -d postgres pgadmin

# Run Docker migrations
docker-migrate:
	@echo "Running migrations in Docker..."
	COMPOSE_PROJECT_NAME=${COMPOSE_PROJECT_NAME} docker compose -f docker/docker-compose.migrations.yml run --rm migrations

# Start the backend API service only
docker-api:
	@echo "Starting application API..."
	COMPOSE_PROJECT_NAME=${COMPOSE_PROJECT_NAME} docker compose -f docker/docker-compose.yml build api
	COMPOSE_PROJECT_NAME=${COMPOSE_PROJECT_NAME} docker compose -f docker/docker-compose.yml up -d api

# Start supporting services (Redis, Vault, monitoring, etc.)
docker-services:
	@echo "Starting supporting services..."
	COMPOSE_PROJECT_NAME=${COMPOSE_PROJECT_NAME} docker compose -f docker/docker-compose.yml up -d redis vault otel-collector jaeger prometheus grafana

# Show logs from Docker services
docker-logs:
	@echo "Showing logs..."
	@COMPOSE_PROJECT_NAME=${COMPOSE_PROJECT_NAME} ./scripts/docker-logs.sh $(filter-out $@,$(MAKECMDGOALS))

# Show and follow logs from Docker services
docker-logs-f:
	@echo "Following logs..."
	@COMPOSE_PROJECT_NAME=${COMPOSE_PROJECT_NAME} ./scripts/docker-logs.sh -f $(filter-out $@,$(MAKECMDGOALS))

# Shutdown all Docker services
docker-down:
	@echo "Stopping all Docker services..."
	COMPOSE_PROJECT_NAME=${COMPOSE_PROJECT_NAME} docker compose -f docker/docker-compose.yml down

# Restart a chosen Docker service
docker-restart:
	@echo "Restarting the chosen Docker service..."
	@COMPOSE_PROJECT_NAME=${COMPOSE_PROJECT_NAME} ./scripts/docker-restart.sh $(filter-out $@,$(MAKECMDGOALS))

# Reload API container
docker-reload-api: swagger
	@echo "Reloading API container..."
	COMPOSE_PROJECT_NAME=${COMPOSE_PROJECT_NAME} docker compose -f docker/docker-compose.yml stop api
	COMPOSE_PROJECT_NAME=${COMPOSE_PROJECT_NAME} docker compose -f docker/docker-compose.yml build api
	COMPOSE_PROJECT_NAME=${COMPOSE_PROJECT_NAME} docker compose -f docker/docker-compose.yml up -d api

# Run migrations up
migrate-up:
	@echo "Running migrations up..."
	go run ./migrations/migrate.go up

# Run migrations down
migrate-down:
	@echo "Running migrations down..."
	go run ./migrations/migrate.go down

# Help target
help:
	@echo "Context-Space-Backend Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Build Targets:"
	@echo "  deps              - Install dependencies (including mockery)"
	@echo "  init              - Initialize project"
	@echo "  clean             - Remove build artifacts"
	@echo "  build             - Build with auto-detected version (or VER=x.x.x for specific version)"
	@echo "  build-with-version - Alias for 'build' (backward compatibility)"
	@echo "  version           - Print version information"
	@echo ""
	@echo "Testing Targets:"
	@echo "  generate-mocks     - Generate mock files using mockery"
	@echo "  clean-mocks        - Remove all generated mock files"
	@echo "  mocks              - Clean and regenerate all mocks"
	@echo "  test-unit          - Run unit tests (Service layer with full mocks)"
	@echo "  test-component     - Run component tests (Infrastructure with mock deps)"
	@echo "  test-integration   - Run integration tests (Cross-layer collaboration)"
	@echo "  test               - Run all layered tests"
	@echo "  test-coverage      - Comprehensive coverage analysis with layer breakdown"
	@echo "  test-coverage-quick - Quick coverage check for development"
	@echo "  lint               - Run linters"
	@echo ""
	@echo "Documentation Targets:"
	@echo "  swagger           - Generate Swagger API documentation"
	@echo ""
	@echo "Docker Targets:"
	@echo "  docker-db         - Start only the database"
	@echo "  docker-migrate    - Run migrations in Docker"
	@echo "  docker-api        - Start only the application"
	@echo "  docker-services   - Start supporting services"
	@echo "  docker-logs       - Show logs from Docker services"
	@echo "  docker-logs-f     - Show and follow logs from Docker services"
	@echo "  docker-down       - Shutdown all Docker services"
	@echo "  docker-restart    - Restart a chosen Docker service"
	@echo "  docker-reload-api - Reload API container"
	@echo ""
	@echo "Migration Targets:"
	@echo "  migrate-up        - Run migrations up"
	@echo "  migrate-down      - Run migrations down"

# Special rule to handle arguments for docker-logs
%:
	@:

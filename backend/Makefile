# =============================================================================
# Makefile - Go Backend
# =============================================================================

# Build variables
APP_NAME := go-copilot
BUILD_DIR := ./bin
MAIN_PATH := ./cmd/api
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION := $(shell go version | cut -d' ' -f3)
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GoVersion=$(GO_VERSION)"

# Docker variables
DOCKER_COMPOSE := docker-compose -f ../docker/docker-compose.yml

# Database variables
DB_DSN ?= postgres://postgres:postgres@localhost:5432/app_dev?sslmode=disable
MIGRATIONS_DIR := ./migrations

# Colors for output
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m

.PHONY: help
.DEFAULT_GOAL := help

# =============================================================================
# Help
# =============================================================================

help: ## Show this help message
	@echo "$(GREEN)Go Copilot Backend - Available Commands$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n\n"} \
		/^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Setup

setup: ## First-time project setup (install tools, copy env files)
	@echo "$(GREEN)Setting up project...$(NC)"
	@cp -n .env.example .env 2>/dev/null || true
	@echo "$(GREEN)Installing Go tools...$(NC)"
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/vektra/mockery/v2@latest
	@echo "$(GREEN)Downloading Go modules...$(NC)"
	go mod download
	go mod verify
	@echo "$(GREEN)Setup complete!$(NC)"

##@ Development

dev: ## Start docker-compose services (PostgreSQL, Redis)
	@echo "$(GREEN)Starting development services...$(NC)"
	$(DOCKER_COMPOSE) up -d postgres redis
	@echo "$(GREEN)Services started!$(NC)"

dev-down: ## Stop docker-compose services
	@echo "$(YELLOW)Stopping development services...$(NC)"
	$(DOCKER_COMPOSE) down
	@echo "$(GREEN)Services stopped!$(NC)"

dev-logs: ## Show docker-compose logs
	$(DOCKER_COMPOSE) logs -f

run: ## Run application in development mode
	@echo "$(GREEN)Starting application...$(NC)"
	go run $(MAIN_PATH)/main.go

run-watch: ## Run application with hot reload (requires air)
	@command -v air > /dev/null 2>&1 || { echo "$(RED)air is not installed. Run: go install github.com/cosmtrek/air@latest$(NC)"; exit 1; }
	air

##@ Build

build: ## Build production binary with version injection
	@echo "$(GREEN)Building $(APP_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(APP_NAME)$(NC)"

build-linux: ## Build for Linux (useful for Docker)
	@echo "$(GREEN)Building $(APP_NAME) for Linux...$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(MAIN_PATH)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(APP_NAME)-linux-amd64$(NC)"

clean: ## Remove build artifacts
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	rm -rf $(BUILD_DIR)
	rm -rf tmp/
	go clean -cache -testcache
	@echo "$(GREEN)Clean complete!$(NC)"

##@ Testing

test: ## Run unit tests with coverage
	@echo "$(GREEN)Running unit tests...$(NC)"
	go test -v -race -cover -coverprofile=coverage.out ./...
	@echo "$(GREEN)Tests complete!$(NC)"

test-short: ## Run unit tests (short mode)
	go test -v -short ./...

test-integration: ## Run integration tests with test database
	@echo "$(GREEN)Running integration tests...$(NC)"
	go test -v -race -tags=integration ./...
	@echo "$(GREEN)Integration tests complete!$(NC)"

test-coverage: test ## Generate and open coverage report
	@echo "$(GREEN)Generating coverage report...$(NC)"
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report: coverage.html$(NC)"

benchmark: ## Run benchmarks
	go test -bench=. -benchmem ./...

##@ Code Quality

lint: ## Run golangci-lint
	@echo "$(GREEN)Running linter...$(NC)"
	@command -v golangci-lint > /dev/null 2>&1 || { echo "$(RED)golangci-lint is not installed. Run: make setup$(NC)"; exit 1; }
	golangci-lint run ./...
	@echo "$(GREEN)Lint complete!$(NC)"

lint-fix: ## Run golangci-lint with auto-fix
	golangci-lint run --fix ./...

fmt: ## Format code (gofmt + goimports)
	@echo "$(GREEN)Formatting code...$(NC)"
	gofmt -s -w .
	@command -v goimports > /dev/null 2>&1 && goimports -w -local github.com/tranvuongduy2003/go-copilot . || true
	@echo "$(GREEN)Format complete!$(NC)"

vet: ## Run go vet
	go vet ./...

check: fmt lint vet test ## Run all checks (format, lint, vet, test)
	@echo "$(GREEN)All checks passed!$(NC)"

##@ Database Migrations

migrate-create: ## Create a new migration (usage: make migrate-create name=create_users_table)
	@if [ -z "$(name)" ]; then echo "$(RED)Error: name is required. Usage: make migrate-create name=migration_name$(NC)"; exit 1; fi
	goose -dir $(MIGRATIONS_DIR) create $(name) sql
	@echo "$(GREEN)Migration created!$(NC)"

migrate-up: ## Run all pending migrations
	@echo "$(GREEN)Running migrations...$(NC)"
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up
	@echo "$(GREEN)Migrations complete!$(NC)"

migrate-down: ## Rollback last migration
	@echo "$(YELLOW)Rolling back last migration...$(NC)"
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down
	@echo "$(GREEN)Rollback complete!$(NC)"

migrate-down-all: ## Rollback all migrations
	@echo "$(RED)Rolling back ALL migrations...$(NC)"
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down-to 0
	@echo "$(GREEN)Rollback complete!$(NC)"

migrate-status: ## Show migration status
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" status

migrate-reset: migrate-down-all migrate-up ## Reset database (rollback all and re-run)
	@echo "$(GREEN)Database reset complete!$(NC)"

##@ Code Generation

generate: ## Run code generation (mocks, wire, etc.)
	@echo "$(GREEN)Running code generation...$(NC)"
	go generate ./...
	@echo "$(GREEN)Generation complete!$(NC)"

mock: ## Generate mocks using mockery
	@echo "$(GREEN)Generating mocks...$(NC)"
	@command -v mockery > /dev/null 2>&1 || { echo "$(RED)mockery is not installed. Run: make setup$(NC)"; exit 1; }
	mockery --all --keeptree --dir=internal --output=internal/mocks
	@echo "$(GREEN)Mocks generated!$(NC)"

wire: ## Run wire for dependency injection
	@echo "$(GREEN)Running wire...$(NC)"
	@command -v wire > /dev/null 2>&1 || { echo "$(RED)wire is not installed. Run: make setup$(NC)"; exit 1; }
	wire ./...
	@echo "$(GREEN)Wire complete!$(NC)"

##@ Docker

docker-build: ## Build Docker image
	docker build -t $(APP_NAME):$(VERSION) -f ../docker/Dockerfile.backend.go ..

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env $(APP_NAME):$(VERSION)

##@ Dependencies

deps: ## Download dependencies
	go mod download

deps-update: ## Update dependencies
	go get -u ./...
	go mod tidy

deps-tidy: ## Tidy dependencies
	go mod tidy
	go mod verify

deps-vendor: ## Vendor dependencies
	go mod vendor

##@ Documentation

docs: ## Generate API documentation
	@echo "$(GREEN)Generating API documentation...$(NC)"
	@command -v swag > /dev/null 2>&1 || go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g $(MAIN_PATH)/main.go -o ./docs
	@echo "$(GREEN)Documentation generated!$(NC)"

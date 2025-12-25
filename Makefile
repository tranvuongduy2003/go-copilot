# =============================================================================
# Makefile - Development and DevOps Commands
# =============================================================================
# Usage: make <target>
# Run `make help` to see all available commands
# =============================================================================

.PHONY: help install dev build test lint format clean docker-* k8s-* tf-* monitoring-*

# Default shell
SHELL := /bin/bash

# Variables
PROJECT_NAME := fullstack-app
DOCKER_REGISTRY := ghcr.io
IMAGE_TAG := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
DOCKER_COMPOSE := docker compose
K8S_NAMESPACE := $(PROJECT_NAME)

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m

# =============================================================================
# Help
# =============================================================================

help: ## Show this help message
	@echo ""
	@echo "$(BLUE)$(PROJECT_NAME) - Development Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""

# =============================================================================
# Development
# =============================================================================

install: ## Install all dependencies
	@echo "$(BLUE)Installing dependencies...$(NC)"
	cd frontend && pnpm install
	cd backend && pnpm install
	@echo "$(GREEN)Dependencies installed!$(NC)"

dev: ## Start development environment
	@echo "$(BLUE)Starting development environment...$(NC)"
	$(DOCKER_COMPOSE) -f docker/docker-compose.yml up -d postgres redis
	@sleep 3
	@make -j2 dev-frontend dev-backend

dev-frontend: ## Start frontend development server
	cd frontend && pnpm dev

dev-backend: ## Start backend development server
	cd backend && pnpm dev

dev-docker: ## Start full development environment with Docker
	$(DOCKER_COMPOSE) -f docker/docker-compose.yml up -d

dev-stop: ## Stop development environment
	$(DOCKER_COMPOSE) -f docker/docker-compose.yml down

dev-logs: ## View development logs
	$(DOCKER_COMPOSE) -f docker/docker-compose.yml logs -f

# =============================================================================
# Build
# =============================================================================

build: build-frontend build-backend ## Build all services

build-frontend: ## Build frontend
	@echo "$(BLUE)Building frontend...$(NC)"
	cd frontend && pnpm build
	@echo "$(GREEN)Frontend built!$(NC)"

build-backend: ## Build backend
	@echo "$(BLUE)Building backend...$(NC)"
	cd backend && pnpm build
	@echo "$(GREEN)Backend built!$(NC)"

# =============================================================================
# Testing
# =============================================================================

test: test-frontend test-backend ## Run all tests

test-frontend: ## Run frontend tests
	@echo "$(BLUE)Running frontend tests...$(NC)"
	cd frontend && pnpm test

test-backend: ## Run backend tests
	@echo "$(BLUE)Running backend tests...$(NC)"
	cd backend && pnpm test

test-e2e: ## Run E2E tests
	@echo "$(BLUE)Running E2E tests...$(NC)"
	cd e2e && pnpm test

test-ci: ## Run tests in CI mode
	@echo "$(BLUE)Running CI tests...$(NC)"
	cd frontend && pnpm test:ci
	cd backend && pnpm test:ci

test-coverage: ## Run tests with coverage
	cd frontend && pnpm test:coverage
	cd backend && pnpm test:coverage

# =============================================================================
# Code Quality
# =============================================================================

lint: lint-frontend lint-backend ## Lint all code

lint-frontend: ## Lint frontend code
	cd frontend && pnpm lint

lint-backend: ## Lint backend code
	cd backend && pnpm lint

lint-fix: ## Fix linting issues
	cd frontend && pnpm lint:fix
	cd backend && pnpm lint:fix

format: ## Format all code
	cd frontend && pnpm format
	cd backend && pnpm format

typecheck: ## Run type checking
	cd frontend && pnpm type-check
	cd backend && pnpm type-check

# =============================================================================
# Docker
# =============================================================================

docker-build: ## Build all Docker images
	@echo "$(BLUE)Building Docker images...$(NC)"
	docker build -f docker/Dockerfile.frontend -t $(PROJECT_NAME)-frontend:$(IMAGE_TAG) .
	docker build -f docker/Dockerfile.backend -t $(PROJECT_NAME)-backend:$(IMAGE_TAG) .
	docker build -f docker/Dockerfile.nginx -t $(PROJECT_NAME)-nginx:$(IMAGE_TAG) .
	@echo "$(GREEN)Docker images built!$(NC)"

docker-push: ## Push Docker images to registry
	@echo "$(BLUE)Pushing Docker images...$(NC)"
	docker tag $(PROJECT_NAME)-frontend:$(IMAGE_TAG) $(DOCKER_REGISTRY)/$(PROJECT_NAME)/frontend:$(IMAGE_TAG)
	docker tag $(PROJECT_NAME)-backend:$(IMAGE_TAG) $(DOCKER_REGISTRY)/$(PROJECT_NAME)/backend:$(IMAGE_TAG)
	docker tag $(PROJECT_NAME)-nginx:$(IMAGE_TAG) $(DOCKER_REGISTRY)/$(PROJECT_NAME)/nginx:$(IMAGE_TAG)
	docker push $(DOCKER_REGISTRY)/$(PROJECT_NAME)/frontend:$(IMAGE_TAG)
	docker push $(DOCKER_REGISTRY)/$(PROJECT_NAME)/backend:$(IMAGE_TAG)
	docker push $(DOCKER_REGISTRY)/$(PROJECT_NAME)/nginx:$(IMAGE_TAG)
	@echo "$(GREEN)Docker images pushed!$(NC)"

docker-up: ## Start production Docker environment
	$(DOCKER_COMPOSE) -f docker/docker-compose.prod.yml up -d

docker-down: ## Stop Docker environment
	$(DOCKER_COMPOSE) -f docker/docker-compose.yml down
	$(DOCKER_COMPOSE) -f docker/docker-compose.prod.yml down

docker-clean: ## Clean Docker resources
	docker system prune -af
	docker volume prune -f

docker-test: ## Run tests in Docker
	$(DOCKER_COMPOSE) -f docker/docker-compose.test.yml up --abort-on-container-exit --exit-code-from test-runner

# =============================================================================
# Database
# =============================================================================

db-migrate: ## Run database migrations
	cd backend && pnpm db:migrate

db-migrate-create: ## Create a new migration
	@read -p "Migration name: " name; \
	cd backend && pnpm db:migrate:create $$name

db-seed: ## Seed the database
	cd backend && pnpm db:seed

db-reset: ## Reset the database
	cd backend && pnpm db:reset

db-studio: ## Open database studio
	cd backend && pnpm db:studio

# =============================================================================
# Kubernetes
# =============================================================================

k8s-apply-staging: ## Apply Kubernetes manifests to staging
	@echo "$(BLUE)Deploying to staging...$(NC)"
	kubectl apply -k k8s/overlays/staging
	@echo "$(GREEN)Deployed to staging!$(NC)"

k8s-apply-prod: ## Apply Kubernetes manifests to production
	@echo "$(YELLOW)Deploying to production...$(NC)"
	kubectl apply -k k8s/overlays/production
	@echo "$(GREEN)Deployed to production!$(NC)"

k8s-delete-staging: ## Delete staging deployment
	kubectl delete -k k8s/overlays/staging

k8s-status: ## Check Kubernetes status
	kubectl get pods,svc,ingress -n $(K8S_NAMESPACE)

k8s-logs: ## View Kubernetes logs
	kubectl logs -f -l app=$(PROJECT_NAME) -n $(K8S_NAMESPACE)

k8s-port-forward: ## Port forward services locally
	kubectl port-forward svc/frontend 3000:80 -n $(K8S_NAMESPACE) &
	kubectl port-forward svc/backend 8080:8080 -n $(K8S_NAMESPACE) &

# =============================================================================
# Terraform
# =============================================================================

tf-init: ## Initialize Terraform
	cd terraform && terraform init

tf-plan-staging: ## Plan Terraform changes for staging
	cd terraform && terraform plan -var-file=environments/staging.tfvars -out=tfplan

tf-plan-prod: ## Plan Terraform changes for production
	cd terraform && terraform plan -var-file=environments/production.tfvars -out=tfplan

tf-apply: ## Apply Terraform changes
	cd terraform && terraform apply tfplan

tf-destroy-staging: ## Destroy staging infrastructure
	cd terraform && terraform destroy -var-file=environments/staging.tfvars

tf-output: ## Show Terraform outputs
	cd terraform && terraform output

# =============================================================================
# Monitoring
# =============================================================================

monitoring-up: ## Start monitoring stack
	$(DOCKER_COMPOSE) -f monitoring/docker-compose.monitoring.yml up -d

monitoring-down: ## Stop monitoring stack
	$(DOCKER_COMPOSE) -f monitoring/docker-compose.monitoring.yml down

monitoring-logs: ## View monitoring logs
	$(DOCKER_COMPOSE) -f monitoring/docker-compose.monitoring.yml logs -f

# =============================================================================
# Security
# =============================================================================

security-scan: ## Run security scans
	@echo "$(BLUE)Running security scans...$(NC)"
	trivy fs --severity HIGH,CRITICAL .
	gitleaks detect --source . --verbose

security-audit: ## Run dependency audit
	cd frontend && pnpm audit
	cd backend && pnpm audit

# =============================================================================
# Cleanup
# =============================================================================

clean: ## Clean build artifacts
	rm -rf frontend/.next frontend/out frontend/node_modules
	rm -rf backend/dist backend/node_modules
	rm -rf coverage
	rm -rf .turbo

clean-all: clean docker-clean ## Clean everything including Docker

# =============================================================================
# Utilities
# =============================================================================

version: ## Show current version
	@echo "$(GREEN)Version: $(IMAGE_TAG)$(NC)"

env-check: ## Check environment setup
	@echo "Checking environment..."
	@command -v node >/dev/null 2>&1 && echo "$(GREEN)Node.js: $$(node -v)$(NC)" || echo "$(RED)Node.js: Not found$(NC)"
	@command -v pnpm >/dev/null 2>&1 && echo "$(GREEN)pnpm: $$(pnpm -v)$(NC)" || echo "$(RED)pnpm: Not found$(NC)"
	@command -v docker >/dev/null 2>&1 && echo "$(GREEN)Docker: $$(docker -v)$(NC)" || echo "$(RED)Docker: Not found$(NC)"
	@command -v kubectl >/dev/null 2>&1 && echo "$(GREEN)kubectl: $$(kubectl version --client --short 2>/dev/null)$(NC)" || echo "$(YELLOW)kubectl: Not found$(NC)"
	@command -v terraform >/dev/null 2>&1 && echo "$(GREEN)Terraform: $$(terraform -v | head -1)$(NC)" || echo "$(YELLOW)Terraform: Not found$(NC)"

generate-types: ## Generate TypeScript types
	cd backend && pnpm generate:types
	cp backend/src/types/api.d.ts frontend/src/types/

docs: ## Generate documentation
	cd backend && pnpm docs:generate

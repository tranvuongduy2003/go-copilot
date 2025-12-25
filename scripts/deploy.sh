#!/bin/bash
# =============================================================================
# Deployment Script
# =============================================================================
# Usage: ./scripts/deploy.sh [staging|production] [--dry-run]
# =============================================================================

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
ENVIRONMENT="${1:-staging}"
DRY_RUN="${2:-}"
DOCKER_REGISTRY="${DOCKER_REGISTRY:-ghcr.io}"
IMAGE_NAME="${IMAGE_NAME:-myorg/myapp}"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Functions
log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
log_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# Validate environment
validate_environment() {
    case "$ENVIRONMENT" in
        staging|production)
            log_info "Deploying to: $ENVIRONMENT"
            ;;
        *)
            log_error "Invalid environment: $ENVIRONMENT. Use 'staging' or 'production'"
            ;;
    esac
}

# Pre-deployment checks
pre_deploy_checks() {
    log_info "Running pre-deployment checks..."

    # Check for uncommitted changes
    if [ -n "$(git status --porcelain)" ]; then
        log_warning "You have uncommitted changes"
        if [ "$ENVIRONMENT" == "production" ]; then
            log_error "Cannot deploy to production with uncommitted changes"
        fi
    fi

    # Verify kubectl access
    if ! kubectl cluster-info >/dev/null 2>&1; then
        log_error "Cannot connect to Kubernetes cluster"
    fi

    # Verify namespace exists
    kubectl get namespace "$ENVIRONMENT" >/dev/null 2>&1 || \
        kubectl create namespace "$ENVIRONMENT"

    log_success "Pre-deployment checks passed!"
}

# Build and push images
build_and_push() {
    log_info "Building Docker images (version: $VERSION)..."

    # Build images
    docker build -f docker/Dockerfile.frontend -t "$DOCKER_REGISTRY/$IMAGE_NAME/frontend:$VERSION" .
    docker build -f docker/Dockerfile.backend -t "$DOCKER_REGISTRY/$IMAGE_NAME/backend:$VERSION" .
    docker build -f docker/Dockerfile.nginx -t "$DOCKER_REGISTRY/$IMAGE_NAME/nginx:$VERSION" .

    if [ "$DRY_RUN" != "--dry-run" ]; then
        log_info "Pushing Docker images..."
        docker push "$DOCKER_REGISTRY/$IMAGE_NAME/frontend:$VERSION"
        docker push "$DOCKER_REGISTRY/$IMAGE_NAME/backend:$VERSION"
        docker push "$DOCKER_REGISTRY/$IMAGE_NAME/nginx:$VERSION"
        log_success "Images pushed!"
    else
        log_warning "Dry run - skipping image push"
    fi
}

# Run database migrations
run_migrations() {
    log_info "Running database migrations..."

    if [ "$DRY_RUN" != "--dry-run" ]; then
        kubectl exec -n "$ENVIRONMENT" deployment/backend -- npm run db:migrate
        log_success "Migrations complete!"
    else
        log_warning "Dry run - skipping migrations"
    fi
}

# Deploy to Kubernetes
deploy() {
    log_info "Deploying to Kubernetes..."

    # Update image tags in kustomization
    cd "k8s/overlays/$ENVIRONMENT"

    # Set images
    kustomize edit set image \
        "frontend=$DOCKER_REGISTRY/$IMAGE_NAME/frontend:$VERSION" \
        "backend=$DOCKER_REGISTRY/$IMAGE_NAME/backend:$VERSION" \
        "nginx=$DOCKER_REGISTRY/$IMAGE_NAME/nginx:$VERSION"

    cd -

    if [ "$DRY_RUN" == "--dry-run" ]; then
        log_info "Dry run - showing what would be applied:"
        kubectl apply -k "k8s/overlays/$ENVIRONMENT" --dry-run=client
    else
        kubectl apply -k "k8s/overlays/$ENVIRONMENT"

        # Wait for rollout
        log_info "Waiting for rollout to complete..."
        kubectl rollout status deployment/frontend -n "$ENVIRONMENT" --timeout=300s
        kubectl rollout status deployment/backend -n "$ENVIRONMENT" --timeout=300s

        log_success "Deployment complete!"
    fi
}

# Health check
health_check() {
    log_info "Running health checks..."

    local max_attempts=30
    local attempt=1

    while [ $attempt -le $max_attempts ]; do
        if kubectl exec -n "$ENVIRONMENT" deployment/backend -- wget -q --spider http://localhost:8080/health; then
            log_success "Health check passed!"
            return 0
        fi
        log_info "Waiting for services... (attempt $attempt/$max_attempts)"
        sleep 10
        ((attempt++))
    done

    log_error "Health check failed after $max_attempts attempts"
}

# Rollback
rollback() {
    log_warning "Rolling back deployment..."
    kubectl rollout undo deployment/frontend -n "$ENVIRONMENT"
    kubectl rollout undo deployment/backend -n "$ENVIRONMENT"
    log_success "Rollback complete!"
}

# Notify
notify() {
    local status="$1"
    local message="Deployment to $ENVIRONMENT: $status (version: $VERSION)"

    # Slack notification (if webhook is configured)
    if [ -n "${SLACK_WEBHOOK_URL:-}" ]; then
        curl -s -X POST -H 'Content-type: application/json' \
            --data "{\"text\":\"$message\"}" \
            "$SLACK_WEBHOOK_URL" || true
    fi

    log_info "$message"
}

# Cleanup
cleanup() {
    local exit_code=$?
    if [ $exit_code -ne 0 ]; then
        log_error "Deployment failed!"
        notify "FAILED"
        if [ "$ENVIRONMENT" == "staging" ]; then
            rollback
        fi
    fi
}

# Main
main() {
    trap cleanup EXIT

    echo ""
    echo "=========================================="
    echo "  Deployment Script"
    echo "  Environment: $ENVIRONMENT"
    echo "  Version: $VERSION"
    if [ "$DRY_RUN" == "--dry-run" ]; then
        echo "  Mode: DRY RUN"
    fi
    echo "=========================================="
    echo ""

    validate_environment
    pre_deploy_checks
    build_and_push
    deploy

    if [ "$DRY_RUN" != "--dry-run" ]; then
        run_migrations
        health_check
        notify "SUCCESS"
    fi

    log_success "Deployment to $ENVIRONMENT completed successfully!"
}

main

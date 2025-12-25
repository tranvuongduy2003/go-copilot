#!/usr/bin/env bash
# =============================================================================
# Development Environment Startup Script
# =============================================================================
# Usage: ./scripts/dev.sh [options]
#   --frontend    Start only frontend
#   --backend     Start only backend
#   --docker      Start full Docker environment
#   --services    Start only infrastructure services (DB, Redis)
#   --help        Show this help message
# =============================================================================

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Project root directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Docker compose file
DOCKER_COMPOSE_FILE="$PROJECT_ROOT/docker/docker-compose.yml"

# =============================================================================
# Helper Functions
# =============================================================================

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_dependencies() {
    local missing=()

    if ! command -v docker &> /dev/null; then
        missing+=("docker")
    fi

    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        missing+=("docker-compose")
    fi

    if ! command -v node &> /dev/null; then
        missing+=("node")
    fi

    if ! command -v go &> /dev/null; then
        missing+=("go")
    fi

    if [ ${#missing[@]} -ne 0 ]; then
        log_error "Missing dependencies: ${missing[*]}"
        exit 1
    fi
}

docker_compose() {
    if docker compose version &> /dev/null; then
        docker compose "$@"
    else
        docker-compose "$@"
    fi
}

wait_for_service() {
    local service=$1
    local max_attempts=${2:-30}
    local attempt=1

    log_info "Waiting for $service to be ready..."

    while [ $attempt -le $max_attempts ]; do
        if docker_compose -f "$DOCKER_COMPOSE_FILE" ps "$service" 2>/dev/null | grep -q "healthy\|running"; then
            log_success "$service is ready"
            return 0
        fi
        sleep 1
        ((attempt++))
    done

    log_error "$service failed to start within $max_attempts seconds"
    return 1
}

# =============================================================================
# Service Functions
# =============================================================================

start_infrastructure() {
    log_info "Starting infrastructure services (PostgreSQL, Redis)..."
    docker_compose -f "$DOCKER_COMPOSE_FILE" up -d postgres redis

    wait_for_service "postgres"
    wait_for_service "redis"

    log_success "Infrastructure services are running"
}

start_frontend() {
    log_info "Starting frontend development server..."
    cd "$PROJECT_ROOT/frontend"

    if [ ! -d "node_modules" ]; then
        log_info "Installing frontend dependencies..."
        if command -v pnpm &> /dev/null; then
            pnpm install
        else
            npm install
        fi
    fi

    if command -v pnpm &> /dev/null; then
        pnpm dev
    else
        npm run dev
    fi
}

start_backend() {
    log_info "Starting backend development server..."
    cd "$PROJECT_ROOT/backend"

    # Check if it's a Go project
    if [ -f "go.mod" ]; then
        log_info "Detected Go backend"

        # Run migrations if goose is available
        if command -v goose &> /dev/null; then
            log_info "Running database migrations..."
            goose -dir migrations/sql postgres "${DATABASE_URL:-postgresql://postgres:postgres@localhost:5432/app_dev}" up || true
        fi

        # Start with air for hot reload if available, otherwise use go run
        if command -v air &> /dev/null; then
            air
        else
            go run cmd/api/main.go
        fi
    else
        # Fallback to Node.js backend
        if [ ! -d "node_modules" ]; then
            log_info "Installing backend dependencies..."
            if command -v pnpm &> /dev/null; then
                pnpm install
            else
                npm install
            fi
        fi

        if command -v pnpm &> /dev/null; then
            pnpm dev
        else
            npm run dev
        fi
    fi
}

start_docker_full() {
    log_info "Starting full Docker development environment..."
    docker_compose -f "$DOCKER_COMPOSE_FILE" up -d
    log_success "All services are running"

    echo ""
    log_info "Services:"
    echo "  Frontend:   http://localhost:3000"
    echo "  Backend:    http://localhost:8080"
    echo "  PostgreSQL: localhost:5432"
    echo "  Redis:      localhost:6379"
    echo ""
    log_info "Run 'docker compose -f docker/docker-compose.yml logs -f' to view logs"
}

start_all() {
    start_infrastructure

    echo ""
    log_info "Starting frontend and backend in parallel..."
    echo ""

    # Start frontend and backend in background
    (start_frontend) &
    FRONTEND_PID=$!

    (start_backend) &
    BACKEND_PID=$!

    # Handle shutdown
    trap "kill $FRONTEND_PID $BACKEND_PID 2>/dev/null; exit" SIGINT SIGTERM

    wait
}

show_help() {
    echo "Development Environment Startup Script"
    echo ""
    echo "Usage: ./scripts/dev.sh [options]"
    echo ""
    echo "Options:"
    echo "  --frontend    Start only frontend development server"
    echo "  --backend     Start only backend development server"
    echo "  --docker      Start full Docker environment"
    echo "  --services    Start only infrastructure services (DB, Redis)"
    echo "  --help        Show this help message"
    echo ""
    echo "Without options, starts infrastructure + frontend + backend"
}

# =============================================================================
# Main
# =============================================================================

main() {
    cd "$PROJECT_ROOT"

    case "${1:-}" in
        --frontend)
            check_dependencies
            start_frontend
            ;;
        --backend)
            check_dependencies
            start_backend
            ;;
        --docker)
            check_dependencies
            start_docker_full
            ;;
        --services)
            check_dependencies
            start_infrastructure
            ;;
        --help|-h)
            show_help
            ;;
        "")
            check_dependencies
            start_all
            ;;
        *)
            log_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
}

main "$@"

#!/bin/bash
# =============================================================================
# Project Setup Script
# =============================================================================
# Usage: ./scripts/setup.sh [--dev|--prod|--ci]
# =============================================================================

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Functions
log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
log_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."

    local missing=()

    command -v node >/dev/null 2>&1 || missing+=("node")
    command -v pnpm >/dev/null 2>&1 || missing+=("pnpm")
    command -v docker >/dev/null 2>&1 || missing+=("docker")

    if [ ${#missing[@]} -ne 0 ]; then
        log_error "Missing prerequisites: ${missing[*]}"
        log_info "Please install the missing tools before continuing."
        exit 1
    fi

    log_success "All prerequisites met!"
}

# Setup environment files
setup_env() {
    log_info "Setting up environment files..."

    if [ ! -f .env ]; then
        cp .env.example .env
        log_success "Created .env from .env.example"
    else
        log_warning ".env already exists, skipping..."
    fi

    if [ ! -f frontend/.env.local ]; then
        cp frontend/.env.example frontend/.env.local 2>/dev/null || echo "NEXT_PUBLIC_API_URL=http://localhost:8080" > frontend/.env.local
        log_success "Created frontend/.env.local"
    fi

    if [ ! -f backend/.env ]; then
        cp backend/.env.example backend/.env 2>/dev/null || cat > backend/.env << EOF
NODE_ENV=development
PORT=8080
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/app_dev
REDIS_URL=redis://localhost:6379
JWT_SECRET=dev-secret-change-in-production-$(openssl rand -hex 32)
EOF
        log_success "Created backend/.env"
    fi
}

# Install dependencies
install_deps() {
    log_info "Installing dependencies..."

    # Frontend
    if [ -d frontend ]; then
        log_info "Installing frontend dependencies..."
        cd frontend && pnpm install && cd ..
    fi

    # Backend
    if [ -d backend ]; then
        log_info "Installing backend dependencies..."
        cd backend && pnpm install && cd ..
    fi

    # E2E tests
    if [ -d e2e ]; then
        log_info "Installing E2E dependencies..."
        cd e2e && pnpm install && cd ..
    fi

    log_success "Dependencies installed!"
}

# Setup git hooks
setup_hooks() {
    log_info "Setting up git hooks..."

    if [ -d .git ]; then
        # Install husky if available
        if command -v pnpm >/dev/null 2>&1; then
            pnpm exec husky install 2>/dev/null || log_warning "Husky not configured"
        fi

        # Create pre-commit hook
        cat > .git/hooks/pre-commit << 'EOF'
#!/bin/sh
# Pre-commit hook

# Run linting
make lint || exit 1

# Run type checking
make typecheck || exit 1

# Check for secrets
if command -v gitleaks >/dev/null 2>&1; then
    gitleaks protect --staged --verbose || exit 1
fi
EOF
        chmod +x .git/hooks/pre-commit
        log_success "Git hooks configured!"
    else
        log_warning "Not a git repository, skipping hooks setup"
    fi
}

# Setup Docker environment
setup_docker() {
    log_info "Setting up Docker environment..."

    # Check if Docker is running
    if ! docker info >/dev/null 2>&1; then
        log_error "Docker is not running. Please start Docker and try again."
        exit 1
    fi

    # Create Docker network if not exists
    docker network create app-network 2>/dev/null || true

    # Pull base images
    log_info "Pulling base images..."
    docker pull node:20-alpine
    docker pull postgres:16-alpine
    docker pull redis:7-alpine

    log_success "Docker environment ready!"
}

# Setup database
setup_database() {
    log_info "Setting up database..."

    # Start PostgreSQL container
    docker compose -f docker/docker-compose.yml up -d postgres

    # Wait for PostgreSQL to be ready
    log_info "Waiting for PostgreSQL..."
    sleep 5

    # Run migrations
    if [ -d backend ]; then
        cd backend
        pnpm db:migrate 2>/dev/null || log_warning "No migrations to run"
        pnpm db:seed 2>/dev/null || log_warning "No seeds to run"
        cd ..
    fi

    log_success "Database setup complete!"
}

# Main setup
main() {
    local mode="${1:-dev}"

    echo ""
    echo "=========================================="
    echo "  Project Setup - Mode: $mode"
    echo "=========================================="
    echo ""

    check_prerequisites

    case "$mode" in
        --dev|dev)
            setup_env
            install_deps
            setup_hooks
            setup_docker
            setup_database
            ;;
        --prod|prod)
            setup_env
            install_deps
            ;;
        --ci|ci)
            install_deps
            ;;
        *)
            log_error "Unknown mode: $mode"
            echo "Usage: $0 [--dev|--prod|--ci]"
            exit 1
            ;;
    esac

    echo ""
    log_success "Setup complete!"
    echo ""
    echo "Next steps:"
    echo "  1. Review and update .env files"
    echo "  2. Run 'make dev' to start development"
    echo "  3. Visit http://localhost:3000"
    echo ""
}

main "$@"

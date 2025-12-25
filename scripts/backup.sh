#!/bin/bash
# =============================================================================
# Database Backup Script
# =============================================================================
# Usage: ./scripts/backup.sh [--restore <backup_file>]
# =============================================================================

set -euo pipefail

# Configuration
BACKUP_DIR="${BACKUP_DIR:-./backups}"
RETENTION_DAYS="${RETENTION_DAYS:-30}"
S3_BUCKET="${S3_BUCKET:-}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_NAME="backup_${TIMESTAMP}"

# Database connection (from environment or defaults)
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-app}"
DB_USER="${DB_USER:-postgres}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
log_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# Create backup directory
setup() {
    mkdir -p "$BACKUP_DIR"
    log_info "Backup directory: $BACKUP_DIR"
}

# Backup PostgreSQL
backup_postgres() {
    log_info "Creating PostgreSQL backup..."

    local backup_file="$BACKUP_DIR/${BACKUP_NAME}_postgres.sql.gz"

    PGPASSWORD="${DB_PASSWORD:-}" pg_dump \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        --format=custom \
        --compress=9 \
        -f "$backup_file"

    log_success "PostgreSQL backup created: $backup_file"
    echo "$backup_file"
}

# Backup Redis
backup_redis() {
    log_info "Creating Redis backup..."

    local redis_host="${REDIS_HOST:-localhost}"
    local redis_port="${REDIS_PORT:-6379}"
    local backup_file="$BACKUP_DIR/${BACKUP_NAME}_redis.rdb"

    # Trigger BGSAVE and wait
    redis-cli -h "$redis_host" -p "$redis_port" BGSAVE

    sleep 5

    # Copy RDB file
    if [ -f "/var/lib/redis/dump.rdb" ]; then
        cp /var/lib/redis/dump.rdb "$backup_file"
        log_success "Redis backup created: $backup_file"
    else
        log_warning "Redis RDB file not found, skipping..."
    fi
}

# Upload to S3
upload_to_s3() {
    if [ -n "$S3_BUCKET" ]; then
        log_info "Uploading backups to S3..."

        aws s3 cp "$BACKUP_DIR/" "s3://$S3_BUCKET/backups/" \
            --recursive \
            --exclude "*" \
            --include "${BACKUP_NAME}*"

        log_success "Backups uploaded to S3!"
    fi
}

# Cleanup old backups
cleanup_old_backups() {
    log_info "Cleaning up backups older than $RETENTION_DAYS days..."

    # Local cleanup
    find "$BACKUP_DIR" -type f -name "backup_*" -mtime +"$RETENTION_DAYS" -delete

    # S3 cleanup (if configured)
    if [ -n "$S3_BUCKET" ]; then
        local cutoff_date=$(date -d "$RETENTION_DAYS days ago" +%Y-%m-%d)
        aws s3 ls "s3://$S3_BUCKET/backups/" | while read -r line; do
            local file_date=$(echo "$line" | awk '{print $1}')
            local file_name=$(echo "$line" | awk '{print $4}')
            if [[ "$file_date" < "$cutoff_date" ]]; then
                aws s3 rm "s3://$S3_BUCKET/backups/$file_name"
                log_info "Deleted old backup: $file_name"
            fi
        done
    fi

    log_success "Cleanup complete!"
}

# Restore from backup
restore() {
    local backup_file="$1"

    if [ ! -f "$backup_file" ]; then
        log_error "Backup file not found: $backup_file"
    fi

    log_warning "This will restore the database from: $backup_file"
    read -p "Are you sure you want to continue? (yes/no): " confirm

    if [ "$confirm" != "yes" ]; then
        log_info "Restore cancelled"
        exit 0
    fi

    log_info "Restoring database..."

    PGPASSWORD="${DB_PASSWORD:-}" pg_restore \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        --clean \
        --if-exists \
        "$backup_file"

    log_success "Database restored from: $backup_file"
}

# List backups
list_backups() {
    log_info "Local backups:"
    ls -lh "$BACKUP_DIR"/backup_* 2>/dev/null || echo "No local backups found"

    if [ -n "$S3_BUCKET" ]; then
        log_info "S3 backups:"
        aws s3 ls "s3://$S3_BUCKET/backups/" 2>/dev/null || echo "No S3 backups found"
    fi
}

# Main
main() {
    case "${1:-backup}" in
        --restore)
            restore "${2:-}"
            ;;
        --list)
            list_backups
            ;;
        --cleanup)
            cleanup_old_backups
            ;;
        backup|*)
            setup
            backup_postgres
            backup_redis
            upload_to_s3
            cleanup_old_backups

            log_success "Backup completed successfully!"
            log_info "Backup name: $BACKUP_NAME"
            ;;
    esac
}

main "$@"

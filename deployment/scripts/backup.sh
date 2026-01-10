#!/bin/bash

# Backup script for FullStack application
# Usage: ./backup.sh

set -e

BACKUP_DIR="/opt/fullstack/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

echo "========================================="
echo "Starting Backup Process"
echo "Timestamp: $TIMESTAMP"
echo "========================================="

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

# Load environment variables
if [ -f "$PROJECT_ROOT/.env.production" ]; then
    export $(cat "$PROJECT_ROOT/.env.production" | grep -v '^#' | xargs)
fi

# Backup PostgreSQL database
echo "Backing up PostgreSQL database..."
docker exec fullstack-postgres-prod pg_dump -U ${DB_USER:-postgres} ${DB_NAME:-fullstack} | gzip > "$BACKUP_DIR/db_backup_$TIMESTAMP.sql.gz"

# Backup Redis data
echo "Backing up Redis data..."
docker exec fullstack-redis-prod redis-cli SAVE
docker cp fullstack-redis-prod:/data/dump.rdb "$BACKUP_DIR/redis_backup_$TIMESTAMP.rdb"

# Backup application files (if needed)
echo "Backing up application files..."
tar -czf "$BACKUP_DIR/app_backup_$TIMESTAMP.tar.gz" \
    -C "$PROJECT_ROOT" \
    --exclude='node_modules' \
    --exclude='dist' \
    --exclude='.git' \
    --exclude='*.log' \
    .

# Remove backups older than 7 days
echo "Cleaning up old backups..."
find "$BACKUP_DIR" -name "*.gz" -mtime +7 -delete
find "$BACKUP_DIR" -name "*.rdb" -mtime +7 -delete

# Calculate backup size
BACKUP_SIZE=$(du -sh "$BACKUP_DIR" | cut -f1)

echo "========================================="
echo "Backup completed successfully!"
echo "========================================="
echo "Backup location: $BACKUP_DIR"
echo "Total backup size: $BACKUP_SIZE"
echo ""
echo "Backup files created:"
echo "  - db_backup_$TIMESTAMP.sql.gz"
echo "  - redis_backup_$TIMESTAMP.rdb"
echo "  - app_backup_$TIMESTAMP.tar.gz"

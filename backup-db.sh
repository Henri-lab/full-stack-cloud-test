#!/bin/bash

# FreeGemini 数据库备份脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
BACKUP_DIR="./backups"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="fullstack"
DB_USER="postgres"
DB_HOST="localhost"
DB_PORT="5432"

echo "🗄️  Starting database backup..."
echo ""

# 创建备份目录
mkdir -p "$BACKUP_DIR"

# 备份文件名
BACKUP_FILE="$BACKUP_DIR/${DB_NAME}_${DATE}.sql"
BACKUP_FILE_GZ="$BACKUP_FILE.gz"

# 执行备份
echo "Backing up database: $DB_NAME"
if PGPASSWORD=postgres pg_dump -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" "$DB_NAME" > "$BACKUP_FILE"; then
    echo -e "${GREEN}✅ Database backup created: $BACKUP_FILE${NC}"

    # 压缩备份文件
    echo "Compressing backup..."
    gzip "$BACKUP_FILE"
    echo -e "${GREEN}✅ Backup compressed: $BACKUP_FILE_GZ${NC}"

    # 显示备份文件大小
    SIZE=$(du -h "$BACKUP_FILE_GZ" | cut -f1)
    echo "Backup size: $SIZE"
else
    echo -e "${RED}❌ Backup failed!${NC}"
    exit 1
fi

# 清理旧备份（保留最近7天）
echo ""
echo "Cleaning up old backups (keeping last 7 days)..."
find "$BACKUP_DIR" -name "*.sql.gz" -mtime +7 -delete
echo -e "${GREEN}✅ Old backups cleaned up${NC}"

# 列出所有备份
echo ""
echo "Available backups:"
ls -lh "$BACKUP_DIR"/*.sql.gz 2>/dev/null || echo "No backups found"

echo ""
echo -e "${GREEN}🎉 Backup completed successfully!${NC}"

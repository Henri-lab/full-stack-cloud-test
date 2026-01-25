#!/bin/bash

# FreeGemini 数据库恢复脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
BACKUP_DIR="./backups"
DB_NAME="fullstack"
DB_USER="postgres"
DB_HOST="localhost"
DB_PORT="5432"

echo "🔄 Database Restore Tool"
echo ""

# 检查备份目录
if [ ! -d "$BACKUP_DIR" ]; then
    echo -e "${RED}❌ Backup directory not found: $BACKUP_DIR${NC}"
    exit 1
fi

# 列出可用备份
echo "Available backups:"
echo ""
BACKUPS=($(ls -t "$BACKUP_DIR"/*.sql.gz 2>/dev/null))

if [ ${#BACKUPS[@]} -eq 0 ]; then
    echo -e "${RED}❌ No backup files found in $BACKUP_DIR${NC}"
    exit 1
fi

# 显示备份列表
for i in "${!BACKUPS[@]}"; do
    BACKUP_FILE="${BACKUPS[$i]}"
    BACKUP_NAME=$(basename "$BACKUP_FILE")
    BACKUP_SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
    BACKUP_DATE=$(stat -f "%Sm" -t "%Y-%m-%d %H:%M:%S" "$BACKUP_FILE" 2>/dev/null || stat -c "%y" "$BACKUP_FILE" 2>/dev/null | cut -d'.' -f1)
    echo "[$i] $BACKUP_NAME ($BACKUP_SIZE) - $BACKUP_DATE"
done

echo ""
echo -n "Select backup to restore (0-$((${#BACKUPS[@]}-1))): "
read SELECTION

# 验证选择
if ! [[ "$SELECTION" =~ ^[0-9]+$ ]] || [ "$SELECTION" -ge ${#BACKUPS[@]} ]; then
    echo -e "${RED}❌ Invalid selection${NC}"
    exit 1
fi

SELECTED_BACKUP="${BACKUPS[$SELECTION]}"
echo ""
echo -e "${YELLOW}⚠️  WARNING: This will replace the current database!${NC}"
echo "Selected backup: $(basename "$SELECTED_BACKUP")"
echo ""
echo -n "Are you sure you want to continue? (yes/no): "
read CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    echo "Restore cancelled."
    exit 0
fi

# 解压备份文件
echo ""
echo "Decompressing backup..."
TEMP_SQL="/tmp/restore_$(date +%s).sql"
gunzip -c "$SELECTED_BACKUP" > "$TEMP_SQL"

# 删除现有数据库
echo "Dropping existing database..."
PGPASSWORD=postgres dropdb -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" "$DB_NAME" --if-exists

# 创建新数据库
echo "Creating new database..."
PGPASSWORD=postgres createdb -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" "$DB_NAME"

# 恢复数据
echo "Restoring database..."
if PGPASSWORD=postgres psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" < "$TEMP_SQL" > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Database restored successfully!${NC}"
else
    echo -e "${RED}❌ Restore failed!${NC}"
    rm -f "$TEMP_SQL"
    exit 1
fi

# 清理临时文件
rm -f "$TEMP_SQL"

echo ""
echo -e "${GREEN}🎉 Restore completed successfully!${NC}"
echo ""
echo "Database: $DB_NAME"
echo "Restored from: $(basename "$SELECTED_BACKUP")"

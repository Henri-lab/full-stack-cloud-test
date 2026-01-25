#!/bin/bash

# FreeGemini 停止脚本

set -e

echo "🛑 Stopping FreeGemini services..."
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 从 PID 文件读取进程 ID
if [ -f "logs/backend.pid" ]; then
    BACKEND_PID=$(cat logs/backend.pid)
    if kill -0 $BACKEND_PID 2>/dev/null; then
        echo "Stopping backend (PID: $BACKEND_PID)..."
        kill $BACKEND_PID
        echo -e "${GREEN}✅ Backend stopped${NC}"
    else
        echo -e "${YELLOW}⚠️  Backend is not running${NC}"
    fi
    rm -f logs/backend.pid
else
    echo -e "${YELLOW}⚠️  No backend PID file found${NC}"
fi

if [ -f "logs/frontend.pid" ]; then
    FRONTEND_PID=$(cat logs/frontend.pid)
    if kill -0 $FRONTEND_PID 2>/dev/null; then
        echo "Stopping frontend (PID: $FRONTEND_PID)..."
        kill $FRONTEND_PID
        echo -e "${GREEN}✅ Frontend stopped${NC}"
    else
        echo -e "${YELLOW}⚠️  Frontend is not running${NC}"
    fi
    rm -f logs/frontend.pid
else
    echo -e "${YELLOW}⚠️  No frontend PID file found${NC}"
fi

# 清理可能残留的进程
echo ""
echo "Cleaning up any remaining processes..."
pkill -f "go run cmd/api/main.go" 2>/dev/null || true
pkill -f "npm run dev" 2>/dev/null || true
pkill -f "vite" 2>/dev/null || true

echo ""
echo -e "${GREEN}🎉 All services stopped successfully!${NC}"

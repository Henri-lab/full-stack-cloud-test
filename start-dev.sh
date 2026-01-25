#!/bin/bash

# FreeGemini 开发环境启动脚本

set -e

echo "🚀 Starting FreeGemini Development Environment..."
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed. Please install Go 1.24 or later.${NC}"
    exit 1
fi

# 检查 Node.js 是否安装
if ! command -v node &> /dev/null; then
    echo -e "${RED}❌ Node.js is not installed. Please install Node.js 18 or later.${NC}"
    exit 1
fi

# 检查 PostgreSQL 是否运行
if ! pg_isready -h localhost -p 5432 &> /dev/null; then
    echo -e "${YELLOW}⚠️  PostgreSQL is not running on localhost:5432${NC}"
    echo "Please start PostgreSQL or update DATABASE_URL in deployment/.env"
    echo ""
fi

# 检查后端依赖
echo "📦 Checking backend dependencies..."
cd backend
if [ ! -d "vendor" ] && [ ! -f "go.sum" ]; then
    echo "Installing Go dependencies..."
    go mod download
fi
cd ..

# 检查前端依赖
echo "📦 Checking frontend dependencies..."
cd frontend
if [ ! -d "node_modules" ]; then
    echo "Installing Node.js dependencies..."
    npm install
fi
cd ..

echo ""
echo -e "${GREEN}✅ All dependencies are ready!${NC}"
echo ""

# 创建日志目录
mkdir -p logs

# 启动后端
echo "🔧 Starting backend server..."
cd backend
go run cmd/api/main.go > ../logs/backend.log 2>&1 &
BACKEND_PID=$!
echo "Backend PID: $BACKEND_PID"
cd ..

# 等待后端启动
echo "⏳ Waiting for backend to start..."
for i in {1..30}; do
    if curl -s http://localhost:8080/api/health > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Backend is ready!${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}❌ Backend failed to start. Check logs/backend.log${NC}"
        kill $BACKEND_PID 2>/dev/null || true
        exit 1
    fi
    sleep 1
done

# 启动前端
echo "🎨 Starting frontend server..."
cd frontend
npm run dev > ../logs/frontend.log 2>&1 &
FRONTEND_PID=$!
echo "Frontend PID: $FRONTEND_PID"
cd ..

# 等待前端启动
echo "⏳ Waiting for frontend to start..."
sleep 3

echo ""
echo -e "${GREEN}🎉 FreeGemini is now running!${NC}"
echo ""
echo "📍 Services:"
echo "   - Frontend: http://localhost:3000"
echo "   - Backend:  http://localhost:8080"
echo "   - API Docs: http://localhost:8080/api/health"
echo ""
echo "📝 Logs:"
echo "   - Backend:  logs/backend.log"
echo "   - Frontend: logs/frontend.log"
echo ""
echo "🛑 To stop all services, run:"
echo "   kill $BACKEND_PID $FRONTEND_PID"
echo ""
echo "💡 Tips:"
echo "   - View backend logs: tail -f logs/backend.log"
echo "   - View frontend logs: tail -f logs/frontend.log"
echo "   - Check API health: curl http://localhost:8080/api/health"
echo ""

# 保存 PIDs 到文件
echo "$BACKEND_PID" > logs/backend.pid
echo "$FRONTEND_PID" > logs/frontend.pid

echo "Press Ctrl+C to stop all services..."
echo ""

# 等待用户中断
trap "echo ''; echo 'Stopping services...'; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null || true; rm -f logs/*.pid; echo 'Services stopped.'; exit 0" INT TERM

# 保持脚本运行
wait

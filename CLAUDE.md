# Claude Code 项目文档

> 此文档供 Claude Code 快速理解项目结构和当前进度

## 项目概述

**FreeGemini** - 一个完整的全栈应用，包含用户认证、任务管理和邮箱账号管理功能。

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | React 18 + TypeScript + Vite + Tailwind CSS |
| 后端 | Go 1.24 + Gin + GORM |
| 数据库 | PostgreSQL 15 |
| 缓存 | Redis 7 |
| 部署 | Docker + Nginx |
| 监控 | Prometheus + Grafana |

## 目录结构

```
fullStack/
├── backend/
│   ├── cmd/api/main.go              # 入口文件
│   └── internal/
│       ├── config/config.go         # 配置管理
│       ├── database/database.go     # 数据库连接
│       ├── handlers/
│       │   ├── auth.go              # 认证处理 (注册/登录/登出)
│       │   └── task.go              # 任务 CRUD
│       ├── middleware/middleware.go # CORS/Auth/Logger 中间件
│       └── models/models.go         # User/Task 模型
├── frontend/
│   └── src/
│       ├── pages/
│       │   ├── Home.tsx             # 首页
│       │   ├── Login.tsx            # 登录页
│       │   ├── Register.tsx         # 注册页
│       │   ├── Tasks.tsx            # 任务管理
│       │   └── Emails.tsx           # 邮箱管理 (含 TOTP)
│       ├── services/api.ts          # Axios API 客户端
│       └── resource/emails.json     # 邮箱数据模板（真实数据在数据库）
└── deployment/
    ├── docker-compose.yml           # 开发环境
    └── docker-compose.prod.yml      # 生产环境
```

## API 端点

### 认证 (无需 Token)
- `POST /api/v1/auth/register` - 注册
- `POST /api/v1/auth/login` - 登录
- `POST /api/v1/auth/logout` - 登出

### 任务 (需要 JWT Token)
- `GET /api/v1/tasks` - 获取所有任务
- `GET /api/v1/tasks/:id` - 获取单个任务
- `POST /api/v1/tasks` - 创建任务
- `PUT /api/v1/tasks/:id` - 更新任务
- `DELETE /api/v1/tasks/:id` - 删除任务

### 邮箱 (需要 JWT Token)
- `GET /api/v1/emails` - 获取所有邮箱
- `GET /api/v1/emails/:id` - 获取单个邮箱
- `POST /api/v1/emails` - 创建邮箱
- `PUT /api/v1/emails/:id` - 更新邮箱
- `DELETE /api/v1/emails/:id` - 删除邮箱

### 健康检查
- `GET /api/health` - 服务状态

## 数据模型

### User
```go
ID, Username, Email, PasswordHash, CreatedAt, UpdatedAt, DeletedAt
```

### Task
```go
ID, Title, Description, Status(open/in_progress/completed), CreatorID, CreatedAt, UpdatedAt, DeletedAt
```

### Email
```go
ID, Main, Password, Deputy, Key2FA, Banned, Price, Sold, NeedRepair, Source, CreatedAt, UpdatedAt, DeletedAt
```

## 安全特性 (已实现)

| 特性 | 状态 | 说明 |
|------|------|------|
| JWT 认证 | ✅ | 7天过期，生产环境必须设置 JWT_SECRET |
| 密码强度 | ✅ | 12位+大小写+数字+特殊字符 |
| 登录限流 | ✅ | 5次失败后封禁IP 15分钟 |
| 注册限流 | ✅ | 复用登录限流器 |
| CORS 限制 | ✅ | 生产环境必须设置 CORS_ORIGIN |
| 密码哈希 | ✅ | bcrypt cost=12 |
| 安全头 | ✅ | X-Frame-Options, X-Content-Type-Options, X-XSS-Protection |

## 环境变量

### 生产环境必需
```bash
ENVIRONMENT=production
JWT_SECRET=<至少32字符的密钥>
CORS_ORIGIN=https://yourdomain.com
DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=require
```

### 开发环境默认值
```bash
ENVIRONMENT=development
JWT_SECRET=<自动生成随机密钥>
CORS_ORIGIN=http://localhost:3000,http://localhost:5173
DATABASE_URL=postgres://postgres:postgres@localhost:5432/fullstack?sslmode=disable
PORT=8080
```

## 最近更新

### 2025-01-25 安全加固
- [x] JWT 密钥：生产环境强制设置，开发环境随机生成
- [x] 密码要求：从8位提升到12位，增加特殊字符要求
- [x] 注册限流：防止批量注册攻击
- [x] CORS：生产环境强制配置，拒绝非白名单来源

### 已有功能
- [x] 用户注册/登录/登出
- [x] JWT Token 认证
- [x] 任务 CRUD
- [x] 邮箱管理界面
- [x] TOTP 动态码生成
- [x] Docker 部署配置
- [x] Prometheus + Grafana 监控

## 开发命令

```bash
# 后端
cd backend && go run cmd/api/main.go

# 前端
cd frontend && npm run dev

# Docker 部署
cd deployment && docker-compose up -d

# 构建检查
cd backend && go build ./...
cd frontend && npm run build
```

## 关键文件位置

| 功能 | 文件 |
|------|------|
| 认证逻辑 | `backend/internal/handlers/auth.go` |
| 密码验证 | `backend/internal/handlers/auth.go:111-155` |
| 限流器 | `backend/internal/handlers/auth.go:17-90` |
| CORS 中间件 | `backend/internal/middleware/middleware.go:13-65` |
| JWT 中间件 | `backend/internal/middleware/middleware.go:96-129` |
| 配置加载 | `backend/internal/config/config.go` |
| 前端 API | `frontend/src/services/api.ts` |
| 邮箱数据模板 | `frontend/src/resource/emails.json` |
| 邮箱 API | `backend/internal/handlers/email.go` |

## 待办事项

- [ ] 添加 refresh token 机制
- [x] 实现邮箱数据后端存储
- [ ] 添加用户角色权限
- [ ] 实现审计日志
- [ ] 添加 2FA 登录支持

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
│       ├── config/config.go         # 配置管理（JWT 自动生成）
│       ├── database/database.go     # 数据库连接和迁移
│       ├── handlers/
│       │   ├── auth.go              # 认证处理 (注册/登录/登出)
│       │   ├── task.go              # 任务 CRUD
│       │   ├── email.go             # 邮箱 CRUD + 批量导入
│       │   └── payment.go           # 支付和 License Key 管理
│       ├── middleware/
│       │   ├── middleware.go        # CORS/Auth/Logger 中间件
│       │   └── license.go           # License Key 验证中间件
│       └── models/models.go         # User/Task/Email/EmailFamily/Payment/LicenseKey 模型
├── frontend/
│   └── src/
│       ├── pages/
│       │   ├── Home.tsx             # 首页
│       │   ├── Login.tsx            # 登录页
│       │   ├── Register.tsx         # 注册页
│       │   ├── Tasks.tsx            # 任务管理
│       │   ├── Emails.tsx           # 邮箱管理 (含 TOTP + 导入 + License Key)
│       │   └── Payment.tsx          # 支付和 License Key 管理
│       ├── services/api.ts          # Axios API 客户端
│       └── resource/emails.json     # 邮箱数据模板（真实数据在数据库）
└── deployment/
    ├── .env                         # 环境变量配置
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
- `GET /api/v1/emails` - 获取所有邮箱（含 family 邮箱）
- `GET /api/v1/emails/:id` - 获取单个邮箱
- `POST /api/v1/emails` - 创建邮箱
- `POST /api/v1/emails/import` - 批量导入邮箱（JSON 文件上传，需 License Key）
- `POST /api/v1/emails/verify` - 批量验证邮箱状态（需 License Key）
- `PUT /api/v1/emails/:id` - 更新邮箱
- `DELETE /api/v1/emails/:id` - 删除邮箱

### 支付和 License Key (需要 JWT Token)
- `GET /api/v1/payments/products` - 获取产品列表
- `POST /api/v1/payments/orders` - 创建支付订单
- `GET /api/v1/payments/orders/:order_no` - 获取订单详情
- `POST /api/v1/payments/notify` - 支付回调（模拟）
- `GET /api/v1/keys` - 获取我的密钥列表
- `POST /api/v1/keys/activate` - 激活密钥
- `POST /api/v1/keys/check` - 检查密钥状态

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
ID, Main, Password, Deputy, Key2FA, Status, Banned, Price, Sold, NeedRepair, Source, Familys, CreatedAt, UpdatedAt, DeletedAt
```

**Status 字段**: unknown, live, verify, dead

### EmailFamily
```go
ID, EmailID, Email, Password, Code, Contact, Issue, DeletedAt
```

### Payment
```go
ID, UserID, OrderNo, Amount, ProductType, QuotaAmount, Status, PaymentMethod, TransactionID, PaidAt, ExpiredAt, CreatedAt, UpdatedAt, DeletedAt
```

**Status 字段**: pending, paid, expired, refunded

### LicenseKey
```go
ID, UserID, PaymentID, KeyCode, ProductType, QuotaTotal, QuotaUsed, Status, ActivatedAt, CreatedAt, UpdatedAt, DeletedAt
```

**Status 字段**: active, exhausted, revoked

**ProductType 字段**: basic, pro, enterprise

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

### 2025-01-25 付费功能 (License Key 系统)
- [x] **按次数付费模式** - 用户购买指定次数的使用额度
  - 基础版：10元 / 100次验证
  - 专业版：30元 / 500次验证
  - 企业版：50元 / 1000次验证
- [x] **在线支付接口** - 支持支付宝/微信支付
  - 订单管理系统
  - 支付回调处理
  - 15分钟订单过期机制
- [x] **License Key 管理**
  - 自动生成格式化密钥（XXXX-XXXX-XXXX-XXXX）
  - 密钥激活和验证
  - 额度追踪和消耗
  - 状态管理（active/exhausted/revoked）
- [x] **功能分级权限**
  - 基础版：邮箱验证
  - 专业版：邮箱验证 + 导入 + 任务管理
  - 企业版：所有功能 + API 访问
- [x] **前端支付页面**
  - 产品选择和对比
  - 收款码展示（支持真实图片）
  - 订单状态追踪
  - 我的密钥管理
- [x] **邮箱验证集成**
  - 验证功能需要 License Key
  - 自动消耗额度
  - 实时显示剩余额度
  - Key 输入和验证界面
- [x] **自动权限生效**
  - 支付成功自动设置当前 Key
  - 前端请求自动携带 `X-License-Key`
  - 邮箱导入与任务管理需有效 Key

### 2025-01-25 批量复制功能
- [x] **批量复制邮箱地址** - 一键复制选中的邮箱到剪贴板
  - 表头显示 "Copy (N)" 按钮（N 是选中数量）
  - 自动换行格式，每个邮箱占一行
  - 成功提示消息，3 秒后自动消失
  - 方便在其他网站（如 gmailver.com）使用
- [x] **改进的选择功能**
  - 表头复选框全选/取消全选
  - 行复选框单独选择
  - 实时显示选中数量

### 2025-01-25 SMTP 邮箱验证功能
- [x] **自己实现的 SMTP 验证** - 无需依赖第三方 API
  - 直接通过 SMTP 协议验证邮箱可用性
  - 支持 Gmail, Outlook, Yahoo 等所有邮箱
  - 多端口尝试（25, 587, 465）
  - STARTTLS 加密支持
  - 多 MX 服务器尝试
- [x] **双验证模式** - 支持 SMTP 和 API 两种验证方式
  - SMTP 验证：无需 key，直接验证（适合少量验证）
  - API 验证：需要从 gmailver.com 获取 key（速度快，推荐）
- [x] **前端验证方法选择**
  - 单选按钮切换验证方式
  - SMTP 模式无需输入 key
  - API 模式显示 key 输入框
- [x] **修复第三方 API 响应解析** - 正确处理 gmailver.com 返回的数据格式
- [x] **SMTP 故障排查文档** - 详细的问题诊断和解决方案

### 2025-01-25 邮箱验证功能
- [x] 邮箱验证状态字段：添加 status 字段（unknown/live/verify/dead）
- [x] 批量验证接口：`POST /api/v1/emails/verify`
  - 支持批量验证邮箱状态
  - 集成第三方 API（gmailver.com）
  - 自动更新数据库状态
- [x] 前端验证功能：
  - Key 输入框（从 gmailver.com 获取）
  - 邮箱选择功能（单选/全选）
  - 验证状态显示（Live/Verify/Dead 徽章）
  - 实时状态更新

### 2025-01-25 邮箱批量导入功能
- [x] EmailFamily 模型：支持每个邮箱关联多个 family 邮箱
- [x] 批量导入接口：`POST /api/v1/emails/import`
  - 支持 JSON 文件上传
  - 事务处理，遇到重复邮箱回滚
  - 同时导入主邮箱和 family 邮箱
- [x] 前端上传组件：Import JSON 按钮
- [x] 导入状态提示：成功/失败消息显示
- [x] 修复 JWT 认证 bug：AuthHandler 现在从配置接收 jwtSecret

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
- [x] 邮箱批量导入（JSON 文件上传）
- [x] 邮箱批量验证（live/verify/dead 状态）
- [x] EmailFamily 关联管理
- [x] TOTP 动态码生成
- [x] Docker 部署配置
- [x] Prometheus + Grafana 监控
- [x] **License Key 付费系统**
- [x] **按次数付费模式**
- [x] **在线支付接口（支付宝/微信）**
- [x] **功能分级权限控制**

## 开发命令

```bash
# 后端（会自动生成 JWT_SECRET）
cd backend && go run cmd/api/main.go

# 前端
cd frontend && npm run dev

# 使用启动脚本（推荐）
./start-dev.sh  # 启动所有服务
./stop-dev.sh   # 停止所有服务

# 数据库备份和恢复
./backup-db.sh   # 备份数据库
./restore-db.sh  # 恢复数据库

# Docker 部署
cd deployment && docker-compose up -d

# 构建检查
cd backend && go build ./...
cd frontend && npm run build

# 数据库迁移（自动执行）
# 启动后端时会自动运行 AutoMigrate
```

## 测试流程

### 1. 启动服务
```bash
# 终端 1 - 启动后端
cd backend
go run cmd/api/main.go
# 输出: WARNING: JWT_SECRET not set, using random secret...
# 输出: Server starting on port 8080

# 终端 2 - 启动前端
cd frontend
npm run dev
# 输出: Local: http://localhost:3000
```

### 2. 测试认证
1. 访问 http://localhost:3000
2. 注册新用户（密码需要 12 位+大小写+数字+特殊字符）
3. 登录成功后跳转到 Tasks 页面

### 3. 测试邮箱导入
1. 登录后访问 Emails 页面
2. 准备 JSON 文件（参考上面的格式）
3. 点击 "Import JSON" 按钮上传
4. 查看导入结果提示
5. 刷新页面查看导入的邮箱数据

### 4. 测试付费功能
1. 登录后点击导航栏的 "License Key"
2. 在 "购买 Key" 标签页选择产品（基础版/专业版/企业版）
3. 点击 "立即购买" 创建订单
4. 扫描收款码或点击 "模拟支付" 按钮
5. 支付成功后自动生成 License Key
6. 在 "我的 Key" 标签页查看已购买的密钥
7. 复制 Key Code 用于邮箱验证功能

### 5. 测试邮箱验证（需要 License Key）
1. 在 Emails 页面点击 "Verify Emails" 按钮
2. 输入 License Key 并点击 "Check Key" 验证
3. 选择验证方法（SMTP 或 API）
4. 勾选要验证的邮箱
5. 点击 "Verify" 按钮开始验证
6. 验证完成后查看状态更新和剩余额度

## 关键文件位置

| 功能 | 文件 |
|------|------|
| 认证逻辑 | `backend/internal/handlers/auth.go` |
| AuthHandler 构造 | `backend/internal/handlers/auth.go:96-100` |
| 密码验证 | `backend/internal/handlers/auth.go:111-155` |
| 限流器 | `backend/internal/handlers/auth.go:17-90` |
| CORS 中间件 | `backend/internal/middleware/middleware.go:13-65` |
| JWT 中间件 | `backend/internal/middleware/middleware.go:96-129` |
| License Key 中间件 | `backend/internal/middleware/license.go` |
| 配置加载 | `backend/internal/config/config.go` |
| JWT 自动生成 | `backend/internal/config/config.go:28-31` |
| 前端 API | `frontend/src/services/api.ts` |
| 邮箱数据模板 | `frontend/src/resource/emails.json` |
| 邮箱 API | `backend/internal/handlers/email.go` |
| 邮箱导入接口 | `backend/internal/handlers/email.go:298-403` |
| 邮箱验证接口 | `backend/internal/handlers/email.go:421-470` |
| SMTP 验证器 | `backend/internal/handlers/email_verify_smtp.go` |
| SMTP 验证逻辑 | `backend/internal/handlers/email_verify_smtp.go:18-95` |
| API 验证逻辑 | `backend/internal/handlers/email.go:473-520` |
| 支付 API | `backend/internal/handlers/payment.go` |
| 支付模型 | `backend/internal/models/models.go:72-95` |
| License Key 模型 | `backend/internal/models/models.go:97-110` |
| 邮箱模型 | `backend/internal/models/models.go:31-58` |
| 前端导入组件 | `frontend/src/pages/Emails.tsx:157-182` |
| 前端验证组件 | `frontend/src/pages/Emails.tsx:205-260` |
| 前端验证方法选择 | `frontend/src/pages/Emails.tsx:360-410` |
| 前端支付页面 | `frontend/src/pages/Payment.tsx` |
| 前端 License Key 输入 | `frontend/src/pages/Emails.tsx:520-575` |

## 待办事项

- [ ] 添加 refresh token 机制
- [x] 实现邮箱数据后端存储
- [ ] 添加用户角色权限
- [ ] 实现审计日志
- [ ] 添加 2FA 登录支持

## 邮箱批量导入 JSON 格式

```json
{
  "emails": [
    {
      "main": "test@gmail.com",
      "password": "TestPass123!",
      "deputy": "backup@gmail.com",
      "key_2FA": "ABCD1234EFGH5678",
      "meta": {
        "banned": false,
        "price": 0,
        "sold": false,
        "need_repair": false,
        "from": "source1"
      },
      "familys": [
        {
          "email": "family1@gmail.com",
          "password": "FamilyPass123!",
          "code": "123456",
          "contact": "qq:123456;phone:13800138000",
          "issue": "正常使用"
        }
      ]
    }
  ]
}
```

### 导入规则
- 遇到重复邮箱（main 字段相同）会报错终止
- 使用事务处理，失败时回滚所有更改
- familys 数组可以为空
- meta 字段可选，不提供则使用默认值

## 常见问题

### License Key 相关

**问题**：邮箱验证时提示 "需要有效的 License Key"

**原因**：邮箱验证功能需要购买并激活 License Key

**解决方案**：
1. 访问 `/payment` 页面购买 License Key
2. 选择合适的产品套餐（基础版/专业版/企业版）
3. 完成支付后获取 Key Code
4. 在 Emails 页面的验证区域输入 Key Code
5. 点击 "Check Key" 验证密钥有效性
6. 验证成功后即可使用邮箱验证功能

**问题**：License Key 额度用尽

**原因**：每次邮箱验证会消耗 1 次额度，当额度用尽时密钥状态变为 exhausted

**解决方案**：
- 购买新的 License Key
- 或升级到更高级别的套餐（更多额度）

### JWT 认证失败
**问题**：登录时返回 500 错误，日志显示 "JWT_SECRET environment variable is required"

**原因**：之前的代码在运行时从环境变量读取 JWT_SECRET，但 Go 不会自动加载 `.env` 文件

**解决方案**（已修复）：
- AuthHandler 现在从配置对象接收 jwtSecret
- config.Load() 会自动生成随机密钥（开发环境）或要求设置（生产环境）
- 不再依赖运行时环境变量

### 邮箱导入失败
**问题**：上传 JSON 文件后返回 409 Conflict

**原因**：JSON 中的邮箱地址（main 字段）已存在于数据库中

**解决方案**：
- 检查数据库中是否已有相同的邮箱
- 修改 JSON 文件中的邮箱地址
- 或先删除数据库中的重复记录

### 前端无法连接后端
**问题**：前端请求返回 CORS 错误

**原因**：CORS 配置不正确

**解决方案**：
- 确保后端在 8080 端口运行
- 确保前端在 3000 端口运行
- 检查 `deployment/.env` 中的 CORS_ORIGIN 配置

## 技术细节

### JWT 认证流程
1. 用户登录 → 后端验证密码
2. 生成 JWT Token（有效期 7 天）
3. 前端存储 Token 到 localStorage
4. 后续请求在 Header 中携带 `Authorization: Bearer <token>`
5. 后端中间件验证 Token 有效性

### 邮箱导入流程
1. 用户选择 JSON 文件
2. 前端使用 FormData 上传文件
3. 后端解析 JSON 内容
4. 检查是否有重复邮箱（main 字段）
5. 使用数据库事务批量插入
6. 同时插入主邮箱和 family 邮箱
7. 失败时回滚所有更改

### 数据库关系
```
User (1) ----< (N) Task
Email (1) ----< (N) EmailFamily
```

- 一个用户可以创建多个任务
- 一个邮箱可以有多个 family 邮箱
- 使用软删除（DeletedAt）保留历史记录

## 性能优化建议

### 已实现
- [x] 数据库索引：email.main (uniqueIndex), email_family.email_id (index)
- [x] 预加载关联：GetEmails 使用 Preload("Familys")
- [x] 事务处理：批量导入使用单个事务
- [x] 前端状态管理：避免不必要的重新渲染

### 待优化
- [ ] 分页查询：邮箱列表数据量大时需要分页
- [ ] 缓存层：使用 Redis 缓存常用数据
- [ ] 连接池：配置数据库连接池参数
- [ ] 静态资源 CDN：生产环境使用 CDN 加速

## 安全建议

### 已实现
- [x] 密码哈希：bcrypt cost=12
- [x] JWT 签名：HMAC-SHA256
- [x] 登录限流：5 次失败封禁 15 分钟
- [x] CORS 限制：生产环境白名单
- [x] SQL 注入防护：GORM 参数化查询
- [x] XSS 防护：React 自动转义

### 待加强
- [ ] HTTPS：生产环境强制使用 HTTPS
- [ ] Rate Limiting：全局 API 限流
- [ ] 审计日志：记录敏感操作
- [ ] 2FA：双因素认证登录
- [ ] 密码重置：邮件验证码重置密码

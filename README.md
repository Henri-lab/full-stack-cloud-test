# FreeGemini

一个完整的全栈应用，包含用户认证、任务管理和邮箱账号管理功能。

## ✨ 特性

- 🔐 **用户认证** - 注册、登录、JWT Token 认证
- 📋 **任务管理** - 完整的 CRUD 操作
- 📧 **邮箱管理** - 邮箱账号管理 + 批量导入
- 👨‍👩‍👧‍👦 **Family 邮箱** - 每个邮箱可关联多个 family 账号
- 🔑 **TOTP 支持** - 动态验证码生成
- 🛡️ **安全加固** - 密码强度验证、登录限流、CORS 保护
- 📦 **批量导入** - JSON 文件批量导入邮箱数据
- 🐳 **Docker 支持** - 一键部署

## 🚀 快速开始

### 前置要求

- Go 1.24+
- Node.js 18+
- PostgreSQL 15+
- (可选) Docker & Docker Compose

### 开发环境

#### 方式 1: 使用启动脚本（推荐）

```bash
# 启动所有服务
./start-dev.sh

# 停止所有服务
./stop-dev.sh
```

#### 方式 2: 手动启动

```bash
# 终端 1 - 启动后端
cd backend
go run cmd/api/main.go

# 终端 2 - 启动前端
cd frontend
npm install
npm run dev
```

访问 http://localhost:3000

### Docker 部署

```bash
cd deployment
docker-compose up -d
```

## 📚 文档

- [项目文档](CLAUDE.md) - 完整的项目结构和技术细节
- [API 文档](API_DOCS.md) - RESTful API 接口说明
- [导入指南](IMPORT_GUIDE.md) - 邮箱批量导入使用指南

## 🏗️ 技术栈

### 后端
- **语言**: Go 1.24
- **框架**: Gin
- **ORM**: GORM
- **数据库**: PostgreSQL 15
- **缓存**: Redis 7
- **认证**: JWT

### 前端
- **框架**: React 18
- **语言**: TypeScript
- **构建工具**: Vite
- **样式**: Tailwind CSS
- **HTTP 客户端**: Axios
- **2FA**: OTPAuth

### 部署
- **容器化**: Docker + Docker Compose
- **反向代理**: Nginx
- **监控**: Prometheus + Grafana

## 📁 项目结构

```
fullStack/
├── backend/                 # Go 后端
│   ├── cmd/api/            # 应用入口
│   └── internal/           # 内部包
│       ├── config/         # 配置管理
│       ├── database/       # 数据库连接
│       ├── handlers/       # HTTP 处理器
│       ├── middleware/     # 中间件
│       └── models/         # 数据模型
├── frontend/               # React 前端
│   └── src/
│       ├── pages/          # 页面组件
│       ├── services/       # API 服务
│       └── resource/       # 静态资源
├── deployment/             # 部署配置
│   ├── .env               # 环境变量
│   ├── docker-compose.yml # Docker 配置
│   └── nginx/             # Nginx 配置
├── logs/                   # 日志文件
├── CLAUDE.md              # 项目文档
├── API_DOCS.md            # API 文档
├── IMPORT_GUIDE.md        # 导入指南
├── test-import.json       # 测试数据
├── start-dev.sh           # 启动脚本
└── stop-dev.sh            # 停止脚本
```

## 🔑 核心功能

### 1. 用户认证

- ✅ 用户注册（密码强度验证）
- ✅ 用户登录（JWT Token）
- ✅ 登录限流（5次失败封禁15分钟）
- ✅ 密码要求：12位+大小写+数字+特殊字符

### 2. 任务管理

- ✅ 创建任务
- ✅ 查看任务列表
- ✅ 更新任务状态
- ✅ 删除任务（软删除）

### 3. 邮箱管理

- ✅ 邮箱 CRUD 操作
- ✅ 批量导入（JSON 文件）
- ✅ Family 邮箱关联
- ✅ TOTP 动态验证码
- ✅ 状态管理（Active/Banned/Sold/Need Repair）

### 4. 批量导入

支持 JSON 文件批量导入邮箱数据：

```json
{
  "emails": [
    {
      "main": "test@gmail.com",
      "password": "TestPass123!",
      "deputy": "backup@gmail.com",
      "key_2FA": "JBSWY3DPEHPK3PXP",
      "meta": {
        "banned": false,
        "price": 10,
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

## 🔒 安全特性

- ✅ **密码哈希**: bcrypt (cost=12)
- ✅ **JWT 认证**: HMAC-SHA256 签名
- ✅ **登录限流**: 5次失败封禁15分钟
- ✅ **CORS 保护**: 生产环境白名单
- ✅ **SQL 注入防护**: GORM 参数化查询
- ✅ **XSS 防护**: React 自动转义
- ✅ **密码强度**: 12位+大小写+数字+特殊字符

## 📊 API 端点

### 认证
- `POST /api/v1/auth/register` - 注册
- `POST /api/v1/auth/login` - 登录
- `POST /api/v1/auth/logout` - 登出

### 任务（需要认证）
- `GET /api/v1/tasks` - 获取所有任务
- `POST /api/v1/tasks` - 创建任务
- `PUT /api/v1/tasks/:id` - 更新任务
- `DELETE /api/v1/tasks/:id` - 删除任务

### 邮箱（需要认证）
- `GET /api/v1/emails` - 获取所有邮箱
- `POST /api/v1/emails` - 创建邮箱
- `POST /api/v1/emails/import` - 批量导入
- `PUT /api/v1/emails/:id` - 更新邮箱
- `DELETE /api/v1/emails/:id` - 删除邮箱

详细 API 文档请查看 [API_DOCS.md](API_DOCS.md)

## 🧪 测试

### 测试导入功能

项目包含测试数据文件 `test-import.json`：

```bash
# 1. 启动服务
./start-dev.sh

# 2. 注册并登录
# 访问 http://localhost:3000

# 3. 进入 Emails 页面

# 4. 点击 "Import JSON" 上传 test-import.json

# 5. 查看导入结果
```

### 测试 API

```bash
# 健康检查
curl http://localhost:8080/api/health

# 注册用户
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"MyPassword123!@#"}'

# 登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"MyPassword123!@#"}'
```

### 邮箱数据导入/导出

从 `frontend/src/resource/emails.json` 导入到数据库：

```bash
cd backend
DATABASE_URL=postgres://postgres:postgres@localhost:5432/fullstack?sslmode=disable \
  go run cmd/seed-emails/main.go
```

从数据库导出为 SQL：

```bash
cd backend
DATABASE_URL=postgres://postgres:postgres@localhost:5432/fullstack?sslmode=disable \
  go run cmd/export-emails/main.go > emails.sql
```

## 🔧 配置

### 环境变量

在 `deployment/.env` 中配置：

```bash
# 数据库
DATABASE_URL=postgres://postgres:postgres@localhost:5432/fullstack?sslmode=disable

# JWT（生产环境必须设置）
JWT_SECRET=your-secret-key-at-least-32-characters

# 服务器
PORT=8080
ENVIRONMENT=development

# CORS（生产环境必须设置）
CORS_ORIGIN=https://yourdomain.com
```

### 开发环境

开发环境会自动生成随机 JWT_SECRET，无需手动配置。

### 生产环境

生产环境必须设置：
- `JWT_SECRET` - 至少32字符
- `CORS_ORIGIN` - 允许的前端域名
- `DATABASE_URL` - 使用 SSL 连接

## 📝 开发日志

### v1.0.0 (2025-01-25)

**新增功能**:
- ✅ 用户认证系统
- ✅ 任务管理功能
- ✅ 邮箱管理功能
- ✅ 邮箱批量导入
- ✅ EmailFamily 关联管理
- ✅ TOTP 动态验证码
- ✅ Docker 部署支持

**安全加固**:
- ✅ 密码强度验证（12位+复杂度）
- ✅ 登录限流保护
- ✅ JWT 认证修复
- ✅ CORS 配置

**Bug 修复**:
- ✅ 修复 JWT 认证 panic 问题
- ✅ 修复 AuthHandler 环境变量依赖

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

## 🙏 致谢

- [Gin](https://github.com/gin-gonic/gin) - Go Web 框架
- [GORM](https://gorm.io/) - Go ORM
- [React](https://react.dev/) - 前端框架
- [Vite](https://vitejs.dev/) - 构建工具
- [Tailwind CSS](https://tailwindcss.com/) - CSS 框架

## 📞 支持

如有问题，请查看：
- [项目文档](CLAUDE.md)
- [API 文档](API_DOCS.md)
- [导入指南](IMPORT_GUIDE.md)

---

Made with ❤️ by FreeGemini Team

## 技术栈

### 前端
- React 18
- Vite
- React Router
- Axios

### 后端
- Go 1.21+
- Gin Web Framework
- GORM
- JWT Authentication
- PostgreSQL

### 基础设施
- Docker & Docker Compose
- Nginx
- PostgreSQL 15
- Redis 7

### 监控
- Prometheus
- Grafana
- Loki

## 项目结构

```
fullStack/
├── frontend/           # React前端应用
├── backend/            # Go后端API
├── deployment/         # 部署配置
│   ├── docker/        # Docker配置文件
│   └── scripts/       # 部署脚本
├── monitoring/         # 监控配置
│   ├── prometheus/
│   ├── grafana/
│   └── loki/
├── docs/              # 项目文档
└── ARCHITECTURE.md    # 架构设计文档
```

## 快速开始

### 前置要求

- Node.js 18+
- Go 1.21+
- Docker & Docker Compose
- Git

### 本地开发

#### 1. 克隆项目

```bash
git clone <repository-url>
cd fullStack
```

#### 2. 启动数据库服务

```bash
cd deployment
docker-compose up -d postgres redis
```

#### 3. 启动后端

```bash
cd backend
go mod download
go run cmd/api/main.go
```

后端将在 http://localhost:8080 运行

#### 4. 启动前端

```bash
cd frontend
npm install
npm run dev
```

前端将在 http://localhost:3000 运行

### Docker部署

#### 开发环境

```bash
cd deployment
docker-compose up -d
```

#### 生产环境

```bash
cd deployment
./scripts/deploy.sh production
```

## API文档

### 认证接口

- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/logout` - 用户登出

### 任务接口

- `GET /api/v1/tasks` - 获取任务列表
- `GET /api/v1/tasks/:id` - 获取任务详情
- `POST /api/v1/tasks` - 创建任务
- `PUT /api/v1/tasks/:id` - 更新任务
- `DELETE /api/v1/tasks/:id` - 删除任务

### 健康检查

- `GET /api/health` - 健康检查

## 环境变量配置

复制 `.env.example` 到 `.env` 并配置以下变量：

```bash
# 数据库配置
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=fullstack
DATABASE_URL=postgres://postgres:password@localhost:5432/fullstack?sslmode=disable

# JWT配置
JWT_SECRET=your-jwt-secret-key

# 服务器配置
PORT=8080
ENVIRONMENT=production

# Grafana配置
GRAFANA_PASSWORD=your-grafana-password
```

## 部署到生产环境

### 1. 初始化服务器

在Debian/Ubuntu服务器上运行：

```bash
sudo bash deployment/scripts/init-server.sh
```

这将安装：
- Docker & Docker Compose
- 防火墙配置
- Fail2ban
- SSL证书工具

### 2. 配置环境变量

```bash
cp .env.example .env.production
# 编辑 .env.production 填入生产环境配置
```

### 3. 部署应用

```bash
cd deployment
./scripts/deploy.sh production
```

### 4. 配置SSL证书

```bash
sudo certbot --nginx -d yourdomain.com
```

### 5. 设置自动备份

添加到crontab：

```bash
# 每天凌晨2点备份
0 2 * * * /opt/fullstack/deployment/scripts/backup.sh
```

## 监控

### Grafana

访问 http://your-server:3000

默认用户名: admin
密码: 在 .env 文件中配置

### Prometheus

访问 http://your-server:9090

## 开发指南

### 前端开发

```bash
cd frontend
npm run dev      # 开发服务器
npm run build    # 生产构建
npm run preview  # 预览生产构建
```

### 后端开发

```bash
cd backend
go run cmd/api/main.go  # 运行开发服务器
go test ./...           # 运行测试
go build -o bin/api cmd/api/main.go  # 构建二进制文件
```

### 数据库迁移

数据库迁移使用GORM自动迁移功能，在应用启动时自动执行。

## 安全最佳实践

1. **更改默认密码** - 修改所有默认密码
2. **使用强JWT密钥** - 至少32个字符的随机字符串
3. **启用HTTPS** - 使用Let's Encrypt配置SSL
4. **定期更新** - 保持系统和依赖包更新
5. **备份数据** - 定期备份数据库
6. **监控日志** - 定期检查应用和系统日志

## 故障排查

### 数据库连接失败

检查数据库是否运行：
```bash
docker-compose ps postgres
```

查看数据库日志：
```bash
docker-compose logs postgres
```

### 后端启动失败

查看后端日志：
```bash
docker-compose logs backend
```

### 前端无法访问API

检查Nginx配置和后端服务状态：
```bash
docker-compose ps
docker-compose logs frontend
```

## 性能优化

1. **数据库索引** - 为常用查询字段添加索引
2. **Redis缓存** - 缓存频繁访问的数据
3. **CDN** - 使用CDN加速静态资源
4. **Gzip压缩** - Nginx已配置Gzip压缩
5. **连接池** - 数据库连接池已配置

## 贡献指南

1. Fork项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启Pull Request

## 许可证

MIT License

## 联系方式

项目链接: [https://github.com/yourusername/fullstack](https://github.com/yourusername/fullstack)

## 致谢

- [React](https://react.dev/)
- [Go](https://go.dev/)
- [Gin](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [Docker](https://www.docker.com/)
- [Prometheus](https://prometheus.io/)
- [Grafana](https://grafana.com/)

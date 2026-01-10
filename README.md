# FullStack Application

一个最小化的全栈应用系统，包含前端React、后端Go、数据库PostgreSQL、自动化运维和监控系统。

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

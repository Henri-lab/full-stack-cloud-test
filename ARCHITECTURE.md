# 最小化全栈系统架构规划

## 项目概述

构建一个最小化的全栈系统，为未来的中欧贸易众包平台打基础。目标是跑通整个技术流程，而不是堆积功能。

## 技术栈选型

### 前端
- **React 18** - UI框架
- **Vite** - 构建工具（比CRA更快）
- **React Router** - 路由管理
- **Axios** - HTTP客户端
- **TailwindCSS** - 样式框架（可选，快速开发）

### 后端
- **Go 1.21+** - 后端语言
- **Gin** - Web框架（轻量高性能）
- **GORM** - ORM框架
- **JWT** - 身份认证
- **Validator** - 数据验证

### 数据库
- **PostgreSQL 15** - 主数据库（生产级，支持复杂查询）
- **Redis 7** - 缓存和会话存储

### 运维部署
- **Docker** - 容器化
- **Docker Compose** - 本地开发环境
- **Nginx** - 反向代理和静态文件服务
- **GitHub Actions** - CI/CD（可选）
- **Shell脚本** - 自动化部署

### 安全和监控
- **Let's Encrypt** - SSL证书
- **Prometheus** - 监控指标收集
- **Grafana** - 可视化监控面板
- **Loki** - 日志聚合（轻量级）
- **Fail2ban** - 防暴力破解

## 最小化业务逻辑

### 核心功能（MVP）
1. **用户系统**
   - 用户注册
   - 用户登录
   - JWT认证

2. **任务管理**（模拟众包平台的核心）
   - 创建任务
   - 查看任务列表
   - 任务详情

3. **健康检查**
   - API健康检查端点
   - 数据库连接检查

## 系统架构图

```
┌─────────────────────────────────────────────────────────────┐
│                         用户浏览器                            │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTPS
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    Nginx (反向代理)                          │
│  - SSL终止                                                   │
│  - 静态文件服务                                               │
│  - 请求路由                                                   │
└──────────┬──────────────────────────┬───────────────────────┘
           │                          │
           │ /api/*                   │ /*
           ▼                          ▼
┌──────────────────────┐    ┌──────────────────────┐
│   Go Backend API     │    │   React Frontend     │
│   (Gin框架)          │    │   (静态文件)          │
│  - RESTful API       │    │  - SPA应用           │
│  - JWT认证           │    └──────────────────────┘
│  - 业务逻辑          │
└──────────┬───────────┘
           │
           ├─────────────┬─────────────┐
           ▼             ▼             ▼
┌──────────────┐  ┌──────────┐  ┌──────────────┐
│ PostgreSQL   │  │  Redis   │  │ Prometheus   │
│ (主数据库)    │  │  (缓存)   │  │ (监控指标)    │
└──────────────┘  └──────────┘  └──────┬───────┘
                                       │
                                       ▼
                              ┌──────────────┐
                              │   Grafana    │
                              │ (监控面板)    │
                              └──────────────┘
```

## 目录结构

```
fullStack/
├── frontend/                 # React前端
│   ├── src/
│   │   ├── components/      # 组件
│   │   ├── pages/           # 页面
│   │   ├── services/        # API服务
│   │   ├── utils/           # 工具函数
│   │   ├── App.jsx
│   │   └── main.jsx
│   ├── public/
│   ├── package.json
│   └── vite.config.js
│
├── backend/                  # Go后端
│   ├── cmd/
│   │   └── api/
│   │       └── main.go      # 入口文件
│   ├── internal/
│   │   ├── handlers/        # HTTP处理器
│   │   ├── models/          # 数据模型
│   │   ├── middleware/      # 中间件
│   │   ├── database/        # 数据库连接
│   │   └── config/          # 配置
│   ├── pkg/                 # 可复用包
│   ├── go.mod
│   └── go.sum
│
├── deployment/               # 部署相关
│   ├── docker/
│   │   ├── Dockerfile.frontend
│   │   ├── Dockerfile.backend
│   │   └── nginx.conf
│   ├── docker-compose.yml
│   ├── docker-compose.prod.yml
│   └── scripts/
│       ├── deploy.sh        # 部署脚本
│       ├── backup.sh        # 备份脚本
│       └── init-server.sh   # 服务器初始化
│
├── monitoring/               # 监控配置
│   ├── prometheus/
│   │   └── prometheus.yml
│   ├── grafana/
│   │   └── dashboards/
│   └── loki/
│       └── loki-config.yml
│
├── docs/                     # 文档
│   ├── API.md               # API文档
│   ├── DEPLOYMENT.md        # 部署文档
│   └── DEVELOPMENT.md       # 开发文档
│
└── README.md
```

## 数据库设计

### users表
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### tasks表
```sql
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    creator_id INTEGER REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'open',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## API设计

### 认证相关
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/logout` - 用户登出

### 任务相关
- `GET /api/v1/tasks` - 获取任务列表
- `GET /api/v1/tasks/:id` - 获取任务详情
- `POST /api/v1/tasks` - 创建任务（需认证）
- `PUT /api/v1/tasks/:id` - 更新任务（需认证）
- `DELETE /api/v1/tasks/:id` - 删除任务（需认证）

### 健康检查
- `GET /api/health` - 健康检查
- `GET /api/metrics` - Prometheus指标

## 部署流程

### 1. 服务器初始化
```bash
# 更新系统
apt update && apt upgrade -y

# 安装必要软件
apt install -y docker.io docker-compose git ufw fail2ban

# 配置防火墙
ufw allow 22
ufw allow 80
ufw allow 443
ufw enable

# 配置Docker
systemctl enable docker
systemctl start docker
```

### 2. 应用部署
```bash
# 克隆代码
git clone <repository-url>
cd fullStack

# 配置环境变量
cp .env.example .env
# 编辑.env文件

# 构建和启动
docker-compose -f deployment/docker-compose.prod.yml up -d

# 配置SSL证书
certbot --nginx -d yourdomain.com
```

### 3. 监控配置
```bash
# 启动监控服务
docker-compose -f monitoring/docker-compose.monitoring.yml up -d

# 访问Grafana
# http://your-server:3000
# 默认用户名/密码: admin/admin
```

## 安全措施

### 应用层
1. **JWT认证** - 所有敏感操作需要认证
2. **密码加密** - 使用bcrypt加密存储
3. **输入验证** - 所有用户输入进行验证
4. **CORS配置** - 限制跨域请求
5. **Rate Limiting** - API请求频率限制

### 系统层
1. **防火墙** - UFW配置，只开放必要端口
2. **Fail2ban** - 防止暴力破解
3. **SSL/TLS** - HTTPS加密传输
4. **定期更新** - 系统和依赖包定期更新
5. **最小权限** - 应用使用非root用户运行

### 数据库层
1. **连接加密** - 使用SSL连接
2. **强密码** - 数据库使用强密码
3. **定期备份** - 自动化数据库备份
4. **访问限制** - 数据库只允许本地访问

## 监控指标

### 应用指标
- API请求数量和延迟
- 错误率
- 活跃用户数
- 数据库连接池状态

### 系统指标
- CPU使用率
- 内存使用率
- 磁盘使用率
- 网络流量

### 业务指标
- 用户注册数
- 任务创建数
- 日活跃用户（DAU）

## 开发流程

### 本地开发
```bash
# 前端开发
cd frontend
npm install
npm run dev

# 后端开发
cd backend
go mod download
go run cmd/api/main.go

# 启动本地数据库
docker-compose up -d postgres redis
```

### 测试
```bash
# 前端测试
npm run test

# 后端测试
go test ./...

# 集成测试
./scripts/integration-test.sh
```

## 扩展路径

当系统稳定运行后，可以考虑以下扩展：

1. **功能扩展**
   - 任务接单功能
   - 支付系统集成
   - 消息通知系统
   - 文件上传功能

2. **技术优化**
   - 引入消息队列（RabbitMQ/Kafka）
   - 实现微服务架构
   - 添加全文搜索（Elasticsearch）
   - CDN加速

3. **运维优化**
   - Kubernetes部署
   - 自动扩缩容
   - 多区域部署
   - 灾难恢复方案

## 预估成本（月度）

### 最小配置
- **云服务器**: $5-10/月（1核2G）
- **域名**: $10-15/年
- **SSL证书**: 免费（Let's Encrypt）
- **总计**: ~$10-15/月

### 推荐配置
- **云服务器**: $20-40/月（2核4G）
- **数据库备份存储**: $5/月
- **CDN流量**: $5-10/月
- **总计**: ~$30-55/月

## 时间规划

### 第一阶段：基础搭建（学习重点）
- 前端React基础结构
- 后端Go API框架
- 数据库设计和连接
- Docker容器化

### 第二阶段：核心功能（业务逻辑）
- 用户认证系统
- 任务CRUD功能
- 前后端联调

### 第三阶段：部署上线（运维实践）
- 服务器配置
- Docker部署
- Nginx配置
- SSL证书

### 第四阶段：监控优化（稳定性）
- Prometheus监控
- Grafana面板
- 日志系统
- 性能优化

## 学习资源

### Go后端
- Go官方文档: https://go.dev/doc/
- Gin框架: https://gin-gonic.com/
- GORM: https://gorm.io/

### React前端
- React官方文档: https://react.dev/
- Vite: https://vitejs.dev/

### 运维部署
- Docker文档: https://docs.docker.com/
- Nginx配置: https://nginx.org/en/docs/

### 监控
- Prometheus: https://prometheus.io/docs/
- Grafana: https://grafana.com/docs/

## 下一步行动

1. 创建项目目录结构
2. 初始化前端React项目
3. 初始化后端Go项目
4. 编写Docker配置文件
5. 实现最小化功能
6. 本地测试
7. 部署到服务器
8. 配置监控

---

**注意**: 这是一个学习导向的架构，重点是跑通流程和理解各个组件的作用。在实际生产环境中，需要根据具体需求进行调整和优化。

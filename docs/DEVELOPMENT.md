# 开发文档

## 本地开发环境设置

### 前置要求

- Node.js 18+ 和 npm
- Go 1.21+
- Docker 和 Docker Compose
- PostgreSQL 15 (可选，可使用Docker)
- Redis 7 (可选，可使用Docker)

### 环境配置

1. **克隆项目**

```bash
git clone <repository-url>
cd fullStack
```

2. **配置环境变量**

```bash
cp .env.example .env.development
```

编辑 `.env.development` 文件，配置本地开发环境变量。

3. **启动数据库服务**

使用Docker启动PostgreSQL和Redis：

```bash
cd deployment
docker-compose up -d postgres redis
```

或者使用本地安装的数据库服务。

### 前端开发

1. **安装依赖**

```bash
cd frontend
npm install
```

2. **启动开发服务器**

```bash
npm run dev
```

前端将在 http://localhost:3000 运行，支持热重载。

3. **构建生产版本**

```bash
npm run build
```

构建产物将输出到 `dist/` 目录。

### 后端开发

1. **安装依赖**

```bash
cd backend
go mod download
```

2. **运行开发服务器**

```bash
go run cmd/api/main.go
```

后端API将在 http://localhost:8080 运行。

3. **构建二进制文件**

```bash
go build -o bin/api cmd/api/main.go
```

## 代码结构

### 前端结构

```
frontend/
├── src/
│   ├── components/     # 可复用组件
│   ├── pages/          # 页面组件
│   │   ├── Home.jsx
│   │   ├── Login.jsx
│   │   ├── Register.jsx
│   │   └── Tasks.jsx
│   ├── services/       # API服务
│   │   └── api.js
│   ├── utils/          # 工具函数
│   ├── App.jsx         # 主应用组件
│   ├── App.css         # 应用样式
│   ├── main.jsx        # 入口文件
│   └── index.css       # 全局样式
├── public/             # 静态资源
├── index.html          # HTML模板
├── vite.config.js      # Vite配置
└── package.json        # 依赖配置
```

### 后端结构

```
backend/
├── cmd/
│   └── api/
│       └── main.go           # 应用入口
├── internal/
│   ├── config/
│   │   └── config.go         # 配置管理
│   ├── database/
│   │   └── database.go       # 数据库连接
│   ├── handlers/
│   │   ├── auth.go           # 认证处理器
│   │   └── task.go           # 任务处理器
│   ├── middleware/
│   │   └── middleware.go     # 中间件
│   └── models/
│       └── models.go         # 数据模型
├── pkg/                      # 可复用包
├── go.mod                    # Go模块定义
└── go.sum                    # 依赖锁定
```

## API开发

### 添加新的API端点

1. **定义数据模型** (internal/models/models.go)

```go
type NewModel struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    Name      string    `gorm:"not null" json:"name"`
    CreatedAt time.Time `json:"created_at"`
}
```

2. **创建处理器** (internal/handlers/newmodel.go)

```go
type NewModelHandler struct {
    db *gorm.DB
}

func NewNewModelHandler(db *gorm.DB) *NewModelHandler {
    return &NewModelHandler{db: db}
}

func (h *NewModelHandler) GetAll(c *gin.Context) {
    // 实现逻辑
}
```

3. **注册路由** (cmd/api/main.go)

```go
newModel := v1.Group("/newmodel")
newModel.Use(middleware.AuthMiddleware(cfg.JWTSecret))
{
    handler := handlers.NewNewModelHandler(db)
    newModel.GET("", handler.GetAll)
}
```

### 添加中间件

在 `internal/middleware/middleware.go` 中添加新的中间件：

```go
func NewMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 中间件逻辑
        c.Next()
    }
}
```

## 前端开发

### 添加新页面

1. **创建页面组件** (src/pages/NewPage.jsx)

```jsx
function NewPage() {
    return (
        <div>
            <h1>New Page</h1>
        </div>
    )
}

export default NewPage
```

2. **添加路由** (src/App.jsx)

```jsx
import NewPage from './pages/NewPage'

// 在Routes中添加
<Route path="/newpage" element={<NewPage />} />
```

### API调用

使用 `src/services/api.js` 中的axios实例：

```jsx
import api from '../services/api'

const fetchData = async () => {
    try {
        const response = await api.get('/endpoint')
        console.log(response.data)
    } catch (error) {
        console.error('Error:', error)
    }
}
```

## 测试

### 后端测试

创建测试文件 `*_test.go`：

```go
package handlers

import (
    "testing"
)

func TestSomething(t *testing.T) {
    // 测试逻辑
}
```

运行测试：

```bash
go test ./...
```

### 前端测试

(待添加测试框架)

## 数据库

### 迁移

GORM自动迁移在应用启动时执行。添加新模型后，在 `internal/database/database.go` 中注册：

```go
func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.User{},
        &models.Task{},
        &models.NewModel{}, // 新模型
    )
}
```

### 查询示例

```go
// 查询所有
var items []models.Item
db.Find(&items)

// 条件查询
db.Where("status = ?", "active").Find(&items)

// 创建
db.Create(&item)

// 更新
db.Model(&item).Updates(map[string]interface{}{"status": "done"})

// 删除
db.Delete(&item)
```

## 调试

### 后端调试

使用Delve调试器：

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug cmd/api/main.go
```

### 前端调试

使用浏览器开发者工具：
- Chrome DevTools
- React Developer Tools扩展

### 日志

后端日志会输出到控制台。在生产环境中，日志会被Docker收集。

查看日志：

```bash
# 开发环境
docker-compose logs -f backend

# 生产环境
docker-compose -f docker-compose.prod.yml logs -f backend
```

## 性能优化

### 前端优化

1. **代码分割** - 使用React.lazy()和Suspense
2. **图片优化** - 使用WebP格式，懒加载
3. **缓存策略** - 合理使用浏览器缓存
4. **减少重渲染** - 使用React.memo, useMemo, useCallback

### 后端优化

1. **数据库索引** - 为常用查询字段添加索引
2. **连接池** - GORM已配置连接池
3. **缓存** - 使用Redis缓存热数据
4. **并发处理** - 使用Go的goroutine

## 常见问题

### 端口冲突

如果端口被占用，修改以下配置：
- 前端: `vite.config.js` 中的 `server.port`
- 后端: `.env` 文件中的 `PORT`

### CORS错误

后端已配置CORS中间件，允许所有来源。生产环境应限制允许的来源。

### 数据库连接失败

检查：
1. 数据库服务是否运行
2. 连接字符串是否正确
3. 防火墙设置

## 代码规范

### Go代码规范

- 使用 `gofmt` 格式化代码
- 遵循 [Effective Go](https://go.dev/doc/effective_go) 指南
- 使用有意义的变量名
- 添加必要的注释

### JavaScript/React规范

- 使用ES6+语法
- 组件使用函数式组件和Hooks
- 使用有意义的组件和变量名
- 保持组件小而专注

## Git工作流

1. 从main分支创建特性分支
2. 开发并提交更改
3. 推送到远程仓库
4. 创建Pull Request
5. 代码审查
6. 合并到main分支

## 部署流程

参见 [DEPLOYMENT.md](DEPLOYMENT.md)

# API文档

## 基础信息

- **Base URL**: `http://localhost:8080/api/v1`
- **认证方式**: JWT Bearer Token
- **Content-Type**: `application/json`

## 认证

大部分API端点需要认证。在请求头中包含JWT token：

```
Authorization: Bearer <your-jwt-token>
```

## 响应格式

### 成功响应

```json
{
  "data": {},
  "message": "Success"
}
```

### 错误响应

```json
{
  "error": "Error message"
}
```

## HTTP状态码

- `200 OK` - 请求成功
- `201 Created` - 资源创建成功
- `400 Bad Request` - 请求参数错误
- `401 Unauthorized` - 未认证或token无效
- `403 Forbidden` - 无权限访问
- `404 Not Found` - 资源不存在
- `409 Conflict` - 资源冲突（如用户已存在）
- `500 Internal Server Error` - 服务器错误

## API端点

### 健康检查

#### GET /api/health

检查API服务健康状态

**请求**
```bash
curl http://localhost:8080/api/health
```

**响应**
```json
{
  "status": "ok",
  "database": "connected"
}
```

---

## 认证接口

### 用户注册

#### POST /api/v1/auth/register

注册新用户

**请求体**
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123"
}
```

**验证规则**
- `username`: 必填，3-50个字符
- `email`: 必填，有效的邮箱格式
- `password`: 必填，至少6个字符

**成功响应** (201 Created)
```json
{
  "message": "User created successfully",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com"
  }
}
```

**错误响应** (409 Conflict)
```json
{
  "error": "User already exists"
}
```

**示例**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

---

### 用户登录

#### POST /api/v1/auth/login

用户登录获取JWT token

**请求体**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**成功响应** (200 OK)
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com"
  }
}
```

**错误响应** (401 Unauthorized)
```json
{
  "error": "Invalid credentials"
}
```

**示例**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

---

### 用户登出

#### POST /api/v1/auth/logout

用户登出（客户端需删除token）

**成功响应** (200 OK)
```json
{
  "message": "Logged out successfully"
}
```

**示例**
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout
```

---

## 任务接口

所有任务接口都需要认证。

### 获取任务列表

#### GET /api/v1/tasks

获取当前用户的所有任务

**请求头**
```
Authorization: Bearer <token>
```

**成功响应** (200 OK)
```json
[
  {
    "id": 1,
    "title": "Complete project documentation",
    "description": "Write comprehensive API documentation",
    "status": "open",
    "creator_id": 1,
    "created_at": "2024-01-05T10:00:00Z",
    "updated_at": "2024-01-05T10:00:00Z"
  },
  {
    "id": 2,
    "title": "Review pull requests",
    "description": "Review and merge pending PRs",
    "status": "in_progress",
    "creator_id": 1,
    "created_at": "2024-01-05T11:00:00Z",
    "updated_at": "2024-01-05T11:30:00Z"
  }
]
```

**示例**
```bash
curl http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer <your-token>"
```

---

### 获取单个任务

#### GET /api/v1/tasks/:id

获取指定ID的任务详情

**路径参数**
- `id`: 任务ID

**成功响应** (200 OK)
```json
{
  "id": 1,
  "title": "Complete project documentation",
  "description": "Write comprehensive API documentation",
  "status": "open",
  "creator_id": 1,
  "created_at": "2024-01-05T10:00:00Z",
  "updated_at": "2024-01-05T10:00:00Z"
}
```

**错误响应** (404 Not Found)
```json
{
  "error": "Task not found"
}
```

**示例**
```bash
curl http://localhost:8080/api/v1/tasks/1 \
  -H "Authorization: Bearer <your-token>"
```

---

### 创建任务

#### POST /api/v1/tasks

创建新任务

**请求体**
```json
{
  "title": "New task title",
  "description": "Task description"
}
```

**验证规则**
- `title`: 必填，1-200个字符
- `description`: 可选

**成功响应** (201 Created)
```json
{
  "id": 3,
  "title": "New task title",
  "description": "Task description",
  "status": "open",
  "creator_id": 1,
  "created_at": "2024-01-05T12:00:00Z",
  "updated_at": "2024-01-05T12:00:00Z"
}
```

**示例**
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "New task title",
    "description": "Task description"
  }'
```

---

### 更新任务

#### PUT /api/v1/tasks/:id

更新指定任务

**路径参数**
- `id`: 任务ID

**请求体**
```json
{
  "title": "Updated title",
  "description": "Updated description",
  "status": "completed"
}
```

**可选字段**
- `title`: 任务标题
- `description`: 任务描述
- `status`: 任务状态 (open, in_progress, completed)

**成功响应** (200 OK)
```json
{
  "id": 1,
  "title": "Updated title",
  "description": "Updated description",
  "status": "completed",
  "creator_id": 1,
  "created_at": "2024-01-05T10:00:00Z",
  "updated_at": "2024-01-05T13:00:00Z"
}
```

**错误响应** (404 Not Found)
```json
{
  "error": "Task not found"
}
```

**示例**
```bash
curl -X PUT http://localhost:8080/api/v1/tasks/1 \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed"
  }'
```

---

### 删除任务

#### DELETE /api/v1/tasks/:id

删除指定任务

**路径参数**
- `id`: 任务ID

**成功响应** (200 OK)
```json
{
  "message": "Task deleted successfully"
}
```

**错误响应** (404 Not Found)
```json
{
  "error": "Task not found"
}
```

**示例**
```bash
curl -X DELETE http://localhost:8080/api/v1/tasks/1 \
  -H "Authorization: Bearer <your-token>"
```

---

## 错误处理

### 常见错误

#### 401 Unauthorized

Token缺失或无效

```json
{
  "error": "Authorization header required"
}
```

或

```json
{
  "error": "Invalid token"
}
```

**解决方案**: 重新登录获取新token

#### 400 Bad Request

请求参数验证失败

```json
{
  "error": "Key: 'RegisterRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"
}
```

**解决方案**: 检查请求参数格式

#### 500 Internal Server Error

服务器内部错误

```json
{
  "error": "Internal server error"
}
```

**解决方案**: 查看服务器日志，联系管理员

---

## 使用示例

### 完整工作流程

1. **注册用户**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

2. **登录获取token**
```bash
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }' | jq -r '.token')
```

3. **创建任务**
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My first task",
    "description": "This is a test task"
  }'
```

4. **获取任务列表**
```bash
curl http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer $TOKEN"
```

5. **更新任务**
```bash
curl -X PUT http://localhost:8080/api/v1/tasks/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed"
  }'
```

6. **删除任务**
```bash
curl -X DELETE http://localhost:8080/api/v1/tasks/1 \
  -H "Authorization: Bearer $TOKEN"
```

---

## 速率限制

目前未实施速率限制。生产环境建议添加速率限制以防止滥用。

## 版本控制

当前API版本: v1

API版本通过URL路径指定: `/api/v1/...`

## 支持

如有API相关问题，请：
1. 查看本文档
2. 检查GitHub Issues
3. 联系技术支持

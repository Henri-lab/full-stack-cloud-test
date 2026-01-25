# FreeGemini API 文档

## 基础信息

- **Base URL**: `http://localhost:8080/api/v1`
- **认证方式**: JWT Bearer Token
- **Content-Type**: `application/json`

## 认证接口

### 注册用户

**POST** `/auth/register`

**请求体**:
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "MyPassword123!@#"
}
```

**密码要求**:
- 至少 12 个字符
- 包含大写字母
- 包含小写字母
- 包含数字
- 包含特殊字符 `!@#$%^&*()_+-=[]{}|;':\",./<>?`

**成功响应** (201):
```json
{
  "message": "User created successfully",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

**错误响应**:
- `400` - 请求参数错误
- `409` - 用户已存在
- `429` - 请求过于频繁（5次失败后封禁15分钟）

---

### 用户登录

**POST** `/auth/login`

**请求体**:
```json
{
  "email": "test@example.com",
  "password": "MyPassword123!@#"
}
```

**成功响应** (200):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

**错误响应**:
- `400` - 请求参数错误
- `401` - 邮箱或密码错误
- `429` - 登录失败次数过多

**限流规则**:
- 5次失败后封禁IP 15分钟
- 成功登录后清除失败记录

---

### 用户登出

**POST** `/auth/logout`

**Headers**:
```
Authorization: Bearer <token>
```

**成功响应** (200):
```json
{
  "message": "Logged out successfully"
}
```

**说明**: JWT 是无状态的，登出主要在客户端删除 token

---

## 任务接口

所有任务接口都需要 JWT 认证。

### 获取所有任务

**GET** `/tasks`

**Headers**:
```
Authorization: Bearer <token>
```

**成功响应** (200):
```json
[
  {
    "id": 1,
    "title": "完成项目文档",
    "description": "编写 API 文档和使用指南",
    "status": "in_progress",
    "creator_id": 1,
    "created_at": "2025-01-25T10:00:00Z",
    "updated_at": "2025-01-25T11:00:00Z"
  }
]
```

---

### 获取单个任务

**GET** `/tasks/:id`

**Headers**:
```
Authorization: Bearer <token>
```

**成功响应** (200):
```json
{
  "id": 1,
  "title": "完成项目文档",
  "description": "编写 API 文档和使用指南",
  "status": "in_progress",
  "creator_id": 1,
  "created_at": "2025-01-25T10:00:00Z",
  "updated_at": "2025-01-25T11:00:00Z"
}
```

**错误响应**:
- `404` - 任务不存在

---

### 创建任务

**POST** `/tasks`

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "title": "新任务",
  "description": "任务描述",
  "status": "open"
}
```

**Status 可选值**:
- `open` - 待处理
- `in_progress` - 进行中
- `completed` - 已完成

**成功响应** (201):
```json
{
  "id": 2,
  "title": "新任务",
  "description": "任务描述",
  "status": "open",
  "creator_id": 1,
  "created_at": "2025-01-25T12:00:00Z",
  "updated_at": "2025-01-25T12:00:00Z"
}
```

---

### 更新任务

**PUT** `/tasks/:id`

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "title": "更新后的标题",
  "description": "更新后的描述",
  "status": "completed"
}
```

**成功响应** (200):
```json
{
  "id": 1,
  "title": "更新后的标题",
  "description": "更新后的描述",
  "status": "completed",
  "creator_id": 1,
  "created_at": "2025-01-25T10:00:00Z",
  "updated_at": "2025-01-25T13:00:00Z"
}
```

---

### 删除任务

**DELETE** `/tasks/:id`

**Headers**:
```
Authorization: Bearer <token>
```

**成功响应** (200):
```json
{
  "message": "Task deleted successfully"
}
```

**说明**: 使用软删除，数据不会真正从数据库删除

---

## 邮箱接口

所有邮箱接口都需要 JWT 认证。

### 获取所有邮箱

**GET** `/emails`

**Headers**:
```
Authorization: Bearer <token>
```

**Query Params**:
- `import_id` (可选) - 过滤指定导入批次的数据

**成功响应** (200):
```json
[
  {
    "id": 1,
    "main": "test@gmail.com",
    "password": "TestPass123!",
    "deputy": "backup@gmail.com",
    "key_2FA": "JBSWY3DPEHPK3PXP",
    "meta": {
      "banned": false,
      "created_at": "2025-01-25T10:00:00Z",
      "updated_at": "2025-01-25T10:00:00Z",
      "price": 10,
      "sold": false,
      "need_repair": false,
      "from": "source1"
    },
    "familys": [
      {
        "id": 1,
        "email": "family1@gmail.com",
        "password": "FamilyPass123!",
        "code": "123456",
        "contact": "qq:123456;phone:13800138000",
        "issue": "正常使用"
      }
    ]
  }
]
```

---

### 获取导入记录列表

**GET** `/emails/imports`

**Headers**:
```
Authorization: Bearer <token>
```

**成功响应** (200):
```json
[
  {
    "id": 12,
    "name": "test-emails.json",
    "created_at": "2025-01-25T10:00:00Z",
    "count": 3
  }
]
```

---

### 获取单个邮箱

**GET** `/emails/:id`

**Headers**:
```
Authorization: Bearer <token>
```

**成功响应** (200):
```json
{
  "id": 1,
  "main": "test@gmail.com",
  "password": "TestPass123!",
  "deputy": "backup@gmail.com",
  "key_2FA": "JBSWY3DPEHPK3PXP",
  "meta": {
    "banned": false,
    "created_at": "2025-01-25T10:00:00Z",
    "updated_at": "2025-01-25T10:00:00Z",
    "price": 10,
    "sold": false,
    "need_repair": false,
    "from": "source1"
  },
  "familys": []
}
```

---

### 创建邮箱

**POST** `/emails`

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "main": "new@gmail.com",
  "password": "NewPass123!",
  "deputy": "backup@gmail.com",
  "key_2FA": "JBSWY3DPEHPK3PXP",
  "meta": {
    "banned": false,
    "price": 10,
    "sold": false,
    "need_repair": false,
    "from": "source1"
  }
}
```

**成功响应** (201):
```json
{
  "id": 2,
  "main": "new@gmail.com",
  "password": "NewPass123!",
  "deputy": "backup@gmail.com",
  "key_2FA": "JBSWY3DPEHPK3PXP",
  "meta": {
    "banned": false,
    "created_at": "2025-01-25T14:00:00Z",
    "updated_at": "2025-01-25T14:00:00Z",
    "price": 10,
    "sold": false,
    "need_repair": false,
    "from": "source1"
  },
  "familys": []
}
```

---

### 批量导入邮箱

**POST** `/emails/import`

**Headers**:
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**请求体**:
- `file`: JSON 文件

**JSON 文件格式**:
```json
{
  "emails": [
    {
      "main": "import1@gmail.com",
      "password": "ImportPass123!",
      "deputy": "backup1@gmail.com",
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
          "contact": "qq:123456",
          "issue": "正常使用"
        }
      ]
    }
  ]
}
```

**成功响应** (200):
```json
{
  "message": "Import successful",
  "imported": 1
}
```

**错误响应**:
- `400` - 文件未上传或 JSON 格式错误
- `409` - 邮箱已存在
- `500` - 导入失败

**导入规则**:
- 遇到重复邮箱（main 字段）会报错终止
- 使用事务处理，失败时回滚所有更改
- familys 数组可以为空

---

### 更新邮箱

**PUT** `/emails/:id`

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "main": "updated@gmail.com",
  "password": "UpdatedPass123!",
  "deputy": "newbackup@gmail.com",
  "key_2FA": "NEWKEY123456",
  "meta": {
    "banned": true,
    "price": 20,
    "sold": true,
    "need_repair": false,
    "from": "source2"
  }
}
```

**成功响应** (200):
```json
{
  "id": 1,
  "main": "updated@gmail.com",
  "password": "UpdatedPass123!",
  "deputy": "newbackup@gmail.com",
  "key_2FA": "NEWKEY123456",
  "meta": {
    "banned": true,
    "created_at": "2025-01-25T10:00:00Z",
    "updated_at": "2025-01-25T15:00:00Z",
    "price": 20,
    "sold": true,
    "need_repair": false,
    "from": "source2"
  },
  "familys": []
}
```

---

### 删除邮箱

**DELETE** `/emails/:id`

**Headers**:
```
Authorization: Bearer <token>
```

**成功响应** (200):
```json
{
  "message": "Email deleted successfully"
}
```

**说明**: 使用软删除，关联的 family 邮箱也会被软删除

---

## 健康检查

### 服务状态

**GET** `/health`

**无需认证**

**成功响应** (200):
```json
{
  "status": "ok",
  "database": "connected"
}
```

---

## 错误响应格式

所有错误响应都遵循以下格式：

```json
{
  "error": "错误描述信息"
}
```

### 常见错误码

| 状态码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 401 | 未认证或认证失败 |
| 403 | 无权限访问 |
| 404 | 资源不存在 |
| 409 | 资源冲突（如邮箱已存在） |
| 429 | 请求过于频繁 |
| 500 | 服务器内部错误 |

---

## 认证流程

### 1. 获取 Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "MyPassword123!@#"
  }'
```

**响应**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {...}
}
```

### 2. 使用 Token 访问受保护接口

```bash
curl -X GET http://localhost:8080/api/v1/emails \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 3. Token 过期

- Token 有效期：7 天
- 过期后需要重新登录获取新 Token
- 前端会自动处理 401 错误并跳转到登录页

---

## 限流规则

### 登录/注册限流

- **规则**: 5次失败后封禁IP 15分钟
- **窗口期**: 10分钟内的失败次数
- **重置**: 成功登录后清除失败记录

**示例**:
```
第1次失败 → 剩余4次
第2次失败 → 剩余3次
第3次失败 → 剩余2次
第4次失败 → 剩余1次
第5次失败 → 封禁15分钟
```

---

## 使用示例

### JavaScript (Axios)

```javascript
import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  headers: {
    'Content-Type': 'application/json'
  }
})

// 添加请求拦截器
api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 登录
const login = async (email, password) => {
  const response = await api.post('/auth/login', { email, password })
  localStorage.setItem('token', response.data.token)
  return response.data
}

// 获取邮箱列表
const getEmails = async () => {
  const response = await api.get('/emails')
  return response.data
}

// 导入邮箱
const importEmails = async (file) => {
  const formData = new FormData()
  formData.append('file', file)
  const response = await api.post('/emails/import', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
  return response.data
}
```

### cURL

```bash
# 登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"MyPassword123!@#"}'

# 获取邮箱列表
curl -X GET http://localhost:8080/api/v1/emails \
  -H "Authorization: Bearer YOUR_TOKEN"

# 导入邮箱
curl -X POST http://localhost:8080/api/v1/emails/import \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@test-import.json"

# 创建任务
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"新任务","description":"任务描述","status":"open"}'
```

---

## 开发建议

### 1. 错误处理

```javascript
try {
  const response = await api.post('/auth/login', { email, password })
  // 处理成功
} catch (error) {
  if (error.response) {
    // 服务器返回错误
    console.error(error.response.data.error)
    if (error.response.status === 429) {
      // 处理限流
    }
  } else {
    // 网络错误
    console.error('Network error')
  }
}
```

### 2. Token 刷新

当前版本使用固定过期时间（7天），建议实现 refresh token 机制：
- Access Token: 短期（1小时）
- Refresh Token: 长期（30天）
- 自动刷新机制

### 3. 批量操作

导入大量邮箱时建议：
- 单次不超过 100 个
- 分批导入
- 显示进度条

---

## 更新日志

### v1.0.0 (2025-01-25)

**新增功能**:
- ✅ 用户认证（注册/登录/登出）
- ✅ 任务管理 CRUD
- ✅ 邮箱管理 CRUD
- ✅ 邮箱批量导入
- ✅ EmailFamily 关联管理
- ✅ 登录限流保护
- ✅ JWT 认证
- ✅ CORS 配置

**安全特性**:
- ✅ 密码强度验证（12位+大小写+数字+特殊字符）
- ✅ bcrypt 密码哈希（cost=12）
- ✅ JWT 签名验证
- ✅ 登录失败限流（5次/15分钟）
- ✅ SQL 注入防护
- ✅ XSS 防护

---

## 支持

如有问题，请查看：
- 项目文档：`CLAUDE.md`
- 导入指南：`IMPORT_GUIDE.md`
- GitHub Issues: [项目地址]

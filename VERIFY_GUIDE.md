# 邮箱验证功能使用指南

## 功能概述

邮箱验证功能允许你批量检测邮箱账号的状态（live/verify/dead），使用第三方 API（gmailver.com）进行验证。

## 验证状态说明

- **Live**: 邮箱可用，状态正常
- **Verify**: 需要验证，可能需要额外操作
- **Dead**: 邮箱不可用或已失效
- **Unknown**: 未验证状态

## 使用步骤

### 1. 获取验证 Key

1. 访问 https://gmailver.com 或 https://etbrower.com
2. 打开浏览器开发者工具（F12）
3. 切换到 Network 标签
4. 刷新页面
5. 查找请求到 `check1.php` 的请求
6. 在请求 payload 中找到 `key` 字段
7. 复制 key 值（例如：`d12da1defe5474edea9a574c7c9ecd98`）

**注意**: Key 每次刷新网站都会变化，需要重新获取。

### 2. 在前端使用验证功能

1. 登录后访问 Emails 页面
2. 点击 "Verify Emails" 按钮
3. 在弹出的输入框中粘贴从 gmailver.com 获取的 key
4. 勾选需要验证的邮箱（可以使用表头的复选框全选）
5. 点击 "Verify (N)" 按钮开始验证
6. 等待验证完成，查看结果

### 3. 查看验证结果

验证完成后：
- 表格中会显示每个邮箱的验证状态（Live/Verify/Dead）
- 状态会自动保存到数据库
- 可以根据状态筛选和管理邮箱

## API 使用

### 验证邮箱接口

**POST** `/api/v1/emails/verify`

**Headers**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**:
```json
{
  "mail": ["email1@gmail.com", "email2@gmail.com"],
  "key": "d12da1defe5474edea9a574c7c9ecd98"
}
```

**成功响应** (200):
```json
{
  "results": [
    {
      "email": "email1@gmail.com",
      "status": "live"
    },
    {
      "email": "email2@gmail.com",
      "status": "dead"
    }
  ],
  "total": 2
}
```

**错误响应**:
- `400` - 请求参数错误（缺少 key 或邮箱列表）
- `500` - 验证失败（第三方 API 错误）

## 使用示例

### cURL 示例

```bash
# 1. 登录获取 token
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"MyPassword123!@#"}' \
  | jq -r '.token')

# 2. 验证邮箱
curl -X POST http://localhost:8080/api/v1/emails/verify \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "mail": ["test1@gmail.com", "test2@gmail.com"],
    "key": "d12da1defe5474edea9a574c7c9ecd98"
  }'
```

### JavaScript 示例

```javascript
import api from './services/api'

// 验证邮箱
const verifyEmails = async (emails, key) => {
  try {
    const response = await api.post('/emails/verify', {
      mail: emails,
      key: key
    })

    console.log(`Verified ${response.data.total} emails`)
    response.data.results.forEach(result => {
      console.log(`${result.email}: ${result.status}`)
    })
  } catch (error) {
    console.error('Verification failed:', error)
  }
}

// 使用
verifyEmails(
  ['test1@gmail.com', 'test2@gmail.com'],
  'd12da1defe5474edea9a574c7c9ecd98'
)
```

## 注意事项

### 1. Key 有效期

- Key 每次刷新 gmailver.com 都会变化
- 旧的 key 会失效
- 建议每次验证前重新获取 key

### 2. 批量验证限制

- 建议每次验证不超过 50 个邮箱
- 大批量验证建议分批进行
- 避免频繁请求导致被限流

### 3. 第三方 API 依赖

- 验证功能依赖 gmailver.com API
- 如果第三方 API 不可用，验证会失败
- 建议实现自己的验证逻辑（见下文）

## 自己实现验证逻辑（推荐）

如果不想依赖第三方 API，可以实现自己的邮箱验证逻辑：

### 方法 1: SMTP 验证

```go
// 伪代码示例
func verifyEmailSMTP(email string) (string, error) {
    // 1. 解析邮箱域名
    domain := strings.Split(email, "@")[1]

    // 2. 查询 MX 记录
    mxRecords, err := net.LookupMX(domain)
    if err != nil {
        return "dead", err
    }

    // 3. 连接 SMTP 服务器
    conn, err := smtp.Dial(mxRecords[0].Host + ":25")
    if err != nil {
        return "dead", err
    }
    defer conn.Close()

    // 4. 验证邮箱是否存在
    // HELO, MAIL FROM, RCPT TO
    // ...

    return "live", nil
}
```

### 方法 2: API 验证

使用其他邮箱验证服务：
- **Hunter.io** - https://hunter.io/email-verifier
- **ZeroBounce** - https://www.zerobounce.net/
- **NeverBounce** - https://neverbounce.com/

### 方法 3: 正则验证 + DNS 检查

```go
func verifyEmailBasic(email string) string {
    // 1. 正则验证格式
    if !isValidEmailFormat(email) {
        return "dead"
    }

    // 2. 检查域名 MX 记录
    domain := strings.Split(email, "@")[1]
    _, err := net.LookupMX(domain)
    if err != nil {
        return "dead"
    }

    return "verify" // 需要进一步验证
}
```

## 常见问题

### Q1: Key 从哪里获取？

**A**: 访问 gmailver.com，打开浏览器开发者工具，在 Network 标签中找到 check1.php 请求，从 payload 中复制 key。

### Q2: 验证失败怎么办？

**A**: 检查以下几点：
1. Key 是否正确且未过期
2. 第三方 API 是否可访问
3. 邮箱格式是否正确
4. 网络连接是否正常

### Q3: 可以自动获取 Key 吗？

**A**: 理论上可以，但需要：
1. 模拟浏览器访问 gmailver.com
2. 解析页面获取 key
3. 可能需要处理反爬虫机制

不推荐这种方式，建议实现自己的验证逻辑。

### Q4: 验证结果准确吗？

**A**: 准确性取决于第三方 API。建议：
1. 定期重新验证
2. 结合多种验证方式
3. 人工复核重要邮箱

### Q5: 如何批量验证所有邮箱？

**A**:
1. 点击表头的复选框全选所有邮箱
2. 输入验证 key
3. 点击 "Verify" 按钮
4. 等待验证完成

## 技术实现

### 后端实现

验证接口位于 `backend/internal/handlers/email.go`:

```go
func (h *EmailHandler) VerifyEmails(c *gin.Context) {
    // 1. 接收请求参数
    var req VerifyEmailRequest
    c.ShouldBindJSON(&req)

    // 2. 调用第三方 API
    resp, err := http.Post("https://gmailver.com/php/check1.php", ...)

    // 3. 解析响应
    var apiResponse map[string]interface{}
    json.Unmarshal(body, &apiResponse)

    // 4. 更新数据库
    for email, status := range apiResponse {
        h.db.Model(&models.Email{}).
            Where("main = ?", email).
            Update("status", status)
    }

    // 5. 返回结果
    c.JSON(200, results)
}
```

### 前端实现

验证功能位于 `frontend/src/pages/Emails.tsx`:

```typescript
const handleVerifyEmails = async () => {
    // 1. 获取选中的邮箱
    const emailsToVerify = Array.from(selectedEmails)
        .map(id => emails.find(e => e.id === id)?.main)

    // 2. 调用验证 API
    const response = await api.post('/emails/verify', {
        mail: emailsToVerify,
        key: verifyKey
    })

    // 3. 更新本地状态
    setEmails(prevEmails => prevEmails.map(email => {
        const result = response.data.results.find(r => r.email === email.main)
        return result ? { ...email, status: result.status } : email
    }))
}
```

## 数据库结构

Email 表新增 `status` 字段：

```sql
ALTER TABLE emails ADD COLUMN status VARCHAR(20) DEFAULT 'unknown';

-- 可选：添加索引
CREATE INDEX idx_emails_status ON emails(status);
```

## 性能优化

### 1. 批量验证优化

```go
// 分批验证，每批 50 个
const batchSize = 50
for i := 0; i < len(emails); i += batchSize {
    end := i + batchSize
    if end > len(emails) {
        end = len(emails)
    }
    batch := emails[i:end]
    verifyBatch(batch, key)
}
```

### 2. 缓存验证结果

```go
// 使用 Redis 缓存验证结果
func getCachedStatus(email string) (string, bool) {
    val, err := redisClient.Get(ctx, "email:status:"+email).Result()
    if err == nil {
        return val, true
    }
    return "", false
}

func cacheStatus(email, status string) {
    redisClient.Set(ctx, "email:status:"+email, status, 24*time.Hour)
}
```

### 3. 异步验证

```go
// 使用 goroutine 异步验证
func verifyEmailsAsync(emails []string, key string) <-chan VerifyResult {
    results := make(chan VerifyResult)

    go func() {
        defer close(results)
        for _, email := range emails {
            status := verifyEmail(email, key)
            results <- VerifyResult{Email: email, Status: status}
        }
    }()

    return results
}
```

## 总结

邮箱验证功能提供了快速批量检测邮箱状态的能力，但依赖第三方 API。建议：

1. **短期使用**: 使用 gmailver.com API + 手动获取 key
2. **长期使用**: 实现自己的 SMTP 验证逻辑
3. **生产环境**: 使用专业的邮箱验证服务（Hunter.io 等）

---

更新时间: 2025-01-25

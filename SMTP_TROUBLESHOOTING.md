# SMTP 验证故障排查指南

## 🔍 常见问题和解决方案

### 问题 1: "EOF" 错误

**错误信息**:
```json
{
  "email": "test@gmail.com",
  "status": "unknown",
  "error": "failed to create SMTP client: EOF"
}
```

**原因**:
1. Gmail 的 SMTP 服务器拒绝了连接
2. 端口 25 被 ISP 或防火墙封禁
3. 服务器要求 TLS 加密但未启用

**解决方案**:

#### 方案 1: 使用改进的 SMTP 验证器（已实现）

最新版本已经支持：
- ✅ 多端口尝试（25, 587, 465）
- ✅ STARTTLS 支持
- ✅ 多个 MX 服务器尝试
- ✅ 更好的错误处理

重启后端服务即可使用新版本。

#### 方案 2: 使用 API 验证（推荐）

如果 SMTP 验证仍然失败，建议使用 API 验证：

```bash
# 1. 访问 https://gmailver.com
# 2. F12 → Network → 找到 check1.php
# 3. 复制 key
# 4. 使用 API 验证
```

**优点**:
- 速度快（0.1-0.5秒/邮箱）
- 不受网络限制
- 准确率高

**缺点**:
- 需要手动获取 key
- Key 会过期

### 问题 2: 所有邮箱返回 "unknown"

**原因**:
1. 网络环境限制（公司网络、云服务器）
2. SMTP 端口被封禁
3. 邮件服务器不支持 RCPT TO 验证

**解决方案**:

#### 检查网络连接

```bash
# 测试能否连接到 Gmail SMTP 服务器
telnet gmail-smtp-in.l.google.com 25

# 或使用 nc
nc -zv gmail-smtp-in.l.google.com 25
```

如果无法连接，说明端口被封禁。

#### 使用 VPN 或代理

某些 ISP 会封禁 25 端口，使用 VPN 可以绕过限制。

#### 切换到 API 验证

API 验证不受端口限制影响。

### 问题 3: Gmail 特定问题

**Gmail 的限制**:
- Gmail 对 SMTP 验证有严格的限流
- 可能会拒绝来自某些 IP 的连接
- 需要 TLS 加密

**解决方案**:

#### 1. 使用 Gmail API（最可靠）

```go
// 需要 OAuth2 认证
// 参考: https://developers.google.com/gmail/api
```

#### 2. 使用第三方验证服务

推荐的服务：
- **Hunter.io** - https://hunter.io/email-verifier
- **ZeroBounce** - https://www.zerobounce.net/
- **NeverBounce** - https://neverbounce.com/

#### 3. 使用 API 验证（gmailver.com）

这是目前最简单可靠的方案。

## 🛠️ 调试步骤

### 1. 检查 MX 记录

```bash
# 查询 Gmail 的 MX 记录
nslookup -type=mx gmail.com

# 或使用 dig
dig gmail.com MX
```

**预期输出**:
```
gmail.com	MX	5 gmail-smtp-in.l.google.com.
gmail.com	MX	10 alt1.gmail-smtp-in.l.google.com.
```

### 2. 测试 SMTP 连接

```bash
# 手动连接到 SMTP 服务器
telnet gmail-smtp-in.l.google.com 25

# 输入以下命令:
HELO example.com
MAIL FROM: <verify@example.com>
RCPT TO: <test@gmail.com>
QUIT
```

**成功响应**:
```
220 mx.google.com ESMTP
250 mx.google.com at your service
250 2.1.0 OK
250 2.1.5 OK  # 邮箱存在
221 2.0.0 closing connection
```

**失败响应**:
```
550 5.1.1 The email account that you tried to reach does not exist.
```

### 3. 检查防火墙

```bash
# 检查出站端口 25 是否开放
sudo iptables -L OUTPUT -n | grep 25

# 或使用 ufw
sudo ufw status
```

### 4. 查看详细日志

在后端代码中添加日志：

```go
func (v *SMTPVerifier) VerifyEmail(email string) (string, error) {
    log.Printf("[SMTP] Starting verification for: %s", email)

    // ... 验证逻辑

    log.Printf("[SMTP] Result for %s: %s (error: %v)", email, status, err)
    return status, err
}
```

## 🔧 配置优化

### 增加超时时间

编辑 `backend/internal/handlers/email_verify_smtp.go`:

```go
func NewSMTPVerifier() *SMTPVerifier {
    return &SMTPVerifier{
        fromEmail: "verify@example.com",
        timeout:   30 * time.Second, // 从 10 秒增加到 30 秒
    }
}
```

### 减少验证延迟

编辑 `backend/internal/handlers/email.go`:

```go
func (h *EmailHandler) verifyEmailsSMTP(emails []string) []VerifyEmailResponse {
    // ...
    time.Sleep(200 * time.Millisecond) // 从 500ms 减少到 200ms
}
```

### 使用代理

```go
func (v *SMTPVerifier) VerifyEmail(email string) (string, error) {
    // 使用 SOCKS5 代理
    dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1080", nil, proxy.Direct)
    if err != nil {
        return "", err
    }

    conn, err := dialer.Dial("tcp", mxHost+":25")
    // ...
}
```

## 📊 不同环境的验证成功率

### 本地开发环境

| 环境 | 成功率 | 说明 |
|------|--------|------|
| 家庭网络 | 30-50% | ISP 可能封禁端口 25 |
| 公司网络 | 10-30% | 防火墙限制 |
| VPN | 60-80% | 取决于 VPN 服务器 |

### 云服务器

| 服务商 | 成功率 | 说明 |
|--------|--------|------|
| AWS EC2 | 20-40% | 默认封禁端口 25 |
| Google Cloud | 0% | 完全封禁端口 25 |
| DigitalOcean | 50-70% | 需要申请解封 |
| Vultr | 60-80% | 相对宽松 |
| 自建服务器 | 80-95% | 取决于 ISP |

### 推荐方案

| 场景 | 推荐方案 | 原因 |
|------|---------|------|
| 本地开发 | API 验证 | 简单可靠 |
| 生产环境 | 第三方服务 | 专业稳定 |
| 少量验证 | API 验证 | 速度快 |
| 大量验证 | 第三方服务 | 批量优惠 |

## 🚀 替代方案

### 方案 1: 使用 API 验证（推荐）

**优点**:
- ✅ 不受网络限制
- ✅ 速度快
- ✅ 准确率高

**缺点**:
- ❌ 需要手动获取 key
- ❌ Key 会过期

**使用方法**:
```bash
# 前端选择 "API Verification"
# 输入从 gmailver.com 获取的 key
```

### 方案 2: 使用第三方验证服务

#### Hunter.io

```bash
curl "https://api.hunter.io/v2/email-verifier?email=test@gmail.com&api_key=YOUR_API_KEY"
```

**价格**: $49/月（1000 次验证）

#### ZeroBounce

```bash
curl "https://api.zerobounce.net/v2/validate?api_key=YOUR_API_KEY&email=test@gmail.com"
```

**价格**: $16/月（2000 次验证）

#### NeverBounce

```bash
curl "https://api.neverbounce.com/v4/single/check" \
  -d "key=YOUR_API_KEY" \
  -d "email=test@gmail.com"
```

**价格**: $8/月（1000 次验证）

### 方案 3: 仅验证域名（快速但不准确）

```go
func VerifyEmailQuick(email string) (string, error) {
    domain := strings.Split(email, "@")[1]

    // 只检查 MX 记录
    _, err := net.LookupMX(domain)
    if err != nil {
        return "dead", err
    }

    return "verify", nil // 域名有效，但不确定邮箱是否存在
}
```

**优点**:
- 速度极快（< 100ms）
- 不受端口限制

**缺点**:
- 只能验证域名，不能验证具体邮箱

### 方案 4: 发送验证邮件

```go
func VerifyEmailBySending(email string) (string, error) {
    // 发送带验证链接的邮件
    // 用户点击链接后标记为 verified

    // 优点: 100% 准确
    // 缺点: 需要用户操作，速度慢
}
```

## 📝 最佳实践

### 1. 混合验证策略

```go
func VerifyEmailSmart(email string) (string, error) {
    // 1. 先检查 MX 记录（快速）
    status, err := VerifyEmailQuick(email)
    if status == "dead" {
        return "dead", err
    }

    // 2. 尝试 SMTP 验证
    status, err = VerifyEmailSMTP(email)
    if err == nil {
        return status, nil
    }

    // 3. 如果 SMTP 失败，使用 API
    status, err = VerifyEmailAPI(email, apiKey)
    return status, err
}
```

### 2. 缓存验证结果

```go
// 使用 Redis 缓存验证结果 24 小时
func GetCachedStatus(email string) (string, bool) {
    val, err := redisClient.Get(ctx, "email:status:"+email).Result()
    if err == nil {
        return val, true
    }
    return "", false
}

func CacheStatus(email, status string) {
    redisClient.Set(ctx, "email:status:"+email, status, 24*time.Hour)
}
```

### 3. 异步验证

```go
// 对于大量邮箱，使用异步验证
func VerifyEmailsAsync(emails []string) <-chan VerifyResult {
    results := make(chan VerifyResult, len(emails))

    go func() {
        defer close(results)
        for _, email := range emails {
            status, err := VerifyEmail(email)
            results <- VerifyResult{Email: email, Status: status, Error: err}
        }
    }()

    return results
}
```

### 4. 限流保护

```go
// 限制每秒验证次数
var verifyLimiter = rate.NewLimiter(rate.Every(time.Second), 10)

func VerifyEmailWithRateLimit(email string) (string, error) {
    if !verifyLimiter.Allow() {
        return "", fmt.Errorf("rate limit exceeded")
    }
    return VerifyEmail(email)
}
```

## 🎯 总结

### SMTP 验证适用场景

✅ **适合**:
- 自建服务器（端口 25 未封禁）
- 少量邮箱验证（< 10 个）
- 不想依赖第三方服务

❌ **不适合**:
- 云服务器（端口 25 被封）
- 大量邮箱验证（> 50 个）
- 需要高准确率

### 推荐方案

| 场景 | 方案 | 原因 |
|------|------|------|
| **开发测试** | API 验证 | 简单快速 |
| **生产环境（少量）** | API 验证 | 可靠稳定 |
| **生产环境（大量）** | 第三方服务 | 专业高效 |
| **预算有限** | MX 记录检查 | 免费但不准确 |

### 当前状态

由于 SMTP 验证在大多数环境下会遇到端口封禁问题，**强烈建议使用 API 验证方式**：

1. 访问 https://gmailver.com
2. 获取验证 key
3. 在前端选择 "API Verification"
4. 输入 key 进行验证

这是目前最可靠的方案！

---

**文档版本**: v1.0.0
**更新时间**: 2025-01-25
**作者**: FreeGemini Team

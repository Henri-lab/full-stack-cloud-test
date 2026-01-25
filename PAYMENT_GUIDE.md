# License Key 付费系统使用指南

## 概述

本系统采用 **按次数付费** 的 License Key 模式，用户购买密钥后可获得指定次数的邮箱验证额度。

## 产品套餐

### 基础版 - ¥10
- **验证次数**: 100 次
- **功能权限**: 邮箱验证
- **适用场景**: 个人用户，少量验证需求

### 专业版 - ¥30
- **验证次数**: 500 次
- **功能权限**: 邮箱验证 + 邮箱导入 + 任务管理
- **适用场景**: 小团队，中等验证需求

### 企业版 - ¥50
- **验证次数**: 1000 次
- **功能权限**: 所有功能 + API 访问 + 优先支持
- **适用场景**: 企业用户，大量验证需求

## 购买流程

### 1. 访问支付页面
- 登录后点击导航栏的 **"License Key"** 按钮
- 或直接访问 `/payment` 路径

### 2. 选择产品
- 在 "购买 Key" 标签页浏览三种产品套餐
- 点击想要购买的产品卡片进行选择
- 选中的产品会高亮显示并出现 "立即购买" 按钮

### 3. 创建订单
- 点击 **"立即购买"** 按钮
- 系统自动生成订单号（格式：ORD + 时间戳 + 随机码）
- 订单有效期为 **15 分钟**

### 4. 完成支付

#### 方式一：扫码支付（生产环境）
1. 页面显示支付宝和微信收款码
2. 使用手机扫描对应的收款码
3. 完成支付后等待系统确认

#### 方式二：模拟支付（开发/测试环境）
1. 点击 **"模拟支付宝支付"** 或 **"模拟微信支付"** 按钮
2. 系统立即处理支付并生成 License Key
3. 3 秒后自动跳转到 "我的 Key" 页面

### 5. 获取 License Key
- 支付成功后系统自动生成 License Key
- Key 格式：`XXXX-XXXX-XXXX-XXXX`（16位十六进制）
- 在 "我的 Key" 标签页查看所有已购买的密钥

## 使用 License Key

### 1. 自动设置当前 Key（推荐）
- 支付成功后系统会自动把新 Key 设为当前使用 Key
- 或在 "我的 Key" 页面点击 **"设为当前"** 按钮
- 设置后前端会自动在请求里携带 `X-License-Key`

### 3. 验证邮箱
1. 输入并验证 License Key
2. 选择验证方法：
   - **SMTP 验证**：无需额外 key，直接验证（较慢但可靠）
   - **API 验证**：需要从 gmailver.com 获取 key（快速）
3. 勾选要验证的邮箱（可多选）
4. 点击 **"Verify"** 按钮开始验证
5. 每验证一个邮箱消耗 **1 次额度**

### 4. 查看剩余额度
- 验证成功后会显示剩余额度
- 在 "我的 Key" 页面查看进度条
- 额度用尽后密钥状态变为 "已用尽"

## License Key 状态

### Active（可用）
- 密钥正常可用
- 还有剩余额度
- 可以继续使用验证功能

### Exhausted（已用尽）
- 额度已全部消耗
- 无法继续使用
- 需要购买新的密钥

### Revoked（已撤销）
- 密钥被管理员撤销
- 无法使用
- 请联系客服

## 功能权限说明

### 基础版权限
- ✅ 邮箱验证（SMTP + API）
- ❌ 邮箱批量导入
- ❌ 任务管理
- ❌ API 访问

### 专业版权限
- ✅ 邮箱验证（SMTP + API）
- ✅ 邮箱批量导入
- ✅ 任务管理
- ❌ API 访问

### 企业版权限
- ✅ 邮箱验证（SMTP + API）
- ✅ 邮箱批量导入
- ✅ 任务管理
- ✅ API 访问
- ✅ 优先客服支持

## API 使用说明

### 获取产品列表
```bash
GET /api/v1/payments/products
Authorization: Bearer <JWT_TOKEN>
```

### 创建订单
```bash
POST /api/v1/payments/orders
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "product_type": "basic" | "pro" | "enterprise"
}
```

### 查询订单
```bash
GET /api/v1/payments/orders/:order_no
Authorization: Bearer <JWT_TOKEN>
```

### 支付回调（模拟）
```bash
POST /api/v1/payments/notify
Content-Type: application/json

{
  "order_no": "ORD20250125123456abcd",
  "transaction_id": "TXN1234567890",
  "payment_method": "alipay" | "wechat"
}
```

### 获取我的密钥
```bash
GET /api/v1/keys
Authorization: Bearer <JWT_TOKEN>
```

### 激活密钥
```bash
POST /api/v1/keys/activate
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "key_code": "XXXX-XXXX-XXXX-XXXX"
}
```

### 检查密钥状态
```bash
POST /api/v1/keys/check
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "key_code": "XXXX-XXXX-XXXX-XXXX"
}
```

### 使用密钥验证邮箱
```bash
POST /api/v1/emails/verify
Authorization: Bearer <JWT_TOKEN>
X-License-Key: XXXX-XXXX-XXXX-XXXX
Content-Type: application/json

{
  "mail": ["test@gmail.com", "test2@gmail.com"],
  "method": "smtp" | "api",
  "key": "gmailver_api_key" // 仅 API 方法需要
}
```

## 自动权限说明

- **邮箱导入**、**任务管理** 等功能需要 License Key
- 前端会自动携带当前 Key（无需手动复制）
- 若未设置 Key，会提示权限不足并引导设置

## 常见问题

### Q: 订单过期了怎么办？
A: 订单有效期为 15 分钟，过期后需要重新创建订单。未支付的订单会自动标记为 expired。

### Q: 可以退款吗？
A: 目前系统支持退款功能，但需要联系管理员处理。退款后密钥状态会变为 revoked。

### Q: 一个密钥可以多个账号使用吗？
A: 不可以。密钥首次使用时会绑定到当前用户账号，其他账号无法使用。

### Q: 额度用完了可以充值吗？
A: 不支持充值。额度用完后需要购买新的 License Key。

### Q: 如何查看历史订单？
A: 目前系统只显示已生成的密钥。如需查看订单历史，请联系管理员。

### Q: 支付失败怎么办？
A:
1. 检查订单是否过期（15分钟有效期）
2. 确认支付金额是否正确
3. 联系客服处理异常订单

### Q: License Key 丢失了怎么办？
A: 在 "我的 Key" 页面可以查看所有已购买的密钥。建议妥善保存密钥。

### Q: 可以转让 License Key 吗？
A: 不支持。密钥激活后绑定到用户账号，无法转让。

## 技术实现

### 数据模型

#### Payment（支付订单）
```go
type Payment struct {
    ID            uint
    UserID        uint
    OrderNo       string    // 订单号
    Amount        int       // 金额（分）
    ProductType   string    // basic/pro/enterprise
    QuotaAmount   int       // 购买的次数额度
    Status        string    // pending/paid/expired/refunded
    PaymentMethod string    // alipay/wechat
    TransactionID string    // 第三方支付流水号
    PaidAt        *time.Time
    ExpiredAt     time.Time // 订单过期时间
}
```

#### LicenseKey（授权密钥）
```go
type LicenseKey struct {
    ID          uint
    UserID      uint
    PaymentID   uint
    KeyCode     string    // 密钥代码
    ProductType string    // basic/pro/enterprise
    QuotaTotal  int       // 总次数
    QuotaUsed   int       // 已使用次数
    Status      string    // active/exhausted/revoked
    ActivatedAt *time.Time
}
```

### 中间件

#### LicenseKeyMiddleware
验证请求是否携带有效的 License Key，并检查功能权限。

#### ConsumeQuota
在请求成功后自动消耗指定数量的额度。

### 安全特性

1. **订单过期机制**：15 分钟未支付自动过期
2. **密钥绑定**：首次使用时绑定用户账号
3. **额度追踪**：每次使用自动扣减额度
4. **状态管理**：自动更新密钥状态（active → exhausted）
5. **权限控制**：基于产品类型的功能分级

## 集成真实支付

### 支付宝集成

1. 申请支付宝开放平台账号
2. 创建应用并获取 APP_ID 和密钥
3. 配置回调地址
4. 修改 `payment.go` 中的支付逻辑：

```go
// 生成支付宝支付链接
func generateAlipayURL(orderNo string, amount int) string {
    // 使用支付宝 SDK 生成支付链接
    // ...
}
```

### 微信支付集成

1. 申请微信商户平台账号
2. 获取商户号和 API 密钥
3. 配置回调地址
4. 修改 `payment.go` 中的支付逻辑：

```go
// 生成微信支付二维码
func generateWechatQRCode(orderNo string, amount int) string {
    // 使用微信支付 SDK 生成二维码
    // ...
}
```

### 回调处理

真实支付平台会在支付成功后调用回调接口：

```go
// PaymentNotify 处理支付回调
func (h *PaymentHandler) PaymentNotify(c *gin.Context) {
    // 1. 验证签名
    // 2. 查询订单
    // 3. 更新订单状态
    // 4. 生成 License Key
    // 5. 返回成功响应
}
```

## 监控和统计

建议添加以下监控指标：

1. **订单统计**
   - 每日订单数量
   - 订单转化率
   - 平均订单金额

2. **密钥使用**
   - 活跃密钥数量
   - 平均额度使用率
   - 密钥过期率

3. **收入统计**
   - 每日/每月收入
   - 产品销售占比
   - 用户复购率

## 后续优化建议

1. **订单管理**
   - 添加订单列表页面
   - 支持订单查询和筛选
   - 导出订单报表

2. **密钥管理**
   - 支持密钥备注
   - 密钥使用历史
   - 密钥分享（企业版）

3. **支付优化**
   - 支持更多支付方式
   - 自动重试失败订单
   - 支付成功通知（邮件/短信）

4. **用户体验**
   - 添加优惠券系统
   - 会员等级制度
   - 推荐奖励机制

5. **安全加固**
   - 支付签名验证
   - 订单防重放
   - 密钥防盗刷

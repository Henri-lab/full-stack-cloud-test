# FreeGemini 邮箱批量导入使用指南

## 快速开始

### 1. 启动服务

```bash
# 终端 1 - 启动后端
cd backend
go run cmd/api/main.go

# 终端 2 - 启动前端
cd frontend
npm run dev
```

### 2. 注册和登录

1. 访问 http://localhost:3000
2. 点击 "Register" 注册新用户
3. 密码要求：至少 12 位，包含大小写字母、数字和特殊字符
   - 示例：`MyPassword123!@#`
4. 注册成功后自动跳转到登录页
5. 登录成功后跳转到 Tasks 页面

### 3. 导入邮箱数据

1. 点击导航栏的 "Emails" 进入邮箱管理页面
2. 点击 "Import JSON" 按钮
3. 选择准备好的 JSON 文件（参考 `test-import.json`）
4. 等待导入完成
5. 查看成功提示：`Import successful (2 emails)`
6. 在 "Select saved dataset..." 中选择刚导入的数据
7. 点击 "Load Dataset" 加载到表格

### 4. 查看邮箱详情

导入成功后，表格会显示：
- **Main Email** - 主邮箱地址
- **Password** - 邮箱密码
- **Deputy Email** - 备用邮箱
- **2FA Key** - 两步验证密钥
- **2FA Code** - 点击 "Generate" 生成动态验证码
- **Status** - 状态标签（Active/Banned/Sold/Need Repair）

## JSON 文件格式说明

### 基本结构

```json
{
  "emails": [
    {
      "main": "必填 - 主邮箱地址",
      "password": "必填 - 邮箱密码",
      "deputy": "可选 - 备用邮箱",
      "key_2FA": "可选 - 2FA 密钥",
      "meta": {
        "banned": false,
        "price": 0,
        "sold": false,
        "need_repair": false,
        "from": "可选 - 来源标识"
      },
      "familys": [
        {
          "email": "必填 - family 邮箱地址",
          "password": "必填 - 密码",
          "code": "可选 - 验证码",
          "contact": "可选 - 联系方式",
          "issue": "可选 - 问题描述"
        }
      ]
    }
  ]
}
```

### 字段说明

#### Email 主对象
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| main | string | ✅ | 主邮箱地址，必须唯一 |
| password | string | ✅ | 邮箱密码 |
| deputy | string | ❌ | 备用邮箱地址 |
| key_2FA | string | ❌ | TOTP 两步验证密钥 |
| meta | object | ❌ | 元数据信息 |
| familys | array | ❌ | 关联的 family 邮箱列表 |

#### Meta 元数据
| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| banned | boolean | false | 是否被封禁 |
| price | integer | 0 | 价格 |
| sold | boolean | false | 是否已售出 |
| need_repair | boolean | false | 是否需要修复 |
| from | string | "" | 来源标识 |

#### Family 邮箱
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| email | string | ✅ | family 邮箱地址 |
| password | string | ✅ | 密码 |
| code | string | ❌ | 验证码 |
| contact | string | ❌ | 所有者联系方式（如：qq:123;phone:456） |
| issue | string | ❌ | 问题描述（如：账户失效、申请售后） |

## 导入规则

### ✅ 成功条件
- JSON 格式正确
- 所有必填字段都已提供
- main 邮箱地址在数据库中不存在
- 邮箱地址格式有效

### ❌ 失败情况
1. **JSON 格式错误**
   - 错误信息：`Invalid JSON format: ...`
   - 解决方案：使用 JSON 验证工具检查格式

2. **邮箱已存在**
   - 错误信息：`Email already exists: xxx@gmail.com`
   - 解决方案：修改 JSON 中的邮箱地址或删除数据库中的记录

3. **没有邮箱数据**
   - 错误信息：`No emails to import`
   - 解决方案：确保 JSON 中 emails 数组不为空

4. **文件未上传**
   - 错误信息：`No file uploaded`
   - 解决方案：确保选择了文件

## 示例场景

### 场景 1：导入单个邮箱（无 family）

```json
{
  "emails": [
    {
      "main": "simple@gmail.com",
      "password": "SimplePass123!",
      "deputy": "",
      "key_2FA": "",
      "familys": []
    }
  ]
}
```

### 场景 2：导入邮箱（带多个 family）

```json
{
  "emails": [
    {
      "main": "main@gmail.com",
      "password": "MainPass123!",
      "deputy": "backup@gmail.com",
      "key_2FA": "JBSWY3DPEHPK3PXP",
      "meta": {
        "banned": false,
        "price": 10,
        "sold": false,
        "need_repair": false,
        "from": "批发商A"
      },
      "familys": [
        {
          "email": "family1@gmail.com",
          "password": "Family1Pass!",
          "code": "123456",
          "contact": "qq:123456;phone:13800138000",
          "issue": "正常使用"
        },
        {
          "email": "family2@gmail.com",
          "password": "Family2Pass!",
          "code": "654321",
          "contact": "wechat:test123",
          "issue": "需要验证邮箱"
        },
        {
          "email": "family3@gmail.com",
          "password": "Family3Pass!",
          "code": "",
          "contact": "telegram:@test",
          "issue": "账户失效，申请售后"
        }
      ]
    }
  ]
}
```

### 场景 3：批量导入多个邮箱

```json
{
  "emails": [
    {
      "main": "batch1@gmail.com",
      "password": "Batch1Pass123!",
      "deputy": "backup1@gmail.com",
      "key_2FA": "JBSWY3DPEHPK3PXP",
      "familys": []
    },
    {
      "main": "batch2@gmail.com",
      "password": "Batch2Pass123!",
      "deputy": "backup2@gmail.com",
      "key_2FA": "HXDMVJECJJWSRB3H",
      "familys": []
    },
    {
      "main": "batch3@gmail.com",
      "password": "Batch3Pass123!",
      "deputy": "backup3@gmail.com",
      "key_2FA": "MFRGGZDFMZTWQ2LK",
      "familys": []
    }
  ]
}
```

## 常见问题

### Q1: 导入后看不到数据？
**A:** 检查以下几点：
1. 是否登录成功（查看浏览器 localStorage 是否有 token）
2. 是否在 Emails 页面（URL 应该是 /emails）
3. 打开浏览器开发者工具，查看 Network 标签是否有 API 请求
4. 查看后端日志是否有错误

### Q2: 如何生成 2FA 密钥？
**A:** 2FA 密钥是 Base32 编码的字符串，可以使用以下方式生成：
- 在线工具：搜索 "TOTP secret generator"
- 使用现有的 2FA 密钥（从 Google Authenticator 等应用导出）
- 留空：导入时可以不提供 key_2FA

### Q3: Contact 字段格式有要求吗？
**A:** 没有严格要求，建议使用分号分隔多个联系方式：
- `qq:123456;phone:13800138000`
- `wechat:test123;telegram:@test`
- `email:test@qq.com`

### Q4: 导入失败会影响数据库吗？
**A:** 不会。导入使用数据库事务，失败时会自动回滚所有更改。

### Q5: 可以更新已存在的邮箱吗？
**A:** 当前版本不支持。如果邮箱已存在，导入会失败。需要先删除旧记录或使用 PUT 接口单独更新。

## 技术细节

### 导入流程

```
用户选择文件
    ↓
前端上传 (FormData)
    ↓
后端接收文件
    ↓
解析 JSON
    ↓
验证格式
    ↓
检查重复 (main 字段)
    ↓
开始事务
    ↓
插入 Email 记录
    ↓
插入 EmailFamily 记录
    ↓
提交事务
    ↓
返回成功响应
```

### 事务保证

所有导入操作在单个数据库事务中执行：
- 如果任何一个邮箱插入失败，整个导入回滚
- 如果任何一个 family 邮箱插入失败，整个导入回滚
- 保证数据一致性

### 性能考虑

- 单次导入建议不超过 100 个邮箱
- 大批量导入建议分批进行
- 每个邮箱的 family 数量建议不超过 10 个

## 下一步

完成邮箱导入后，你可以：
1. 点击 "Generate" 生成 TOTP 动态验证码
2. 点击 "Copy" 按钮快速复制邮箱、密码等信息
3. 使用搜索框过滤邮箱
4. 查看邮箱状态（Active/Banned/Sold/Need Repair）

## 支持

如有问题，请查看：
- 项目文档：`CLAUDE.md`
- 后端日志：运行 `go run cmd/api/main.go` 的终端输出
- 前端日志：浏览器开发者工具 Console 标签

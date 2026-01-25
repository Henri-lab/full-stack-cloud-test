# 项目计划书

## 📈 性能指标
- 导入速度: ~100 邮箱/秒
- 数据库查询: <50ms（带索引）
- JWT 验证: <5ms
- 前端加载: <2s（首次）

## 🔒 安全检查清单
- ✅ 密码哈希（bcrypt cost=12）
- ✅ JWT 签名验证
- ✅ 登录限流（5次/15分钟）
- ✅ CORS 白名单
- ✅ SQL 注入防护
- ✅ XSS 防护
- ✅ 密码强度验证

## 📝 待优化项目与可执行计划

### 1) 功能增强

**目标**: 完成邮箱管理核心体验升级。

**可执行计划**
1. 分页查询（邮箱列表）
   - 后端: `GET /api/v1/emails` 支持 `page`、`page_size` 参数，返回 `total`、`items`。
   - 数据库: 为 `emails` 表添加索引（`main`、`deputy`、`created_at`）。
   - 前端: 加分页组件与页码状态，支持跳页与每页数量。
2. 搜索优化（全文搜索）
   - 后端: 增加 `q` 参数，支持 `main`、`deputy` 模糊搜索。
   - 数据库: 使用 `GIN`/`tsvector`（Postgres）或 `ILIKE` + 索引。
   - 前端: 搜索节流（300ms）并提示匹配条数。
3. 批量操作（批量删除/更新）
   - 后端: 批量接口 `POST /api/v1/emails/batch`（ids + action）。
   - 前端: 表格勾选与批量操作栏。
4. 导出功能（JSON/CSV）
   - 后端: `GET /api/v1/emails/export?format=json|csv`。
   - 前端: 下载按钮与进度提示。

### 2) 性能优化

**目标**: 降低延迟、提升列表渲染稳定性。

**可执行计划**
1. Redis 缓存
   - 缓存邮箱列表与常用统计（Active/Banned 等）。
   - 缓存键加入分页与搜索参数。
2. 数据库连接池调优
   - 调整 GORM 连接池：`SetMaxOpenConns`、`SetMaxIdleConns`、`SetConnMaxLifetime`。
3. 前端代码分割
   - Emails 页面懒加载；分离图表/重组件。
4. 图片懒加载
   - 若未来加入图片资源，统一使用 `loading="lazy"`。

### 3) 安全加固

**目标**: 强化认证与审计能力。

**可执行计划**
1. Refresh Token
   - 增加 refresh token 表与接口（rotate & revoke）。
2. 2FA 登录
   - 用户级 2FA 开关与备份码。
3. 审计日志
   - 记录关键操作与 IP/UA。
4. API 限流
   - 在网关或中间件添加速率限制。

### 4) 用户体验

**目标**: 提升可用性和多端体验。

**可执行计划**
1. 暗黑模式
   - 使用 Tailwind `dark` 模式切换。
2. 多语言支持
   - i18n 方案（如 `react-i18next`）。
3. 移动端适配
   - 表格移动端优化（卡片布局）。
4. 离线支持
   - Service Worker 缓存与离线提示页。

## 里程碑建议
- M1: 分页 + 搜索（2-4 天）
- M2: 批量操作 + 导出（3-5 天）
- M3: Redis 缓存 + 连接池调优（2-3 天）
- M4: Refresh Token + 审计日志（4-7 天）

## 账号平台需求落地计划（可配置/稳健）

### 业务规则建议
- 订阅改为有效期（Subscription），控制账号池“访问权限”
- 功能消耗型操作（验证/导出/绑定）可选配额（Quota）
- 家庭组需要解绑后再绑定，可选冷却期
- 临时账号默认限制同一用户同时占用 1 个，可按档位扩展
- 纯独享账号提供“再次查看凭证”入口，需审计/二次确认

### 可配置策略设计
- 规则配置化：订阅时长、家庭组容量、临时占用上限、冷却期
- 权限分层：订阅门槛 + 功能权限（可与 License/Quota 结合）
- 特性开关：可灰度启用（如“再次查看凭证”）

### 数据模型建议（主表 + 子表）
- accounts（统一账号表：type, main, password, key_2fa, status…）
- temporary_usages（account_id, user_id, started_at, expires_at, returned_at）
- exclusive_purchases（account_id, user_id, payment_id, purchased_at）
- family_groups（account_id, capacity）
- family_bindings（family_group_id, user_id, member_email, member_password_encrypted, created_at）
- subscriptions（user_id, plan, expires_at, status）
- audit_logs（user_id, action, target_id, metadata, created_at）

### 接口设计建议（可扩展）
- GET /accounts（默认脱敏，支持 type 过滤与分页）
- POST /accounts/temporary/claim | /release
- POST /accounts/exclusive/purchase
- POST /accounts/family/bind | /unbind
- GET /accounts/exclusive/:id/credentials（仅授权用户）
- GET /subscriptions/me | POST /subscriptions/renew

### 稳健性实施步骤
1) 抽离规则到配置表/环境变量（避免代码硬改）
2) 接口参数化（分页/过滤/策略开关）
3) 核心流程加审计日志
4) 权限中间件复用（订阅、License、Quota）
5) 版本化 API（v1/v2）逐步演进

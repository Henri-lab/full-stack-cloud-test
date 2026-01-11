# Go 企业级项目实战训练指南（以本项目为例）

> 面向：有一点 Go 基础、刚接触企业级 Web 项目 / 微服务的新手  
> 目标：通过一套可重复练习的步骤，学会在“真实”项目里定位、分析、解决问题，尤其是高并发和生产事故相关的问题。

---

## 1. 准备环境

本项目后端结构（简化）：

- 入口：`backend/cmd/api/main.go`
- 配置：`backend/internal/config/config.go`
- 数据库：`backend/internal/database/database.go`
- HTTP：Gin，路由：`/api/health`、`/api/v1/auth/*`、`/api/v1/tasks/*`
- ORM：GORM + PostgreSQL

推荐使用的运行方式：

```bash
cd backend
go run ./cmd/api
```

或者使用 Docker：

```bash
cd deployment
docker-compose -f docker-compose.yml up backend postgres redis
```

在开始所有训练前，确保：

- `/api/health` 能返回 `{"status":"ok","database":"connected"}`；
- 可以通过前端完成注册 / 登录 / 创建任务。

---

## 2. 基础训练：日志 + 错误排查

### 2.1 读懂项目结构和启动流程

练习：

1. 打开 `backend/cmd/api/main.go`，找出：
   - 配置加载在哪里 (`config.Load`)；
   - 数据库连接 / 迁移在哪里 (`database.Connect` / `database.Migrate`)；
   - 路由分组 `/api/v1`、`/auth`、`/tasks` 在哪里注册。
2. 改动一行日志，例如启动时打印当前环境：

```go
log.Printf("Server starting on port %s, env=%s", port, cfg.Environment)
```

3. 重新运行服务，观察控制台输出，确认修改生效。

### 2.2 模拟常见错误并排查

1. **模拟数据库配置错误**
   - 把 `.env.development` 或 `.env` 中的 `DB_PASSWORD` 改错；
   - 启动后端，观察错误日志；
   - 修复密码，再启动一次。
2. **模拟路由 404 / 校验错误**
   - 用 curl 或 Postman 调用不存在的路由 `/api/v1/tasks/xxx`；
   - 观察 Gin 返回什么状态码 / 错误信息；
   - 在 handler 中加上更友好的错误返回（例如自定义错误结构）。

目标：习惯通过日志 + HTTP 返回码来定位问题。

---

## 3. 并发基础：请求处理和数据库

这个项目是典型的“无共享状态”Web 服务：每个请求在自己的 Goroutine 里处理，核心共享资源是数据库连接池。

### 3.1 阅读数据库连接代码

练习：

1. 打开 `backend/internal/database/database.go`（如果还没看过）：
   - 找出 `gorm.Open` 的调用；
   - 看看是否设置了连接池参数（最大连接数、空闲连接数、超时等）。
2. 尝试增加连接池配置（示例）：

```go
sqlDB, err := db.DB()
if err != nil {
    return nil, err
}
sqlDB.SetMaxOpenConns(25)
sqlDB.SetMaxIdleConns(25)
sqlDB.SetConnMaxLifetime(5 * time.Minute)
```

3. 在开发环境跑起来，确认没有编译错误，再继续下面的并发测试。

---

## 4. 高并发模拟：压测入门

目标：对项目某个接口（比如 `/api/v1/tasks`）进行高并发请求，观察服务表现、发现瓶颈。

### 4.1 选择一个压测工具

推荐工具（任选其一，建议从简单的开始）：

- `hey`：简单易用，适合入门。
- `wrk`：更强大，Lua 脚本扩展。
- `k6`：脚本化测试（JavaScript），适合复杂场景。

下面以 `hey` 为例（其他工具可自行类比）。

安装（本机）：

```bash
go install github.com/rakyll/hey@latest
```

确认出现在 `$GOPATH/bin` 后：

```bash
hey -h
```

### 4.2 压测健康检查接口

先从轻量的 `/api/health` 开始：

```bash
hey -n 1000 -c 50 http://localhost:8080/api/health
```

- `-n 1000`：总请求数 1000；
- `-c 50`：并发数 50。

观察输出：

- QPS（Requests/sec）；
- 平均延迟、P95、P99；
- 有无非 200 状态码。

练习：

1. 尝试不同并发数：`-c 10`、`-c 100`、`-c 200`。
2. 记录一个简单表格：并发数、QPS、P95 延迟的变化。

### 4.3 压测带数据库的接口

选择一个读多写少的接口，例如获取任务列表 `/api/v1/tasks`：

1. 先在前端 / Postman 登录，拿到 JWT token。
2. 用 `hey` 压测时带上 Header（用 `-H`）：

```bash
hey -n 1000 -c 50 \
  -H "Authorization: Bearer <你的token>" \
  http://localhost:8080/api/v1/tasks
```

练习：

1. 比较 `/api/health` 和 `/api/v1/tasks` 的延迟和 QPS 差异。
2. 在数据库里增加任务数量（比如 1、100、1000 条），再分别压测，观察查询慢多少。

> 思考：当任务数量更多时，是否需要加索引？可以在 PostgreSQL 里对常用查询字段创建索引，然后再压测比较效果。

---

## 5. 并发问题排查：race / 泄漏 / 超时

### 5.1 使用竞争检测（race detector）

虽然当前项目主要共享资源是数据库，但我们依然可以通过演练来学习竞争检测。

1. 在某个 handler 中（例如 `taskHandler`），故意添加一个共享全局变量，并在多个请求中修改它。
2. 使用 `-race` 运行服务：

```bash
cd backend
go run -race ./cmd/api
```

3. 再用 `hey` 对该接口高并发请求，观察终端是否打印 data race 报告。

练习目标：理解 data race 是什么、如何利用 `-race` 找出问题代码行。

### 5.2 Goroutine 泄漏和超时处理

练习思路：

1. 在某个接口里，模拟一个耗时很久的操作（例如 `time.Sleep(10 * time.Second)`）。
2. 不设置超时的情况下压测，观察：
   - 客户端体验（大量请求堆积、超时）；
   - 服务器资源占用（Goroutine 数量上升）。
3. 然后在 Gin Handler 或数据库操作中引入 Context 超时：

```go
ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
defer cancel()

// 在数据库操作或外部调用中传入 ctx
```

4. 再次压测，观察是否在超过 2 秒后快速返回错误，避免无限堆积。

---

## 6. 性能分析：pprof 入门（可选进阶）

诉求：找到 CPU / 内存热点，了解系统在高并发下的瓶颈。

步骤概览：

1. 在 `main.go` 中引入 `net/http/pprof`，单独开一个 debug 端口：

```go
import (
    "net/http"
    _ "net/http/pprof"
)

func main() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()

    // 原来的 Gin 启动逻辑...
}
```

2. 高并发压测某接口（例如 `/api/v1/tasks`）。
3. 期间用 `go tool pprof` 拉取数据：

```bash
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

4. 在 pprof 交互界面中使用：
   - `top`：看最耗时函数；
   - `web`：生成调用图（需要 graphviz）。

练习目标：能看懂哪个函数最耗 CPU / 哪段代码分配内存最多。

---

## 7. 线上问题演练：从报警到修复的完整流程

可以模拟一次“生产事故”，按企业流程演练：

1. 人为制造一个 bug：
   - 例如在某个接口里对 nil 指针解引用，或 SQL 写错。
2. 在 Docker Compose 的生产环境起服务，前端访问该接口，让它报错。
3. 通过：
   - Docker 日志：`docker-compose -f deployment/docker-compose.prod.yml logs -f backend`
   - 应用日志（Gin / 自己的 log.Printf）
   - `/api/health` 状态
   分析问题。
4. 在本地修复 bug，写好单元测试 / 集成测试。
5. 使用 GitHub Actions 部署到服务器（已经配置好的 workflow），验证修复。

目标：体验一遍“问题出现 → 分析日志 → 定位代码 → 本地重现 → 修复 + 测试 → 部署”的完整闭环。

---

## 8. 推荐练习路线总结

按阶段反复练习，每次只加一点难度：

1. **第一轮**
   - 本地跑服务；
   - 用 `hey` 压 `/api/health`；
   - 改一点日志。
2. **第二轮**
   - 压 `/api/v1/tasks`；
   - 调整数据库连接池；
   - 手动在 PostgreSQL 里加数据、加索引。
3. **第三轮**
   - 引入 `-race`，刻意制造并发问题并修复；
   - 引入 Context 超时；
   - 用 pprof 看一次 CPU profile。
4. **第四轮**
   - 在 Docker Compose + 服务器上复现一次“生产事故”；
   - 用 GitHub Actions 部署修复版本。

只要你把上面的每一步都至少做一遍，就已经拥有了“在企业级 Go Web 项目里处理真实问题”的完整链路经验，而不仅仅是“能写 Go 语法”。  
后续可以在此基础上继续扩展：日志结构化、分布式 tracing、熔断 / 限流、微服务拆分等。


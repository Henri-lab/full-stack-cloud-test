# 部署文档

## 部署架构

本项目支持多种部署方式：
1. Docker Compose部署（推荐）
2. 手动部署
3. Kubernetes部署（未来支持）

## 服务器要求

### 最小配置
- CPU: 1核
- 内存: 2GB
- 存储: 20GB
- 操作系统: Debian 11+ / Ubuntu 20.04+

### 推荐配置
- CPU: 2核
- 内存: 4GB
- 存储: 50GB SSD
- 操作系统: Debian 11+ / Ubuntu 20.04+

## 部署步骤

### 1. 服务器初始化

在全新的Debian/Ubuntu服务器上运行初始化脚本：

```bash
# 下载项目
git clone <repository-url> /opt/fullstack
cd /opt/fullstack

# 运行初始化脚本
sudo bash deployment/scripts/init-server.sh
```

初始化脚本会安装：
- Docker和Docker Compose
- 防火墙(UFW)
- Fail2ban
- Certbot (SSL证书工具)
- 其他必要工具

### 2. 配置环境变量

```bash
cd /opt/fullstack
cp .env.example .env.production
```

编辑 `.env.production` 文件，配置生产环境变量：

```bash
# 数据库配置
DB_USER=postgres
DB_PASSWORD=<生成强密码>
DB_NAME=fullstack
DATABASE_URL=postgres://postgres:<密码>@postgres:5432/fullstack?sslmode=disable

# JWT配置 (至少32个字符的随机字符串)
JWT_SECRET=<生成随机密钥>

# 服务器配置
PORT=8080
ENVIRONMENT=production

# Grafana配置
GRAFANA_PASSWORD=<设置Grafana密码>
```

生成随机密钥：
```bash
openssl rand -base64 32
```

### 3. 配置域名

在域名DNS设置中，添加A记录指向服务器IP：

```
A    @    your-server-ip
A    www  your-server-ip
```

### 4. 部署应用

```bash
cd /opt/fullstack/deployment
./scripts/deploy.sh production
```

部署脚本会：
1. 拉取最新代码
2. 构建Docker镜像
3. 启动所有服务
4. 运行数据库迁移

### 5. 配置SSL证书

使用Let's Encrypt获取免费SSL证书：

```bash
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com
```

Certbot会自动配置Nginx并设置自动续期。

### 6. 验证部署

检查所有服务是否正常运行：

```bash
cd /opt/fullstack/deployment
docker-compose -f docker-compose.prod.yml ps
```

访问以下URL验证：
- 前端: https://yourdomain.com
- API健康检查: https://yourdomain.com/api/health
- Grafana: http://yourdomain.com:3000
- Prometheus: http://yourdomain.com:9090

## 更新部署

### 更新应用代码

```bash
cd /opt/fullstack
git pull origin main
cd deployment
./scripts/deploy.sh production
```

### 更新单个服务

```bash
cd /opt/fullstack/deployment

# 更新后端
docker-compose -f docker-compose.prod.yml up -d --build backend

# 更新前端
docker-compose -f docker-compose.prod.yml up -d --build frontend
```

## 备份和恢复

### 自动备份

设置定时备份任务：

```bash
# 编辑crontab
crontab -e

# 添加以下行（每天凌晨2点备份）
0 2 * * * /opt/fullstack/deployment/scripts/backup.sh
```

### 手动备份

```bash
cd /opt/fullstack/deployment
./scripts/backup.sh
```

备份文件保存在 `/opt/fullstack/backups/`

### 恢复数据库

```bash
# 恢复PostgreSQL
gunzip < /opt/fullstack/backups/db_backup_TIMESTAMP.sql.gz | \
docker exec -i fullstack-postgres-prod psql -U postgres -d fullstack

# 恢复Redis
docker cp /opt/fullstack/backups/redis_backup_TIMESTAMP.rdb \
fullstack-redis-prod:/data/dump.rdb
docker-compose -f docker-compose.prod.yml restart redis
```

## 监控

### Grafana

访问: http://your-server:3000

默认登录：
- 用户名: admin
- 密码: 在.env.production中配置的GRAFANA_PASSWORD

首次登录后：
1. 添加Prometheus数据源
2. 导入预配置的dashboard
3. 设置告警规则

### Prometheus

访问: http://your-server:9090

查看指标和配置告警规则。

### 日志查看

```bash
# 查看所有服务日志
docker-compose -f docker-compose.prod.yml logs -f

# 查看特定服务日志
docker-compose -f docker-compose.prod.yml logs -f backend
docker-compose -f docker-compose.prod.yml logs -f frontend
```

## 性能优化

### 数据库优化

1. **添加索引**

```sql
-- 为常用查询字段添加索引
CREATE INDEX idx_tasks_creator_id ON tasks(creator_id);
CREATE INDEX idx_tasks_status ON tasks(status);
```

2. **配置连接池**

在后端代码中已配置GORM连接池。

3. **定期维护**

```bash
# 进入数据库容器
docker exec -it fullstack-postgres-prod psql -U postgres -d fullstack

# 运行VACUUM
VACUUM ANALYZE;
```

### Nginx优化

编辑 `deployment/docker/nginx.conf`：

```nginx
# 增加worker进程
worker_processes auto;

# 增加连接数
events {
    worker_connections 2048;
}

# 启用缓存
proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=my_cache:10m;
```

### 应用优化

1. **启用Redis缓存** - 缓存频繁访问的数据
2. **使用CDN** - 加速静态资源
3. **数据库读写分离** - 使用主从复制
4. **水平扩展** - 增加后端实例

## 安全加固

### 1. 防火墙配置

```bash
# 只开放必要端口
sudo ufw allow 22/tcp   # SSH
sudo ufw allow 80/tcp   # HTTP
sudo ufw allow 443/tcp  # HTTPS
sudo ufw enable
```

### 2. SSH安全

编辑 `/etc/ssh/sshd_config`：

```
PermitRootLogin no
PasswordAuthentication no
PubkeyAuthentication yes
```

重启SSH服务：
```bash
sudo systemctl restart sshd
```

### 3. Fail2ban配置

Fail2ban已在初始化脚本中安装。检查状态：

```bash
sudo fail2ban-client status
```

### 4. 定期更新

```bash
# 更新系统包
sudo apt update && sudo apt upgrade -y

# 更新Docker镜像
cd /opt/fullstack/deployment
docker-compose -f docker-compose.prod.yml pull
docker-compose -f docker-compose.prod.yml up -d
```

### 5. 密钥管理

- 使用强密码（至少16个字符）
- 定期轮换JWT密钥
- 不要在代码中硬编码密钥
- 使用环境变量管理敏感信息

## 故障排查

### 服务无法启动

1. 检查Docker服务状态
```bash
sudo systemctl status docker
```

2. 查看容器日志
```bash
docker-compose -f docker-compose.prod.yml logs
```

3. 检查端口占用
```bash
sudo netstat -tulpn | grep LISTEN
```

### 数据库连接失败

1. 检查数据库容器状态
```bash
docker-compose -f docker-compose.prod.yml ps postgres
```

2. 测试数据库连接
```bash
docker exec -it fullstack-postgres-prod psql -U postgres -d fullstack
```

3. 检查环境变量配置

### 内存不足

1. 检查内存使用
```bash
free -h
docker stats
```

2. 增加swap空间
```bash
sudo fallocate -l 4G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

### SSL证书问题

1. 检查证书状态
```bash
sudo certbot certificates
```

2. 手动续期
```bash
sudo certbot renew
```

## 回滚

如果新版本出现问题，可以快速回滚：

```bash
# 1. 切换到上一个版本
cd /opt/fullstack
git checkout <previous-commit-hash>

# 2. 重新部署
cd deployment
./scripts/deploy.sh production

# 3. 恢复数据库（如果需要）
# 使用最近的备份恢复
```

## 扩展部署

### 负载均衡

使用Nginx作为负载均衡器，配置多个后端实例：

```nginx
upstream backend {
    server backend1:8080;
    server backend2:8080;
    server backend3:8080;
}

server {
    location /api/ {
        proxy_pass http://backend;
    }
}
```

### 数据库主从复制

配置PostgreSQL主从复制以提高读性能和可用性。

### 容器编排

对于大规模部署，考虑使用Kubernetes进行容器编排。

## 成本估算

### 云服务器（月度）

- **基础配置**: $5-10/月 (1核2G)
- **推荐配置**: $20-40/月 (2核4G)
- **高性能配置**: $80-160/月 (4核8G)

### 其他成本

- 域名: $10-15/年
- SSL证书: 免费 (Let's Encrypt)
- 备份存储: $5-10/月
- CDN: $5-20/月（可选）

## 支持

如遇到部署问题，请：
1. 查看日志文件
2. 检查GitHub Issues
3. 联系技术支持

## 检查清单

部署前检查：
- [ ] 服务器满足最低要求
- [ ] 域名DNS已配置
- [ ] 环境变量已正确配置
- [ ] 防火墙规则已设置
- [ ] SSL证书已配置
- [ ] 备份策略已设置
- [ ] 监控系统已配置
- [ ] 日志系统正常工作

部署后验证：
- [ ] 所有服务正常运行
- [ ] 前端可以访问
- [ ] API响应正常
- [ ] 数据库连接正常
- [ ] SSL证书有效
- [ ] 监控数据正常收集
- [ ] 备份任务正常执行

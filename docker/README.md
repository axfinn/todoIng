# 🐳 TodoIng Docker 部署方案

欢迎使用 TodoIng 的 Docker 部署方案！本目录提供了多种生产就绪的 Docker Compose 配置，适用于不同的部署场景和环境需求。

## � 快速开始

### 1. 选择部署方案

```bash
# 🔵 推荐：使用 Golang 后端
./deploy.sh golang up

# 🟢 使用 Node.js 后端  
./deploy.sh nodejs up

# 🛠️ 开发环境 (包含开发工具)
./deploy.sh dev up --profile golang
```

### 2. 访问应用
- **前端应用**: http://localhost
- **API 接口**: http://localhost:5004/api (Golang) 或 http://localhost:5001/api (Node.js)
- **API 文档**: http://localhost:5004/swagger/

### 3. 停止服务
```bash
./deploy.sh golang down
```

---

## 📋 可用部署方案

| 方案 | 文件 | 适用场景 | 端口 |
|------|------|----------|------|
| 🔵 **Golang** | `docker-compose.golang.yml` | 生产推荐，性能优秀 | 5004, 9090 |
| 🟢 **Node.js** | `docker-compose.nodejs.yml` | 传统部署，快速上手 | 5001 |
| 🛠️ **开发环境** | `docker-compose.dev-full.yml` | 开发调试，热重载 | 3000, 5001/5004 |
| 🚀 **生产环境** | `docker-compose.prod.yml` | 生产级，高可用 | 80, 443 |
| 🏗️ **微服务** | `docker-compose.microservices.yml` | 大型项目，服务解耦 | 8000 |

---

## 🔧 配置说明

### 环境变量配置

1. **复制配置模板**:
```bash
cp .env.example .env
```

2. **编辑配置文件**:
```bash
# 必须配置的变量
JWT_SECRET=your_super_secret_jwt_key_change_in_production
MONGO_ROOT_PASSWORD=your_secure_mongo_password
EMAIL_USER=your-email@gmail.com
EMAIL_PASS=your-email-app-password

# 可选配置
DISABLE_REGISTRATION=false
ENABLE_CAPTCHA=true
DEFAULT_USERNAME=admin
DEFAULT_EMAIL=admin@example.com
```

### Docker 镜像配置

默认使用预构建的远程镜像，加快部署速度：

```bash
# 🐳 Docker 镜像配置 (默认使用远程镜像)
GOLANG_BACKEND_IMAGE=axiu/todoing-go:latest      # Golang 后端镜像
NODEJS_BACKEND_IMAGE=axiu/todoing:latest         # Node.js 后端镜像  
FRONTEND_IMAGE=axiu/todoing-frontend:latest      # React 前端镜像
```

**镜像使用策略**:

- **远程镜像** (默认): 快速启动，无需本地构建
- **本地构建**: 使用 `--build` 选项或设置镜像环境变量为空

```bash
# 使用远程镜像 (推荐)
./deploy.sh golang up

# 强制本地构建
./deploy.sh golang up --build

# 或设置环境变量使用本地构建
export GOLANG_BACKEND_IMAGE=""
./deploy.sh golang up
```

### 域名和 SSL 配置 (生产环境)

```bash
# 生产环境变量
DOMAIN=your-domain.com
SSL_EMAIL=admin@your-domain.com
CORS_ORIGINS=https://your-domain.com
```

---

## 📖 详细使用指南

### 🔵 Golang 方案 (推荐)

最佳性能和功能完整性，包含 gRPC 支持和监控。

```bash
# 基础启动
./deploy.sh golang up

# 包含 gRPC 服务
./deploy.sh golang up --profile grpc

# 包含监控服务 (Prometheus + Grafana)
./deploy.sh golang up --profile monitoring

# 查看服务状态
./deploy.sh golang ps

# 查看后端日志
./deploy.sh golang logs backend-golang --follow
```

**访问地址**:
- 应用: http://localhost
- API: http://localhost:5004/api
- Swagger: http://localhost:5004/swagger/
- gRPC: localhost:9090 (需要 --profile grpc)
- Grafana: http://localhost:3001 (admin/admin123, 需要 --profile monitoring)

### 🟢 Node.js 方案

传统 JavaScript 全栈方案，适合快速部署。

```bash
# 启动服务
./deploy.sh nodejs up

# 查看日志
./deploy.sh nodejs logs backend-nodejs --follow
```

**访问地址**:
- 应用: http://localhost
- API: http://localhost:5001/api
- 数据库: localhost:27017

### 🛠️ 开发环境

包含完整的开发工具链，支持热重载和调试。

```bash
# 启动 Golang 开发环境
./deploy.sh dev up --profile golang

# 启动 Node.js 开发环境  
./deploy.sh dev up --profile nodejs

# 启动测试数据库
./deploy.sh dev up --profile testing

# 进入容器调试
./deploy.sh dev exec backend-golang-dev bash
```

**开发工具**:
- 前端开发: <http://localhost:3000> (Vite 热重载)
- Mongo Express: <http://localhost:8081> (admin/admin123)
- Redis Commander: <http://localhost:8082> (admin/admin123)
- MailHog: <http://localhost:8025> (邮件测试)
- API 文档: <http://localhost:4000>

### 🚀 生产环境

企业级部署，包含负载均衡、SSL、监控和备份。

```bash
# 基础生产部署
./deploy.sh prod up

# 包含数据库副本集
./deploy.sh prod up --profile replica

# 包含 ELK 日志收集
./deploy.sh prod up --profile logging

# 执行数据备份
./deploy.sh prod --profile backup run backup
```

**生产特性**:
- 负载均衡 (多后端实例)
- SSL/HTTPS 支持
- 数据库主从复制
- Redis 主从复制
- Prometheus + Grafana 监控
- ELK 日志收集 (可选)
- 自动备份

### 🏗️ 微服务方案

适合大型项目和团队协作的微服务架构。

```bash
# 启动微服务
./deploy.sh micro up

# 执行数据库迁移
./deploy.sh micro --profile migration run migrator
```

**微服务组件**:
- API 网关: Kong (<http://localhost:8000>)
- 服务发现: Consul (<http://localhost:8500>)
- 消息队列: RabbitMQ (<http://localhost:15672>)
- 对象存储: MinIO (<http://localhost:9001>)
- 链路追踪: Jaeger (<http://localhost:16686>)

**独立服务**:
- 认证服务 (端口: 5001)
- 任务服务 (端口: 5002)
- 报告服务 (端口: 5003)
- 通知服务 (端口: 5004)
- 文件服务 (端口: 5005)

---

## 🛠️ 运维指南

### 服务管理命令

```bash
# 查看所有服务状态
./deploy.sh [方案] ps

# 重启特定服务
./deploy.sh [方案] restart [服务名]

# 查看服务日志
./deploy.sh [方案] logs [服务名] --follow

# 进入容器
./deploy.sh [方案] exec [服务名] bash

# 强制重新构建
./deploy.sh [方案] up --build

# 更新镜像
./deploy.sh [方案] pull
```

### 数据管理

```bash
# 备份数据库
docker exec todoing_mongodb_golang mongodump --authenticationDatabase admin -u admin -p password123 --out /backup

# 恢复数据库
docker exec todoing_mongodb_golang mongorestore --authenticationDatabase admin -u admin -p password123 /backup

# 清理无用数据
./deploy.sh [方案] down --remove
```

### 监控和调试

```bash
# 查看资源使用情况
docker stats

# 查看网络连接
docker network ls

# 查看数据卷
docker volume ls

# 清理系统
docker system prune -a
```

---

## 🔒 安全最佳实践

### 1. 密码安全

- 使用强密码 (至少 16 位，包含大小写字母、数字、特殊字符)
- 定期更换密码
- 不要在代码中硬编码密码

### 2. 网络安全

- 生产环境启用 HTTPS
- 配置防火墙规则
- 限制数据库访问

### 3. 数据保护

- 定期备份数据
- 加密敏感数据
- 配置访问日志

### 4. 容器安全

- 定期更新镜像
- 使用非 root 用户运行容器
- 扫描镜像漏洞

---

## 📊 性能优化

### 1. 资源配置

```yaml
# 在 docker-compose.yml 中添加资源限制
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 2G
        reservations:
          cpus: '1.0'
          memory: 1G
```

### 2. 缓存策略

- Redis 缓存配置
- Nginx 静态文件缓存
- 数据库查询优化

### 3. 负载均衡

- 多实例部署
- 健康检查配置
- 连接池优化

---

## 🐛 故障排查

### 常见问题

**1. 端口冲突**

```bash
# 查看端口占用
lsof -i :5004

# 修改端口配置
vim docker-compose.golang.yml
```

**2. 数据库连接失败**

```bash
# 检查数据库状态
./deploy.sh golang ps
./deploy.sh golang logs mongodb-golang

# 检查网络连接
docker network inspect docker_todoing_golang_network
```

**3. 内存不足**

```bash
# 查看系统资源
docker stats
free -h

# 清理无用容器
docker system prune -a
```

**4. 服务启动失败**

```bash
# 查看详细日志
./deploy.sh golang logs [服务名] --follow

# 重新构建镜像
./deploy.sh golang build [服务名]
```

### 健康检查

所有服务都配置了健康检查，可通过以下方式查看：

```bash
# 查看健康状态
docker-compose -f docker-compose.golang.yml ps

# 手动执行健康检查
curl http://localhost:5004/health
```

---

## 📞 获取帮助

### 1. 查看部署脚本帮助

```bash
./deploy.sh --help
```

### 2. 查看 Docker Compose 帮助

```bash
docker-compose --help
```

### 3. 常用调试命令

```bash
# 查看容器详情
docker inspect [容器名]

# 查看容器进程
docker top [容器名]

# 进入容器
docker exec -it [容器名] bash
```

### 4. 社区支持

- GitHub Issues: [TodoIng Issues](https://github.com/axfinn/todoIng/issues)
- 文档: 查看项目根目录的 README.md
- 示例: 参考 `docker/` 目录下的配置文件

---

## 📝 更新日志

### v1.0.0 (2025-08-10)

- ✨ 初始发布
- 🔵 Golang 部署方案
- 🟢 Node.js 部署方案
- 🛠️ 开发环境方案
- 🚀 生产环境方案
- 🏗️ 微服务方案
- 📖 完整文档和脚本

---

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](../LICENSE) 文件了解详情。


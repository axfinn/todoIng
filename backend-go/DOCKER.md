# 🐳 TodoIng Golang 后端 - Docker 使用指南

本文档介绍如何使用 Docker 构建、运行和部署 TodoIng 的 Golang 后端服务。

## 📦 Docker 镜像信息

### 官方镜像
- **镜像名称**: `axiu/todoing-go`
- **Docker Hub**: https://hub.docker.com/r/axiu/todoing-go
- **支持架构**: `linux/amd64`, `linux/arm64`

### 标签说明
| 标签 | 描述 | 用途 |
|------|------|------|
| `latest` | 最新稳定版本 | 生产环境推荐 |
| `v1.0.0` | 具体版本号 | 生产环境固定版本 |
| `dev` | 开发分支版本 | 开发测试环境 |
| `YYYYMMDD-{sha}` | 每日构建版本 | CI/CD 环境 |

---

## 🚀 快速开始

### 1. 直接运行官方镜像

```bash
# 启动 Golang 后端服务
docker run -d \
  --name todoing-go-backend \
  -p 5004:5004 \
  -e MONGO_URI="mongodb://admin:password@localhost:27017/todoing?authSource=admin" \
  -e JWT_SECRET="your_super_secret_jwt_key_change_in_production" \
  axiu/todoing-go:latest

# 查看日志
docker logs -f todoing-go-backend

# 健康检查
curl http://localhost:5004/health
```

### 2. 使用 Docker Compose (推荐)

```bash
# 进入项目 docker 目录
cd /path/to/todoIng/docker

# 一键启动完整服务 (包含数据库)
./deploy.sh golang up

# 访问服务
echo "API: http://localhost:5004/api"
echo "Swagger: http://localhost:5004/swagger/"
```

---

## 🔨 本地构建

### 构建方式 1: 使用 Dockerfile

```bash
# 进入 Golang 后端目录
cd backend-go

# 构建镜像
docker build -t todoing-go:local .

# 运行本地构建的镜像
docker run -d -p 5004:5004 \
  -e MONGO_URI="mongodb://localhost:27017/todoing" \
  -e JWT_SECRET="your_jwt_secret" \
  todoing-go:local
```

### 构建方式 2: 使用 Docker Compose

```bash
# 在 docker 目录中本地构建
cd docker

# 设置环境变量使用本地构建
export GOLANG_BACKEND_IMAGE=""

# 或使用 --build 选项强制构建
./deploy.sh golang up --build
```

### 构建方式 3: 使用 Make 命令

```bash
# 进入 Golang 后端目录
cd backend-go

# 使用 Make 构建 Docker 镜像
make docker-build

# 运行构建的镜像
make docker-run

# 查看所有 Docker 相关命令
make help | grep docker
```

---

## ⚙️ 环境变量配置

### 必需环境变量

```bash
# 数据库连接
MONGO_URI=mongodb://admin:password@localhost:27017/todoing?authSource=admin

# JWT 密钥 (生产环境必须修改)
JWT_SECRET=your_super_secret_jwt_key_change_in_production_min_32_chars
```

### 可选环境变量

```bash
# 服务配置
PORT=5004                          # 服务端口
GIN_MODE=release                   # Gin 模式 (debug/release)
LOG_LEVEL=info                     # 日志级别 (debug/info/warn/error)

# JWT 配置
JWT_EXPIRY=24h                     # JWT 过期时间

# Redis 配置 (可选)
REDIS_URL=redis://:password@localhost:6379

# 功能开关
DISABLE_REGISTRATION=false         # 禁用注册功能
ENABLE_CAPTCHA=true               # 启用验证码
ENABLE_EMAIL_VERIFICATION=false   # 启用邮件验证

# 默认管理员账户
DEFAULT_USERNAME=admin
DEFAULT_EMAIL=admin@example.com
DEFAULT_PASSWORD=admin123

# 邮件配置 (如果启用邮件功能)
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USER=your-email@gmail.com
EMAIL_PASS=your-email-app-password

# CORS 配置
CORS_ORIGINS=http://localhost:3000,http://localhost:5173

# 监控配置
ENABLE_METRICS=true               # 启用 Prometheus 指标
METRICS_PORT=9090                 # 指标端口
```

---

## 🏗️ 多阶段构建详解

### Dockerfile 结构

```dockerfile
# 阶段 1: 构建阶段
FROM golang:1.23-alpine AS builder
# - 安装构建依赖
# - 下载 Go 模块
# - 编译 Go 应用

# 阶段 2: 运行阶段  
FROM alpine:3.18
# - 最小化运行环境
# - 非 root 用户
# - 安全配置
```

### 构建优化特性

- **多阶段构建**: 减小最终镜像大小 (< 20MB)
- **构建缓存**: 支持 Go 模块缓存和构建缓存
- **多架构支持**: 自动构建 AMD64 和 ARM64 版本
- **安全配置**: 非 root 用户运行，最小权限

---

## 🔍 健康检查和监控

### 健康检查端点

```bash
# 基础健康检查
curl http://localhost:5004/health
# 响应: {"status":"ok","timestamp":"2025-08-10T12:00:00Z"}

# 详细健康检查 (包含依赖检查)
curl http://localhost:5004/health/detailed
# 响应: {"status":"ok","dependencies":{"mongodb":"ok","redis":"ok"}}

# 就绪状态检查
curl http://localhost:5004/ready
# 响应: {"ready":true}
```

### Prometheus 监控指标

```bash
# 访问监控指标 (如果启用)
curl http://localhost:9090/metrics

# 常用指标
# - todoing_http_requests_total: HTTP 请求总数
# - todoing_http_request_duration_seconds: 请求耗时
# - todoing_database_connections: 数据库连接数
# - todoing_active_users: 活跃用户数
```

### Docker 容器监控

```bash
# 查看容器资源使用情况
docker stats todoing-go-backend

# 查看容器日志
docker logs -f todoing-go-backend

# 查看容器进程
docker top todoing-go-backend

# 进入容器调试
docker exec -it todoing-go-backend sh
```

---

## 🐛 故障排查

### 常见问题

#### 1. 容器启动失败

```bash
# 查看启动日志
docker logs todoing-go-backend

# 检查环境变量
docker exec todoing-go-backend env | grep -E "(MONGO|JWT|PORT)"

# 检查端口占用
lsof -i :5004
```

#### 2. 数据库连接失败

```bash
# 测试数据库连接
docker exec todoing-go-backend sh -c "echo 'db.runCommand(\"ping\")' | mongosh $MONGO_URI"

# 检查网络连接
docker network ls
docker network inspect [network_name]
```

#### 3. 内存或 CPU 问题

```bash
# 查看资源使用
docker stats todoing-go-backend

# 设置资源限制
docker run -d --name todoing-go-backend \
  --memory=512m \
  --cpus=1.0 \
  -p 5004:5004 \
  axiu/todoing-go:latest
```

### 调试模式

```bash
# 启用调试模式
docker run -d \
  -e GIN_MODE=debug \
  -e LOG_LEVEL=debug \
  -p 5004:5004 \
  axiu/todoing-go:latest

# 启用开发模式 (包含调试端口)
docker run -d \
  -p 5004:5004 \
  -p 2345:2345 \
  -e DEBUG=true \
  axiu/todoing-go:dev
```

---

## 🚀 部署最佳实践

### 生产环境部署

```bash
# 1. 使用固定版本标签
docker pull axiu/todoing-go:v1.0.0

# 2. 设置资源限制
docker run -d --name todoing-go-prod \
  --memory=1g \
  --cpus=2.0 \
  --restart=unless-stopped \
  -p 5004:5004 \
  axiu/todoing-go:v1.0.0

# 3. 使用 Docker Compose 生产配置
cd docker
./deploy.sh prod up --profile replica
```

### 安全配置

```bash
# 1. 非 root 用户运行 (镜像已配置)
docker run --user 1000:1000 axiu/todoing-go:latest

# 2. 只读文件系统
docker run --read-only \
  --tmpfs /tmp \
  --tmpfs /var/log \
  axiu/todoing-go:latest

# 3. 限制网络访问
docker run --network custom-network axiu/todoing-go:latest
```

### 高可用部署

```bash
# 1. 使用 Docker Swarm
docker service create \
  --name todoing-go-service \
  --replicas 3 \
  --publish 5004:5004 \
  axiu/todoing-go:latest

# 2. 使用 Kubernetes
kubectl apply -f k8s/todoing-go-deployment.yaml

# 3. 使用负载均衡器
docker run -d --name nginx-lb \
  -p 80:80 \
  -v ./nginx.conf:/etc/nginx/nginx.conf \
  nginx:alpine
```

---

## 📊 性能优化

### 容器优化

```bash
# 1. 调整 Go 运行时参数
docker run -d \
  -e GOGC=100 \
  -e GOMAXPROCS=2 \
  -e GOMEMLIMIT=1GiB \
  axiu/todoing-go:latest

# 2. 优化垃圾回收
docker run -d \
  -e GOGC=50 \
  -e GODEBUG=madvdontneed=1 \
  axiu/todoing-go:latest
```

### 数据库连接池

```bash
# 设置数据库连接池参数
docker run -d \
  -e MONGO_MAX_POOL_SIZE=100 \
  -e MONGO_MIN_POOL_SIZE=10 \
  -e MONGO_MAX_IDLE_TIME=300s \
  axiu/todoing-go:latest
```

---

## 🔄 CI/CD 集成

### GitHub Actions 使用

```yaml
# 在 .github/workflows/docker-publish.yml 中
- name: Build and push Golang backend
  uses: docker/build-push-action@v5
  with:
    context: ./backend-go
    platforms: linux/amd64,linux/arm64
    push: true
    tags: axiu/todoing-go:latest,axiu/todoing-go:${{ github.ref_name }}
```

### 自动化测试

```bash
# 运行集成测试
docker run --rm \
  -e TEST_MODE=true \
  -e MONGO_URI="mongodb://mongo-test:27017/test" \
  axiu/todoing-go:latest \
  go test ./...

# 使用测试专用镜像
docker build --target test -t todoing-go-test .
docker run --rm todoing-go-test
```

---

## 📚 相关文档

- [项目主文档](../README.md)
- [Docker 部署指南](../docker/README.md)
- [API 文档](./docs/api.md)
- [开发指南](./docs/development.md)
- [配置说明](./docs/configuration.md)

---

## 🤝 贡献指南

### 镜像构建贡献

1. **Fork 项目**
2. **创建功能分支**: `git checkout -b feature/docker-improvement`
3. **修改 Dockerfile**: 优化构建过程或添加新功能
4. **测试构建**: `make docker-build && make docker-test`
5. **提交 PR**: 描述改进内容和测试结果

### 文档贡献

欢迎改进本文档，包括但不限于：
- 添加新的使用场景
- 完善故障排查指南
- 补充性能优化建议
- 更新最佳实践

---

## 📞 获取帮助

- **GitHub Issues**: [提交问题](https://github.com/axfinn/todoIng/issues)
- **Docker Hub**: [镜像页面](https://hub.docker.com/r/axiu/todoing-go)
- **文档**: 查看项目根目录的完整文档
- **社区**: 加入讨论和交流

---

*最后更新: 2025-08-10*

# 🚀 TodoIng Backend - Go 版本

<div align="center">

**高性能的 Go 语言任务管理系统后端服务**

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org/)
[![gRPC](https://img.shields.io/badge/gRPC-Ready-green.svg)](https://grpc.io/)
[![MongoDB](https://img.shields.io/badge/MongoDB-5.0+-green.svg)](https://mongodb.com/)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://docker.com/)
[![API Docs](https://img.shields.io/badge/API-Swagger-orange.svg)](http://localhost:5004/swagger/)

</div>

## ✨ 功能概述

TodoIng Go 后端是 Node.js 版本的高性能重构，采用现代化的 Go 技术栈，提供：

### 🎯 核心功能
- ✅ **用户认证系统** - JWT 令牌认证，邮箱验证码
- ✅ **任务管理** - 完整的 CRUD 操作，导入导出功能
- ✅ **报表生成** - 日报、周报、月报自动生成
- ✅ **AI 润色** - 集成 OpenAI 的报表优化功能
- ✅ **图形验证码** - 防机器人注册登录保护
- ✅ **邮箱验证** - 完整的邮箱验证码系统

### 🔌 API 架构
- ✅ **RESTful API** - 标准的 HTTP JSON API
- ✅ **gRPC 服务** - 高性能的二进制协议
- ✅ **Swagger 文档** - 自动生成的 API 文档
- ✅ **中间件支持** - 认证、日志、错误处理
- ✅ **结构化日志** - 完整的请求追踪

## 🏗️ 技术架构

### 核心技术栈
```
📦 技术选型
├── 🌐 Web 框架: Gorilla Mux
├── 🗄️ 数据库: MongoDB + Official Driver
├── 🔐 认证: JWT + Bcrypt
├── 📡 RPC: gRPC + Protobuf
├── 📚 文档: Swagger/OpenAPI
├── 📧 邮件: Go-Mail
├── 🎯 验证码: 自研图形验证码
└── 🐳 部署: Docker + Docker Compose
```

### 项目结构
```
backend-go/
├── 📁 cmd/                    # 应用程序入口
│   ├── api/                   # HTTP API 服务器
│   └── grpc/                  # gRPC 服务器
├── 📁 internal/               # 内部代码
│   ├── api/                   # HTTP 处理器
│   ├── auth/                  # 认证服务
│   ├── captcha/               # 验证码服务
│   ├── config/                # 配置管理
│   ├── email/                 # 邮件服务
│   ├── models/                # 数据模型
│   └── services/              # 业务逻辑
├── 📁 api/proto/v1/           # Protobuf 定义
├── 📁 docs/                   # API 文档
├── 📁 tools/                  # 开发工具
│   └── generate_complete_api.go # 文档生成器
├── 📁 scripts/                # 构建脚本
├── 🐳 Dockerfile             # 容器镜像
├── 🔧 Makefile               # 构建配置
└── 📦 go.mod                 # Go 模块定义
```

## 🚀 快速开始

### 📋 环境要求
- **Go** 1.23 或更高版本
- **MongoDB** 5.0 或更高版本
- **Protocol Buffers** 编译器 (可选，用于 gRPC)

### 🛠️ 开发环境搭建

#### 1. 克隆并准备项目
```bash
# 克隆项目
git clone https://github.com/axfinn/todoIng.git
cd todoIng/backend-go

# 安装依赖
make deps

# 生成 API 文档
make docs
```

#### 2. 配置环境变量
```bash
# 复制环境变量模板
cp .env.example .env

# 编辑配置文件
vim .env
```

#### 3. 启动数据库
```bash
# 使用 Docker 启动 MongoDB
cd .. && docker-compose -f docker-compose.dev.yml up mongodb -d

# 或使用本地 MongoDB
mongod --dbpath ./data/db
```

#### 4. 运行服务
```bash
# 方式1: 使用 Makefile
make run

# 方式2: 直接运行
go run ./cmd/api/main.go

# 方式3: 构建后运行
make build && ./server
```

#### 5. 验证服务
```bash
# 健康检查
curl http://localhost:5004/health

# 查看 API 文档
open http://localhost:5004/swagger/
```

### 🐳 Docker 部署

> 📖 **详细的 Docker 使用指南**: 查看 [DOCKER.md](./DOCKER.md) 获取完整的镜像使用文档

#### 快速启动 (使用官方镜像)

```bash
# 方式1: 直接运行官方镜像
docker run -d \
  --name todoing-go-backend \
  -p 5004:5004 \
  -e MONGO_URI="mongodb://admin:password@localhost:27017/todoing?authSource=admin" \
  -e JWT_SECRET="your_super_secret_jwt_key" \
  axiu/todoing-go:latest

# 方式2: 使用完整的 Docker Compose (推荐)
cd ../docker
./deploy.sh golang up
```

#### 本地构建镜像

```bash
# 构建镜像
docker build -t todoing-backend:local .

# 或使用 Makefile
make docker-build

# 或使用 Docker Compose 本地构建
cd ../docker
./deploy.sh golang up --build
```

#### 镜像信息

- **官方镜像**: `axiu/todoing-go:latest`
- **支持架构**: `linux/amd64`, `linux/arm64`
- **Docker Hub**: [axiu/todoing-go](https://hub.docker.com/r/axiu/todoing-go)
- **自动更新**: 每次发布都会自动构建新镜像

## 🔧 开发工具

### Makefile 命令
```bash
# 📚 查看所有可用命令
make help

# 🏗️ 开发命令
make deps              # 下载依赖
make build             # 构建 HTTP 服务器
make build-grpc        # 构建 gRPC 服务器
make run               # 运行 HTTP 服务器
make run-grpc          # 运行 gRPC 服务器

# 📖 文档相关
make docs              # 生成完整 API 文档
make docs-dev          # 快速文档更新
make docs-clean        # 清理文档文件
make docs-serve        # 生成文档并启动服务

# 🧪 测试相关
make test              # 运行所有测试
make test-unit         # 运行单元测试
make test-coverage     # 生成覆盖率报告

# 🔍 代码质量
make fmt               # 格式化代码
make lint              # 代码检查
make vet               # 代码静态分析

# 🐳 Docker 相关
make docker-build      # 构建 Docker 镜像
make docker-run        # 运行 Docker 容器

# 🧹 清理命令
make clean             # 清理构建文件
make clean-proto       # 清理生成的 Proto 文件
```

### API 文档生成

系统提供了强大的 API 文档自动生成功能：

```bash
# 生成完整 API 文档 (包含 REST + gRPC)
make docs

# 输出示例:
# 🚀 生成 TodoIng 完整 API 文档...
# ✅ 完整 API 文档生成成功!
# 📄 文件位置: docs/api_complete.json
# 📊 包含接口: 16 个 REST API + 5 个 gRPC API
# 📋 数据模型: 26 个
# 🔗 访问地址: http://localhost:5004/swagger/
```

#### 文档包含内容
- **16个 REST API 接口**：完整的 HTTP 接口规范
- **5个 gRPC API 接口**：Protobuf 定义的高性能接口
- **26个数据模型**：详细的请求/响应类型定义
- **完整示例**：每个接口都包含示例数据

## 🔌 API 接口

### REST API 端点

#### 🔐 认证模块
```
POST   /api/auth/register           # 用户注册
POST   /api/auth/login              # 用户登录  
GET    /api/auth/me                 # 获取用户信息
POST   /api/auth/send-email-code    # 发送注册验证码
POST   /api/auth/send-login-email-code # 发送登录验证码
GET    /api/auth/captcha            # 获取图形验证码
POST   /api/auth/verify-captcha     # 验证图形验证码
```

#### 📋 任务管理
```
GET    /api/tasks                   # 获取任务列表
POST   /api/tasks                   # 创建新任务
GET    /api/tasks/{id}              # 获取任务详情
PUT    /api/tasks/{id}              # 更新任务
DELETE /api/tasks/{id}              # 删除任务
GET    /api/tasks/export/all        # 导出所有任务
POST   /api/tasks/import            # 批量导入任务
```

#### 📊 报表管理
```
GET    /api/reports                 # 获取报表列表
POST   /api/reports/generate        # 生成新报表
GET    /api/reports/{id}            # 获取报表详情
DELETE /api/reports/{id}            # 删除报表
POST   /api/reports/{id}/polish     # AI 润色报表
GET    /api/reports/{id}/export/{format} # 导出报表 (pdf/excel/word)
```

### gRPC 服务

#### 认证服务 (AuthService)
```protobuf
service AuthService {
  rpc Login(AuthLoginRequest) returns (AuthLoginResponse);
  rpc Register(AuthRegisterRequest) returns (AuthRegisterResponse);
  rpc GetUserInfo(GetUserInfoRequest) returns (GetUserInfoResponse);
}
```

#### 任务服务 (TaskService)
```protobuf
service TaskService {
  rpc CreateTask(CreateTaskRequest) returns (CreateTaskResponse);
  rpc ListTasks(ListTasksRequest) returns (ListTasksResponse);
}
```

#### 报表服务 (ReportService)
```protobuf
service ReportService {
  rpc GenerateReport(GenerateReportRequest) returns (GenerateReportResponse);
}
```

#### 验证码服务 (CaptchaService)
```protobuf
service CaptchaService {
  rpc GenerateCaptcha(GenerateCaptchaRequest) returns (GenerateCaptchaResponse);
  rpc VerifyCaptcha(VerifyCaptchaRequest) returns (VerifyCaptchaResponse);
}
```

## ⚙️ 配置管理

### 环境变量
```bash
# 🗄️ 数据库配置
MONGODB_URI=mongodb://localhost:27017/todoing
DB_NAME=todoing

# 🔐 JWT 认证配置
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRES_IN=168h  # 7 天

# 📧 邮件服务配置
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password

# 🤖 OpenAI 配置 (AI 润色功能)
OPENAI_API_KEY=sk-your-openai-api-key
OPENAI_MODEL=gpt-3.5-turbo

# 🌐 服务配置
HTTP_PORT=5004          # HTTP 服务端口
GRPC_PORT=9000         # gRPC 服务端口
LOG_LEVEL=info         # 日志级别

# 🔧 功能开关
ENABLE_CAPTCHA=true              # 启用图形验证码
ENABLE_EMAIL_VERIFICATION=true  # 启用邮箱验证
DISABLE_REGISTRATION=false      # 禁用注册功能
DEBUG_MODE=false               # 调试模式
```

### 功能开关说明
| 环境变量 | 类型 | 默认值 | 说明 |
|----------|------|--------|------|
| `ENABLE_CAPTCHA` | boolean | `false` | 启用图形验证码功能 |
| `ENABLE_EMAIL_VERIFICATION` | boolean | `false` | 启用邮箱验证码功能 |
| `DISABLE_REGISTRATION` | boolean | `false` | 禁用用户注册功能 |
| `DEBUG_MODE` | boolean | `false` | 启用详细调试日志 |
| `CORS_ENABLED` | boolean | `true` | 启用跨域请求支持 |

## 🧪 测试

### 运行测试
```bash
# 运行所有测试
make test

# 运行单元测试
make test-unit

# 运行转换层测试
make test-convert

# 生成覆盖率报告
make test-coverage
```

### API 测试示例
```bash
# 1. 健康检查
curl http://localhost:5004/health

# 2. 获取验证码
curl http://localhost:5004/api/auth/captcha

# 3. 用户注册
curl -X POST http://localhost:5004/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123456",
    "nickname": "测试用户",
    "emailCode": "123456",
    "emailCodeId": "code-id-123"
  }'

# 4. 用户登录
curl -X POST http://localhost:5004/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123456"
  }'

# 5. 创建任务 (需要 JWT Token)
curl -X POST http://localhost:5004/api/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "完成项目文档",
    "description": "编写详细的API文档",
    "priority": "high",
    "dueDate": "2024-12-31T23:59:59Z"
  }'
```

## 📊 性能指标

### 系统性能
- **启动时间**: < 3 秒
- **API 响应时间**: 平均 < 100ms
- **内存占用**: < 50MB (空载)
- **并发处理**: 1000+ 请求/秒
- **gRPC 性能**: 比 REST API 快 30-50%

### 基准测试
```bash
# HTTP API 压力测试
ab -n 10000 -c 100 http://localhost:5004/health

# gRPC 性能测试 (需要 ghz 工具)
ghz --insecure --proto ./api/proto/v1/auth.proto \
    --call AuthService.Login \
    -d '{"email":"test@example.com","password":"123456"}' \
    localhost:9000
```

## 🚀 生产部署

### Docker 部署
```bash
# 构建生产镜像
docker build -t todoing-backend:prod -f Dockerfile .

# 运行容器
docker run -d \
  --name todoing-backend \
  -p 5004:5004 \
  -e MONGODB_URI=mongodb://mongo:27017/todoing \
  todoing-backend:prod
```

### Kubernetes 部署
```yaml
# k8s-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: todoing-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: todoing-backend
  template:
    metadata:
      labels:
        app: todoing-backend
    spec:
      containers:
      - name: backend
        image: todoing-backend:prod
        ports:
        - containerPort: 5004
        env:
        - name: MONGODB_URI
          value: "mongodb://mongo-service:27017/todoing"
```

### 云平台部署建议
- **AWS**: ECS + RDS DocumentDB
- **Google Cloud**: Cloud Run + Firestore
- **Azure**: Container Instances + Cosmos DB
- **阿里云**: ACK + MongoDB Atlas

## 🔍 监控和日志

### 结构化日志
系统使用结构化日志记录，便于分析和监控：

```json
{
  "timestamp": "2024-08-10T12:00:00Z",
  "level": "info",
  "service": "todoing-backend",
  "method": "POST",
  "path": "/api/tasks",
  "status": 201,
  "duration": "45ms",
  "user_id": "user123",
  "request_id": "req-abc123"
}
```

### 健康检查端点
```bash
# 基础健康检查
GET /health

# 详细健康状态
GET /health/detailed

# 数据库连接状态
GET /health/db
```

## 🤝 贡献指南

### 开发流程
1. **Fork 项目** 并创建功能分支
2. **遵循代码规范**: `gofmt`, `golint`, `go vet`
3. **编写测试** 确保代码覆盖率 > 80%
4. **更新文档** 包括 API 变更
5. **提交 PR** 并等待 Review

### 代码规范
```bash
# 格式化代码
make fmt

# 静态检查
make vet

# 代码规范检查
make lint

# 运行所有检查
make fmt vet lint test
```

## 📝 更新日志

查看 [CHANGELOG.md](./CHANGELOG.md) 了解详细的版本更新历史。

## 🔗 相关链接

- **项目主页**: https://github.com/axfinn/todoIng
- **API 文档**: http://localhost:5004/swagger/
- **问题反馈**: https://github.com/axfinn/todoIng/issues
- **gRPC 文档**: [Proto 文件](./api/proto/v1/)

## 📄 许可证

本项目采用 [MIT 许可证](../LICENSE)。

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给我们一个 Star！**

Made with ❤️ in Go | [返回主项目](../)

</div>


# todoIng - 现代化任务管理系统

![todoIng Dashboard](img/dashboard.png)

一个使用 React、TypeScript、Node.js、MongoDB 和 Docker 构建的现代化任务管理系统。

## 核心特性

- **任务生命周期管理** - 完整的任务状态跟踪（创建、进行中、暂停、完成、取消）
- **变更历史追踪** - 记录任务的每一次变更，支持随时回溯任务状态
- **可视化历史展示** - 以图形化方式展示任务变更历史，类似 Git 提交历史
- **Web 图形化界面** - 直观的用户界面，便于操作和查看任务状态
- **实时同步** - 多设备间任务状态实时同步
- **数据统计** - 提供任务完成情况的统计和分析
- **Docker 自动构建** - 使用 GitHub Actions 实现的自动化 Docker 镜像构建和部署

## 文档

详细的项目设计文档请查看 [docs](./docs) 目录：

- [项目概述](../README.md)
- [系统架构与技术设计](docs/technical-design.md)
- [UI/UX 设计](docs/ui-ux-design.md)
- [API 设计](docs/api-design.md)
- [数据库设计](docs/database-design.md)
- [实现技术方案](docs/implementation-plan.md)

开发相关文档请查看 [docs/development](./docs/development) 目录：

- [开发计划](docs/development/development-plan.md)

## 技术栈

- 前端: React 18, TypeScript, Redux, Bootstrap 5, React Router v6
- 后端: Node.js, Express, Mongoose
- 数据库: MongoDB
- 部署: Docker, Docker Compose
- 其他: i18next (多语言), Axios (HTTP 客户端)

## 快速开始

### 使用 Docker Compose 部署（推荐）

```bash
# 克隆项目
git clone https://github.com/axfinn/todoIng.git
cd todoIng

# 启动服务
docker-compose up -d

# 访问应用
# 前端: http://localhost
# 后端 API: http://localhost:5001/api
```

### 使用 Docker Hub 镜像

```bash
# 拉取镜像
docker pull axiu/todoing:latest
docker pull axiu/todoing-frontend:latest

# 运行服务
docker-compose -f docker-compose.local.yml up -d
```

#### 开发环境部署
```bash
# 克隆项目
git clone <repository-url>
cd todoIng

# 构建并启动开发环境
docker-compose -f docker-compose.dev.yml up -d

# 应用将在以下地址可用：
# 前端: http://localhost:3000
# 后端API: http://localhost:5001
```

### 方法二：手动部署

#### 后端设置
```bash
# 进入后端目录
cd backend

# 安装依赖
npm install

# 创建.env文件并配置环境变量
cp .env.example .env
# 编辑.env文件，设置MONGO_URI和JWT_SECRET

# 启动后端服务
npm start
# 或者开发模式
npm run dev
```

#### 前端设置
```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 创建.env文件并配置环境变量
cp .env.example .env
# 编辑.env文件，设置VITE_API_URL

# 启动前端开发服务器
npm run dev

# 或者构建生产版本
npm run build
```

## Docker 镜像

预构建的 Docker 镜像可在 Docker Hub 获取:

- 后端镜像: [axiu/todoing](https://hub.docker.com/r/axiu/todoing)
- 前端镜像: [axiu/todoing-frontend](https://hub.docker.com/r/axiu/todoing-frontend)

## Docker 自动构建

本项目使用 GitHub Actions 实现 Docker 镜像的自动构建和推送。每当有新的 Git 标签创建时，GitHub Actions 会自动构建 Docker 镜像并推送到 Docker Hub。

### 配置自动构建

要为你的 fork 配置自动构建，需要在 GitHub 仓库中设置以下 Secrets:

1. `DOCKERHUB_USERNAME` - 你的 Docker Hub 用户名
2. `DOCKERHUB_TOKEN` - 你的 Docker Hub 访问令牌

生成 Docker Hub 访问令牌的步骤:
1. 登录到 [Docker Hub](https://hub.docker.com/)
2. 进入 Account Settings（账户设置）
3. 点击 Security（安全）选项卡
4. 点击 "New Access Token"（新建访问令牌）
5. 为令牌添加描述（例如："GitHub Actions"）
6. 选择适当的权限（通常选择 Read & Write）
7. 点击 "Generate"（生成）
8. 复制生成的令牌并将其作为 `DOCKERHUB_TOKEN` Secret 添加到 GitHub

## 赞助作者

如果你觉得这个项目对你有帮助，欢迎赞助作者一杯咖啡！

<div style="display: flex; justify-content: center; gap: 20px; margin: 20px 0;">
  <div style="text-align: center;">
    <h4>支付宝</h4>
    <img src="./img/alipay.JPG" alt="支付宝收款码" width="200">
  </div>
  <div style="text-align: center;">
    <h4>微信</h4>
    <img src="./img/wxpay.JPG" alt="微信收款码" width="200">
  </div>
</div>
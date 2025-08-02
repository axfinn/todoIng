# todoIng

todoIng 是一个创新的任务管理系统，它不仅提供基本的任务管理功能，还为每个任务提供完整的生命周期追踪和变更历史记录，就像 Git 管理代码变更一样管理任务。

## 核心特性

- **任务生命周期管理** - 完整的任务状态跟踪（创建、进行中、暂停、完成、取消）
- **变更历史追踪** - 记录任务的每一次变更，支持随时回溯任务状态
- **可视化历史展示** - 以图形化方式展示任务变更历史，类似 Git 提交历史
- **Web 图形化界面** - 直观的用户界面，便于操作和查看任务状态
- **实时同步** - 多设备间任务状态实时同步
- **数据统计** - 提供任务完成情况的统计和分析

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

### 前端
- React 19
- TypeScript
- Redux Toolkit
- React Router v7
- Bootstrap 5
- Bootstrap Icons
- i18next (多语言支持)
- Axios (HTTP客户端)

### 后端
- Node.js
- Express.js
- MongoDB
- Mongoose
- JWT 认证
- bcryptjs (密码加密)
- express-validator (数据验证)

## 快速开始

### 方法一：使用Docker（推荐）

#### 生产环境部署
```bash
# 克隆项目
git clone <repository-url>
cd todoIng

# 构建并启动所有服务
docker-compose up -d

# 应用将在以下地址可用：
# 前端: http://localhost
# 后端API: http://localhost:5001
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
# 核心模块开发计划

## 1. 概述

本文档详细描述了 todoIng 项目核心模块的开发计划，包括用户认证、任务管理和基础架构。这些模块是整个系统的基础，需要优先开发和稳定实现。

## 2. 功能边界

核心模块包含以下功能：

1. 用户认证系统（注册、登录、JWT 认证）
2. 基础任务管理（创建、读取、更新、删除任务）
3. 系统基础架构（项目结构、配置管理、错误处理）

这些功能是系统运行的基础，其他所有模块都依赖于这些核心功能。

## 3. 技术实现

### 3.1 后端技术栈
- Node.js 运行时环境
- Express.js Web 框架
- MongoDB 数据库
- Mongoose ODM
- JWT 用于身份验证
- bcryptjs 用于密码加密

### 3.2 前端技术栈
- React 18 前端框架
- TypeScript 类型系统
- Redux Toolkit 状态管理
- React Router 路由管理
- Axios HTTP 客户端

## 4. 开发任务分解

### 4.1 第一阶段：项目初始化和基础架构 (Week 1, Day 1-2)

#### 后端任务
1. 项目结构初始化
   - 创建项目目录结构
   - 初始化 package.json
   - 安装核心依赖包
   - 配置 ESLint 和 Prettier

2. 基础架构搭建
   - 配置 Express.js 应用
   - 实现基础路由
   - 配置环境变量管理
   - 实现日志记录系统

3. 数据库配置
   - 配置 MongoDB 连接
   - 实现连接池管理
   - 实现连接错误处理

#### 前端任务
1. 项目初始化
   - 使用 Vite 或 Create React App 初始化项目
   - 配置 TypeScript
   - 安装核心依赖包
   - 配置 ESLint 和 Prettier

2. 基础架构搭建
   - 创建目录结构
   - 配置 React Router
   - 实现基础页面路由
   - 配置全局样式

### 4.2 第二阶段：用户认证系统 (Week 1, Day 3-5)

#### 后端任务
1. 用户模型实现
   - 设计 User Schema
   - 实现用户注册验证逻辑
   - 实现密码加密存储
   - 实现用户数据访问层

2. 认证接口开发
   - 实现用户注册接口 (/api/auth/register)
   - 实现用户登录接口 (/api/auth/login)
   - 实现 JWT Token 生成和验证
   - 实现用户信息获取接口 (/api/auth/me)

3. 认证中间件
   - 实现 JWT 验证中间件
   - 实现权限检查中间件
   - 实现错误处理中间件

#### 前端任务
1. 认证页面开发
   - 实现登录页面组件
   - 实现注册页面组件
   - 实现表单验证
   - 实现错误提示

2. 认证状态管理
   - 实现 Redux 认证切片
   - 实现用户状态持久化
   - 实现认证路由保护
   - 实现自动登出功能

### 4.3 第三阶段：基础任务管理 (Week 2)

#### 后端任务
1. 任务模型实现
   - 设计 Task Schema
   - 实现任务状态枚举
   - 实现任务优先级枚举
   - 实现任务数据访问层

2. 任务管理接口
   - 实现创建任务接口 (/api/tasks)
   - 实现获取任务列表接口 (/api/tasks)
   - 实现获取任务详情接口 (/api/tasks/:id)
   - 实现更新任务接口 (/api/tasks/:id)
   - 实现删除任务接口 (/api/tasks/:id)

3. 任务验证和权限
   - 实现任务数据验证
   - 实现任务所有权验证
   - 实现任务状态转换规则

#### 前端任务
1. 任务管理界面
   - 实现任务列表页面
   - 实现任务详情页面
   - 实现任务创建/编辑表单
   - 实现任务状态变更组件

2. 任务状态管理
   - 实现 Redux 任务切片
   - 实现任务数据获取和缓存
   - 实现任务操作反馈
   - 实现任务搜索和筛选

## 5. 模块交互设计

### 5.1 用户认证与任务管理交互
```
用户认证模块 ←→ 任务管理模块
     ↓              ↓
   JWT Token    用户ID关联
```

用户在访问任务管理功能时，必须通过认证模块验证身份，任务数据将与用户ID关联。

### 5.2 前后端交互
```
前端 ←→ 后端 API ←→ 数据库
认证 ←→ 任务管理
```

所有数据交互通过 RESTful API 进行，确保前后端解耦。

## 6. 数据模型

### 6.1 用户模型 (User)
```javascript
{
  _id: ObjectId,
  username: String,
  email: String,
  password: String,  // 加密存储
  createdAt: Date,
  updatedAt: Date
}
```

### 6.2 任务模型 (Task)
```javascript
{
  _id: ObjectId,
  title: String,
  description: String,
  status: String,     // created, in-progress, paused, completed, cancelled
  priority: String,   // low, medium, high
  assignee: ObjectId, // reference to User
  createdBy: ObjectId, // reference to User
  createdAt: Date,
  updatedAt: Date
}
```

## 7. API 接口设计

### 7.1 认证相关接口
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `GET /api/auth/me` - 获取当前用户信息

### 7.2 任务管理接口
- `GET /api/tasks` - 获取任务列表
- `POST /api/tasks` - 创建新任务
- `GET /api/tasks/:id` - 获取任务详情
- `PUT /api/tasks/:id` - 更新任务
- `DELETE /api/tasks/:id` - 删除任务

## 8. 错误处理

### 8.1 后端错误处理
- 统一错误响应格式
- HTTP 状态码规范使用
- 详细错误日志记录
- 错误信息国际化支持

### 8.2 前端错误处理
- 全局错误拦截
- 用户友好的错误提示
- 错误重试机制
- 优雅降级处理

## 9. 测试计划

### 9.1 单元测试
- 用户模型测试
- 任务模型测试
- 认证服务测试
- 任务服务测试

### 9.2 集成测试
- 认证接口测试
- 任务管理接口测试
- 数据验证测试
- 权限控制测试

## 10. 部署考虑

### 10.1 环境配置
- 开发环境配置
- 测试环境配置
- 生产环境配置

### 10.2 安全考虑
- 环境变量管理
- 密码加密存储
- JWT Token 安全
- CORS 配置

## 11. 后续扩展

核心模块为后续功能模块提供基础支撑：

1. **历史追踪模块** - 基于用户和任务模型扩展
2. **团队管理模块** - 扩展用户关系和任务分配
3. **统计分析模块** - 基于任务数据进行聚合分析
4. **移动端支持** - 复用后端 API 接口

核心模块的稳定性和可扩展性直接影响整个系统的质量。
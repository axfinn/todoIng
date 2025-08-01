# 技术设计

## 系统架构

todoIng 采用前后端分离的架构模式，支持任务生命周期追踪和变更历史记录：

```
┌─────────────────┐    HTTP API    ┌──────────────────┐
│   Web Browser   │◄──────────────►│  Node.js Server  │
└─────────────────┘                └──────────────────┘
                                            │
                                     MongoDB│
                                            ▼
                                  ┌──────────────────┐
                                  │     Database     │
                                  └──────────────────┘
```

## 核心概念

### 任务生命周期
每个任务都有完整的生命周期，包括：
- 创建 (Created)
- 进行中 (In Progress)
- 暂停 (Paused)
- 完成 (Completed)
- 取消 (Cancelled)

### 变更追踪
系统将追踪任务的所有变更，包括：
- 状态变更
- 内容更新
- 分配变更
- 优先级调整

## 扩展功能架构

### 用户模块扩展
- 用户个人资料管理
- 用户偏好设置
- 用户成就系统
- 用户活跃度统计

### 回溯模块
- 任务状态快照功能
- 时间点任务状态恢复
- 任务历史对比功能

### 生命历程
- 用户所有任务的总时间线视图
- 任务完成里程碑展示
- 个人成长轨迹分析

### 事项总结
- 任务完成情况统计
- 个人/团队绩效分析
- 任务模式识别和建议

### 团队管理
- 团队和组织结构管理
- 团队任务分配和跟踪
- 团队绩效统计和分析

## 前端技术设计

### 技术栈
- React 18.x
- TypeScript
- React Router v6
- Redux Toolkit 状态管理
- Axios HTTP 客户端
- D3.js 或 Chart.js 用于历史数据可视化

### 项目结构
```
src/
├── components/        # 可复用组件
├── pages/             # 页面组件
├── store/             # Redux 状态管理
├── services/          # API 服务层
├── utils/             # 工具函数
├── hooks/             # 自定义 hooks
├── styles/            # 全局样式
├── types/             # TypeScript 类型定义
└── App.tsx            # 根组件
```

### 核心组件设计

#### 1. 任务列表组件
展示所有任务及其当前状态

#### 2. 任务详情组件
展示任务详细信息和变更历史时间线

#### 3. 任务编辑组件
用于创建和编辑任务

#### 4. 历史时间线组件
以图形化方式展示任务变更历史，类似Git提交历史

#### 5. 状态统计面板
展示任务状态分布和统计信息

### 扩展组件设计

#### 1. 团队管理组件
- 团队列表和详情展示
- 团队成员管理
- 团队任务视图

#### 2. 用户个人资料组件
- 用户信息展示和编辑
- 个人成就展示
- 个人任务统计

#### 3. 生命历程时间线组件
- 用户所有任务的时间线视图
- 里程碑标记和展示
- 成长轨迹可视化

#### 4. 数据分析和总结组件
- 任务完成情况图表
- 性能统计面板
- 个性化建议展示

### 状态管理
使用 Redux Toolkit 进行全局状态管理，主要管理以下状态：
- 用户认证状态
- 任务列表数据
- 任务详情及变更历史
- 团队信息
- UI 状态（加载状态、错误信息等）

### 路由设计
- `/` - 主页/仪表板
- `/tasks` - 任务列表
- `/tasks/:id` - 任务详情和历史
- `/tasks/new` - 创建新任务
- `/tasks/:id/edit` - 编辑任务
- `/teams` - 团队列表
- `/teams/:id` - 团队详情
- `/profile` - 用户个人资料
- `/timeline` - 生命历程时间线
- `/analytics` - 数据分析和可视化
- `/login` - 登录页
- `/register` - 注册页

## 后端技术设计

### 技术栈
- Node.js
- Express.js
- MongoDB + Mongoose
- JWT 认证
- bcryptjs 密码加密

### 项目结构
```
src/
├── controllers/       # 控制器层
├── models/            # 数据模型
├── routes/            # 路由定义
├── middleware/        # 中间件
├── utils/             # 工具函数
├── config/            # 配置文件
└── app.js             # 应用入口
```

### API 设计规范
- RESTful API 风格
- JSON 数据格式
- 统一的响应格式
- JWT Token 认证

### 数据模型设计

#### 用户模型 (User)
- `_id`: ObjectId
- `username`: String
- `email`: String
- `password`: String (hashed)
- `profile`: Object (头像、简介等个人信息)
- `preferences`: Object (用户偏好设置)
- `achievements`: Array (用户成就)
- `createdAt`: Date
- `updatedAt`: Date

#### 团队模型 (Team)
- `_id`: ObjectId
- `name`: String
- `description`: String
- `members`: Array of Objects (成员及角色)
- `createdAt`: Date
- `updatedAt`: Date

#### 任务模型 (Task)
- `_id`: ObjectId
- `title`: String
- `description`: String
- `status`: Enum (created, in-progress, paused, completed, cancelled)
- `priority`: Enum (low, medium, high)
- `assignee`: ObjectId (reference to User)
- `team`: ObjectId (reference to Team, 可选)
- `tags`: [String]
- `dueDate`: Date
- `createdAt`: Date
- `updatedAt`: Date
- `createdBy`: ObjectId (reference to User)

#### 任务历史记录模型 (TaskHistory)
- `_id`: ObjectId
- `taskId`: ObjectId (reference to Task)
- `field`: String (变更的字段)
- `oldValue`: Mixed (变更前的值)
- `newValue`: Mixed (变更后的值)
- `changedBy`: ObjectId (reference to User)
- `changeType`: String (status-change, update, assign, etc.)
- `timestamp`: Date
- `comment`: String (可选的变更说明)

#### 用户统计模型 (UserStatistics)
- `_id`: ObjectId (reference to User)
- `tasksCompleted`: Number
- `tasksCreated`: Number
- `completionRate`: Number
- `activeDays`: Array of Dates
- `lastActive`: Date
- `achievements`: Array

#### 任务总结模型 (TaskSummary)
- `_id`: ObjectId
- `userId`: ObjectId (reference to User)
- `period`: String (统计周期: daily, weekly, monthly)
- `tasksCompleted`: Number
- `tasksCreated`: Number
- `averageCompletionTime`: Number
- `insights`: Array (系统生成的洞察)
- `generatedAt`: Date

### 核心API端点

#### 认证相关
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录

#### 任务相关
- `GET /api/tasks` - 获取任务列表
- `POST /api/tasks` - 创建新任务
- `GET /api/tasks/:id` - 获取任务详情
- `PUT /api/tasks/:id` - 更新任务
- `DELETE /api/tasks/:id` - 删除任务
- `GET /api/tasks/:id/history` - 获取任务变更历史

#### 团队相关
- `GET /api/teams` - 获取团队列表
- `POST /api/teams` - 创建新团队
- `GET /api/teams/:id` - 获取团队详情
- `PUT /api/teams/:id` - 更新团队信息
- `DELETE /api/teams/:id` - 删除团队
- `POST /api/teams/:id/members` - 添加团队成员
- `DELETE /api/teams/:id/members/:userId` - 移除团队成员

#### 用户相关
- `GET /api/users/:id` - 获取用户详情
- `PUT /api/users/:id` - 更新用户信息
- `GET /api/users/:id/statistics` - 获取用户统计信息

#### 历史记录相关
- `GET /api/history` - 获取所有历史记录
- `GET /api/history/:id` - 获取特定历史记录详情
- `GET /api/history/timeline` - 获取历史时间线

#### 统计和总结相关
- `GET /api/analytics/tasks` - 获取任务分析数据
- `GET /api/summaries/tasks` - 获取任务总结
- `GET /api/users/:id/summaries` - 获取用户任务总结

## 数据库设计

### MongoDB 设计考虑
1. 使用嵌套结构存储任务和历史记录以提高查询性能
2. 为常用查询字段创建索引（status, assignee, dueDate等）
3. 实现数据分页以处理大量任务和历史记录
4. 为团队和用户扩展信息设计合理的数据结构

## 安全设计
- JWT Token 认证
- 密码加密存储
- 请求验证和清理
- CORS 配置
- 速率限制
- 团队权限控制

## 性能优化
- 数据库索引优化
- API 响应缓存
- 前端数据缓存
- 分页加载大数据集
- 团队和用户数据的懒加载
# 项目实现技术方案

## 1. 概述

todoIng 是一个创新的任务管理系统，它不仅提供基本的任务管理功能，还为每个任务提供完整的生命周期追踪和变更历史记录，就像 Git 管理代码变更一样管理任务。本文档详细描述了项目的实现技术方案。

## 2. 技术栈选择

### 2.1 前端技术栈
- **React 18** - 用于构建用户界面的 JavaScript 库
- **TypeScript** - 为 JavaScript 添加静态类型检查
- **Redux Toolkit** - 状态管理解决方案
- **React Router v6** - 声明式路由管理
- **Axios** - HTTP 客户端
- **D3.js** - 数据可视化库，用于历史数据展示
- **Styled Components** - CSS-in-JS 解决方案
- **Formik + Yup** - 表单处理和验证

### 2.2 后端技术栈
- **Node.js** - JavaScript 运行时环境
- **Express.js** - Web 应用框架
- **MongoDB** - NoSQL 数据库
- **Mongoose** - MongoDB 对象建模工具
- **JWT** - JSON Web Token 认证
- **Socket.IO** - 实时通信支持
- **Jest** - JavaScript 测试框架

### 2.3 开发工具
- **Webpack** - 模块打包工具
- **Babel** - JavaScript 编译器
- **ESLint** - 代码静态分析工具
- **Prettier** - 代码格式化工具
- **Docker** - 容器化部署
- **GitHub Actions** - CI/CD 工具

## 3. 系统架构设计

### 3.1 整体架构
```
┌─────────────────┐    HTTP API    ┌──────────────────┐
│   Web Browser   │◄──────────────►│  Node.js Server  │
└─────────────────┘    WebSocket   └──────────────────┘
                                            │
                                     MongoDB│
                                            ▼
                                  ┌──────────────────┐
                                  │     Database     │
                                  └──────────────────┘
```

### 3.2 前端架构
```
src/
├── components/        # 可复用UI组件
├── pages/             # 页面组件
├── store/             # Redux状态管理
├── services/          # API服务层
├── utils/             # 工具函数
├── hooks/             # 自定义hooks
├── styles/            # 全局样式
├── types/             # TypeScript类型定义
├── assets/            # 静态资源
└── App.tsx            # 根组件
```

### 3.3 后端架构
```
src/
├── controllers/       # 控制器层
├── models/            # 数据模型
├── routes/            # 路由定义
├── middleware/        # 中间件
├── services/          # 业务逻辑层
├── utils/             # 工具函数
├── config/            # 配置文件
├── validators/        # 数据验证
└── app.js             # 应用入口
```

## 4. 核心功能实现方案

### 4.1 任务生命周期管理

#### 4.1.1 状态管理
任务状态包括：创建(Created)、进行中(In Progress)、暂停(Paused)、完成(Completed)、取消(Cancelled)。

```javascript
// 任务状态枚举
const TaskStatus = {
  CREATED: 'created',
  IN_PROGRESS: 'in-progress',
  PAUSED: 'paused',
  COMPLETED: 'completed',
  CANCELLED: 'cancelled'
};

// 状态转换规则
const statusTransitions = {
  [TaskStatus.CREATED]: [TaskStatus.IN_PROGRESS, TaskStatus.CANCELLED],
  [TaskStatus.IN_PROGRESS]: [TaskStatus.PAUSED, TaskStatus.COMPLETED, TaskStatus.CANCELLED],
  [TaskStatus.PAUSED]: [TaskStatus.IN_PROGRESS, TaskStatus.CANCELLED],
  [TaskStatus.COMPLETED]: [],
  [TaskStatus.CANCELLED]: []
};
```

#### 4.1.2 状态变更追踪
每次状态变更都会在 [task_histories](file:///Volumes/M20/code/docs/todoIng/docs/database-design.md) 集合中创建记录：

```javascript
const createStatusChangeHistory = async (taskId, oldStatus, newStatus, userId) => {
  return await TaskHistory.create({
    taskId,
    field: 'status',
    oldValue: oldStatus,
    newValue: newStatus,
    changedBy: userId,
    changeType: 'status-change',
    timestamp: new Date()
  });
};
```

### 4.2 变更历史追踪

#### 4.2.1 历史记录创建
在每次任务更新时自动创建历史记录：

```javascript
// 任务更新中间件
const trackTaskChanges = async (req, res, next) => {
  const taskId = req.params.id;
  const updates = req.body;
  const userId = req.user._id;
  
  // 获取任务更新前的状态
  const oldTask = await Task.findById(taskId);
  
  // 继续处理更新请求
  next();
  
  // 请求处理完成后创建历史记录
  setImmediate(async () => {
    const newTask = await Task.findById(taskId);
    for (const [key, value] of Object.entries(updates)) {
      if (oldTask[key] !== newTask[key]) {
        await TaskHistory.create({
          taskId,
          field: key,
          oldValue: oldTask[key],
          newValue: newTask[key],
          changedBy: userId,
          changeType: 'update',
          timestamp: new Date()
        });
      }
    }
  });
};
```

#### 4.2.2 任务快照功能
定期创建任务快照以支持快速恢复：

```javascript
const createTaskSnapshot = async (taskId) => {
  const task = await Task.findById(taskId).populate('assignee createdBy');
  return await TaskHistory.create({
    taskId,
    changedBy: task.createdBy,
    changeType: 'snapshot',
    timestamp: new Date(),
    snapshot: task
  });
};
```

### 4.3 团队管理功能

#### 4.3.1 团队权限控制
基于角色的访问控制(RBAC)：

```javascript
const TeamRoles = {
  OWNER: 'owner',
  ADMIN: 'admin',
  MEMBER: 'member'
};

const checkTeamPermission = (team, userId, requiredRole) => {
  const member = team.members.find(m => m.user.toString() === userId.toString());
  if (!member) return false;
  
  const roleHierarchy = {
    [TeamRoles.MEMBER]: 1,
    [TeamRoles.ADMIN]: 2,
    [TeamRoles.OWNER]: 3
  };
  
  return roleHierarchy[member.role] >= roleHierarchy[requiredRole];
};
```

### 4.4 系统配置管理

#### 4.4.1 注册控制
系统支持通过环境变量禁用用户注册功能，适用于私有部署场景。

```javascript
// 检查是否禁用注册
const isRegistrationDisabled = () => {
  return process.env.DISABLE_REGISTRATION === 'true';
};

// 在注册路由中应用控制
router.post('/register', (req, res) => {
  if (isRegistrationDisabled()) {
    return res.status(403).json({ msg: 'Registration is disabled' });
  }
  // 正常注册逻辑
});
```

#### 4.4.2 默认用户
系统支持配置默认用户，在应用启动时自动创建，方便私有部署场景的初始访问。

```javascript
// 创建默认用户
const createDefaultUser = async () => {
  // 检查是否配置了默认用户
  if (process.env.DEFAULT_USERNAME && process.env.DEFAULT_PASSWORD) {
    // 检查用户是否已存在
    const existingUser = await User.findOne({
      $or: [
        { username: process.env.DEFAULT_USERNAME },
        { email: process.env.DEFAULT_EMAIL }
      ]
    });
    
    // 如果不存在则创建
    if (!existingUser) {
      const user = new User({
        username: process.env.DEFAULT_USERNAME,
        email: process.env.DEFAULT_EMAIL,
        password: process.env.DEFAULT_PASSWORD
      });
      
      const salt = await bcrypt.genSalt(10);
      user.password = await bcrypt.hash(user.password, salt);
      await user.save();
      console.log('Default user created successfully');
    }
  }
};
```

#### 4.4.3 登录验证码
系统支持可选的登录验证码功能，增强系统安全性。

```javascript
// 验证码检查中间件
const checkCaptcha = (req, res, next) => {
  // 检查是否启用了验证码功能
  if (process.env.ENABLE_CAPTCHA !== 'true') {
    return next();
  }
  
  const { captcha } = req.body;
  // 验证验证码逻辑
  if (!captcha || !isValidCaptcha(captcha)) {
    return res.status(400).json({ msg: 'Invalid captcha' });
  }
  
  next();
};

// 在登录路由中应用
router.post('/login', checkCaptcha, (req, res) => {
  // 正常登录逻辑
});

// 生成验证码路由
router.get('/captcha', (req, res) => {
  // 生成验证码逻辑
  const captcha = generateCaptcha();
  // 保存验证码到会话或缓存中
  req.session.captcha = captcha.text;
  // 返回验证码图像
  res.json({ image: captcha.image });
});
```

## 5. 数据库设计实现

### 5.1 MongoDB 模式定义

#### 5.1.1 用户模型
```
const userSchema = new mongoose.Schema({
  username: { type: String, required: true, unique: true },
  email: { type: String, required: true, unique: true },
  password: { type: String, required: true },
  profile: {
    avatar: String,
    firstName: String,
    lastName: String,
    bio: String,
    skills: [String],
    location: String
  },
  preferences: {
    theme: { type: String, default: 'light' },
    notifications: {
      email: { type: Boolean, default: true },
      push: { type: Boolean, default: true }
    },
    language: { type: String, default: 'zh-CN' }
  },
  achievements: [{
    id: String,
    name: String,
    description: String,
    earnedAt: Date
  }]
}, { timestamps: true });
```

#### 5.1.2 任务模型
```
const taskSchema = new mongoose.Schema({
  title: { type: String, required: true },
  description: String,
  status: { 
    type: String, 
    enum: ['created', 'in-progress', 'paused', 'completed', 'cancelled'],
    default: 'created'
  },
  priority: { 
    type: String, 
    enum: ['low', 'medium', 'high'],
    default: 'medium'
  },
  assignee: { type: mongoose.Schema.Types.ObjectId, ref: 'User' },
  team: { type: mongoose.Schema.Types.ObjectId, ref: 'Team' },
  createdBy: { type: mongoose.Schema.Types.ObjectId, ref: 'User', required: true },
  tags: [String],
  dueDate: Date,
  startDate: Date,
  completedDate: Date,
  estimatedHours: Number,
  actualHours: Number,
  attachments: [{
    name: String,
    url: String,
    size: Number,
    type: String
  }]
}, { timestamps: true });
```

#### 5.1.3 团队模型
```
const teamSchema = new mongoose.Schema({
  name: { type: String, required: true },
  description: String,
  avatar: String,
  members: [{
    user: { type: mongoose.Schema.Types.ObjectId, ref: 'User', required: true },
    role: { 
      type: String, 
      enum: ['owner', 'admin', 'member'],
      default: 'member'
    },
    joinedAt: { type: Date, default: Date.now }
  }]
}, { timestamps: true });
```

## 6. API 实现方案

### 6.1 RESTful API 设计
遵循 RESTful 设计原则，使用标准 HTTP 方法：

```javascript
// 任务相关路由
app.get('/api/tasks', taskController.getTasks);
app.post('/api/tasks', taskController.createTask);
app.get('/api/tasks/:id', taskController.getTask);
app.put('/api/tasks/:id', taskController.updateTask);
app.delete('/api/tasks/:id', taskController.deleteTask);
app.get('/api/tasks/:id/history', taskController.getTaskHistory);

// 团队相关路由
app.get('/api/teams', teamController.getTeams);
app.post('/api/teams', teamController.createTeam);
app.get('/api/teams/:id', teamController.getTeam);
app.put('/api/teams/:id', teamController.updateTeam);
app.delete('/api/teams/:id', teamController.deleteTeam);
app.post('/api/teams/:id/members', teamController.addMember);
app.delete('/api/teams/:id/members/:userId', teamController.removeMember);
```

### 6.2 WebSocket 实现实时更新
使用 Socket.IO 实现实时通信：

```javascript
// 服务端
io.on('connection', (socket) => {
  // 加入房间
  socket.on('join', (room) => {
    socket.join(room);
  });
  
  // 任务更新事件
  socket.on('task.updated', (data) => {
    io.to(`task-${data.taskId}`).emit('task.updated', data);
  });
});

// 客户端
const socket = io();
socket.emit('join', `task-${taskId}`);

socket.on('task.updated', (data) => {
  // 更新本地状态
  dispatch(updateTask(data));
});
```

## 7. 前端实现方案

### 7.1 状态管理 (Redux Toolkit)
```javascript
// 任务切片
const taskSlice = createSlice({
  name: 'tasks',
  initialState: {
    tasks: [],
    loading: false,
    error: null
  },
  reducers: {
    setTasks: (state, action) => {
      state.tasks = action.payload;
    },
    addTask: (state, action) => {
      state.tasks.push(action.payload);
    },
    updateTask: (state, action) => {
      const index = state.tasks.findIndex(t => t._id === action.payload._id);
      if (index !== -1) {
        state.tasks[index] = action.payload;
      }
    }
  }
});
```

### 7.2 组件设计
```typescript
// 任务卡片组件
interface TaskCardProps {
  task: Task;
  onStatusChange: (taskId: string, status: TaskStatus) => void;
  onEdit: (taskId: string) => void;
  onDelete: (taskId: string) => void;
}

const TaskCard: React.FC<TaskCardProps> = ({ task, onStatusChange, onEdit, onDelete }) => {
  return (
    <div className="task-card">
      <h3>{task.title}</h3>
      <p>{task.description}</p>
      <div className="task-meta">
        <span className={`status ${task.status}`}>{task.status}</span>
        <span className={`priority ${task.priority}`}>{task.priority}</span>
      </div>
      <div className="task-actions">
        <button onClick={() => onEdit(task._id)}>编辑</button>
        <button onClick={() => onDelete(task._id)}>删除</button>
      </div>
    </div>
  );
};
```

### 7.3 历史时间线可视化
使用 D3.js 实现历史时间线：

```javascript
// 历史时间线组件
const HistoryTimeline: React.FC<{ history: TaskHistory[] }> = ({ history }) => {
  useEffect(() => {
    // 使用 D3.js 绘制时间线
    const svg = d3.select('#timeline')
      .append('svg')
      .attr('width', 800)
      .attr('height', 400);
    
    // 绘制时间线和节点
    // ...
  }, [history]);
  
  return <div id="timeline"></div>;
};
```

## 8. 安全实现方案

### 8.1 认证和授权
```javascript
// JWT 认证中间件
const authenticate = async (req, res, next) => {
  const token = req.header('Authorization')?.replace('Bearer ', '');
  
  if (!token) {
    return res.status(401).json({ error: 'Access denied' });
  }
  
  try {
    const decoded = jwt.verify(token, process.env.JWT_SECRET);
    req.user = await User.findById(decoded.userId);
    next();
  } catch (error) {
    res.status(401).json({ error: 'Invalid token' });
  }
};

// 权限检查中间件
const authorize = (roles = []) => {
  return (req, res, next) => {
    if (!req.user) {
      return res.status(401).json({ error: 'Access denied' });
    }
    
    if (roles.length && !roles.includes(req.user.role)) {
      return res.status(403).json({ error: 'Insufficient permissions' });
    }
    
    next();
  };
};
```

### 8.2 数据验证
```javascript
// 使用 Joi 进行数据验证
const taskValidationSchema = Joi.object({
  title: Joi.string().min(1).max(100).required(),
  description: Joi.string().max(1000),
  status: Joi.string().valid('created', 'in-progress', 'paused', 'completed', 'cancelled'),
  priority: Joi.string().valid('low', 'medium', 'high'),
  assignee: Joi.string().pattern(/^[0-9a-fA-F]{24}$/),
  dueDate: Joi.date().greater('now')
});
```

## 9. 性能优化方案

### 9.1 数据库优化
1. 创建合适的索引
2. 使用聚合管道进行复杂查询
3. 实施分页加载大数据集

```javascript
// 分页查询任务
const getTasks = async (query, page = 1, limit = 10) => {
  const skip = (page - 1) * limit;
  
  const tasks = await Task.find(query)
    .skip(skip)
    .limit(limit)
    .sort({ createdAt: -1 });
    
  const total = await Task.countDocuments(query);
  
  return {
    tasks,
    pagination: {
      page,
      limit,
      total,
      pages: Math.ceil(total / limit)
    }
  };
};
```

### 9.2 前端优化
1. 实施虚拟滚动处理大量数据
2. 使用 React.memo 优化组件渲染
3. 实施代码分割和懒加载

```typescript
// 虚拟滚动实现
const VirtualizedTaskList: React.FC<{ tasks: Task[] }> = ({ tasks }) => {
  const rowHeight = 100;
  const visibleCount = Math.ceil(window.innerHeight / rowHeight);
  const [scrollTop, setScrollTop] = useState(0);
  
  const startIndex = Math.floor(scrollTop / rowHeight);
  const endIndex = Math.min(startIndex + visibleCount, tasks.length);
  
  const visibleTasks = tasks.slice(startIndex, endIndex);
  
  return (
    <div 
      className="task-list"
      onScroll={(e) => setScrollTop(e.currentTarget.scrollTop)}
    >
      <div style={{ height: tasks.length * rowHeight }}>
        <div style={{ transform: `translateY(${startIndex * rowHeight}px)` }}>
          {visibleTasks.map(task => (
            <TaskCard key={task._id} task={task} />
          ))}
        </div>
      </div>
    </div>
  );
};
```

## 10. 部署方案

### 10.1 Docker 化部署
```
# Dockerfile
FROM node:18-alpine

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

EXPOSE 3000

CMD ["npm", "start"]
```

### 10.2 Nginx 反向代理配置
```
server {
    listen 80;
    server_name todoing.example.com;
    
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
    
    location /ws {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

## 11. 测试方案

### 11.1 单元测试
```javascript
// 任务服务测试
describe('Task Service', () => {
  it('should create a new task', async () => {
    const taskData = {
      title: 'Test Task',
      description: 'Test Description',
      createdBy: 'user123'
    };
    
    const task = await taskService.createTask(taskData);
    
    expect(task.title).toBe('Test Task');
    expect(task.status).toBe('created');
  });
  
  it('should not allow invalid status transition', async () => {
    const task = await Task.create({
      title: 'Test Task',
      status: 'completed',
      createdBy: 'user123'
    });
    
    await expect(taskService.updateTask(task._id, { status: 'in-progress' }))
      .rejects
      .toThrow('Invalid status transition');
  });
});
```

### 11.2 集成测试
```javascript
// API 集成测试
describe('Task API', () => {
  it('should get tasks with pagination', async () => {
    // 创建测试数据
    await Task.create([
      { title: 'Task 1', createdBy: 'user123' },
      { title: 'Task 2', createdBy: 'user123' },
      { title: 'Task 3', createdBy: 'user123' }
    ]);
    
    const response = await request(app)
      .get('/api/tasks')
      .query({ page: 1, limit: 2 });
    
    expect(response.status).toBe(200);
    expect(response.body.data.tasks).toHaveLength(2);
    expect(response.body.data.pagination.total).toBe(3);
  });
});
```

## 12. 监控和日志

### 12.1 日志记录
```javascript
// 使用 Winston 记录日志
const logger = winston.createLogger({
  level: 'info',
  format: winston.format.combine(
    winston.format.timestamp(),
    winston.format.json()
  ),
  transports: [
    new winston.transports.File({ filename: 'error.log', level: 'error' }),
    new winston.transports.File({ filename: 'combined.log' })
  ]
});
```

### 12.2 性能监控
```javascript
// API 响应时间监控中间件
const responseTimeMiddleware = (req, res, next) => {
  const start = Date.now();
  
  res.on('finish', () => {
    const duration = Date.now() - start;
    logger.info(`${req.method} ${req.path} - ${duration}ms`);
    
    // 如果响应时间过长，发送告警
    if (duration > 5000) {
      logger.warn(`Slow request: ${req.method} ${req.path} - ${duration}ms`);
    }
  });
  
  next();
};
```

## 13. 后续扩展计划

### 13.1 第一阶段：核心功能完善
- 完善任务生命周期管理
- 优化历史记录追踪功能
- 实现基础团队管理功能

### 13.2 第二阶段：高级功能开发
- 实现生命历程时间线
- 开发数据分析和总结功能
- 增强用户个人资料系统

### 13.3 第三阶段：性能和扩展
- 实施更高级的性能优化
- 增加移动端支持
- 实现离线功能

通过以上技术方案的实施，todoIng 将成为一个功能完整、性能优秀、易于扩展的任务管理系统，能够满足个人和团队的任务管理需求，并提供独特的任务生命周期追踪功能。
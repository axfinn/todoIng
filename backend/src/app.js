const express = require('express');
const cors = require('cors');
require('dotenv').config();

const connectDB = require('./config/db');
const authRoutes = require('./routes/auth');
const taskRoutes = require('./routes/tasks');
const authMiddleware = require('./middleware/auth');
const errorMiddleware = require('./middleware/error');
const loggerMiddleware = require('./middleware/logger');

const app = express();

// 连接数据库
connectDB();

// 中间件
app.use(loggerMiddleware);
app.use(cors());
app.use(express.json());

// 路由
app.use('/api/auth', authRoutes);
app.use('/api/tasks', authMiddleware, taskRoutes);

// 错误处理中间件（必须在路由之后注册）
app.use(errorMiddleware);

// 获取端口配置，默认为5000
const PORT = process.env.PORT || 5000;

app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});
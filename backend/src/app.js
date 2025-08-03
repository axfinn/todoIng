const express = require('express');
const connectDB = require('./config/db');
const dotenv = require('dotenv');
const path = require('path');
const User = require('./models/User');
const bcrypt = require('bcryptjs');
const cors = require('cors');

// Load env vars
dotenv.config();

const app = express();

// Connect Database
connectDB();

// 初始化默认用户
const initializeDefaultUser = async () => {
  try {
    // 检查是否配置了默认用户
    if (process.env.DEFAULT_USERNAME && process.env.DEFAULT_PASSWORD && process.env.DEFAULT_EMAIL) {
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
  } catch (err) {
    console.error('Error initializing default user:', err.message);
  }
};

// 调用初始化函数
initializeDefaultUser();

const authRoutes = require('./routes/auth');
const taskRoutes = require('./routes/tasks');
const authMiddleware = require('./middleware/auth');
const errorMiddleware = require('./middleware/error');
const loggerMiddleware = require('./middleware/logger');

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
const express = require('express');
const cors = require('cors');
const dotenv = require('dotenv');
const connectDB = require('./config/db');
const authRoutes = require('./routes/auth');
const taskRoutes = require('./routes/tasks');
const errorHandler = require('./middleware/error');
const User = require('./models/User');
const bcrypt = require('bcryptjs');

// Load env vars
dotenv.config({ path: './.env' });

// Connect to database
connectDB();

// 创建默认用户
const createDefaultUser = async () => {
  // 检查是否配置了默认用户
  if (process.env.DEFAULT_USERNAME && process.env.DEFAULT_PASSWORD) {
    try {
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
      } else {
        console.log('Default user already exists');
      }
    } catch (err) {
      console.error('Error creating default user:', err.message);
    }
  }
};

const app = express();

// Body parser
app.use(express.json());

// Enable cors
app.use(cors());

// Routes
app.use('/api/auth', authRoutes);
app.use('/api/tasks', taskRoutes);

// Error handler
app.use(errorHandler);

// Validate and use port from environment or default to 5001
const PORT = process.env.PORT ? parseInt(process.env.PORT, 10) : 5001;

const server = app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
  // 创建默认用户
  createDefaultUser();
});

// Handle unhandled promise rejections
process.on('unhandledRejection', (err, promise) => {
  console.log(`Error: ${err.message}`);
  // Close server & exit process
  server.close(() => {
    process.exit(1);
  });
});
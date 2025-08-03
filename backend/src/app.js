const express = require('express');
const cors = require('cors');
const dotenv = require('dotenv');
const connectDB = require('./config/db');
const authRoutes = require('./routes/auth');
const taskRoutes = require('./routes/tasks');
const errorHandler = require('./middleware/error');

// Load env vars
dotenv.config({ path: './.env' });

// Connect to database
connectDB();

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

// Check if port is valid
if (isNaN(PORT) || PORT < 1024 || PORT > 65535) {
  console.error('Invalid port configuration. Using default port 5001');
  process.env.PORT = '5001';
}

app.listen(PORT, '0.0.0.0', () => {
  console.log(`Server running on port ${PORT}`);
  console.log(`Environment: ${process.env.NODE_ENV || 'development'}`);
});
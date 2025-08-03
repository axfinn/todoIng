const mongoose = require('mongoose');

const connectDB = async () => {
  try {
    // 移除已弃用的选项 useNewUrlParser 和 useUnifiedTopology
    const conn = await mongoose.connect(process.env.MONGO_URI);
    
    console.log(`MongoDB Connected: ${conn.connection.host}`);
  } catch (error) {
    console.error(`Error: ${error.message}`);
    process.exit(1);
  }
};

module.exports = connectDB;
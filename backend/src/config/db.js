const mongoose = require('mongoose');
require('dotenv').config();

const connectDB = async () => {
  try {
    // 优先使用环境变量中的MONGO_URI，否则使用默认值
    const conn = await mongoose.connect(process.env.MONGO_URI || 'mongodb://localhost:27017/todoing', {
      useNewUrlParser: true,
      useUnifiedTopology: true,
    });

    console.log(`MongoDB Connected: ${conn.connection.host}`);
  } catch (err) {
    console.error(`Error: ${err.message}`);
    process.exit(1);
  }
};

module.exports = connectDB;
const mongoose = require('mongoose');

const connectDB = async () => {
  try {
    // 尝试连接MongoDB，使用docker-compose中定义的服务名称
    const conn = await mongoose.connect(process.env.MONGO_URI, {
      useNewUrlParser: true,
      useUnifiedTopology: true,
    });

    console.log(`MongoDB Connected: ${conn.connection.host}`);
  } catch (err) {
    console.error(`Error: ${err.message}`);
    // 如果连接失败，5秒后重试
    setTimeout(connectDB, 5000);
  }
};

module.exports = connectDB;
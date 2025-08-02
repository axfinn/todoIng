const fs = require('fs');
const path = require('path');

// 创建日志目录（如果不存在）
const logDirectory = path.join(__dirname, '../logs');
if (!fs.existsSync(logDirectory)) {
  fs.mkdirSync(logDirectory, { recursive: true });
}

const logger = (req, res, next) => {
  const now = new Date();
  const timestamp = now.toISOString();
  const logEntry = `[${timestamp}] ${req.method} ${req.url} - IP: ${req.ip}\n`;
  
  // 输出到控制台
  console.log(logEntry.trim());
  
  // 写入日志文件
  const logFile = path.join(logDirectory, `${now.toISOString().split('T')[0]}.log`);
  fs.appendFileSync(logFile, logEntry);
  
  next();
};

module.exports = logger;
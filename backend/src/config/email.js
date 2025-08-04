const nodemailer = require('nodemailer');
const crypto = require('crypto');

// 存储邮箱验证码 (生产环境建议使用Redis)
const emailCodeStore = new Map();

// 邮箱验证码配置
const EMAIL_CODE_CONFIG = {
  length: 6, // 验证码长度
  lifetime: 10 * 60 * 1000, // 10分钟有效期
  maxAttempts: 3 // 最大尝试次数
};

// 生成邮箱验证码
function generateEmailCode(length = EMAIL_CODE_CONFIG.length) {
  return crypto.randomBytes(length).toString('hex').slice(0, length).toUpperCase();
}

// 生成邮箱验证码ID
function generateEmailCodeId() {
  return crypto.randomBytes(16).toString('hex');
}

// 清理过期的邮箱验证码
function cleanupExpiredEmailCodes() {
  const now = Date.now();
  for (const [id, code] of emailCodeStore.entries()) {
    if (!code.expiresAt || code.expiresAt < now) {
      emailCodeStore.delete(id);
    }
  }
}

// 启动定时器定期清理过期验证码
setInterval(cleanupExpiredEmailCodes, 60000); // 每分钟清理一次

// 创建邮件发送器
function createEmailTransporter() {
  // 从环境变量获取邮件配置
  const config = {
    host: process.env.EMAIL_HOST,
    port: process.env.EMAIL_PORT || 587,
    secure: process.env.EMAIL_SECURE === 'true', // true for 465, false for other ports
    auth: {
      user: process.env.EMAIL_USER,
      pass: process.env.EMAIL_PASS
    }
  };

  // 检查必要配置是否存在
  if (!config.host || !config.auth.user || !config.auth.pass) {
    console.warn('Email configuration not complete. Email verification will not work.');
    return null;
  }

  return nodemailer.createTransport(config);
}

// 发送验证码邮件
async function sendVerificationEmail(to, code) {
  const transporter = createEmailTransporter();
  
  if (!transporter) {
    throw new Error('Email transporter not configured');
  }

  const mailOptions = {
    from: process.env.EMAIL_FROM || process.env.EMAIL_USER,
    to,
    subject: 'TodoIng 邮箱验证码',
    html: `
      <div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #333;">TodoIng 邮箱验证码</h2>
        <p>您好！</p>
        <p>您正在注册 TodoIng 账号，您的邮箱验证码是：</p>
        <div style="background-color: #f5f5f5; padding: 15px; text-align: center; border-radius: 5px; margin: 20px 0;">
          <span style="font-size: 24px; font-weight: bold; color: #e67e22; letter-spacing: 5px;">${code}</span>
        </div>
        <p>验证码有效期为10分钟，请尽快完成注册。</p>
        <p style="color: #999; font-size: 12px;">如果您没有进行此操作，请忽略此邮件。</p>
      </div>
    `
  };

  try {
    await transporter.sendMail(mailOptions);
    return true;
  } catch (error) {
    console.error('Failed to send verification email:', error);
    throw new Error('Failed to send verification email');
  }
}

module.exports = {
  emailCodeStore,
  EMAIL_CODE_CONFIG,
  generateEmailCode,
  generateEmailCodeId,
  cleanupExpiredEmailCodes,
  sendVerificationEmail
};
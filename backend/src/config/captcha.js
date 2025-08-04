const crypto = require('crypto');

// 存储验证码 (生产环境建议使用Redis)
const captchaStore = new Map();

// 验证码配置
const CAPTCHA_CONFIG = {
  length: 6, // 验证码长度
  lifetime: 5 * 60 * 1000, // 5分钟有效期
};

// 生成验证码文本
function generateCaptchaText(length = CAPTCHA_CONFIG.length) {
  const chars = 'ABCDEFGHJKMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789';
  let result = '';
  for (let i = 0; i < length; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  return result;
}

// 生成验证码ID
function generateCaptchaId() {
  return crypto.randomBytes(16).toString('hex');
}

// 清理过期的验证码
function cleanupExpiredCaptchas() {
  const now = Date.now();
  for (const [id, captcha] of captchaStore.entries()) {
    // 如果验证码已过期或没有设置过期时间
    if (!captcha.expiresAt || captcha.expiresAt < now) {
      captchaStore.delete(id);
    }
  }
}

// 启动一个定时器定期清理过期验证码
setInterval(cleanupExpiredCaptchas, 60000); // 每分钟清理一次

module.exports = {
  captchaStore,
  CAPTCHA_CONFIG,
  generateCaptchaText,
  generateCaptchaId,
  cleanupExpiredCaptchas
};
const express = require('express');
const router = express.Router();
const bcrypt = require('bcryptjs');
const jwt = require('jsonwebtoken');
const { check, validationResult } = require('express-validator');
const User = require('../models/User');

// 用于存储验证码的内存存储（实际应用中应使用Redis等）
const captchaStore = new Map();

// 验证码配置
const CAPTCHA_CONFIG = {
  lifetime: 5 * 60 * 1000, // 验证码有效期：5分钟
  length: 6,               // 验证码长度
  chars: '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ' // 验证码字符集
};

// 验证码对象接口
// {
//   text: string,      // 验证码文本
//   expiresAt: number, // 过期时间戳
// }

// 生成随机文本验证码
function generateCaptchaText(length = CAPTCHA_CONFIG.length) {
  let result = '';
  for (let i = 0; i < length; i++) {
    result += CAPTCHA_CONFIG.chars.charAt(Math.floor(Math.random() * CAPTCHA_CONFIG.chars.length));
  }
  return result;
}

// 生成唯一且更安全的验证码ID
function generateCaptchaId() {
  const randomString = Math.random().toString(36).substr(2, 10);
  const timestamp = Date.now().toString(36);
  return `captcha_${timestamp}_${randomString}`;
}

// 定期清理过期CAPTCHAs
function cleanupExpiredCaptchas() {
  const now = Date.now();
  for (const [id, captcha] of captchaStore.entries()) {
    // 如果验证码已过期或没有设置过期时间
    if (!captcha.expiresAt || captcha.expiresAt < now) {
      captchaStore.delete(id);
    }
  }
}

// 启动一个定时器定期清理过期CAPTCHA
setInterval(cleanupExpiredCaptchas, 60 * 1000); // 每分钟清理一次

// 生成SVG格式的验证码图片
function generateCaptchaSVG(text) {
  const width = 150;
  const height = 50;
  
  let svg = `<svg width="${width}" height="${height}" xmlns="http://www.w3.org/2000/svg">`;
  svg += `<rect width="100%" height="100%" fill="#f0f0f0"/>`;
  
  // 添加干扰线
  for (let i = 0; i < 5; i++) {
    const x1 = Math.random() * width;
    const y1 = Math.random() * height;
    const x2 = Math.random() * width;
    const y2 = Math.random() * height;
    svg += `<line x1="${x1}" y1="${y1}" x2="${x2}" y2="${y2}" stroke="#aaa" stroke-width="1"/>`;
  }
  
  // 添加干扰点
  for (let i = 0; i < 20; i++) {
    const cx = Math.random() * width;
    const cy = Math.random() * height;
    const r = Math.random() * 2;
    svg += `<circle cx="${cx}" cy="${cy}" r="${r}" fill="#888"/>`;
  }
  
  // 添加文字（每个字符单独处理以增加干扰效果）
  const fontSize = 24;
  const charWidth = width / text.length;
  for (let i = 0; i < text.length; i++) {
    const x = charWidth * i + charWidth / 2;
    const y = height / 2 + fontSize / 3;
    const rotation = (Math.random() - 0.5) * 30; // 随机旋转角度
    svg += `<text x="${x}" y="${y}" font-family="Arial" font-size="${fontSize}" fill="#333" font-weight="bold" transform="rotate(${rotation} ${x} ${y})">${text[i]}</text>`;
  }
  
  svg += '</svg>';
  
  return svg;
}

// Register User
router.post(
  '/register',
  [
    check('username', 'Username is required').not().isEmpty(),
    check('email', 'Please include a valid email').isEmail(),
    check(
      'password',
      'Please enter a password with 6 or more characters'
    ).isLength({ min: 6 }),
  ],
  async (req, res) => {
    // 检查是否禁用注册功能
    if (process.env.DISABLE_REGISTRATION === 'true') {
      return res.status(403).json({ msg: 'Registration is disabled' });
    }

    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(400).json({ errors: errors.array() });
    }

    const { username, email, password, captcha, captchaId } = req.body;

    try {
      // 检查是否启用了验证码功能
      if (process.env.ENABLE_CAPTCHA === 'true') {
        // 验证验证码逻辑
        if (!captcha) {
          return res.status(400).json({ msg: 'Captcha is required' });
        }
        
        // 检查是否提供了验证码ID
        if (!captchaId) {
          return res.status(400).json({ msg: 'Captcha ID is required' });
        }
        
        // 验证验证码是否正确
        const storedCaptcha = captchaStore.get(captchaId);
        if (!storedCaptcha) {
          return res.status(400).json({ msg: 'Invalid or expired captcha' });
        }
        
        if (storedCaptcha.text !== captcha.toUpperCase()) {
          return res.status(400).json({ msg: 'Invalid captcha' });
        }
        
        // 验证成功后删除验证码，防止重复使用
        captchaStore.delete(captchaId);
      }

      let user = await User.findOne({ email });
      if (user) {
        return res.status(400).json({ msg: 'User already exists' });
      }

      user = new User({
        username,
        email,
        password,
      });

      const salt = await bcrypt.genSalt(10);
      user.password = await bcrypt.hash(password, salt);

      await user.save();

      const payload = {
        user: {
          id: user.id,
        },
      };

      jwt.sign(
        payload,
        process.env.JWT_SECRET,
        { expiresIn: '1h' },
        (err, token) => {
          if (err) throw err;
          res.json({ token });
        }
      );
    } catch (err) {
      console.error(err.message);
      res.status(500).send('Server Error');
    }
  });

// Login User
router.post(
  '/login',
  [
    check('email', 'Please include a valid email').isEmail(),
    check('password', 'Password is required').exists(),
  ],
  async (req, res) => {
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(400).json({ errors: errors.array() });
    }

    const { email, password, captcha, captchaId } = req.body;

    try {
      // 检查是否启用了验证码功能
      if (process.env.ENABLE_CAPTCHA === 'true') {
        // 验证验证码逻辑
        if (!captcha) {
          return res.status(400).json({ msg: 'Captcha is required' });
        }
        
        // 检查是否提供了验证码ID
        if (!captchaId) {
          return res.status(400).json({ msg: 'Captcha ID is required' });
        }
        
        // 验证验证码是否正确
        const storedCaptcha = captchaStore.get(captchaId);
        if (!storedCaptcha) {
          return res.status(400).json({ msg: 'Invalid or expired captcha' });
        }
        
        if (storedCaptcha.text !== captcha.toUpperCase()) {
          return res.status(400).json({ msg: 'Invalid captcha' });
        }
        
        // 验证成功后删除验证码，防止重复使用
        captchaStore.delete(captchaId);
      }

      let user = await User.findOne({ email });
      if (!user) {
        return res.status(400).json({ msg: 'Invalid Credentials' });
      }

      const isMatch = await bcrypt.compare(password, user.password);
      if (!isMatch) {
        return res.status(400).json({ msg: 'Invalid Credentials' });
      }

      const payload = {
        user: {
          id: user.id,
        },
      };

      jwt.sign(
        payload,
        process.env.JWT_SECRET,
        { expiresIn: '1h' },
        (err, token) => {
          if (err) throw err;
          res.json({ token });
        }
      );
    } catch (err) {
      console.error(err.message);
      res.status(500).send('Server Error');
    }
  }
);

// 生成验证码图片的路由
router.get('/captcha', (req, res) => {
  try {
    // 检查是否启用了验证码功能
    if (process.env.ENABLE_CAPTCHA !== 'true') {
      return res.status(400).json({ msg: 'Captcha is not enabled' });
    }
    
    // 生成随机验证码文本
    const captchaText = generateCaptchaText(6);
    
    // 生成验证码图片
    const captchaImage = generateCaptchaSVG(captchaText);
    
    // 生成安全的验证码ID
    const captchaId = generateCaptchaId();
    
    // 存储验证码并设置过期时间
    captchaStore.set(captchaId, {
      text: captchaText,
      createdAt: Date.now(),
      expiresAt: Date.now() + CAPTCHA_CONFIG.lifetime
    });
    
    // 定期清理过期验证码（实际应用中可以使用setInterval或Redis过期机制）
    if (Math.random() < 0.1) { // 10%的概率触发清理
      cleanupExpiredCaptchas();
    }
    
    // 返回验证码图片和ID
    res.json({ 
      image: `data:image/svg+xml;base64,${Buffer.from(captchaImage).toString('base64')}`,
      id: captchaId
    });
  } catch (err) {
    console.error(err.message);
    res.status(500).send('Server Error');
  }
});

module.exports = router;
const bcrypt = require('bcryptjs');
const jwt = require('jsonwebtoken');
const { body, validationResult, check } = require('express-validator');
const User = require('../models/User');
const { generateCaptchaText, generateCaptchaId, captchaStore, CAPTCHA_CONFIG, cleanupExpiredCaptchas } = require('../config/captcha');
const { emailCodeStore, EMAIL_CODE_CONFIG, generateEmailCode, generateEmailCodeId, cleanupExpiredEmailCodes, sendVerificationEmail } = require('../config/email');
const auth = require('../middleware/auth');

const router = require('express').Router();

// SVG验证码生成函数
function generateCaptchaSVG(text) {
  const width = 150;
  const height = 50;
  
  let svg = `<svg width="${width}" height="${height}" xmlns="http://www.w3.org/2000/svg">`;
  
  // 添加背景
  svg += `<rect width="100%" height="100%" fill="#f0f0f0"/>`;
  
  // 添加干扰线
  for (let i = 0; i < 5; i++) {
    const x1 = Math.random() * width;
    const y1 = Math.random() * height;
    const x2 = Math.random() * width;
    const y2 = Math.random() * height;
    const color = `#${Math.floor(Math.random()*16777215).toString(16)}`;
    svg += `<line x1="${x1}" y1="${y1}" x2="${x2}" y2="${y2}" stroke="${color}" stroke-width="1"/>`;
  }
  
  // 添加验证码文本
  const fontSize = 24;
  const x = (width - text.length * fontSize * 0.6) / 2;
  const y = height / 2 + fontSize / 3;
  
  for (let i = 0; i < text.length; i++) {
    const charX = x + i * fontSize * 0.7;
    const charY = y + (Math.random() - 0.5) * 10;
    const rotation = (Math.random() - 0.5) * 30;
    const color = `#${Math.floor(Math.random()*16777215).toString(16)}`;
    
    svg += `<text x="${charX}" y="${charY}" font-size="${fontSize}" fill="${color}" transform="rotate(${rotation} ${charX} ${charY})" font-family="Arial">${text[i]}</text>`;
  }
  
  svg += '</svg>';
  
  return svg;
}

// @route   GET api/auth/me
// @desc    Get authenticated user
// @access  Private
router.get('/me', auth, async (req, res) => {
  try {
    const user = await User.findById(req.user.id).select('-password');
    res.json(user);
  } catch (err) {
    console.error(err.message);
    res.status(500).send('Server Error');
  }
});

// @route   POST api/auth/send-email-code
// @desc    Send email verification code
// @access  Public
router.post(
  '/send-email-code',
  check('email', 'Please include a valid email').isEmail(),
  async (req, res) => {
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(400).json({ errors: errors.array() });
    }

    const { email } = req.body;

    try {
      // 检查用户是否已存在
      let user = await User.findOne({ email });
      if (!user) {
        return res.status(400).json({ msg: 'User does not exist' });
      }

      // 生成邮箱验证码
      const emailCode = generateEmailCode();
      const emailCodeId = generateEmailCodeId();

      // 存储验证码并设置过期时间
      emailCodeStore.set(emailCodeId, {
        email,
        code: emailCode,
        createdAt: Date.now(),
        expiresAt: Date.now() + EMAIL_CODE_CONFIG.lifetime,
        attempts: 0
      });

      // 发送验证码邮件
      await sendVerificationEmail(email, emailCode);

      // 定期清理过期验证码
      if (Math.random() < 0.1) { // 10%的概率触发清理
        cleanupExpiredEmailCodes();
      }

      res.json({ 
        msg: 'Verification code sent successfully',
        id: emailCodeId
      });
    } catch (err) {
      console.error(err.message);
      res.status(500).send('Server Error');
    }
  }
);

// @route   POST api/auth/send-login-email-code
// @desc    Send email verification code (for login)
// @access  Public
router.post(
  '/send-login-email-code',
  check('email', 'Please include a valid email').isEmail(),
  async (req, res) => {
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(400).json({ errors: errors.array() });
    }

    const { email } = req.body;

    try {
      // 检查用户是否存在
      let user = await User.findOne({ email });
      if (!user) {
        return res.status(400).json({ msg: 'User does not exist' });
      }

      // 生成邮箱验证码
      const emailCode = generateEmailCode();
      const emailCodeId = generateEmailCodeId();

      // 存储验证码并设置过期时间
      emailCodeStore.set(emailCodeId, {
        email,
        code: emailCode,
        createdAt: Date.now(),
        expiresAt: Date.now() + EMAIL_CODE_CONFIG.lifetime,
        attempts: 0
      });

      // 发送验证码邮件
      await sendVerificationEmail(email, emailCode);

      // 定期清理过期验证码
      if (Math.random() < 0.1) { // 10%的概率触发清理
        cleanupExpiredEmailCodes();
      }

      res.json({ 
        msg: 'Login verification code sent successfully',
        id: emailCodeId
      });
    } catch (err) {
      console.error(err.message);
      res.status(500).send('Server Error');
    }
  }
);

// @route   POST api/auth/verify-captcha
// @desc    Verify captcha
// @access  Public
router.post(
  '/verify-captcha',
  async (req, res) => {
    const { captcha, captchaId } = req.body;

    try {
      // 检查是否启用了验证码功能
      if (process.env.ENABLE_CAPTCHA !== 'true') {
        return res.status(400).json({ msg: 'Captcha is not enabled' });
      }

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
      
      res.json({ msg: 'Captcha verified successfully' });
    } catch (err) {
      console.error(err.message);
      res.status(500).send('Server Error');
    }
  }
);

// @route   POST api/auth/register
// @desc    Register user
// @access  Public
router.post(
  '/register',
  check('username', 'Name is required').notEmpty(), // 修复字段验证，使用username而不是name
  check('email', 'Please include a valid email').isEmail(),
  check(
    'password',
    'Please enter a password with 6 or more characters'
  ).isLength({ min: 6 }),
  async (req, res) => {
    // 检查是否禁用注册功能
    if (process.env.DISABLE_REGISTRATION === 'true') {
      return res.status(403).json({ msg: 'Registration is disabled' });
    }

    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(400).json({ errors: errors.array() });
    }

    const { username, email, password, emailCode, emailCodeId } = req.body;

    try {
      // 检查是否启用了邮箱验证码功能
      if (process.env.ENABLE_EMAIL_VERIFICATION === 'true') {
        // 验证邮箱验证码逻辑
        if (!emailCode) {
          return res.status(400).json({ msg: 'Email verification code is required' });
        }
        
        // 检查是否提供了邮箱验证码ID
        if (!emailCodeId) {
          return res.status(400).json({ msg: 'Email verification code ID is required' });
        }
        
        // 验证验证码是否正确
        const storedEmailCode = emailCodeStore.get(emailCodeId);
        if (!storedEmailCode) {
          return res.status(400).json({ msg: 'Invalid or expired email verification code' });
        }
        
        // 检查邮箱是否匹配
        if (storedEmailCode.email !== email) {
          return res.status(400).json({ msg: 'Email does not match the verification code' });
        }
        
        // 检查尝试次数
        if (storedEmailCode.attempts >= EMAIL_CODE_CONFIG.maxAttempts) {
          emailCodeStore.delete(emailCodeId);
          return res.status(400).json({ msg: 'Too many attempts. Please request a new verification code.' });
        }
        
        // 增加尝试次数
        storedEmailCode.attempts += 1;
        emailCodeStore.set(emailCodeId, storedEmailCode);
        
        // 验证验证码是否正确
        if (storedEmailCode.code !== emailCode.toUpperCase()) {
          return res.status(400).json({ msg: 'Invalid email verification code' });
        }
        
        // 验证成功后删除验证码，防止重复使用
        emailCodeStore.delete(emailCodeId);
      }

      let user = await User.findOne({ email });
      if (user) {
        return res.status(400).json({ msg: 'User already exists' });
      }

      user = new User({
        username, // 使用username而不是name
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

// @route   POST api/auth/login
// @desc    Authenticate user & get token (support both password and email code login)
// @access  Public
router.post(
  '/login',
  [
    check('email', 'Please include a valid email').isEmail(),
  ],
  async (req, res) => {
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(400).json({ errors: errors.array() });
    }

    const { email, password, captcha, captchaId, emailCode, emailCodeId } = req.body;

    try {
      let user = await User.findOne({ email });
      if (!user) {
        return res.status(400).json({ msg: 'Invalid Credentials' });
      }

      // 检查是否启用了邮箱验证码登录
      const isEmailCodeLogin = process.env.ENABLE_EMAIL_VERIFICATION === 'true' && emailCode && emailCodeId;
      
      // 如果提供了邮箱验证码，则验证邮箱验证码登录
      if (isEmailCodeLogin) {
        // 验证邮箱验证码逻辑
        const storedEmailCode = emailCodeStore.get(emailCodeId);
        if (!storedEmailCode) {
          return res.status(400).json({ msg: 'Invalid or expired email verification code' });
        }
        
        // 检查邮箱是否匹配
        if (storedEmailCode.email !== email) {
          return res.status(400).json({ msg: 'Email does not match the verification code' });
        }
        
        // 检查尝试次数
        if (storedEmailCode.attempts >= EMAIL_CODE_CONFIG.maxAttempts) {
          emailCodeStore.delete(emailCodeId);
          return res.status(400).json({ msg: 'Too many attempts. Please request a new verification code.' });
        }
        
        // 增加尝试次数
        storedEmailCode.attempts += 1;
        emailCodeStore.set(emailCodeId, storedEmailCode);
        
        // 验证验证码是否正确
        if (storedEmailCode.code !== emailCode.toUpperCase()) {
          return res.status(400).json({ msg: 'Invalid email verification code' });
        }
        
        // 验证成功后删除验证码，防止重复使用
        emailCodeStore.delete(emailCodeId);
      } 
      // 否则验证密码登录
      else {
        // 检查是否提供了密码
        if (!password) {
          return res.status(400).json({ msg: 'Password is required' });
        }
        
        const isMatch = await bcrypt.compare(password, user.password);
        if (!isMatch) {
          return res.status(400).json({ msg: 'Invalid Credentials' });
        }
      }

      // 检查是否启用了验证码功能
      // 但邮箱验证码登录时不需要图片验证码
      if (process.env.ENABLE_CAPTCHA === 'true' && !isEmailCodeLogin) {
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
        
        // 检查验证码是否过期
        if (storedCaptcha.expiresAt < Date.now()) {
          captchaStore.delete(captchaId);
          return res.status(400).json({ msg: 'Captcha has expired' });
        }
        
        // 输出调试信息
        console.log('验证码验证信息:', {
          inputCaptcha: captcha,
          storedCaptcha: storedCaptcha.text,
          inputUppercase: captcha.toUpperCase(),
          storedUppercase: storedCaptcha.text.toUpperCase(),
          matchResult: storedCaptcha.text.toUpperCase() === captcha.toUpperCase()
        });
        
        // 统一转换为大写进行比较
        if (storedCaptcha.text.toUpperCase() !== captcha.toUpperCase()) {
          return res.status(400).json({ msg: 'Invalid captcha' });
        }
        
        // 验证成功后删除验证码，防止重复使用
        captchaStore.delete(captchaId);
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

// @route   GET api/auth/captcha
// @desc    Generate captcha image
// @access  Public
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
    
    // 存储验证码并设置过期时间 (转换为大写以统一处理)
    captchaStore.set(captchaId, {
      text: captchaText.toUpperCase(),
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
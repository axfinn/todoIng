const axios = require('axios');

// 配置axios基础URL
const api = axios.create({
  baseURL: 'http://localhost:5001/api'
});

async function testCaptcha() {
  try {
    // 获取验证码
    console.log('获取验证码...');
    const captchaResponse = await api.get('/auth/captcha');
    console.log('验证码响应:', captchaResponse.data);
    
    // 尝试使用错误的验证码注册
    console.log('\n尝试使用错误的验证码注册...');
    const registerResponse = await api.post('/auth/register', {
      username: 'testuser',
      email: 'test@example.com',
      password: 'password123',
      captcha: 'INVALID',
      captchaId: captchaResponse.data.id
    });
    
    console.log('注册成功:', registerResponse.data);
  } catch (error) {
    if (error.response) {
      console.log('注册失败(预期):', error.response.data);
    } else {
      console.log('请求错误:', error.message);
    }
  }
  
  try {
    // 获取验证码
    console.log('\n获取验证码...');
    const captchaResponse = await api.get('/auth/captcha');
    console.log('验证码响应:', captchaResponse.data);
    
    // 使用正确的验证码注册
    console.log('\n尝试使用正确的验证码注册...');
    const registerResponse = await api.post('/auth/register', {
      username: 'testuser2',
      email: 'test2@example.com',
      password: 'password123',
      captcha: captchaResponse.data.text,
      captchaId: captchaResponse.data.id
    });
    
    console.log('注册成功:', registerResponse.data);
  } catch (error) {
    if (error.response) {
      console.log('注册失败:', error.response.data);
    } else {
      console.log('请求错误:', error.message);
    }
  }
}

testCaptcha();
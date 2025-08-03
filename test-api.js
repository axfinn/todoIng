const axios = require('axios');

// 测试API
async function testAPI() {
  try {
    // 1. 获取验证码
    console.log('1. 获取验证码...');
    const captchaRes = await axios.get('http://localhost:5001/api/auth/captcha');
    console.log('验证码ID:', captchaRes.data.id);
    
    // 2. 登录
    console.log('\n2. 登录...');
    const loginRes = await axios.post('http://localhost:5001/api/auth/login', {
      email: 'admin@example.com',
      password: 'password123',
      captcha: 'CGG9Z', // 这里需要手动输入验证码
      captchaId: captchaRes.data.id
    });
    
    console.log('登录成功，Token:', loginRes.data.token);
    
    // 3. 测试获取报告列表
    console.log('\n3. 获取报告列表...');
    const reportsRes = await axios.get('http://localhost:5001/api/reports', {
      headers: {
        'Authorization': `Bearer ${loginRes.data.token}`
      }
    });
    
    console.log('报告列表:', reportsRes.data);
    
  } catch (error) {
    console.error('错误:', error.response ? error.response.data : error.message);
  }
}

testAPI();
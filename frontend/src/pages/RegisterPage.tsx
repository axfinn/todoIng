import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { useTranslation } from 'react-i18next';
import type { AppDispatch, RootState } from '../app/store';
import { registerUser } from '../features/auth/authSlice';
import api from '../config/api';

// 为调试信息添加样式
const debugStyle = {
  border: '1px solid #e9ecef',
  padding: '1rem',
  marginBottom: '1rem',
  backgroundColor: '#f8f9fa',
  fontFamily: 'monospace',
  fontSize: '0.9rem'
};

interface FormData {
  username: string;
  email: string;
  password: string;
  password2: string;
  emailCode: string;
}

const RegisterPage: React.FC = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const dispatch = useDispatch<AppDispatch>();
  
  const { isAuthenticated, isLoading, error } = useSelector((state: RootState) => state.auth);
  
  const [formData, setFormData] = useState<FormData>({
    username: '',
    email: '',
    password: '',
    password2: '',
    captcha: '',
    emailCode: ''
  });
  
  const [captchaData, setCaptchaData] = useState<{ image: string; id: string } | null>(null);
  const [isCaptchaFetching, setIsCaptchaFetching] = useState(false);
  const [emailCodeData, setEmailCodeData] = useState<{ id: string } | null>(null);
  const [isEmailCodeSending, setIsEmailCodeSending] = useState(false);
  const [emailCodeSent, setEmailCodeSent] = useState(false);
  const [countdown, setCountdown] = useState(0);
  const [captchaVerified, setCaptchaVerified] = useState(false);
  
  const { username, email, password, password2, emailCode } = formData;
  
  // 检查环境变量
  const isCaptchaEnabled = import.meta.env.VITE_ENABLE_CAPTCHA === 'true';
  const isEmailVerificationEnabled = import.meta.env.VITE_ENABLE_EMAIL_VERIFICATION === 'true';
  const isRegistrationDisabled = import.meta.env.VITE_DISABLE_REGISTRATION === 'true';
  
  // 倒计时效果
  useEffect(() => {
    let timer: NodeJS.Timeout;
    if (countdown > 0) {
      timer = setTimeout(() => setCountdown(countdown - 1), 1000);
    }
    return () => {
      if (timer) clearTimeout(timer);
    };
  }, [countdown]);
  
// Removed captcha related effect
  
  // 检查认证状态
  useEffect(() => {
    if (isAuthenticated) {
      navigate('/dashboard');
    }
  }, [isAuthenticated, navigate]);
  
  // 如果注册被禁用，重定向到登录页面
  if (isRegistrationDisabled) {
    return (
      <div className="container py-5">
        <div className="row justify-content-center">
          <div className="col-md-6 col-lg-5">
            <div className="card shadow-sm">
              <div className="card-body p-5">
                <div className="text-center mb-4">
                  <h2 className="fw-bold">{t('auth.register.title')}</h2>
                </div>
                
                <div className="alert alert-warning" role="alert">
                  <h4 className="alert-heading">{t('auth.register.disabledTitle')}</h4>
                  <p>{t('auth.register.disabledMessage')}</p>
                  <hr />
                  <Link to="/login" className="btn btn-primary">
                    {t('auth.register.login')}
                  </Link>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  // 发送邮箱验证码
  const sendEmailCode = async () => {
    if (!isEmailVerificationEnabled || !email || !/\S+@\S+\.\S+/.test(email)) {
      return;
    }
    
    if (countdown > 0) return;
    
    setIsEmailCodeSending(true);
    try {
      const res = await api.post('/auth/send-email-code', { email });
      setEmailCodeData({
        id: res.data.id
      });
      setEmailCodeSent(true);
      setCountdown(60); // 60秒倒计时
    } catch (err: any) {
      console.error(err);
      alert(err.response?.data?.msg || 'Failed to send verification code');
    } finally {
      setIsEmailCodeSending(false);
    }
  };

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    // 前端表单验证
    if (password !== password2) {
      console.warn('密码不匹配:', { password, confirmPassword: password2 });
      alert(t('auth.register.passwordMismatch'));
      return;
    }
    
    // 检查邮箱格式
    const emailRegex = /\S+@\S+\.\S+/;
    if (!emailRegex.test(email)) {
      console.warn('无效邮箱格式:', { email });
      alert(t('auth.register.invalidEmail'));
      return;
    }
    
    // 检查邮箱验证码是否已发送
    if (isEmailVerificationEnabled && !emailCodeSent) {
      console.warn('邮箱验证码未发送');
      alert(t('auth.register.emailCodeRequired'));
      return;
    }
    
    // 构建用户数据
    const userData: Record<string, unknown> = {
      username,
      email,
      password
    };

    // 添加邮箱验证码数据
    if (isEmailVerificationEnabled && emailCodeData) {
      userData.emailCode = emailCode;
      userData.emailCodeId = emailCodeData.id;
    }
    
    // 调试信息
    console.log('注册数据:', {
      ...userData,
      password: '*****', // 隐藏密码
      timestamp: new Date().toISOString()
    });
    
    // 发送注册请求
    try {
      dispatch(registerUser(userData));
    } catch (error) {
      console.error('注册请求异常:', error);
      alert(t('auth.register.errors.registrationFailed'));
    }
  };

  // 验证图片验证码
  const verifyCaptcha = async () => {
    if (!captcha || !captchaData) {
      alert(t('auth.register.enterCaptcha'));
      return;
    }
    
    try {
      // 发送一个简单的请求来验证验证码
      await api.post('/auth/verify-captcha', {
        captcha,
        captchaId: captchaData.id
      });
      
      // 验证成功
      setCaptchaVerified(true);
      alert(t('auth.register.captchaVerified'));
    } catch (err: any) {
      console.error(err);
      alert(err.response?.data?.msg || 'Failed to verify captcha');
      refreshCaptcha(); // 刷新验证码
    }
  };

  if (isLoading) {
    return (
      <div className="container py-5">
        <div className="d-flex justify-content-center align-items-center" style={{ height: '70vh' }}>
          <div className="text-center">
            <div className="spinner-border text-primary" role="status">
              <span className="visually-hidden">{t('common.loading')}</span>
            </div>
            <p className="mt-3 text-muted">{t('auth.register.pleaseWait')}</p>
          </div>
        </div>
      </div>
    );
  }
  
  // 添加调试信息显示区域
  const renderDebugInfo = () => {
    if (import.meta.env.MODE !== 'development') return null;
    
    return (
      <div style={debugStyle}>
        <h5 className="mb-3">调试信息 (仅开发环境显示)</h5>
        <pre className="mb-0" style={{ overflowX: 'auto' }}>
          {JSON.stringify({
            formState: formData,
            emailVerificationEnabled: isEmailVerificationEnabled,
            emailCodeSent: emailCodeSent,
            countdown: countdown,
            emailCodeData: emailCodeData,
            registrationDisabled: isRegistrationDisabled
          }, null, 2)}
        </pre>
      </div>
    );
  };

  return (
    <div className="container py-5">
      <div className="row justify-content-center">
        <div className="col-md-6 col-lg-5">
          <div className="card shadow-sm">
            <div className="card-body p-5">
              <div className="text-center mb-4">
                <h2 className="fw-bold">{t('auth.register.title')}</h2>
              </div>

              {error && (
                <div className="alert alert-danger alert-dismissible fade show" role="alert">
                  <strong>{t('auth.register.error')}</strong> {error}
                  <button
                    type="button"
                    className="btn-close"
                    data-bs-dismiss="alert"
                    aria-label="Close"
                    onClick={() => dispatch({ type: 'auth/clearError' })}
                  ></button>
                </div>
              )}
              
              {/* 显示调试信息 */}
              {renderDebugInfo()}

              <form onSubmit={onSubmit}>
                <div className="mb-3">
                  <label htmlFor="username" className="form-label">{t('auth.register.name')}</label>
                  <input
                    type="text"
                    className="form-control"
                    id="username"
                    name="username"
                    value={username}
                    onChange={onChange}
                    required
                  />
                </div>

                <div className="mb-3">
                  <label htmlFor="email" className="form-label">{t('auth.register.email')}</label>
                  <input
                    type="email"
                    className="form-control"
                    id="email"
                    name="email"
                    value={email}
                    onChange={onChange}
                    required
                  />
                </div>
                
                {/* 邮箱验证码输入框 */}
              {isEmailVerificationEnabled && (
                <div className="mb-3">
                  <label htmlFor="emailCode" className="form-label">{t('auth.register.emailCode')}</label>
                  <div className="input-group">
                    <input
                      type="text"
                      className="form-control"
                      id="emailCode"
                      name="emailCode"
                      value={emailCode}
                      onChange={onChange}
                      required
                      placeholder={t('auth.register.enterEmailCode')}
                    />
                    <button 
                      type="button"
                      className="btn btn-outline-secondary"
                      onClick={sendEmailCode}
                      disabled={isEmailCodeSending || countdown > 0 || !email || !/\S+@\S+\.\S+/.test(email)}
                    >
                      {isEmailCodeSending ? (
                        t('common.loading')
                      ) : countdown > 0 ? (
                        t('auth.register.resendCode', { seconds: countdown })
                      ) : (
                        t('auth.register.getCode')
                      )}
                    </button>
                  </div>
                  {emailCodeSent && (
                    <div className="form-text text-success">
                      {t('auth.register.emailCodeSent')}
                    </div>
                  )}
                </div>
              )}

              <div className="mb-3">
                <label htmlFor="password" className="form-label">{t('auth.register.password')}</label>
                <input
                  type="password"
                  className="form-control"
                  id="password"
                  name="password"
                  value={password}
                  onChange={onChange}
                  required
                />
              </div>

                <div className="mb-3">
                  <label htmlFor="password2" className="form-label">{t('auth.register.confirmPassword')}</label>
                  <input
                    type="password"
                    className="form-control"
                    id="password2"
                    name="password2"
                    value={password2}
                    onChange={onChange}
                    required
                  />
                </div>


                <div className="d-grid">
                  <button 
                    type="submit" 
                    className="btn btn-primary btn-lg" 
                    disabled={
                      isLoading || 
                      (isEmailVerificationEnabled && !emailCodeSent)
                    }
                  >
                    {isLoading ? t('common.loading') : t('auth.register.submit')}
                  </button>
                </div>
              </form>

              <div className="text-center mt-4">
                <p className="mb-0">
                  {t('auth.register.haveAccount')} <Link to="/login" className="text-decoration-none">{t('auth.register.login')}</Link>
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default RegisterPage;
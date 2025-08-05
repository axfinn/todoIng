import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useNavigate, Link } from 'react-router-dom';
import { loginUser } from '../features/auth/authSlice';
import { useTranslation } from 'react-i18next';
import type { AppDispatch, RootState } from '../app/store';
import api from '../config/api';

const LoginPage: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  const navigate = useNavigate();
  const { t } = useTranslation();
  const { isAuthenticated, isLoading, error } = useSelector((state: RootState) => state.auth);

  const [backgroundImage, setBackgroundImage] = useState('');

  // 背景图片数组
  const backgroundImages = [
    'https://picsum.photos/1920/1080?random=1',
    'https://picsum.photos/1920/1080?random=2',
    'https://picsum.photos/1920/1080?random=3'
  ];

  // 设置随机背景图片
  useEffect(() => {
    const randomIndex = Math.floor(Math.random() * backgroundImages.length);
    setBackgroundImage(backgroundImages[randomIndex]);
  }, []);

  const [formData, setFormData] = useState({
    email: '',
    password: '',
    captcha: '',
    emailCode: ''
  });

  const [captchaData, setCaptchaData] = useState<{ image: string; id: string } | null>(null);
  const [isCaptchaFetching, setIsCaptchaFetching] = useState(false);
  const [emailCodeData, setEmailCodeData] = useState<{ id: string } | null>(null);
  const [isEmailCodeSending, setIsEmailCodeSending] = useState(false);
  const [countdown, setCountdown] = useState(0);
  const [loginMethod, setLoginMethod] = useState<'password' | 'emailCode'>('password');

  // 从环境变量获取功能开关状态
  const isCaptchaEnabled = import.meta.env.VITE_ENABLE_CAPTCHA === 'true';
  const isEmailVerificationEnabled = import.meta.env.VITE_ENABLE_EMAIL_VERIFICATION === 'true';

  // 获取验证码
  const getCaptcha = async () => {
    if (!isCaptchaEnabled) return;
    
    setIsCaptchaFetching(true);
    try {
      const res = await api.get('/auth/captcha');
      setCaptchaData({
        image: res.data.image,
        id: res.data.id
      });
    } catch (err) {
      console.error('获取验证码失败:', err);
    } finally {
      setIsCaptchaFetching(false);
    }
  };

  // 发送邮箱验证码
  const sendEmailCode = async () => {
    if (!isEmailVerificationEnabled || !formData.email || !/\S+@\S+\.\S+/.test(formData.email)) {
      return;
    }
    
    if (countdown > 0) return;
    
    setIsEmailCodeSending(true);
    try {
      const res = await api.post('/auth/send-login-email-code', { email: formData.email });
      setEmailCodeData({
        id: res.data.id
      });
      setCountdown(60); // 60秒倒计时
    } catch (err: any) {
      console.error(err);
      alert(err.response?.data?.msg || 'Failed to send verification code');
    } finally {
      setIsEmailCodeSending(false);
    }
  };

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

  useEffect(() => {
    if (isAuthenticated) {
      navigate('/dashboard');
    }
  }, [isAuthenticated, navigate]);

  useEffect(() => {
    // 如果启用了验证码功能且使用密码登录，则获取验证码
    if (isCaptchaEnabled && loginMethod === 'password') {
      getCaptcha();
    }
  }, [isCaptchaEnabled, loginMethod]);

  const { email, password, captcha, emailCode } = formData;

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    // 准备登录数据
    const userData: Record<string, string> = { email };
    
    // 根据登录方式添加相应数据
    if (loginMethod === 'password') {
      userData.password = password;
    } else if (loginMethod === 'emailCode' && emailCodeData) {
      userData.emailCode = emailCode;
      userData.emailCodeId = emailCodeData.id;
    }
    
    // 如果启用了验证码且有验证码ID，并且使用密码登录，则添加验证码相关数据
    if (isCaptchaEnabled && captchaData && loginMethod === 'password') {
      userData.captcha = captcha;
      userData.captchaId = captchaData.id;
    }
    
    dispatch(loginUser(userData));
  };

  if (isLoading) {
    return (
      <div className="container py-5">
        <div className="d-flex justify-content-center align-items-center" style={{ height: '70vh' }}>
          <div className="text-center">
            <div className="spinner-border text-primary" role="status">
              <span className="visually-hidden">{t('common.loading')}</span>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div 
      className="d-flex flex-column min-vh-100"
      style={{
        backgroundImage: `url(${backgroundImage})`,
        backgroundSize: 'cover',
        backgroundPosition: 'center',
        backgroundRepeat: 'no-repeat',
        backgroundAttachment: 'fixed'
      }}
    >
      <div className="container py-5 flex-grow-1 d-flex align-items-center">
        <div className="row justify-content-center w-100">
          <div className="col-md-6 col-lg-5">
            <div className="card shadow-lg" style={{ 
              backgroundColor: 'rgba(255, 255, 255, 0.8)',
              borderRadius: '10px',
              backdropFilter: 'blur(10px)'
            }}>
              <div className="card-body p-5">
                <div className="text-center mb-4">
                  <h2 className="fw-bold">{t('auth.login.title')}</h2>
                </div>
              
                {error && (
                  <div className="alert alert-danger" role="alert">
                    {error}
                  </div>
                )}
              
                <form onSubmit={onSubmit}>
                  <div className="mb-3">
                    <label htmlFor="email" className="form-label">{t('auth.login.email')}</label>
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
                  
                  {/* 登录方式切换 */}
                  {isEmailVerificationEnabled && (
                    <div className="mb-3">
                      <div className="btn-group w-100" role="group">
                        <input
                          type="radio"
                          className="btn-check"
                          name="loginMethod"
                          id="passwordLogin"
                          checked={loginMethod === 'password'}
                          onChange={() => setLoginMethod('password')}
                        />
                        <label className="btn btn-outline-primary" htmlFor="passwordLogin">
                          {t('auth.login.passwordLogin')}
                        </label>
                        
                        <input
                          type="radio"
                          className="btn-check"
                          name="loginMethod"
                          id="emailCodeLogin"
                          checked={loginMethod === 'emailCode'}
                          onChange={() => setLoginMethod('emailCode')}
                        />
                        <label className="btn btn-outline-primary" htmlFor="emailCodeLogin">
                          {t('auth.login.emailCodeLogin')}
                        </label>
                      </div>
                    </div>
                  )}
                  
                  {/* 密码输入框 */}
                  {loginMethod === 'password' && (
                    <div className="mb-3">
                      <label htmlFor="password" className="form-label">{t('auth.login.password')}</label>
                      <input
                        type="password"
                        className="form-control"
                        id="password"
                        name="password"
                        value={password}
                        onChange={onChange}
                        required={loginMethod === 'password'}
                      />
                    </div>
                  )}
                  
                  {/* 邮箱验证码输入框 */}
                  {loginMethod === 'emailCode' && isEmailVerificationEnabled && (
                    <div className="mb-3">
                      <label htmlFor="emailCode" className="form-label">{t('auth.login.emailCode')}</label>
                      <div className="input-group">
                        <input
                          type="text"
                          className="form-control"
                          id="emailCode"
                          name="emailCode"
                          value={emailCode}
                          onChange={onChange}
                          required={loginMethod === 'emailCode'}
                          placeholder={t('auth.login.enterEmailCode')}
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
                            t('auth.login.resendCode', { seconds: countdown })
                          ) : (
                            t('auth.login.getCode')
                          )}
                        </button>
                      </div>
                    </div>
                  )}
                  
                  {/* 验证码输入框 - 仅在密码登录时显示 */}
                  {isCaptchaEnabled && loginMethod === 'password' && (
                    <div className="mb-3 position-relative">
                      <label htmlFor="captcha" className="form-label">{t('auth.login.captcha')}</label>
                      <div className="input-group">
                        <input
                          type="text"
                          className="form-control"
                          id="captcha"
                          name="captcha"
                          value={captcha}
                          onChange={onChange}
                          required
                        />
                        <button 
                          type="button" 
                          className="btn btn-outline-secondary"
                          onClick={getCaptcha}
                          disabled={isCaptchaFetching}
                        >
                          {isCaptchaFetching ? t('common.loading') : t('auth.login.refreshCaptcha')}
                        </button>
                      </div>
                      <div className="mt-2 d-flex justify-content-center">
                        {isCaptchaFetching ? (
                          <div className="captcha-placeholder">{t('common.loading')}</div>
                        ) : (
                          <>
                            {captchaData?.image ? (
                              <img 
                                src={captchaData.image} 
                                alt={t('auth.login.captcha')} 
                                className="captcha-image"
                                onClick={getCaptcha}
                                style={{ cursor: 'pointer', maxHeight: '50px' }}
                              />
                            ) : (
                              <div className="captcha-placeholder">{t('auth.login.captcha')}</div>
                            )}
                          </>
                        )}
                      </div>
                    </div>
                  )}
                  
                  <div className="d-grid">
                    <button type="submit" className="btn btn-primary btn-lg" disabled={isLoading}>
                      {isLoading ? t('common.loading') : t('auth.login.submit')}
                    </button>
                  </div>
                </form>
                
                <div className="text-center mt-4">
                  <p className="mb-0">
                    {t('auth.login.noAccount')}{' '}
                    <Link to="/register" className="text-decoration-none">
                      {t('auth.login.register')}
                    </Link>
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
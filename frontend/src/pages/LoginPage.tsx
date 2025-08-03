import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useNavigate, Link } from 'react-router-dom';
import { loginUser } from '../features/auth/authSlice';
import { useTranslation } from 'react-i18next';
import type { AppDispatch, RootState } from '../app/store';

const LoginPage: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  const navigate = useNavigate();
  const { t } = useTranslation();
  const { isAuthenticated, isLoading, error } = useSelector((state: RootState) => state.auth);

  const [formData, setFormData] = useState({
    email: '',
    password: '',
    captcha: ''
  });

  const [captchaImage, setCaptchaImage] = useState<string | null>(null);
  const [captchaId, setCaptchaId] = useState<string | null>(null);

  // 从环境变量获取验证码功能开关状态
  const isCaptchaEnabled = import.meta.env.VITE_ENABLE_CAPTCHA === 'true';

  // 获取验证码
  const getCaptcha = async () => {
    try {
      const response = await fetch('/api/auth/captcha');
      const data = await response.json();
      if (data.image && data.id) {
        setCaptchaImage(data.image);
        setCaptchaId(data.id);
      }
    } catch (error) {
      console.error('获取验证码失败:', error);
    }
  };

  useEffect(() => {
    if (isAuthenticated) {
      navigate('/dashboard');
    }
  }, [isAuthenticated, navigate]);

  useEffect(() => {
    // 如果启用了验证码功能，则获取验证码
    if (isCaptchaEnabled) {
      getCaptcha();
    }
  }, [isCaptchaEnabled]);

  const { email, password, captcha } = formData;

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    // 准备登录数据
    const userData: Record<string, string> = { email, password };
    
    // 如果启用了验证码且有验证码ID，则添加验证码相关数据
    if (isCaptchaEnabled && captchaId) {
      userData.captcha = captcha;
      userData.captchaId = captchaId;
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
    <div className="container py-5">
      <div className="row justify-content-center">
        <div className="col-md-6 col-lg-5">
          <div className="card shadow-sm">
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
                
                <div className="mb-3">
                  <label htmlFor="password" className="form-label">{t('auth.login.password')}</label>
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
                
                {/* 验证码输入框 - 与注册页面保持一致 */}
                {isCaptchaEnabled && (
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
                        disabled={isLoading}
                      >
                        {t('auth.login.refreshCaptcha')}
                      </button>
                    </div>
                    {captchaImage && (
                      <div className="mt-2 text-center">
                        <img 
                          src={captchaImage} 
                          alt={t('auth.login.captcha')} 
                          className="img-fluid rounded"
                          style={{ maxHeight: '80px', cursor: 'pointer' }}
                          onClick={getCaptcha}
                        />
                      </div>
                    )}
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
              
              <div className="text-center mt-4">
                <div className="d-flex align-items-center justify-content-center">
                  <i className="bi bi-github me-2"></i>
                  <a href="https://github.com/axfinn/todoIng" target="_blank" rel="noopener noreferrer" className="text-decoration-none">
                    Fork me on GitHub
                  </a>
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
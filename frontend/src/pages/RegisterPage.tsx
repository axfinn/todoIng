import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useNavigate, Link } from 'react-router-dom';
import { registerUser } from '../features/auth/authSlice';
import { useTranslation } from 'react-i18next';
import type { AppDispatch, RootState } from '../app/store';

const RegisterPage: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  const navigate = useNavigate();
  const { t } = useTranslation();
  const { isAuthenticated, isLoading, error } = useSelector((state: RootState) => state.auth);
  
  // 从环境变量获取注册功能开关状态
  const isRegistrationDisabled = import.meta.env.VITE_DISABLE_REGISTRATION === 'true';

  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    password2: '',
    captcha: ''
  });

  const [captchaImage, setCaptchaImage] = useState<string | null>(null);

  // 从环境变量获取验证码功能开关状态
  const isCaptchaEnabled = import.meta.env.VITE_ENABLE_CAPTCHA === 'true';


  useEffect(() => {
    if (isAuthenticated) {
      navigate('/dashboard');
    }
  }, [isAuthenticated, navigate]);


  const { username, email, password, password2, captcha } = formData;

  // 获取验证码
  const getCaptcha = async () => {
    try {
      const response = await fetch('/api/auth/captcha');
      const data = await response.json();
      if (response.ok) {
        setCaptchaImage(data.image);
      } else {
        console.error('Failed to get captcha:', data.msg);
      }
    } catch (err) {
      console.error('Error getting captcha:', err);
    }
  };

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (password !== password2) {
      alert(t('auth.register.passwordMismatch'));
      return;
    }
    
    // 准备注册数据
    const registerData: Record<string, string> = { username, email, password, captcha };
    
    dispatch(registerUser(registerData));
  };

  // 如果注册功能被禁用，显示提示信息
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
                <h2 className="fw-bold">{t('auth.register.title')}</h2>
              </div>

              {error && (
                <div className="alert alert-danger" role="alert">
                  {error}
                </div>
              )}

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

                {/* 验证码输入框 - 与登录页面保持一致 */}
                {isCaptchaEnabled && (
                  <div className="mb-3 position-relative">
                    <label htmlFor="captcha" className="form-label">{t('auth.register.captcha')}</label>
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
                        {isLoading ? t('common.loading') : t('auth.register.refreshCaptcha')}
                      </button>
                    </div>
                    {captchaImage && (
                      <div className="mt-2 text-center">
                        <img 
                          src={captchaImage} 
                          alt={t('auth.register.captcha')} 
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
                    {isLoading ? t('common.loading') : t('auth.register.submit')}
                  </button>
                </div>
              </form>

              <div className="text-center mt-4">
                <p className="mb-0">
                  {t('auth.register.haveAccount')} <Link to="/login" className="text-decoration-none">{t('auth.register.login')}</Link>
                </p>
              </div>

              <div className="text-center mt-4">
                <div className="d-flex align-items-center justify-content-center">
                  <i className="bi bi-github me-2"></i>
                  <a 
                    href="https://github.com/axfinn/todoIng" 
                    target="_blank" 
                    rel="noopener noreferrer" 
                    className="text-decoration-none"
                  >
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

export default RegisterPage;
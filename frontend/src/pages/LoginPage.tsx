import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { useTranslation } from 'react-i18next';
import type { RootState, AppDispatch } from '../app/store';
import { loginUser } from '../features/auth/authSlice';
import { Link } from 'react-router-dom';

const LoginPage: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  const navigate = useNavigate();
  const { t } = useTranslation();
  const { isAuthenticated, isLoading, error } = useSelector((state: RootState) => state.auth);

  const [formData, setFormData] = useState({
    email: '',
    password: '',
    captcha: '',
  });

  // 检查是否启用验证码功能
  const isCaptchaEnabled = process.env.REACT_APP_ENABLE_CAPTCHA === 'true';

  useEffect(() => {
    if (isAuthenticated) {
      navigate('/dashboard');
    }
  }, [isAuthenticated, navigate]);

  const { email, password, captcha } = formData;

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    // 准备登录数据
    const loginData: { email: string; password: string; captcha?: string } = {
      email,
      password
    };
    
    // 如果启用了验证码，则添加验证码到请求数据中
    if (isCaptchaEnabled) {
      loginData.captcha = captcha;
    }
    
    dispatch(loginUser(loginData));
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
                        // TODO: Add get captcha logic
                      >
                        {t('auth.login.getCaptcha')}
                      </button>
                    </div>
                  </div>
                )}
                
                <div className="d-grid">
                  <button type="submit" className="btn btn-primary btn-lg" disabled={isLoading}>
                    {t('auth.login.submit')}
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
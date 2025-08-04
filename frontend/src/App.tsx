import React, { useState, useEffect } from 'react';
import { Routes, Route, Link } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { useTranslation } from 'react-i18next';
import type { RootState, AppDispatch } from './app/store';
import { logout } from './features/auth/authSlice';

import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import DashboardPage from './pages/DashboardPage';
import ReportsPage from './pages/ReportsPage';
import ProtectedRoute from './components/ProtectedRoute';

const App: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  const { isAuthenticated } = useSelector((state: RootState) => state.auth);
  const { t, i18n } = useTranslation();
  const [currentLanguage, setCurrentLanguage] = useState('en');
  const [githubStats, setGithubStats] = useState({ stars: 0, forks: 0 });
  const [isGitHubLoading, setIsGitHubLoading] = useState(false);
  const [githubError, setGithubError] = useState(false);
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

  useEffect(() => {
    // 初始化时设置当前语言
    setCurrentLanguage(i18n.language);
  }, [i18n.language]);

  useEffect(() => {
    // 获取GitHub项目统计信息
    setIsGitHubLoading(true);
    setGithubError(false);
    
    fetch('https://api.github.com/repos/axfinn/todoIng')
      .then(response => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.json();
      })
      .then(data => {
        setGithubStats({
          stars: data.stargazers_count || 0,
          forks: data.forks_count || 0
        });
        setIsGitHubLoading(false);
      })
      .catch(error => {
        console.error('Failed to fetch GitHub stats:', error);
        setIsGitHubLoading(false);
        setGithubError(true);
      });
  }, []);

  const handleLogout = () => {
    dispatch(logout());
  };

  const changeLanguage = (lng: string) => {
    i18n.changeLanguage(lng);
    setCurrentLanguage(lng);
  };

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
      <header>
        <nav className="navbar navbar-expand-lg navbar-dark bg-primary shadow-sm">
          <div className="container">
            <Link className="navbar-brand fw-bold d-flex align-items-center" to="/">
              <i className="bi bi-check2-circle me-2"></i>
              todoIng
            </Link>
            <button
              className="navbar-toggler"
              type="button"
              data-bs-toggle="collapse"
              data-bs-target="#navbarNav"
              aria-controls="navbarNav"
              aria-expanded="false"
              aria-label="Toggle navigation"
            >
              <span className="navbar-toggler-icon"></span>
            </button>
            <div className="collapse navbar-collapse" id="navbarNav">
              <ul className="navbar-nav me-auto mb-2 mb-lg-0">
                <li className="nav-item">
                  <Link className="nav-link" to="/">
                    {t('nav.home')}
                  </Link>
                </li>
                {isAuthenticated && (
                  <li className="nav-item">
                    <Link className="nav-link" to="/dashboard">
                      {t('nav.dashboard')}
                    </Link>
                  </li>
                )}
                {isAuthenticated && (
                  <li className="nav-item">
                    <Link className="nav-link" to="/reports">
                      {t('nav.reports')}
                    </Link>
                  </li>
                )}
              </ul>
              <ul className="navbar-nav mb-2 mb-lg-0">
                <li className="nav-item dropdown">
                  <a className="btn btn-outline-light dropdown-toggle me-2" href="#" role="button" data-bs-toggle="dropdown" aria-expanded="false">
                    {currentLanguage === 'en' ? 'English' : '中文'}
                  </a>
                  <ul className="dropdown-menu">
                    <li><a className="dropdown-item" href="#" onClick={(e) => { e.preventDefault(); changeLanguage('en'); }}>English</a></li>
                    <li><hr className="dropdown-divider" /></li>
                    <li><a className="dropdown-item" href="#" onClick={(e) => { e.preventDefault(); changeLanguage('zh'); }}>中文</a></li>
                  </ul>
                </li>
                {isAuthenticated ? (
                  <li className="nav-item">
                    <button className="btn btn-outline-light" onClick={handleLogout}>
                      <i className="bi bi-box-arrow-right me-1"></i>
                      {t('nav.logout')}
                    </button>
                  </li>
                ) : (
                  <>
                    <li className="nav-item">
                      <Link className="nav-link" to="/register">
                        {t('nav.register')}
                      </Link>
                    </li>
                    <li className="nav-item">
                      <Link className="nav-link" to="/login">
                        {t('nav.login')}
                      </Link>
                    </li>
                  </>
                )}
              </ul>
            </div>
          </div>
        </nav>
      </header>
      <main className="flex-grow-1">
        <Routes>
          <Route path="/" element={
            <div className="container py-5">
              <div className="row justify-content-center">
                <div className="col-md-8 text-center" style={{ 
                  backgroundColor: 'rgba(255, 255, 255, 0.8)',
                  borderRadius: '10px',
                  padding: '2rem',
                  boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)'
                }}>
                  <h1 className="display-4 fw-bold mb-4">{t('home.title')}</h1>
                  <p className="lead mb-4">
                    {t('home.description')}
                  </p>
                  {!isAuthenticated && (
                    <div className="d-grid gap-3 d-sm-flex justify-content-sm-center">
                      <Link to="/register" className="btn btn-primary btn-lg px-4 gap-3">
                        {t('home.getStarted')}
                      </Link>
                      <Link to="/login" className="btn btn-outline-primary btn-lg px-4">
                        {t('home.login')}
                      </Link>
                    </div>
                  )}
                  {isAuthenticated && (
                    <div className="d-grid gap-3 d-sm-flex justify-content-sm-center">
                      <Link to="/dashboard" className="btn btn-primary btn-lg px-4 gap-3">
                        {t('nav.dashboard')}
                      </Link>
                    </div>
                  )}
                  
                  <div className="mt-5">
                    <div className="d-flex align-items-center justify-content-center mb-3">
                      <i className="bi bi-github me-2"></i>
                      <a href="https://github.com/axfinn/todoIng" target="_blank" rel="noopener noreferrer" className="text-decoration-none">
                        Fork me on GitHub
                      </a>
                    </div>
                    <div className="d-flex justify-content-center gap-3">
                      <>
                        <button 
                          className="btn btn-outline-dark btn-sm d-flex align-items-center"
                          onClick={() => window.open('https://github.com/axfinn/todoIng', '_blank')}
                        >
                          <i className="bi bi-star-fill me-1"></i> 
                          <span>Star</span>
                          {githubStats.stars > 0 && (
                            <span className="badge bg-secondary ms-1">{githubStats.stars}</span>
                          )}
                        </button>
                        <span className="d-flex align-items-center">
                          <i className="bi bi-git me-1"></i> 
                          {githubStats.forks > 0 && (
                            <span className="badge bg-secondary">{githubStats.forks}</span>
                          )}
                        </span>
                      </>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          } />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/dashboard" element={
            <ProtectedRoute>
              <DashboardPage />
            </ProtectedRoute>
          } />
          <Route path="/reports" element={
            <ProtectedRoute>
              <ReportsPage />
            </ProtectedRoute>
          } />
        </Routes>
      </main>
      <footer className="bg-light py-3 mt-auto">
        <div className="container">
          <div className="text-center text-muted">
            <div className="d-flex align-items-center justify-content-center">
              <i className="bi bi-github me-2"></i>
              <a href="https://github.com/axfinn/todoIng" target="_blank" rel="noopener noreferrer" className="text-decoration-none text-muted">
                Fork me on GitHub
              </a>
            </div>
            <div className="mt-2">
              &copy; {new Date().getFullYear()} {t('footer.copyright')}
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default App;
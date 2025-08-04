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

  useEffect(() => {
    // 初始化时设置当前语言
    setCurrentLanguage(i18n.language);
  }, [i18n.language]);

  const handleLogout = () => {
    dispatch(logout());
  };

  const changeLanguage = (lng: string) => {
    i18n.changeLanguage(lng);
    setCurrentLanguage(lng);
  };

  return (
    <div className="App min-vh-100 d-flex flex-column">
      <nav className="navbar navbar-expand-lg navbar-dark bg-primary shadow-sm">
        <div className="container-fluid">
          <Link className="navbar-brand fw-bold d-flex align-items-center" to="/">
            <i className="bi bi-check2-circle me-2"></i>
            <span className="d-none d-sm-inline">todoIng</span>
            <span className="d-inline d-sm-none">todo</span>
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
                <>
                  <li className="nav-item">
                    <Link className="nav-link" to="/dashboard">
                      {t('nav.dashboard')}
                    </Link>
                  </li>
                  <li className="nav-item">
                    <Link className="nav-link" to="/reports">
                      {t('nav.reports')}
                    </Link>
                  </li>
                </>
              )}
            </ul>
            <ul className="navbar-nav d-flex align-items-center mb-2 mb-lg-0">
              <li className="nav-item dropdown me-2 mb-2 mb-lg-0">
                <a 
                  className="btn btn-sm btn-outline-light dropdown-toggle w-100" 
                  href="#"
                  role="button" 
                  data-bs-toggle="dropdown" 
                  aria-expanded="false"
                >
                  {currentLanguage === 'en' ? 'EN' : '中'}
                </a>
                <ul className="dropdown-menu">
                  <li>
                    <a 
                      className="dropdown-item" 
                      href="#"
                      onClick={(e) => {
                        e.preventDefault();
                        changeLanguage('en');
                      }}
                    >
                      English
                    </a>
                  </li>
                  <li>
                    <hr className="dropdown-divider" />
                  </li>
                  <li>
                    <a 
                      className="dropdown-item" 
                      href="#"
                      onClick={(e) => {
                        e.preventDefault();
                        changeLanguage('zh');
                      }}
                    >
                      中文
                    </a>
                  </li>
                </ul>
              </li>
              {!isAuthenticated ? (
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
              ) : (
                <li className="nav-item">
                  <button className="btn btn-sm btn-outline-light d-flex align-items-center" onClick={handleLogout}>
                    <i className="bi bi-box-arrow-right me-1"></i>
                    <span className="d-none d-sm-inline">{t('nav.logout')}</span>
                    <span className="d-inline d-sm-none">{t('nav.logoutShort')}</span>
                  </button>
                </li>
              )}
            </ul>
          </div>
        </div>
      </nav>
      <main className="flex-grow-1">
        <Routes>
          <Route path="/" element={
            <div className="container py-5">
              <div className="row justify-content-center">
                <div className="col-md-8 text-center">
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
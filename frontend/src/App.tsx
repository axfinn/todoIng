import React from 'react';
import { Routes, Route, Link } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import type { RootState, AppDispatch } from './app/store';
import { logout } from './features/auth/authSlice';

import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import DashboardPage from './pages/DashboardPage';

const App: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  const { isAuthenticated } = useSelector((state: RootState) => state.auth);

  const handleLogout = () => {
    dispatch(logout());
  };

  return (
    <div className="App min-vh-100 d-flex flex-column">
      <nav className="navbar navbar-expand-lg navbar-dark bg-primary shadow-sm">
        <div className="container">
          <Link className="navbar-brand fw-bold d-flex align-items-center" to="/">
            <i className="bi bi-check2-circle me-2"></i>
            todoIng
          </Link>
          <div className="collapse navbar-collapse">
            <ul className="navbar-nav me-auto mb-2 mb-lg-0">
              <li className="nav-item">
                <Link className="nav-link" to="/">
                  Home
                </Link>
              </li>
              {isAuthenticated && (
                <li className="nav-item">
                  <Link className="nav-link" to="/dashboard">
                    Dashboard
                  </Link>
                </li>
              )}
            </ul>
            <ul className="navbar-nav mb-2 mb-lg-0">
              {!isAuthenticated ? (
                <>
                  <li className="nav-item">
                    <Link className="nav-link" to="/register">
                      Register
                    </Link>
                  </li>
                  <li className="nav-item">
                    <Link className="nav-link" to="/login">
                      Login
                    </Link>
                  </li>
                </>
              ) : (
                <li className="nav-item">
                  <button className="btn btn-outline-light" onClick={handleLogout}>
                    <i className="bi bi-box-arrow-right me-1"></i>
                    Logout
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
                  <h1 className="display-4 fw-bold mb-4">Welcome to todoIng!</h1>
                  <p className="lead mb-4">
                    Your modern task management solution. Organize your work, boost your productivity, 
                    and achieve your goals with our intuitive task management system.
                  </p>
                  {!isAuthenticated && (
                    <div className="d-grid gap-3 d-sm-flex justify-content-sm-center">
                      <Link to="/register" className="btn btn-primary btn-lg px-4 gap-3">
                        Get Started
                      </Link>
                      <Link to="/login" className="btn btn-outline-primary btn-lg px-4">
                        Login
                      </Link>
                    </div>
                  )}
                  {isAuthenticated && (
                    <div className="d-grid gap-3 d-sm-flex justify-content-sm-center">
                      <Link to="/dashboard" className="btn btn-primary btn-lg px-4 gap-3">
                        Go to Dashboard
                      </Link>
                    </div>
                  )}
                </div>
              </div>
            </div>
          } />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/dashboard" element={<DashboardPage />} />
        </Routes>
      </main>
      <footer className="bg-light py-3 mt-auto">
        <div className="container">
          <div className="text-center text-muted">
            &copy; {new Date().getFullYear()} todoIng - Task Management System
          </div>
        </div>
      </footer>
    </div>
  );
};

export default App;
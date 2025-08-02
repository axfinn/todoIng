import React, { useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { loginUser } from '../features/auth/authSlice';
import type { AppDispatch, RootState } from '../app/store';

const LoginPage: React.FC = () => {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
  });

  const [errors, setErrors] = useState({
    email: '',
    password: '',
  });

  const { email, password } = formData;
  const dispatch = useDispatch<AppDispatch>();
  const { isLoading, error } = useSelector((state: RootState) => state.auth);

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
    setErrors({ ...errors, [e.target.name]: '' }); // Clear error on change
  };

  const validateForm = () => {
    const newErrors = { email: '', password: '' };
    let isValid = true;

    if (!email.trim()) {
      newErrors.email = 'Email is required';
      isValid = false;
    } else if (!/^[\w-.]+@([\w-]+\.)+[\w-]{2,4}$/.test(email)) {
      newErrors.email = 'Please enter a valid email';
      isValid = false;
    }

    if (!password.trim()) {
      newErrors.password = 'Password is required';
      isValid = false;
    }

    setErrors(newErrors);
    return isValid;
  };

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (validateForm()) {
      dispatch(loginUser({ email, password }));
    }
  };

  return (
    <div className="container py-5">
      <div className="row justify-content-center">
        <div className="col-md-6 col-lg-5">
          <div className="card shadow-lg border-0 rounded-3">
            <div className="card-body p-5">
              <div className="text-center mb-5">
                <div className="d-inline-block bg-primary rounded-circle p-3 mb-4">
                  <i className="bi bi-box-arrow-in-right text-white fs-1"></i>
                </div>
                <h2 className="card-title">Welcome Back</h2>
                <p className="text-muted">Sign in to your account</p>
              </div>
              
              {isLoading && (
                <div className="d-flex justify-content-center mb-4">
                  <div className="spinner-border text-primary" role="status">
                    <span className="visually-hidden">Loading...</span>
                  </div>
                </div>
              )}
              
              {error && (
                <div className="alert alert-danger alert-dismissible fade show" role="alert">
                  <i className="bi bi-exclamation-triangle-fill me-2"></i>
                  {error || 'Login failed'}
                  <button type="button" className="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
                </div>
              )}
              
              <form onSubmit={onSubmit} noValidate>
                <div className="form-group mb-4">
                  <div className="input-group">
                    <span className="input-group-text">
                      <i className="bi bi-envelope"></i>
                    </span>
                    <input
                      type="email"
                      className={`form-control form-control-lg ${errors.email ? 'is-invalid' : ''}`}
                      placeholder="Email Address"
                      name="email"
                      value={email}
                      onChange={onChange}
                      required
                    />
                  </div>
                  {errors.email && <div className="invalid-feedback d-block">{errors.email}</div>}
                </div>
                
                <div className="form-group mb-4">
                  <div className="input-group">
                    <span className="input-group-text">
                      <i className="bi bi-lock"></i>
                    </span>
                    <input
                      type="password"
                      className={`form-control form-control-lg ${errors.password ? 'is-invalid' : ''}`}
                      placeholder="Password"
                      name="password"
                      value={password}
                      onChange={onChange}
                      minLength={6}
                      required
                    />
                  </div>
                  {errors.password && <div className="invalid-feedback d-block">{errors.password}</div>}
                </div>
                
                <div className="d-grid mb-4">
                  <button 
                    type="submit" 
                    className="btn btn-primary btn-lg rounded-pill" 
                    disabled={isLoading}
                  >
                    {isLoading ? (
                      <>
                        <span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                        Signing in...
                      </>
                    ) : (
                      <>
                        <i className="bi bi-box-arrow-in-right me-2"></i>
                        Sign In
                      </>
                    )}
                  </button>
                </div>
                
                <div className="text-center">
                  <p className="mb-0 text-muted">
                    Don't have an account?{' '}
                    <a href="/register" className="text-decoration-none">
                      Sign up here
                    </a>
                  </p>
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
import React, { useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { registerUser } from '../features/auth/authSlice';
import type { AppDispatch, RootState } from '../app/store';

const RegisterPage: React.FC = () => {
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    password2: '',
  });

  const [errors, setErrors] = useState({
    username: '',
    email: '',
    password: '',
    password2: '',
  });

  const { username, email, password, password2 } = formData;
  const dispatch = useDispatch<AppDispatch>();
  const { isLoading, error } = useSelector((state: RootState) => state.auth);

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
    setErrors({ ...errors, [e.target.name]: '' }); // Clear error on change
  };

  const validateForm = () => {
    const newErrors = { username: '', email: '', password: '', password2: '' };
    let isValid = true;

    if (!username.trim()) {
      newErrors.username = 'Username is required';
      isValid = false;
    }

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
    } else if (password.length < 6) {
      newErrors.password = 'Password must be at least 6 characters';
      isValid = false;
    }

    if (!password2.trim()) {
      newErrors.password2 = 'Confirm password is required';
      isValid = false;
    } else if (password !== password2) {
      newErrors.password2 = 'Passwords do not match';
      isValid = false;
    }

    setErrors(newErrors);
    return isValid;
  };

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (validateForm()) {
      dispatch(registerUser({ username, email, password }));
    }
  };

  return (
    <div className="container py-5">
      <div className="row justify-content-center">
        <div className="col-md-6 col-lg-5">
          <div className="card shadow-lg border-0 rounded-3">
            <div className="card-body p-5">
              <div className="text-center mb-5">
                <div className="d-inline-block bg-success rounded-circle p-3 mb-4">
                  <i className="bi bi-person-plus text-white fs-1"></i>
                </div>
                <h2 className="card-title">Create Account</h2>
                <p className="text-muted">Join us today</p>
              </div>
              
              {isLoading && (
                <div className="d-flex justify-content-center mb-4">
                  <div className="spinner-border text-success" role="status">
                    <span className="visually-hidden">Loading...</span>
                  </div>
                </div>
              )}
              
              {error && (
                <div className="alert alert-danger alert-dismissible fade show" role="alert">
                  <i className="bi bi-exclamation-triangle-fill me-2"></i>
                  {error || 'Registration failed'}
                  <button type="button" className="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
                </div>
              )}
              
              <form onSubmit={onSubmit} noValidate>
                <div className="form-group mb-4">
                  <div className="input-group">
                    <span className="input-group-text">
                      <i className="bi bi-person"></i>
                    </span>
                    <input
                      type="text"
                      className={`form-control form-control-lg ${errors.username ? 'is-invalid' : ''}`}
                      placeholder="Username"
                      name="username"
                      value={username}
                      onChange={onChange}
                      required
                    />
                  </div>
                  {errors.username && <div className="invalid-feedback d-block">{errors.username}</div>}
                </div>
                
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
                
                <div className="form-group mb-4">
                  <div className="input-group">
                    <span className="input-group-text">
                      <i className="bi bi-lock-fill"></i>
                    </span>
                    <input
                      type="password"
                      className={`form-control form-control-lg ${errors.password2 ? 'is-invalid' : ''}`}
                      placeholder="Confirm Password"
                      name="password2"
                      value={password2}
                      onChange={onChange}
                      minLength={6}
                      required
                    />
                  </div>
                  {errors.password2 && <div className="invalid-feedback d-block">{errors.password2}</div>}
                </div>
                
                <div className="d-grid mb-4">
                  <button 
                    type="submit" 
                    className="btn btn-success btn-lg rounded-pill" 
                    disabled={isLoading}
                  >
                    {isLoading ? (
                      <>
                        <span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                        Creating Account...
                      </>
                    ) : (
                      <>
                        <i className="bi bi-person-plus me-2"></i>
                        Create Account
                      </>
                    )}
                  </button>
                </div>
                
                <div className="text-center">
                  <p className="mb-0 text-muted">
                    Already have an account?{' '}
                    <a href="/login" className="text-decoration-none">
                      Sign in here
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

export default RegisterPage;
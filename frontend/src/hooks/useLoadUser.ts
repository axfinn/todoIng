import { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { loadUser, selectAuth } from '../features/auth/authSlice';
import type { AppDispatch } from '../app/store';

export const useLoadUser = () => {
  const dispatch = useDispatch<AppDispatch>();
  const { isAuthenticated, token } = useSelector(selectAuth);

  useEffect(() => {
    // 如果有token但还没有用户信息，则加载用户信息
    if (token && !isAuthenticated) {
      dispatch(loadUser());
    }
  }, [dispatch, token, isAuthenticated]);
};
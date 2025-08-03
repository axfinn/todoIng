import { createSlice, createAsyncThunk, type PayloadAction } from '@reduxjs/toolkit';
import api from '../../config/api';

interface AuthState {
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  user: Record<string, unknown> | null;
  error: string | null;
}

const token = localStorage.getItem('token');
const initialState: AuthState = {
  token: token,
  isAuthenticated: !!token,
  isLoading: false,
  user: null,
  error: null,
};

interface AuthResponse {
  token: string;
}

export const registerUser = createAsyncThunk<AuthResponse, Record<string, unknown>, { rejectValue: string }>(
  'auth/registerUser',
  async (userData: Record<string, unknown>, { rejectWithValue }) => {
    try {
      const res = await api.post('/auth/register', userData);
      return res.data as AuthResponse;
    } catch (err: any) {
      if (err.response && err.response.data) {
        return rejectWithValue(err.response.data.msg || 'Registration failed');
      }
      return rejectWithValue('Registration failed');
    }
  }
);

export const loginUser = createAsyncThunk<AuthResponse, Record<string, unknown>, { rejectValue: string }>(
  'auth/loginUser',
  async (userData: Record<string, unknown>, { rejectWithValue }) => {
    try {
      const res = await api.post('/auth/login', userData);
      return res.data as AuthResponse;
    }
    catch (err: any) {
      if (err.response && err.response.data) {
        return rejectWithValue(err.response.data.msg || 'Login failed');
      }
      return rejectWithValue('Login failed');
    }
  }
);

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    logout: (state) => {
      localStorage.removeItem('token');
      state.token = null;
      state.isAuthenticated = false;
      state.user = null;
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(registerUser.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(registerUser.fulfilled, (state, action: PayloadAction<AuthResponse>) => {
        localStorage.setItem('token', action.payload.token);
        state.token = action.payload.token;
        state.isAuthenticated = true;
        state.isLoading = false;
        state.error = null;
      })
      .addCase(registerUser.rejected, (state, action: PayloadAction<any>) => {
        localStorage.removeItem('token');
        state.token = null;
        state.isAuthenticated = false;
        state.isLoading = false;
        state.error = action.payload as string;
      })
      .addCase(loginUser.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(loginUser.fulfilled, (state, action: PayloadAction<AuthResponse>) => {
        localStorage.setItem('token', action.payload.token);
        state.token = action.payload.token;
        state.isAuthenticated = true;
        state.isLoading = false;
        state.error = null;
      })
      .addCase(loginUser.rejected, (state, action: PayloadAction<any>) => {
        localStorage.removeItem('token');
        state.token = null;
        state.isAuthenticated = false;
        state.isLoading = false;
        state.error = action.payload as string;
      });
  },
});

export const { logout } = authSlice.actions;
export default authSlice.reducer;
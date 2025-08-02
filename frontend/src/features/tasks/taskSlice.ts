import { createSlice, createAsyncThunk, type PayloadAction } from '@reduxjs/toolkit';
import axios from 'axios';
import API_BASE_URL from '../../config/api';

interface Task {
  _id: string;
  title: string;
  description?: string;
  status: 'To Do' | 'In Progress' | 'Done';
  priority: 'Low' | 'Medium' | 'High';
  assignee?: string; // User ID
  createdBy: string; // User ID
  createdAt: string;
  updatedAt: string;
}

interface TaskState {
  tasks: Task[];
  isLoading: boolean;
  error: string | null;
}

const initialState: TaskState = {
  tasks: [],
  isLoading: false,
  error: null,
};

// Async Thunks for Task operations
export const fetchTasks = createAsyncThunk<Task[], void, { rejectValue: string }>(
  'tasks/fetchTasks',
  async (_, { rejectWithValue }) => {
    try {
      const res = await axios.get(`${API_BASE_URL}/tasks`);
      return res.data as Task[];
    } catch (err: any) {
      if (err.response && err.response.data) {
        return rejectWithValue(err.response.data.msg || 'Failed to fetch tasks');
      }
      return rejectWithValue('Failed to fetch tasks');
    }
  }
);

export const createTask = createAsyncThunk<Task, Partial<Task>, { rejectValue: string }>(
  'tasks/createTask',
  async (taskData: Partial<Task>, { rejectWithValue }) => {
    try {
      const res = await axios.post(`${API_BASE_URL}/tasks`, taskData);
      return res.data as Task;
    } catch (err: any) {
      if (err.response && err.response.data) {
        return rejectWithValue(err.response.data.msg || 'Failed to create task');
      }
      return rejectWithValue('Failed to create task');
    }
  }
);

export const updateTask = createAsyncThunk<Task, { id: string; taskData: Partial<Task> }, { rejectValue: string }>(
  'tasks/updateTask',
  async ({ id, taskData }: { id: string; taskData: Partial<Task> }, { rejectWithValue }) => {
    try {
      const res = await axios.put(`${API_BASE_URL}/tasks/${id}`, taskData);
      return res.data as Task;
    } catch (err: any) {
      if (err.response && err.response.data) {
        return rejectWithValue(err.response.data.msg || 'Failed to update task');
      }
      return rejectWithValue('Failed to update task');
    }
  }
);

export const deleteTask = createAsyncThunk<string, string, { rejectValue: string }>(
  'tasks/deleteTask',
  async (id: string, { rejectWithValue }) => {
    try {
      await axios.delete(`${API_BASE_URL}/tasks/${id}`);
      return id;
    } catch (err: any) {
      if (err.response && err.response.data) {
        return rejectWithValue(err.response.data.msg || 'Failed to delete task');
      }
      return rejectWithValue('Failed to delete task');
    }
  }
);

const taskSlice = createSlice({
  name: 'tasks',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchTasks.pending, (state) => {
        state.isLoading = true;
      })
      .addCase(fetchTasks.fulfilled, (state, action: PayloadAction<Task[]>) => {
        state.isLoading = false;
        state.tasks = action.payload;
      })
      .addCase(fetchTasks.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      .addCase(createTask.fulfilled, (state, action: PayloadAction<Task>) => {
        state.tasks.push(action.payload);
      })
      .addCase(updateTask.fulfilled, (state, action: PayloadAction<Task>) => {
        const index = state.tasks.findIndex((task) => task._id === action.payload._id);
        if (index !== -1) {
          state.tasks[index] = action.payload;
        }
      })
      .addCase(deleteTask.fulfilled, (state, action: PayloadAction<string>) => {
        state.tasks = state.tasks.filter((task) => task._id !== action.payload);
      });
  },
});

export type { Task };
export default taskSlice.reducer;
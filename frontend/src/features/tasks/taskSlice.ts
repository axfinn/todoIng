import { createSlice, createAsyncThunk, type PayloadAction } from '@reduxjs/toolkit';
import api from '../../config/api';

export interface TaskComment {
  text: string;
  createdBy?: string;
  createdAt: string;
}

export interface Task {
  _id: string;
  title: string;
  description: string;
  status: 'To Do' | 'In Progress' | 'Done';
  priority: 'Low' | 'Medium' | 'High';
  assignee?: string;
  comments?: TaskComment[];
  createdAt: string;
  updatedAt: string;
  deadline?: string | null;
  scheduledDate?: string | null;
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
      const res = await api.get('/tasks');
      return res.data as Task[];
    } catch (err: any) {
      if (err.response && err.response.data) {
        return rejectWithValue(err.response.data.msg || 'Failed to fetch tasks');
      }
      return rejectWithValue('Failed to fetch tasks');
    }
  }
);

export const createTask = createAsyncThunk<Task, Omit<Task, '_id' | 'createdAt' | 'updatedAt' | 'comments'>, { rejectValue: string }>(
  'tasks/createTask',
  async (taskData: Omit<Task, '_id' | 'createdAt' | 'updatedAt' | 'comments'>, { rejectWithValue }) => {
    try {
      const res = await api.post('/tasks', taskData);
      return res.data as Task;
    } catch (err: any) {
      if (err.response && err.response.data) {
        return rejectWithValue(err.response.data.msg || 'Failed to create task');
      }
      return rejectWithValue('Failed to create task');
    }
  }
);

export const updateTask = createAsyncThunk<Task, Partial<Task> & { _id: string }, { rejectValue: string }>(
  'tasks/updateTask',
  async (taskData: Partial<Task> & { _id: string }, { rejectWithValue }) => {
    try {
      const { _id, ...taskUpdate } = taskData;
      const res = await api.put(`/tasks/${_id}`, taskUpdate);
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
  async (taskId: string, { rejectWithValue }) => {
    try {
      await api.delete(`/tasks/${taskId}`);
      return taskId;
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
        state.error = action.error.message || 'Failed to fetch tasks';
      })
      .addCase(createTask.fulfilled, (state, action: PayloadAction<Task>) => {
        state.tasks.unshift(action.payload);
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

export type { TaskState };
export default taskSlice.reducer;
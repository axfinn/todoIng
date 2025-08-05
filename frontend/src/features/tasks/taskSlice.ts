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

// Define the specific fields we want to allow updating
interface UpdateTaskFields {
  title?: string;
  description?: string;
  status?: 'To Do' | 'In Progress' | 'Done';
  assignee?: string;
  deadline?: string | null;
  scheduledDate?: string | null;
}

export const updateTask = createAsyncThunk<Task, { _id: string } & UpdateTaskFields, { rejectValue: string }>(
  'tasks/updateTask',
  async (taskData: { _id: string } & UpdateTaskFields, { rejectWithValue }) => {
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

// Async thunk for exporting tasks
export const exportTasks = createAsyncThunk<{ data: string, filename: string }, void, { rejectValue: string }>(
  'tasks/exportTasks',
  async (_, { rejectWithValue }) => {
    try {
      const response = await api.get('/tasks/export/all', {
        responseType: 'blob'
      });
      
      // Get filename from response headers
      const contentDisposition = response.headers['content-disposition'];
      let filename = `todoing-backup-${new Date().toISOString().slice(0, 10)}.json`;
      
      // Improved filename parsing with better edge case handling
      if (contentDisposition) {
        // First try to match UTF-8 encoded filename
        const utf8FilenameRegex = /filename\*=UTF-8''([\w%\-\.]+)/i;
        const asciiFilenameRegex = /filename="?([^"]+)"?/i;
        
        const utf8Matches = contentDisposition.match(utf8FilenameRegex);
        if (utf8Matches && utf8Matches[1]) {
          filename = decodeURIComponent(utf8Matches[1]);
        } else {
          // Fall back to ASCII filename
          const asciiMatches = contentDisposition.match(asciiFilenameRegex);
          if (asciiMatches && asciiMatches[1]) {
            filename = asciiMatches[1];
          }
        }
      }
      
      // Convert blob to text
      const data = await response.data.text();
      
      return {
        data,
        filename
      };
    } catch (err: any) {
      if (err.response && err.response.data) {
        return rejectWithValue(err.response.data.msg || 'Failed to export tasks');
      }
      return rejectWithValue('Failed to export tasks');
    }
  }
);

// Async thunk for importing tasks
export const importTasks = createAsyncThunk<{ imported: number, errors: any[] }, File, { rejectValue: string }>(
  'tasks/importTasks',
  async (file: File, { rejectWithValue }) => {
    try {
      const formData = new FormData();
      formData.append('file', file);
      
      // First read the file content
      const fileContent = await file.text();
      const tasksData = JSON.parse(fileContent);
      
      // Send the data to the server
      const response = await api.post('/tasks/import', { tasks: tasksData });
      
      return {
        imported: response.data.imported,
        errors: response.data.errors
      };
    } catch (err: any) {
      if (err.response && err.response.data) {
        return rejectWithValue(err.response.data.msg || 'Failed to import tasks');
      }
      return rejectWithValue('Failed to import tasks');
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
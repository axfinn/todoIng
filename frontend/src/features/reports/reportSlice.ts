import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';
import type { RootState } from '../../app/store';
import api from '../../config/api';

// 报告接口
export interface Report {
  _id: string;
  userId: string;
  type: 'daily' | 'weekly' | 'monthly';
  period: string;
  title: string;
  content: string;
  polishedContent: string | null;
  tasks: string[];
  statistics: {
    totalTasks: number;
    completedTasks: number;
    inProgressTasks: number;
    overdueTasks: number;
    completionRate: number;
  };
  createdAt: string;
  updatedAt: string;
}

// 报告生成请求接口
export interface GenerateReportRequest {
  type: 'daily' | 'weekly' | 'monthly';
  period: string;
  startDate: string;
  endDate: string;
}

// AI润色请求接口
export interface PolishReportRequest {
  apiKey: string;
  model?: string;
  apiUrl?: string;
  provider?: string;
}

// 状态接口
interface ReportState {
  reports: Report[];
  currentReport: Report | null;
  loading: boolean;
  error: string | null;
}

// 初始状态
const initialState: ReportState = {
  reports: [],
  currentReport: null,
  loading: false,
  error: null,
}

// 异步action - 获取所有报告
export const fetchReports = createAsyncThunk<
  Report[],
  void,
  { state: RootState, rejectValue: string }
>(
  'reports/fetchReports',
  async (_, { rejectWithValue }) => {
    try {
      const res = await api.get('/reports');
      return res.data;
    } catch (err: any) {
      return rejectWithValue(err.response?.data?.msg || 'Failed to fetch reports');
    }
  }
);

// 异步action - 获取特定报告
export const fetchReport = createAsyncThunk<
  Report,
  string,
  { state: RootState, rejectValue: string }
>(
  'reports/fetchReport',
  async (reportId, { rejectWithValue }) => {
    try {
      const res = await api.get(`/reports/${reportId}`);
      return res.data;
    } catch (err: any) {
      return rejectWithValue(err.response?.data?.msg || 'Failed to fetch report');
    }
  }
);

// 异步action - 生成报告
export const generateReport = createAsyncThunk<
  Report,
  GenerateReportRequest,
  { state: RootState, rejectValue: string }
>(
  'reports/generateReport',
  async (reportData, { rejectWithValue }) => {
    try {
      const res = await api.post('/reports/generate', reportData);
      return res.data;
    } catch (err: any) {
      return rejectWithValue(err.response?.data?.msg || 'Failed to generate report');
    }
  }
);

// 异步action - AI润色报告
export const polishReport = createAsyncThunk<
  Report,
  { reportId: string, polishData: PolishReportRequest },
  { state: RootState, rejectValue: string }
>(
  'reports/polishReport',
  async ({ reportId, polishData }, { rejectWithValue }) => {
    try {
      const res = await api.post(`/reports/${reportId}/polish`, polishData);
      return res.data;
    } catch (err: any) {
      return rejectWithValue(err.response?.data?.msg || 'Failed to polish report');
    }
  }
);

// 异步action - 删除报告
export const deleteReport = createAsyncThunk<
  string, // 返回被删除的报告ID
  string, // 报告ID参数
  { state: RootState, rejectValue: string }
>(
  'reports/deleteReport',
  async (reportId, { rejectWithValue }) => {
    try {
      await api.delete(`/reports/${reportId}`);
      return reportId; // 返回被删除的报告ID
    } catch (err: any) {
      return rejectWithValue(err.response?.data?.msg || 'Failed to delete report');
    }
  }
);

// 异步action - 导出报告
export const exportReport = createAsyncThunk<
  { data: string, filename: string },
  { reportId: string, format: string },
  { state: RootState, rejectValue: string }
>(
  'reports/exportReport',
  async ({ reportId, format }, { rejectWithValue }) => {
    try {
      const response = await api.get(`/reports/${reportId}/export/${format}`, {
        responseType: 'blob'
      });
      
      // 从响应头中获取文件名
      const contentDisposition = response.headers['content-disposition'];
      let filename = `report.${format}`;
      if (contentDisposition) {
        const filenameMatch = contentDisposition.match(/filename="?(.+)"?/);
        if (filenameMatch && filenameMatch[1]) {
          filename = filenameMatch[1];
        }
      }
      
      // 将blob转换为文本
      const data = await response.data.text();
      
      return {
        data,
        filename
      };
    } catch (err: any) {
      return rejectWithValue(err.response?.data?.msg || 'Failed to export report');
    }
  }
);

// 创建切片
export const reportSlice = createSlice({
  name: 'reports',
  initialState,
  reducers: {
    clearCurrentReport: (state) => {
      state.currentReport = null;
    },
    clearError: (state) => {
      state.error = null;
    }
  },
  extraReducers: (builder) => {
    // 获取所有报告
    builder
      .addCase(fetchReports.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchReports.fulfilled, (state, action: PayloadAction<Report[]>) => {
        state.loading = false;
        state.reports = action.payload;
      })
      .addCase(fetchReports.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload || 'Failed to fetch reports';
      })
    
    // 获取特定报告
    .addCase(fetchReport.pending, (state) => {
      state.loading = true;
      state.error = null;
    })
    .addCase(fetchReport.fulfilled, (state, action: PayloadAction<Report>) => {
      state.loading = false;
      state.currentReport = action.payload;
    })
    .addCase(fetchReport.rejected, (state, action) => {
      state.loading = false;
      state.error = action.payload || 'Failed to fetch report';
    })
    
    // 生成报告
    .addCase(generateReport.pending, (state) => {
      state.loading = true;
      state.error = null;
    })
    .addCase(generateReport.fulfilled, (state, action: PayloadAction<Report>) => {
      state.loading = false;
      state.reports.unshift(action.payload);
      state.currentReport = action.payload;
    })
    .addCase(generateReport.rejected, (state, action) => {
      state.loading = false;
      state.error = action.payload || 'Failed to generate report';
    })
    
    // AI润色报告
    .addCase(polishReport.pending, (state) => {
      state.loading = true;
      state.error = null;
    })
    .addCase(polishReport.fulfilled, (state, action: PayloadAction<Report>) => {
      state.loading = false;
      // 更新当前报告的润色内容
      if (state.currentReport) {
        state.currentReport.polishedContent = action.payload.polishedContent;
      }
      // 更新报告列表中的润色内容
      const index = state.reports.findIndex(report => report._id === action.payload._id);
      if (index !== -1) {
        state.reports[index].polishedContent = action.payload.polishedContent;
      }
    })
    .addCase(polishReport.rejected, (state, action) => {
      state.loading = false;
      state.error = action.payload || 'Failed to polish report';
    })
    
    // 删除报告
    .addCase(deleteReport.pending, (state) => {
      state.loading = true;
      state.error = null;
    })
    .addCase(deleteReport.fulfilled, (state, action: PayloadAction<string>) => {
      state.loading = false;
      // 从报告列表中移除被删除的报告
      state.reports = state.reports.filter(report => report._id !== action.payload);
      // 如果当前查看的报告被删除，清除当前报告
      if (state.currentReport && state.currentReport._id === action.payload) {
        state.currentReport = null;
      }
    })
    .addCase(deleteReport.rejected, (state, action) => {
      state.loading = false;
      state.error = action.payload || 'Failed to delete report';
    })
    
    // 导出报告
    .addCase(exportReport.pending, (state) => {
      state.loading = true;
      state.error = null;
    })
    .addCase(exportReport.fulfilled, (state) => {
      state.loading = false;
    })
    .addCase(exportReport.rejected, (state, action) => {
      state.loading = false;
      state.error = action.payload || 'Failed to export report';
    });
  },
});

// 导出actions
export const { clearCurrentReport, clearError } = reportSlice.actions;

// 导出selectors
export const selectReports = (state: RootState) => state.reports.reports;
export const selectCurrentReport = (state: RootState) => state.reports.currentReport;
export const selectReportsLoading = (state: RootState) => state.reports.loading;
export const selectReportsError = (state: RootState) => state.reports.error;

export default reportSlice.reducer;
import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useDispatch, useSelector } from 'react-redux';
import { 
  fetchReports, 
  generateReport, 
  fetchReport,
  polishReport,
  exportReport,
  deleteReport,
  selectReports,
  selectCurrentReport,
  selectReportsLoading,
  selectReportsError
} from '../features/reports/reportSlice';
import type { AppDispatch } from '../app/store';

const ReportsPage: React.FC = () => {
  const { t } = useTranslation();
  const dispatch = useDispatch<AppDispatch>();
  
  const reports = useSelector(selectReports);
  const currentReport = useSelector(selectCurrentReport);
  const loading = useSelector(selectReportsLoading);
  const error = useSelector(selectReportsError);
  
  const [reportType, setReportType] = useState<'daily' | 'weekly' | 'monthly'>('daily');
  const [reportDate, setReportDate] = useState<string>(new Date().toISOString().split('T')[0]);
  const [apiKey, setApiKey] = useState<string>('');
  const [apiUrl, setApiUrl] = useState<string>('https://api.openai.com/v1/chat/completions');
  const [provider, setProvider] = useState<string>('openai');
  const [model, setModel] = useState<string>('gpt-3.5-turbo');
  const [showGenerateForm, setShowGenerateForm] = useState<boolean>(true);
  const [previewMode, setPreviewMode] = useState<'text' | 'preview'>('text');
  
  // 获取所有报告
  useEffect(() => {
    dispatch(fetchReports());
  }, [dispatch]);
  
  // 处理生成报告
  const handleGenerateReport = () => {
    const startDate = new Date(reportDate);
    let endDate = new Date(reportDate);
    
    switch (reportType) {
      case 'daily':
        endDate.setHours(23, 59, 59, 999);
        break;
      case 'weekly':
        // 设置为周日
        const day = startDate.getDay();
        const diff = startDate.getDate() - day;
        startDate.setDate(diff);
        endDate = new Date(startDate);
        endDate.setDate(startDate.getDate() + 6);
        endDate.setHours(23, 59, 59, 999);
        break;
      case 'monthly':
        startDate.setDate(1);
        endDate = new Date(startDate.getFullYear(), startDate.getMonth() + 1, 0);
        endDate.setHours(23, 59, 59, 999);
        break;
    }
    
    const period = reportDate;
    
    dispatch(generateReport({
      type: reportType,
      period,
      startDate: startDate.toISOString(),
      endDate: endDate.toISOString()
    }));
  };
  
  // 处理选择报告
  const handleSelectReport = (reportId: string) => {
    dispatch(fetchReport(reportId));
  };
  
  // 处理AI润色
  const handlePolishReport = () => {
    if (!currentReport) return;
    if (!apiKey) {
      alert('Please enter your API key');
      return;
    }
    
    dispatch(polishReport({
      reportId: currentReport._id,
      polishData: {
        apiKey,
        model,
        apiUrl,
        provider
      }
    }));
  };
  
  // 处理导出报告
  const handleExportReport = (format: string) => {
    if (!currentReport) return;
    
    dispatch(exportReport({
      reportId: currentReport._id,
      format
    })).then((action: any) => {
      if (exportReport.fulfilled.match(action)) {
        const { data, filename } = action.payload;
        const blob = new Blob([data], { type: format === 'md' ? 'text/markdown' : 'text/plain' });
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        window.URL.revokeObjectURL(url);
      }
    });
  };
  
  // 删除报告
  const handleDeleteReport = (reportId: string) => {
    if (window.confirm(t('reports.deleteConfirm'))) {
      dispatch(deleteReport(reportId)).then((action) => {
        if (deleteReport.fulfilled.match(action)) {
          // 成功删除后刷新报告列表
          dispatch(fetchReports());
          alert(t('reports.deleteSuccess'));
        } else if (deleteReport.rejected.match(action)) {
          // 显示错误信息
          alert(action.payload || t('reports.deleteFailed'));
        }
      });
    }
  };
  
  // 渲染Markdown内容
  const renderMarkdown = (content: string) => {
    if (!content) return '';
    
    return content
      .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
      .replace(/\*(.*?)\*/g, '<em>$1</em>')
      .replace(/^# (.*$)/gm, '<h1>$1</h1>')
      .replace(/^## (.*$)/gm, '<h2>$1</h2>')
      .replace(/^### (.*$)/gm, '<h3>$1</h3>')
      .replace(/^#### (.*$)/gm, '<h4>$1</h4>')
      .replace(/^- (.*$)/gm, '<li>$1</li>')
      .replace(/<li>(.*?)<\/li>/gs, '<ul>$&</ul>')
      .replace(/^\d+\. (.*$)/gm, '<li>$1</li>')
      .replace(/<li>(.*?)<\/li>/gs, '<ol>$&</ol>')
      .replace(/\n/g, '<br />');
  };
  
  const reportsLoading = loading;
  const reportsError = error;
  const isGenerating = loading;
  const isPolishing = loading;
  
  return (
    <div className="container-fluid">
      <div className="row">
        <div className="col-12">
          <h2 className="my-4">{t('reports.title')}</h2>
        </div>
      </div>
      
      {/* 报告生成表单 */}
      <div className="row mb-4">
        <div className="col-12">
          <div className="card">
            <div className="card-header">
              <h5 className="mb-0">{t('reports.generateNew')}</h5>
            </div>
            <div className="card-body">
              <div className="row">
                <div className="col-md-3 mb-3">
                  <label className="form-label">{t('reports.type')}</label>
                  <select 
                    className="form-select"
                    value={reportType}
                    onChange={(e) => setReportType(e.target.value as any)}
                  >
                    <option value="daily">{t('reports.daily')}</option>
                    <option value="weekly">{t('reports.weekly')}</option>
                    <option value="monthly">{t('reports.monthly')}</option>
                  </select>
                </div>
                <div className="col-md-3 mb-3">
                  <label className="form-label">{t('reports.date')}</label>
                  <input
                    type="date"
                    className="form-control"
                    value={reportDate}
                    onChange={(e) => setReportDate(e.target.value)}
                  />
                </div>
                <div className="col-md-3 mb-3 d-flex align-items-end">
                  <button
                    className="btn btn-primary"
                    onClick={handleGenerateReport}
                    disabled={isGenerating}
                  >
                    {isGenerating ? (
                      <>
                        <span className="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span>
                        {t('common.loading')}
                      </>
                    ) : (
                      t('reports.generate')
                    )}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <div className="row">
        {/* 左侧报告列表 - 响应式调整 */}
        <div className="col-md-4 reports-sidebar">
          <div className="card">
            <div className="card-header d-flex justify-content-between align-items-center">
              <span>{t('reports.list')}</span>
            </div>
            <div className="card-body p-0">
              {reportsLoading ? (
                <div className="text-center p-3">
                  <div className="spinner-border" role="status">
                    <span className="visually-hidden">{t('common.loading')}</span>
                  </div>
                </div>
              ) : reportsError ? (
                <div className="alert alert-danger m-3">{reportsError}</div>
              ) : reports && reports.length > 0 ? (
                <ul className="list-group list-group-flush">
                  {reports.map(report => (
                    <li 
                      key={report._id} 
                      className={`list-group-item d-flex justify-content-between align-items-start ${currentReport && currentReport._id === report._id ? 'active' : ''}`}
                    >
                      <div 
                        className="flex-grow-1 cursor-pointer"
                        onClick={() => handleSelectReport(report._id)}
                      >
                        <div className="fw-bold">
                          {report.title}
                        </div>
                        <div className="small text-muted">
                          {new Date(report.createdAt).toLocaleDateString()}
                        </div>
                      </div>
                      <button 
                        className="btn btn-outline-danger btn-sm ms-2"
                        onClick={(e) => {
                          e.stopPropagation();
                          handleDeleteReport(report._id);
                        }}
                      >
                        <i className="bi bi-trash"></i>
                      </button>
                    </li>
                  ))}
                </ul>
              ) : (
                <div className="p-3 text-center text-muted">
                  {t('reports.noReports')}
                </div>
              )}
            </div>
          </div>
        </div>
        
        {/* 右侧报告详情 - 响应式调整 */}
        <div className="col-md-8">
          {currentReport ? (
            <div className="card">
              <div className="card-header d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-2 flex-wrap">
                <h5 className="mb-2 mb-md-0">{currentReport.title}</h5>
                <div className="d-flex flex-wrap gap-2">
                  <div className="btn-group btn-group-sm" role="group">
                    <button 
                      type="button" 
                      className={`btn btn-outline-secondary ${previewMode === 'text' ? 'active' : ''}`}
                      onClick={() => setPreviewMode('text')}
                    >
                      {t('reports.textPreview')}
                    </button>
                    <button 
                      type="button" 
                      className={`btn btn-outline-secondary ${previewMode === 'preview' ? 'active' : ''}`}
                      onClick={() => setPreviewMode('preview')}
                    >
                      {t('reports.preview')}
                    </button>
                  </div>
                  <div className="btn-group btn-group-sm" role="group">
                    <button 
                      type="button" 
                      className="btn btn-outline-primary"
                      onClick={() => handleExportReport('text')}
                    >
                      <i className="bi bi-filetype-txt me-1"></i>
                      <span className="d-none d-sm-inline">{t('reports.exportText')}</span>
                    </button>
                    <button 
                      type="button" 
                      className="btn btn-outline-primary"
                      onClick={() => handleExportReport('markdown')}
                    >
                      <i className="bi bi-filetype-md me-1"></i>
                      <span className="d-none d-sm-inline">{t('reports.exportMarkdown')}</span>
                    </button>
                  </div>
                </div>
              </div>
              
              <div className="card-body">
                {/* 统计信息 - 响应式调整 */}
                <div className="row mb-4">
                  <div className="col-12">
                    <h6>{t('reports.statistics')}</h6>
                    <div className="row">
                      <div className="col-sm-6 col-md-4 mb-2">
                        <div className="card bg-light">
                          <div className="card-body p-2 text-center">
                            <div className="small text-muted">{t('reports.totalTasks')}</div>
                            <div className="fw-bold">{currentReport.statistics.totalTasks}</div>
                          </div>
                        </div>
                      </div>
                      <div className="col-sm-6 col-md-4 mb-2">
                        <div className="card bg-light">
                          <div className="card-body p-2 text-center">
                            <div className="small text-muted">{t('reports.completedTasks')}</div>
                            <div className="fw-bold">{currentReport.statistics.completedTasks}</div>
                          </div>
                        </div>
                      </div>
                      <div className="col-sm-6 col-md-4 mb-2">
                        <div className="card bg-light">
                          <div className="card-body p-2 text-center">
                            <div className="small text-muted">{t('reports.inProgressTasks')}</div>
                            <div className="fw-bold">{currentReport.statistics.inProgressTasks}</div>
                          </div>
                        </div>
                      </div>
                      <div className="col-sm-6 col-md-4 mb-2">
                        <div className="card bg-light">
                          <div className="card-body p-2 text-center">
                            <div className="small text-muted">{t('reports.overdueTasks')}</div>
                            <div className="fw-bold">{currentReport.statistics.overdueTasks}</div>
                          </div>
                        </div>
                      </div>
                      <div className="col-sm-6 col-md-4 mb-2">
                        <div className="card bg-light">
                          <div className="card-body p-2 text-center">
                            <div className="small text-muted">{t('reports.completionRate')}</div>
                            <div className="fw-bold">{currentReport.statistics.completionRate}%</div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                
                {/* AI润色功能 - 响应式调整 */}
                <div className="row mb-4">
                  <div className="col-12">
                    <div className="card">
                      <div className="card-body">
                        <h6 className="card-title">{t('reports.aiPolish')}</h6>
                        
                        {/* AI服务提供商选择 */}
                        <div className="mb-3">
                          <label className="form-label">{t('reports.aiProvider')}</label>
                          <select
                            className="form-select"
                            value={provider}
                            onChange={(e) => setProvider(e.target.value)}
                          >
                            <option value="openai">OpenAI</option>
                            <option value="custom">Custom AI Service</option>
                          </select>
                        </div>
                        
                        {/* 模型选择 (仅OpenAI) */}
                        {provider === 'openai' && (
                          <div className="mb-3">
                            <label className="form-label">{t('reports.model')}</label>
                            <select
                              className="form-select"
                              value={model}
                              onChange={(e) => setModel(e.target.value)}
                            >
                              <option value="gpt-3.5-turbo">GPT-3.5 Turbo</option>
                              <option value="gpt-4">GPT-4</option>
                              <option value="gpt-4-turbo">GPT-4 Turbo</option>
                            </select>
                          </div>
                        )}
                        
                        {/* API URL输入 */}
                        <div className="mb-3">
                          <label className="form-label">{t('reports.apiUrl')}</label>
                          <input
                            type="text"
                            className="form-control"
                            placeholder={provider === 'openai' 
                              ? 'https://api.openai.com/v1/chat/completions' 
                              : 'https://your-custom-ai-api.com/endpoint'}
                            value={apiUrl}
                            onChange={(e) => setApiUrl(e.target.value)}
                          />
                        </div>
                        
                        {/* API Key输入 */}
                        <div className="mb-3">
                          <label className="form-label">{t('reports.apiKey')}</label>
                          <input
                            type="password"
                            className="form-control"
                            placeholder={t('reports.enterApiKey')}
                            value={apiKey}
                            onChange={(e) => setApiKey(e.target.value)}
                          />
                        </div>
                        
                        <button 
                          className="btn btn-success d-flex align-items-center justify-content-center"
                          onClick={handlePolishReport}
                          disabled={loading}
                        >
                          {loading ? (
                            <>
                              <span className="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span>
                              {t('reports.polishing')}
                            </>
                          ) : (
                            <>
                              <i className="bi bi-stars me-1"></i>
                              {t('reports.polish')}
                            </>
                          )}
                        </button>
                        <div className="form-text mt-2">
                          {t('reports.apiKeyNote')}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                
                {/* 报告内容 - 响应式调整 */}
                <div className="row">
                  <div className="col-12">
                    <h6>{t('reports.content')}</h6>
                    {previewMode === 'preview' ? (
                      <div className="border rounded p-3 bg-light">
                        <div dangerouslySetInnerHTML={{ __html: renderMarkdown(currentReport.polishedContent || currentReport.content) }} />
                      </div>
                    ) : (
                      <div className="border rounded p-3 bg-light">
                        <pre className="mb-0" style={{ whiteSpace: 'pre-wrap', fontFamily: 'inherit' }}>
                          {currentReport.polishedContent || currentReport.content}
                        </pre>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            </div>
          ) : (
            <div className="card">
              <div className="card-body text-center">
                <p className="text-muted">{t('reports.selectReport')}</p>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default ReportsPage;
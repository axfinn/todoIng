import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useTranslation } from 'react-i18next';
import { fetchTasks, deleteTask, createTask, updateTask, exportTasks, importTasks } from '../features/tasks/taskSlice';
import { generateCalendarICS, downloadICSFile } from '../utils/calendarUtils';
import type { RootState, AppDispatch } from '../app/store';
import type { Task } from '../features/tasks/taskSlice';

const DashboardPage: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  const { tasks, isLoading, error } = useSelector((state: RootState) => state.tasks);
  const { t, i18n } = useTranslation();

  // 获取指定天数后的日期字符串
  const getDateAfterDays = (days: number): string => {
    const date = new Date();
    date.setDate(date.getDate() + days);
    return date.toISOString().split('T')[0];
  };

  // 快捷设置日期的函数
  const setQuickDate = (field: 'deadline' | 'scheduledDate', days: number) => {
    setNewTask({
      ...newTask,
      [field]: getDateAfterDays(days)
    });
  };

  const [newTask, setNewTask] = useState({
    title: '',
    description: '',
    status: 'To Do' as 'To Do' | 'In Progress' | 'Done',
    priority: 'Medium' as 'Low' | 'Medium' | 'High',
    assignee: '',
    deadline: new Date().toISOString().split('T')[0],
    scheduledDate: new Date().toISOString().split('T')[0],
  });

  const [showCreateModal, setShowCreateModal] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);
  const [commentText, setCommentText] = useState<{[key: string]: string}>({});
  const [filterStatus, setFilterStatus] = useState<string>('All');
  const [filterPriority, setFilterPriority] = useState<string>('All');
  const [sortOrder, setSortOrder] = useState<string>('newest');
  const [expandedTaskId, setExpandedTaskId] = useState<string | null>(null); // 用于跟踪展开的任务
  const [githubStats, setGithubStats] = useState({ stars: 0, forks: 0 });
  const [importFile, setImportFile] = useState<File | null>(null);

  // 计算各类任务的数量
  const todoTasksCount = tasks.filter(task => task.status === 'To Do').length;
  const inProgressTasksCount = tasks.filter(task => task.status === 'In Progress').length;
  const doneTasksCount = tasks.filter(task => task.status === 'Done').length;

  useEffect(() => {
    dispatch(fetchTasks());
    
    // 获取GitHub项目统计信息
    fetch('https://api.github.com/repos/axfinn/todoIng')
      .then(response => response.json())
      .then(data => {
        setGithubStats({
          stars: data.stargazers_count || 0,
          forks: data.forks_count || 0
        });
      })
      .catch(error => {
        console.error('Failed to fetch GitHub stats:', error);
      });
  }, [dispatch]);

  const handleDelete = (id: string) => {
    if (window.confirm(t('dashboard.confirmDelete'))) {
      dispatch(deleteTask(id));
    }
  };

  const handleExportCalendar = () => {
    const icsContent = generateCalendarICS(tasks);
    downloadICSFile(icsContent, 'todoing-tasks.ics');
  };

  const handleExportTasks = () => {
    dispatch(exportTasks()).then((action: any) => {
      if (exportTasks.fulfilled.match(action)) {
        const { data, filename } = action.payload;
        const blob = new Blob([data], { type: 'application/json' });
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

  const handleImportTasks = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      dispatch(importTasks(file)).then((action: any) => {
        if (importTasks.fulfilled.match(action)) {
          const { imported, errors } = action.payload;
          alert(t('dashboard.importSuccess', { count: imported, errors: errors.length }));
          // 刷新任务列表
          dispatch(fetchTasks());
        } else if (importTasks.rejected.match(action)) {
          alert(t('dashboard.importFailed') + ': ' + action.payload);
        }
      });
    }
  };

  const handleNewTaskChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setNewTask({ ...newTask, [name]: value });
  };

  const handleCreateTask = (e: React.FormEvent) => {
    e.preventDefault();
    
    // 验证截止日期必须在计划日期之后
    if (newTask.deadline && newTask.scheduledDate) {
      const deadline = new Date(newTask.deadline);
      const scheduledDate = new Date(newTask.scheduledDate);
      if (deadline < scheduledDate) {
        alert(t('dashboard.deadlineBeforeScheduledDate'));
        return;
      }
    }
    
    // 处理日期字段，如果为空则设为null
    const taskToCreate = {
      ...newTask,
      deadline: newTask.deadline || null,
      scheduledDate: newTask.scheduledDate || null,
    };
    
    dispatch(createTask(taskToCreate));
    setNewTask({
      title: '',
      description: '',
      status: 'To Do',
      priority: 'Medium',
      assignee: '',
      deadline: '',
      scheduledDate: '',
    });
    setShowCreateModal(false);
  };

  const handleEdit = (task: Task) => {
    setEditingTask({
      ...task,
      deadline: task.deadline ? task.deadline.split('T')[0] : '',
      scheduledDate: task.scheduledDate ? task.scheduledDate.split('T')[0] : '',
    });
  };

  const handleUpdateTask = (e: React.FormEvent) => {
    e.preventDefault();
    if (editingTask) {
      // 验证截止日期必须在计划日期之后
      if (editingTask.deadline && editingTask.scheduledDate) {
        const deadline = new Date(editingTask.deadline);
        const scheduledDate = new Date(editingTask.scheduledDate);
        if (deadline < scheduledDate) {
          alert(t('dashboard.deadlineBeforeScheduledDate'));
          return;
        }
      }
      
      // 处理日期字段，如果为空则设为null
      const taskToUpdate = {
        ...editingTask,
        deadline: editingTask.deadline || null,
        scheduledDate: editingTask.scheduledDate || null,
      };
      
      dispatch(updateTask(taskToUpdate));
      setEditingTask(null);
    }
  };

  const handleEditChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    if (editingTask) {
      const { name, value } = e.target;
      setEditingTask({ ...editingTask, [name]: value });
    }
  };

  const handleAddComment = (taskId: string) => {
    const comment = commentText[taskId];
    if (comment?.trim()) {
      // 查找当前任务
      const task = tasks.find(t => t._id === taskId);
      if (task) {
        // 创建新的评论对象
        const newComment = {
          text: comment,
          createdAt: new Date().toISOString()
        };
        
        // 更新任务的评论列表
        const updatedComments = [...(task.comments || []), newComment];
        
        // 更新任务，将状态设为"In Progress"
        dispatch(updateTask({ 
          ...task, 
          status: 'In Progress',
          comments: updatedComments 
        }));
        
        // 清除输入框
        setCommentText(prev => ({
          ...prev,
          [taskId]: ''
        }));
      }
    }
  };

  const handleCommentChange = (taskId: string, text: string) => {
    setCommentText({ ...commentText, [taskId]: text });
  };

  // 快速更新任务状态
  const updateTaskStatus = (task: Task, status: 'To Do' | 'In Progress' | 'Done') => {
    const updatedTask = { ...task, status };
    dispatch(updateTask(updatedTask));
  };

  // 快速筛选任务
  const setFilter = (status: string) => {
    setFilterStatus(status);
  };

  // 检查任务是否临近截止日期
  const isTaskNearDeadline = (task: Task) => {
    if (!task.deadline) return false;
    
    const deadline = new Date(task.deadline);
    const now = new Date();
    const diffTime = deadline.getTime() - now.getTime();
    const diffDays = diffTime / (1000 * 3600 * 24);
    
    // 如果截止日期已过或在1天内，则标记为临近截止日期
    return diffDays <= 1 && diffDays >= 0;
  };

  // 检查任务是否已过截止日期
  const isTaskOverdue = (task: Task) => {
    if (!task.deadline) return false;
    
    const deadline = new Date(task.deadline);
    const now = new Date();
    
    return deadline < now;
  };

  const filteredAndSortedTasks = tasks
    .filter(task => filterStatus === 'All' || task.status === filterStatus)
    .filter(task => filterPriority === 'All' || task.priority === filterPriority)
    .sort((a, b) => {
      // 首先按是否有截止日期排序（有截止日期的优先）
      const aHasDeadline = a.deadline ? 1 : 0;
      const bHasDeadline = b.deadline ? 1 : 0;
      
      if (aHasDeadline !== bHasDeadline) {
        return bHasDeadline - aHasDeadline;
      }
      
      // 如果都有截止日期或都没有，按截止日期排序（早的在前）
      if (a.deadline && b.deadline) {
        return new Date(a.deadline).getTime() - new Date(b.deadline).getTime();
      }
      
      // 按创建时间排序
      if (sortOrder === 'newest') {
        return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
      } else {
        return new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime();
      }
    });

  const getStatusClass = (status: string) => {
    switch (status) {
      case 'To Do':
        return 'bg-secondary';
      case 'In Progress':
        return 'bg-warning text-dark';
      case 'Done':
        return 'bg-success';
      default:
        return 'bg-secondary';
    }
  };

  const getPriorityClass = (priority: string) => {
    switch (priority) {
      case 'Low':
        return 'bg-info';
      case 'Medium':
        return 'bg-warning text-dark';
      case 'High':
        return 'bg-danger';
      default:
        return 'bg-secondary';
    }
  };

  // 获取截止日期徽章样式
  const getDeadlineClass = (task: Task) => {
    if (isTaskOverdue(task)) {
      return 'bg-danger';
    } else if (isTaskNearDeadline(task)) {
      return 'bg-warning text-dark';
    }
    return 'bg-secondary';
  };

  const toggleTaskDetails = (taskId: string) => {
    setExpandedTaskId(expandedTaskId === taskId ? null : taskId);
  };

  const translateStatus = (status: string) => {
    switch (status) {
      case 'To Do':
        return t('status.todo');
      case 'In Progress':
        return t('status.inProgress');
      case 'Done':
        return t('status.done');
      default:
        return status;
    }
  };

  const translatePriority = (priority: string) => {
    switch (priority) {
      case 'Low':
        return t('priority.low');
      case 'Medium':
        return t('priority.medium');
      case 'High':
        return t('priority.high');
      default:
        return priority;
    }
  };

  // GitHub star操作
  const handleStarRepo = () => {
    window.open('https://github.com/axfinn/todoIng', '_blank');
  };

  if (isLoading) {
    return (
      <div className="container py-5">
        <div className="d-flex justify-content-center align-items-center" style={{ height: '70vh' }}>
          <div className="text-center">
            <div className="spinner-border text-primary" role="status">
              <span className="visually-hidden">{t('common.loading')}</span>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="container py-4">
      <div className="row">
        <div className="col-12">
          <h2 className="mb-4">{t('dashboard.title')}</h2>
          
          {/* 创建任务按钮和GitHub信息 - 响应式调整 */}
          <div className="d-flex flex-column flex-md-row justify-content-between align-items-md-center mb-4 gap-3">
            <div className="d-flex flex-wrap gap-2">
              <button 
                className="btn btn-primary" 
                onClick={() => setShowCreateModal(true)}
              >
                <i className="bi bi-plus-lg me-1"></i>
                {t('dashboard.newTask')}
              </button>
              
              <button 
                className="btn btn-outline-secondary"
                onClick={handleExportCalendar}
                title={t('dashboard.exportCalendar')}
              >
                <i className="bi bi-calendar-plus me-1"></i>
                <span className="d-none d-sm-inline">{t('dashboard.exportCalendar')}</span>
                <span className="d-inline d-sm-none">{t('dashboard.exportCalendarShort')}</span>
              </button>
              
              <button 
                className="btn btn-outline-info"
                onClick={handleExportTasks}
                title={t('dashboard.exportTasks')}
              >
                <i className="bi bi-download me-1"></i>
                <span className="d-none d-sm-inline">{t('dashboard.exportTasks')}</span>
                <span className="d-inline d-sm-none">{t('dashboard.exportTasksShort')}</span>
              </button>
              
              <label className="btn btn-outline-success mb-0">
                <i className="bi bi-upload me-1"></i>
                <span className="d-none d-sm-inline">{t('dashboard.importTasks')}</span>
                <span className="d-inline d-sm-none">{t('dashboard.importTasksShort')}</span>
                <input 
                  type="file" 
                  accept=".json" 
                  onChange={handleImportTasks} 
                  style={{ display: 'none' }} 
                />
              </label>
            </div>
            
            <div className="d-flex align-items-center">
              <i className="bi bi-github me-2"></i>
              <button 
                className="btn btn-outline-dark btn-sm me-3 d-flex align-items-center"
                onClick={handleStarRepo}
              >
                <i className="bi bi-star-fill me-1"></i> 
                <span className="d-none d-sm-inline">Star</span>
              </button>
              <div className="d-flex">
                <span className="badge bg-secondary me-2 d-flex align-items-center">
                  <i className="bi bi-star me-1"></i> <span className="d-none d-sm-inline"> {githubStats.stars}</span>
                </span>
                <span className="badge bg-secondary d-flex align-items-center">
                  <i className="bi bi-git me-1"></i> <span className="d-none d-sm-inline"> {githubStats.forks}</span>
                </span>
              </div>
            </div>
          </div>

          {/* 快速筛选按钮 - 响应式调整 */}
          <div className="d-flex flex-wrap gap-2 mb-4">
            <button 
              className={`btn ${filterStatus === 'All' ? 'btn-primary' : 'btn-outline-primary'} d-flex align-items-center`}
              onClick={() => setFilter('All')}
            >
              {t('dashboard.allTasks')} <span className="badge bg-white text-primary ms-1">{tasks.length}</span>
            </button>
            <button 
              className={`btn ${filterStatus === 'To Do' ? 'btn-secondary' : 'btn-outline-secondary'} d-flex align-items-center`}
              onClick={() => setFilter('To Do')}
            >
              {t('status.todo')} <span className="badge bg-white text-secondary ms-1">{todoTasksCount}</span>
            </button>
            <button 
              className={`btn ${filterStatus === 'In Progress' ? 'btn-warning' : 'btn-outline-warning'} d-flex align-items-center`}
              onClick={() => setFilter('In Progress')}
            >
              {t('status.inProgress')} <span className="badge bg-white text-warning ms-1">{inProgressTasksCount}</span>
            </button>
            <button 
              className={`btn ${filterStatus === 'Done' ? 'btn-success' : 'btn-outline-success'} d-flex align-items-center`}
              onClick={() => setFilter('Done')}
            >
              {t('status.done')} <span className="badge bg-white text-success ms-1">{doneTasksCount}</span>
            </button>
          </div>

          {/* 筛选和排序控件 - 响应式调整 */}
          <div className="row mb-4">
            <div className="col-md-4 col-sm-6 mb-3">
              <label htmlFor="filterStatus" className="form-label">{t('dashboard.filter.status')}</label>
              <select
                className="form-select"
                id="filterStatus"
                value={filterStatus}
                onChange={(e) => setFilterStatus(e.target.value)}
              >
                <option value="All">{t('dashboard.filter.all')}</option>
                <option value="To Do">{t('status.todo')}</option>
                <option value="In Progress">{t('status.inProgress')}</option>
                <option value="Done">{t('status.done')}</option>
              </select>
            </div>
            <div className="col-md-4 col-sm-6 mb-3">
              <label htmlFor="filterPriority" className="form-label">{t('dashboard.filter.priority')}</label>
              <select
                className="form-select"
                id="filterPriority"
                value={filterPriority}
                onChange={(e) => setFilterPriority(e.target.value)}
              >
                <option value="All">{t('dashboard.filter.all')}</option>
                <option value="Low">{t('priority.low')}</option>
                <option value="Medium">{t('priority.medium')}</option>
                <option value="High">{t('priority.high')}</option>
              </select>
            </div>
            <div className="col-md-4 col-sm-6 mb-3">
              <label htmlFor="sortOrder" className="form-label">{t('dashboard.sort.label')}</label>
              <select
                className="form-select"
                id="sortOrder"
                value={sortOrder}
                onChange={(e) => setSortOrder(e.target.value)}
              >
                <option value="newest">{t('dashboard.sort.newest')}</option>
                <option value="oldest">{t('dashboard.sort.oldest')}</option>
              </select>
            </div>
          </div>

          {/* 任务列表 */}
          <div className="card">
            <div className="card-header">
              <h5 className="mb-0">{t('dashboard.taskList')}</h5>
            </div>
            <div className="card-body">
              {error && (
                <div className="alert alert-danger" role="alert">
                  {error}
                </div>
              )}

              {filteredAndSortedTasks.length === 0 ? (
                <div className="text-center py-5">
                  <p className="mb-0">{t('dashboard.noTasks')}</p>
                </div>
              ) : (
                <div className="row">
                  {filteredAndSortedTasks.map((task) => (
                    <div key={task._id} className="col-12 mb-3">
                      <div className={`card ${isTaskOverdue(task) ? 'border-danger' : ''}`}>
                        <div className="card-body">
                          <div className="d-flex flex-column flex-md-row justify-content-between align-items-md-start gap-2">
                            <div>
                              <h5 className="mb-1">{task.title}</h5>
                              <p className="mb-2 text-muted">
                                {task.description && task.description.length > 100 
                                  ? `${task.description.substring(0, 100)}...` 
                                  : task.description}
                              </p>
                              <div className="d-flex flex-wrap gap-2 mb-2">
                                <span className={`badge ${getStatusClass(task.status)}`}>
                                  {translateStatus(task.status)}
                                </span>
                                <span className={`badge ${getPriorityClass(task.priority)}`}>
                                  {translatePriority(task.priority)}
                                </span>
                                {task.deadline && (
                                  <span className={`badge ${getDeadlineClass(task)}`}>
                                    <i className="bi bi-calendar me-1"></i>
                                    {new Date(task.deadline).toLocaleDateString()}
                                    {(isTaskNearDeadline(task) || isTaskOverdue(task)) && (
                                      <span className="ms-1">
                                        {isTaskOverdue(task) ? t('dashboard.overdue') : t('dashboard.dueSoon')}
                                      </span>
                                    )}
                                  </span>
                                )}
                              </div>
                              <p className="mb-0 small text-muted">
                                <i className="bi bi-clock me-1"></i>
                                {t('dashboard.created')}: {new Date(task.createdAt).toLocaleDateString(i18n.language)}
                              </p>
                            </div>
                            <div className="d-flex flex-wrap gap-2">
                              {/* 快速状态更新按钮 */}
                              {task.status !== 'Done' && (
                                <button
                                  className="btn btn-sm btn-success"
                                  onClick={() => updateTaskStatus(task, 'Done')}
                                  title={t('dashboard.markAsDone')}
                                >
                                  <i className="bi bi-check-circle"></i>
                                  <span className="d-none d-sm-inline ms-1">{t('dashboard.done')}</span>
                                </button>
                              )}
                              {task.status !== 'In Progress' && task.status !== 'Done' && (
                                <button
                                  className="btn btn-sm btn-warning"
                                  onClick={() => updateTaskStatus(task, 'In Progress')}
                                  title={t('dashboard.markAsInProgress')}
                                >
                                  <i className="bi bi-arrow-right-circle"></i>
                                  <span className="d-none d-sm-inline ms-1">{t('dashboard.inProgress')}</span>
                                </button>
                              )}
                              <button
                                className="btn btn-sm btn-outline-primary"
                                onClick={() => toggleTaskDetails(task._id)}
                              >
                                {expandedTaskId === task._id 
                                  ? t('dashboard.hideDetails') 
                                  : t('dashboard.showDetails')}
                              </button>
                              <button
                                className="btn btn-sm btn-outline-secondary"
                                onClick={() => handleEdit(task)}
                              >
                                <i className="bi bi-pencil"></i>
                                <span className="d-none d-sm-inline ms-1">{t('dashboard.edit')}</span>
                              </button>
                              <button
                                className="btn btn-sm btn-outline-danger"
                                onClick={() => handleDelete(task._id)}
                              >
                                <i className="bi bi-trash"></i>
                                <span className="d-none d-sm-inline ms-1">{t('dashboard.delete')}</span>
                              </button>
                            </div>
                          </div>

                          {/* 任务详情和评论 */}
                          {(expandedTaskId === task._id) && (
                            <div className="mt-3 pt-3 border-top">
                              <h6>{t('dashboard.comments')}</h6>
                              <div className="timeline">
                                {/* 任务创建事件 */}
                                <div className="d-flex mb-3">
                                  <div className="flex-shrink-0">
                                    <div className="rounded-circle bg-success d-flex align-items-center justify-content-center" 
                                         style={{ width: '32px', height: '32px' }}>
                                      <i className="bi bi-plus text-white"></i>
                                    </div>
                                  </div>
                                  <div className="flex-grow-1 ms-3">
                                    <div className="card">
                                      <div className="card-body py-2 px-3">
                                        <p className="mb-0">
                                          <strong>{t('dashboard.created')}</strong>
                                        </p>
                                        <small className="text-muted">
                                          {new Date(task.createdAt).toLocaleString(i18n.language)}
                                        </small>
                                      </div>
                                    </div>
                                  </div>
                                </div>

                                {/* 评论列表 */}
                                {task.comments && task.comments.map((comment, index) => (
                                  <div className="d-flex mb-3" key={index}>
                                    <div className="flex-shrink-0">
                                      <div className="rounded-circle bg-primary d-flex align-items-center justify-content-center" 
                                           style={{ width: '32px', height: '32px' }}>
                                        <i className="bi bi-chat text-white"></i>
                                      </div>
                                    </div>
                                    <div className="flex-grow-1 ms-3">
                                      <div className="card">
                                        <div className="card-body py-2 px-3">
                                          <p className="mb-1">{comment.text}</p>
                                          <small className="text-muted">
                                            {new Date(comment.createdAt).toLocaleString(i18n.language)}
                                          </small>
                                        </div>
                                      </div>
                                    </div>
                                  </div>
                                ))}
                              </div>

                              {/* 添加评论 */}
                              <div className="mt-3">
                                <div className="input-group">
                                  <input
                                    type="text"
                                    className="form-control"
                                    placeholder={t('dashboard.addComment')}
                                    value={commentText[task._id] || ''}
                                    onChange={(e) => handleCommentChange(task._id, e.target.value)}
                                    onKeyPress={(e) => {
                                      if (e.key === 'Enter') {
                                        handleAddComment(task._id);
                                      }
                                    }}
                                  />
                                  <button
                                    className="btn btn-outline-primary"
                                    type="button"
                                    onClick={() => handleAddComment(task._id)}
                                  >
                                    {t('dashboard.addCommentButton')}
                                  </button>
                                </div>
                              </div>
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* 创建任务模态框 */}
      {showCreateModal && (
        <div className="modal show d-block" tabIndex={-1} style={{ backgroundColor: 'rgba(0,0,0,0.5)' }}>
          <div className="modal-dialog">
            <div className="modal-content">
              <div className="modal-header">
                <h5 className="modal-title">{t('dashboard.newTask')}</h5>
                <button
                  type="button"
                  className="btn-close"
                  onClick={() => setShowCreateModal(false)}
                ></button>
              </div>
              <form onSubmit={handleCreateTask}>
                <div className="modal-body">
                  <div className="mb-3">
                    <label htmlFor="title" className="form-label">{t('dashboard.taskTitle')}</label>
                    <input
                      type="text"
                      className="form-control"
                      id="title"
                      name="title"
                      value={newTask.title}
                      onChange={handleNewTaskChange}
                      required
                    />
                  </div>
                  <div className="mb-3">
                    <label htmlFor="assignee" className="form-label">{t('dashboard.assignee')}</label>
                    <input
                      type="text"
                      className="form-control"
                      id="assignee"
                      name="assignee"
                      value={newTask.assignee}
                      onChange={handleNewTaskChange}
                    />
                  </div>
                  <div className="mb-3">
                    <label htmlFor="status" className="form-label">{t('dashboard.status')}</label>
                    <select
                      className="form-select"
                      id="status"
                      name="status"
                      value={newTask.status}
                      onChange={handleNewTaskChange}
                    >
                      <option value="To Do">{t('status.todo')}</option>
                      <option value="In Progress">{t('status.inProgress')}</option>
                      <option value="Done">{t('status.done')}</option>
                    </select>
                  </div>
                  <div className="mb-3">
                    <label htmlFor="priority" className="form-label">{t('dashboard.priority')}</label>
                    <select
                      className="form-select"
                      id="priority"
                      name="priority"
                      value={newTask.priority}
                      onChange={handleNewTaskChange}
                    >
                      <option value="Low">{t('priority.low')}</option>
                      <option value="Medium">{t('priority.medium')}</option>
                      <option value="High">{t('priority.high')}</option>
                    </select>
                  </div>
                  <div className="mb-3">
                    <label htmlFor="deadline" className="form-label">{t('dashboard.deadline')}</label>
                    <input
                      type="date"
                      className="form-control"
                      id="deadline"
                      name="deadline"
                      value={newTask.deadline}
                      onChange={handleNewTaskChange}
                    />
                    <div className="mt-1">
                      <button 
                        type="button" 
                        className="btn btn-sm btn-outline-primary me-1"
                        onClick={() => setQuickDate('deadline', 0)}
                      >
                        {t('dashboard.today')}
                      </button>
                      <button 
                        type="button" 
                        className="btn btn-sm btn-outline-primary me-1"
                        onClick={() => setQuickDate('deadline', 1)}
                      >
                        {t('dashboard.tomorrow')}
                      </button>
                      <button 
                        type="button" 
                        className="btn btn-sm btn-outline-primary me-1"
                        onClick={() => setQuickDate('deadline', 2)}
                      >
                        {t('dashboard.dayAfterTomorrow')}
                      </button>
                      <button 
                        type="button" 
                        className="btn btn-sm btn-outline-primary me-1"
                        onClick={() => setQuickDate('deadline', 7)}
                      >
                        {t('dashboard.nextWeek')}
                      </button>
                      <button 
                        type="button" 
                        className="btn btn-sm btn-outline-primary"
                        onClick={() => setQuickDate('deadline', 30)}
                      >
                        {t('dashboard.nextMonth')}
                      </button>
                    </div>
                  </div>
                  <div className="mb-3">
                    <label htmlFor="scheduledDate" className="form-label">{t('dashboard.scheduledDate')}</label>
                    <input
                      type="date"
                      className="form-control"
                      id="scheduledDate"
                      name="scheduledDate"
                      value={newTask.scheduledDate}
                      onChange={handleNewTaskChange}
                    />
                    <div className="mt-1">
                      <button 
                        type="button" 
                        className="btn btn-sm btn-outline-primary me-1"
                        onClick={() => setQuickDate('scheduledDate', 0)}
                      >
                        {t('dashboard.today')}
                      </button>
                      <button 
                        type="button" 
                        className="btn btn-sm btn-outline-primary me-1"
                        onClick={() => setQuickDate('scheduledDate', 1)}
                      >
                        {t('dashboard.tomorrow')}
                      </button>
                      <button 
                        type="button" 
                        className="btn btn-sm btn-outline-primary me-1"
                        onClick={() => setQuickDate('scheduledDate', 2)}
                      >
                        {t('dashboard.dayAfterTomorrow')}
                      </button>
                      <button 
                        type="button" 
                        className="btn btn-sm btn-outline-primary me-1"
                        onClick={() => setQuickDate('scheduledDate', 7)}
                      >
                        {t('dashboard.nextWeek')}
                      </button>
                      <button 
                        type="button" 
                        className="btn btn-sm btn-outline-primary"
                        onClick={() => setQuickDate('scheduledDate', 30)}
                      >
                        {t('dashboard.nextMonth')}
                      </button>
                    </div>
                  </div>
                  <div className="mb-3">
                    <label htmlFor="description" className="form-label">{t('dashboard.description')}</label>
                    <textarea
                      className="form-control"
                      id="description"
                      name="description"
                      rows={3}
                      value={newTask.description}
                      onChange={handleNewTaskChange}
                    ></textarea>
                  </div>
                </div>
                <div className="modal-footer">
                  <button
                    type="button"
                    className="btn btn-secondary"
                    onClick={() => setShowCreateModal(false)}
                  >
                    {t('common.cancel')}
                  </button>
                  <button type="submit" className="btn btn-primary">
                    {t('dashboard.createTask')}
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}

      {/* 编辑任务模态框 */}
      {editingTask && (
        <div className="modal show d-block" tabIndex={-1} style={{ backgroundColor: 'rgba(0,0,0,0.5)' }}>
          <div className="modal-dialog">
            <div className="modal-content">
              <div className="modal-header">
                <h5 className="modal-title">{t('dashboard.edit')}</h5>
                <button
                  type="button"
                  className="btn-close"
                  onClick={() => setEditingTask(null)}
                ></button>
              </div>
              <form onSubmit={handleUpdateTask}>
                <div className="modal-body">
                  <div className="mb-3">
                    <label htmlFor="editTitle" className="form-label">{t('dashboard.taskTitle')}</label>
                    <input
                      type="text"
                      className="form-control"
                      id="editTitle"
                      name="title"
                      value={editingTask.title}
                      onChange={handleEditChange}
                      required
                    />
                  </div>
                  <div className="mb-3">
                    <label htmlFor="editDescription" className="form-label">{t('dashboard.description')}</label>
                    <textarea
                      className="form-control"
                      id="editDescription"
                      name="description"
                      rows={3}
                      value={editingTask.description}
                      onChange={handleEditChange}
                    ></textarea>
                  </div>
                  <div className="mb-3">
                    <label htmlFor="editStatus" className="form-label">{t('dashboard.status')}</label>
                    <select
                      className="form-select"
                      id="editStatus"
                      name="status"
                      value={editingTask.status}
                      onChange={handleEditChange}
                    >
                      <option value="To Do">{t('status.todo')}</option>
                      <option value="In Progress">{t('status.inProgress')}</option>
                      <option value="Done">{t('status.done')}</option>
                    </select>
                  </div>
                  <div className="mb-3">
                    <label htmlFor="editPriority" className="form-label">{t('dashboard.priority')}</label>
                    <select
                      className="form-select"
                      id="editPriority"
                      name="priority"
                      value={editingTask.priority}
                      onChange={handleEditChange}
                    >
                      <option value="Low">{t('priority.low')}</option>
                      <option value="Medium">{t('priority.medium')}</option>
                      <option value="High">{t('priority.high')}</option>
                    </select>
                  </div>
                  <div className="mb-3">
                    <label htmlFor="editDeadline" className="form-label">{t('dashboard.deadline')}</label>
                    <input
                      type="date"
                      className="form-control"
                      id="editDeadline"
                      name="deadline"
                      value={editingTask.deadline || ''}
                      onChange={handleEditChange}
                    />
                    <div className="mt-1">
                      <button 
                        type="button" 
                        className="btn btn-sm btn-outline-primary me-1"
                        onClick={() => editingTask && setEditingTask({
                          ...editingTask,
                          deadline: getDateAfterDays(0)
                        })}
                      >
                        {t('dashboard.today')}
                      </button>
                      <button 
                        type="button" 
                        className="btn btn-sm btn-outline-primary me-1"
                        onClick={() => editingTask && setEditingTask({
                          ...editingTask,
                          deadline: getDateAfterDays(1)
                        })}
                      >
                        {t('dashboard.tomorrow')}
                      </button>
                      <button 
                        type="button" 
                        className="btn btn-sm btn-outline-primary me-1"
                        onClick={() => editingTask && setEditingTask({
                          ...editingTask,
                          deadline: getDateAfterDays(2)
                        })}
                      >
                        {t('dashboard.dayAfterTomorrow')}
                      </button>
                      <button 
                        type="button" 
                        className="btn btn-sm btn-outline-primary me-1"
                        onClick={() => editingTask && setEditingTask({
                          ...editingTask,
                          deadline: getDateAfterDays(7)
                        })}
                      >
                        {t('dashboard.nextWeek')}
                      </button>
                      <button 
                        type="button" 
                        className="btn btn-sm btn-outline-primary"
                        onClick={() => editingTask && setEditingTask({
                          ...editingTask,
                          deadline: getDateAfterDays(30)
                        })}
                      >
                        {t('dashboard.nextMonth')}
                      </button>
                    </div>
                  </div>
                  <div className="mb-3">
                    <label htmlFor="editScheduledDate" className="form-label">{t('dashboard.scheduledDate')}</label>
                    <input
                      type="date"
                      className="form-control"
                      id="editScheduledDate"
                      name="scheduledDate"
                      value={editingTask.scheduledDate || ''}
                      onChange={handleEditChange}
                    />
                    <div className="mt-1">
                      <button 
                        type="button" 
                        className="btn btn-sm btn-outline-primary me-1"
                        onClick={() => editingTask && setEditingTask({
                          ...editingTask,
                          scheduledDate: getDateAfterDays(0)
                        })}
                      >
                        {t('dashboard.today')}
                      </button>
                      <button 
                        type="button" 
                        className="btn btn-sm btn-outline-primary me-1"
                        onClick={() => editingTask && setEditingTask({
                          ...editingTask,
                          scheduledDate: getDateAfterDays(1)
                        })}
                      >
                        {t('dashboard.tomorrow')}
                      </button>
                      <button 
                        type="button" 
                        className="btn btn-sm btn-outline-primary me-1"
                        onClick={() => editingTask && setEditingTask({
                          ...editingTask,
                          scheduledDate: getDateAfterDays(2)
                        })}
                      >
                        {t('dashboard.dayAfterTomorrow')}
                      </button>
                      <button 
                        type="button" 
                        className="btn btn-sm btn-outline-primary me-1"
                        onClick={() => editingTask && setEditingTask({
                          ...editingTask,
                          scheduledDate: getDateAfterDays(7)
                        })}
                      >
                        {t('dashboard.nextWeek')}
                      </button>
                      <button 
                        type="button" 
                        className="btn btn-sm btn-outline-primary"
                        onClick={() => editingTask && setEditingTask({
                          ...editingTask,
                          scheduledDate: getDateAfterDays(30)
                        })}
                      >
                        {t('dashboard.nextMonth')}
                      </button>
                    </div>
                  </div>
                  <div className="mb-3">
                    <label htmlFor="editAssignee" className="form-label">{t('dashboard.assignee')}</label>
                    <input
                      type="text"
                      className="form-control"
                      id="editAssignee"
                      name="assignee"
                      value={editingTask.assignee || ''}
                      onChange={handleEditChange}
                    />
                  </div>
                </div>
                <div className="modal-footer">
                  <button
                    type="button"
                    className="btn btn-secondary"
                    onClick={() => setEditingTask(null)}
                  >
                    {t('common.cancel')}
                  </button>
                  <button type="submit" className="btn btn-primary">
                    {t('common.save')}
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default DashboardPage;
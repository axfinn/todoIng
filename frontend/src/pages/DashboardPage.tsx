export type TaskComment = {
  text: string;
  createdAt: string;
};
import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useTranslation } from 'react-i18next';
import { fetchTasks, deleteTask, createTask, updateTask } from '../features/tasks/taskSlice';
import type { RootState, AppDispatch } from '../app/store';
import type { Task } from '../features/tasks/taskSlice';

const DashboardPage: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  const { tasks, isLoading, error } = useSelector((state: RootState) => state.tasks);
  const { t, i18n } = useTranslation();

  const [newTask, setNewTask] = useState({
    title: '',
    description: '',
    status: 'To Do' as 'To Do' | 'In Progress' | 'Done',
    priority: 'Medium' as 'Low' | 'Medium' | 'High',
    assignee: '',
  });

  const [showCreateModal, setShowCreateModal] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);
  const [commentText, setCommentText] = useState<{[key: string]: string}>({});
  const [filterStatus, setFilterStatus] = useState<string>('All');
  const [filterPriority, setFilterPriority] = useState<string>('All');
  const [sortOrder, setSortOrder] = useState<string>('newest');
  const [expandedTaskId, setExpandedTaskId] = useState<string | null>(null); // 用于跟踪展开的任务
  const [githubStats, setGithubStats] = useState({ stars: 0, forks: 0 });

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

  // 快速更新任务状态
  const updateTaskStatus = (task: Task, status: 'To Do' | 'In Progress' | 'Done') => {
    const updatedTask = { ...task, status };
    dispatch(updateTask(updatedTask));
  };

  const handleNewTaskChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    setNewTask({ ...newTask, [e.target.name]: e.target.value });
  };

  const handleCreateTask = (e: React.FormEvent) => {
    e.preventDefault();
    dispatch(createTask(newTask));
    setNewTask({
      title: '',
      description: '',
      status: 'To Do',
      priority: 'Medium',
      assignee: '',
    });
  };

  const handleEdit = (task: Task) => {
    setEditingTask(task);
  };

  const handleUpdateTask = (e: React.FormEvent) => {
    e.preventDefault();
    if (editingTask) {
      dispatch(updateTask(editingTask));
      setEditingTask(null);
    }
  };

  const handleEditChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    if (editingTask) {
      setEditingTask({ ...editingTask, [e.target.name]: e.target.value });
    }
  };

  const handleAddComment = (taskId: string) => {
    if (commentText[taskId]?.trim()) {
      const updatedTask = tasks.find(task => task._id === taskId);
      if (updatedTask) {
        const newComment = {
          text: commentText[taskId],
          createdAt: new Date().toISOString()
        };
        
        const updatedComments = [...(updatedTask.comments || []), newComment];
        
        dispatch(updateTask({ ...updatedTask, comments: updatedComments }));
        
        setCommentText({ ...commentText, [taskId]: '' });
      }
    }
  };

  const handleCommentChange = (taskId: string, text: string) => {
    setCommentText({ ...commentText, [taskId]: text });
  };

  const filteredAndSortedTasks = tasks
    .filter(task => filterStatus === 'All' || task.status === filterStatus)
    .filter(task => filterPriority === 'All' || task.priority === filterPriority)
    .sort((a, b) => {
      if (sortOrder === 'newest') {
        return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
      } else {
        return new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime();
      }
    });

  const getStatusBadgeClass = (status: string) => {
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

  const getPriorityBadgeClass = (priority: string) => {
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
          
          {/* 创建任务按钮和GitHub信息 */}
          <div className="d-flex justify-content-between align-items-center mb-4">
            <button 
              className="btn btn-primary" 
              onClick={() => setShowCreateModal(true)}
            >
              <i className="bi bi-plus-lg me-1"></i>
              {t('dashboard.newTask')}
            </button>
            
            <div className="d-flex align-items-center">
              <i className="bi bi-github me-2"></i>
              <button 
                className="btn btn-outline-dark btn-sm me-3 d-flex align-items-center"
                onClick={() => window.open('https://github.com/axfinn/todoIng', '_blank')}
              >
                <i className="bi bi-star-fill me-1"></i> 
                <span>Star</span>
              </button>
              <div className="d-flex">
                <span className="badge bg-secondary me-2">
                  <i className="bi bi-star me-1"></i> {githubStats.stars}
                </span>
                <span className="badge bg-secondary">
                  <i className="bi bi-git me-1"></i> {githubStats.forks}
                </span>
              </div>
            </div>
          </div>

          {/* 快速筛选按钮 */}
          <div className="d-flex gap-2 mb-4">
            <button 
              className={`btn ${filterStatus === 'All' ? 'btn-primary' : 'btn-outline-primary'}`}
              onClick={() => setFilterStatus('All')}
            >
              {t('dashboard.allTasks')} <span className="badge bg-white text-primary ms-1">{tasks.length}</span>
            </button>
            <button 
              className={`btn ${filterStatus === 'To Do' ? 'btn-secondary' : 'btn-outline-secondary'}`}
              onClick={() => setFilterStatus('To Do')}
            >
              {t('status.todo')} <span className="badge bg-white text-secondary ms-1">{tasks.filter(task => task.status === 'To Do').length}</span>
            </button>
            <button 
              className={`btn ${filterStatus === 'In Progress' ? 'btn-warning' : 'btn-outline-warning'}`}
              onClick={() => setFilterStatus('In Progress')}
            >
              {t('status.inProgress')} <span className="badge bg-white text-warning ms-1">{tasks.filter(task => task.status === 'In Progress').length}</span>
            </button>
            <button 
              className={`btn ${filterStatus === 'Done' ? 'btn-success' : 'btn-outline-success'}`}
              onClick={() => setFilterStatus('Done')}
            >
              {t('status.done')} <span className="badge bg-white text-success ms-1">{tasks.filter(task => task.status === 'Done').length}</span>
            </button>
          </div>

          {/* 筛选和排序控件 */}
          <div className="row mb-4">
            <div className="col-md-4 mb-3">
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
            <div className="col-md-4 mb-3">
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
            <div className="col-md-4 mb-3">
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
                      <div className="card">
                        <div className="card-body">
                          <div className="d-flex justify-content-between align-items-start">
                            <div>
                              <h5 className="card-title">{task.title}</h5>
                              <p className="card-text text-muted">{task.description}</p>
                              <div className="d-flex gap-2 mb-2">
                                <span className={`badge ${getStatusBadgeClass(task.status)}`}>
                                  {translateStatus(task.status)}
                                </span>
                                <span className={`badge ${getPriorityBadgeClass(task.priority)}`}>
                                  {translatePriority(task.priority)}
                                </span>
                              </div>
                              {task.assignee && (
                                <p className="mb-1">
                                  <small className="text-muted">
                                    <i className="bi bi-person me-1"></i>
                                    {task.assignee}
                                  </small>
                                </p>
                              )}
                              <p className="mb-1">
                                <small className="text-muted">
                                  <i className="bi bi-calendar me-1"></i>
                                  {t('dashboard.created')}: {new Date(task.createdAt).toLocaleString(i18n.language)}
                                </small>
                              </p>
                              <p className="mb-1">
                                <small className="text-muted">
                                  <i className="bi bi-arrow-repeat me-1"></i>
                                  {t('dashboard.updated')}: {new Date(task.updatedAt).toLocaleString(i18n.language)}
                                </small>
                              </p>
                            </div>
                            <div className="d-flex gap-2">
                              {/* 快速状态更新按钮 */}
                              {task.status !== 'Done' && (
                                <button
                                  className="btn btn-sm btn-success"
                                  onClick={() => updateTaskStatus(task, 'Done')}
                                  title={t('dashboard.markAsDone')}
                                >
                                  <i className="bi bi-check-circle"></i>
                                </button>
                              )}
                              {task.status !== 'In Progress' && task.status !== 'Done' && (
                                <button
                                  className="btn btn-sm btn-warning"
                                  onClick={() => updateTaskStatus(task, 'In Progress')}
                                  title={t('dashboard.markAsInProgress')}
                                >
                                  <i className="bi bi-arrow-right-circle"></i>
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
                                {t('dashboard.edit')}
                              </button>
                              <button
                                className="btn btn-sm btn-outline-danger"
                                onClick={() => handleDelete(task._id)}
                              >
                                {t('dashboard.delete')}
                              </button>
                            </div>
                          </div>

                          {/* 快速状态更新按钮 */}
                          {task.status !== 'Done' && (
                            <button
                              className="btn btn-sm btn-success"
                              onClick={() => updateTaskStatus(task, 'Done')}
                              title={t('dashboard.markAsDone')}
                            >
                              <i className="bi bi-check-circle"></i>
                            </button>
                          )}
                          {task.status !== 'In Progress' && task.status !== 'Done' && (
                            <button
                              className="btn btn-sm btn-warning"
                              onClick={() => updateTaskStatus(task, 'In Progress')}
                              title={t('dashboard.markAsInProgress')}
                            >
                              <i className="bi bi-arrow-right-circle"></i>
                            </button>
                          )}

                          {/* 任务详情和评论 */}
                          {(expandedTaskId === task._id || task.status === 'In Progress') && (
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
export type TaskComment = {
  text: string;
  createdAt: string;
};
import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { fetchTasks, deleteTask, createTask, updateTask } from '../features/tasks/taskSlice';
import type { RootState, AppDispatch } from '../app/store';
import type { Task } from '../features/tasks/taskSlice';

const DashboardPage: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  const { tasks, isLoading, error } = useSelector((state: RootState) => state.tasks);

  const [newTask, setNewTask] = useState({
    title: '',
    description: '',
    status: 'To Do' as 'To Do' | 'In Progress' | 'Done',
    priority: 'Medium' as 'Low' | 'Medium' | 'High',
    assignee: '',
  });

  const [editingTask, setEditingTask] = useState<Task | null>(null);
  const [commentText, setCommentText] = useState<{[key: string]: string}>({});
  const [filterStatus, setFilterStatus] = useState<string>('All');
  const [filterPriority, setFilterPriority] = useState<string>('All');
  const [sortOrder, setSortOrder] = useState<string>('newest');

  useEffect(() => {
    dispatch(fetchTasks());
  }, [dispatch]);

  const handleDelete = (id: string) => {
    if (window.confirm('Are you sure you want to delete this task?')) {
      dispatch(deleteTask(id));
    }
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

  const handleEditTask = (task: Task) => {
    setEditingTask(task);
  };

  const handleUpdateTask = () => {
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

  const handleStatusChange = (taskId: string, newStatus: 'To Do' | 'In Progress' | 'Done', comment?: string) => {
    const taskToUpdate = tasks.find(task => task._id === taskId);
    if (taskToUpdate) {
      const updatedTask = { ...taskToUpdate, status: newStatus };
      
      // 如果有备注，则添加到任务中
      if (comment) {
        const newComment = {
          text: comment,
          createdAt: new Date().toISOString(),
        };
        
        updatedTask.comments = updatedTask.comments ? [...updatedTask.comments, newComment] : [newComment];
      }
      
      dispatch(updateTask(updatedTask));
      
      // 清除备注输入
      setCommentText(prev => {
        const newCommentText = { ...prev };
        delete newCommentText[taskId];
        return newCommentText;
      });
    }
  };

  const handleAddComment = (taskId: string) => {
    const comment = commentText[taskId];
    if (comment) {
      handleStatusChange(taskId, 'In Progress', comment); // 默认更新为"In Progress"状态并添加备注
    }
  };

  const handleCommentChange = (taskId: string, text: string) => {
    setCommentText(prev => ({
      ...prev,
      [taskId]: text
    }));
  };

  const filteredAndSortedTasks = tasks
    .filter((task: Task) => {
      if (filterStatus !== 'All' && task.status !== filterStatus) return false;
      if (filterPriority !== 'All' && task.priority !== filterPriority) return false;
      return true;
    })
    .sort((a: Task, b: Task) => {
      if (sortOrder === 'newest') {
        return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
      } else if (sortOrder === 'oldest') {
        return new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime();
      } else if (sortOrder === 'priority') {
        const priorityOrder: { [key: string]: number } = { 
          'High': 3,
          'Medium': 2,
          'Low': 1
        };
        return priorityOrder[b.priority] - priorityOrder[a.priority];
      }
      return 0;
    });

  // Get task counts by status
  const taskCounts = {
    total: tasks.length,
    todo: tasks.filter((task: Task) => task.status === 'To Do').length,
    inProgress: tasks.filter((task: Task) => task.status === 'In Progress').length,
    done: tasks.filter((task: Task) => task.status === 'Done').length,
  };

  if (isLoading && tasks.length === 0) return (
    <div className="container py-5">
      <div className="d-flex justify-content-center align-items-center" style={{ height: '70vh' }}>
        <div className="text-center">
          <div className="spinner-border text-primary" role="status">
            <span className="visually-hidden">Loading tasks...</span>
          </div>
          <p className="mt-3">Loading your tasks...</p>
        </div>
      </div>
    </div>
  );
  
  if (error && tasks.length === 0) return (
    <div className="container py-5">
      <div className="alert alert-danger" role="alert">
        <h4 className="alert-heading">Error!</h4>
        <p>Error: {error}</p>
        <hr />
        <button className="btn btn-outline-danger" onClick={() => dispatch(fetchTasks())}>
          Try Again
        </button>
      </div>
    </div>
  );

  return (
    <div className="container py-5">
      <div className="row">
        <div className="col-12">
          <div className="d-flex justify-content-between align-items-center mb-4">
            <h1 className="fw-bold mb-0">
              <i className="bi bi-list-check me-2"></i>
              My Tasks
            </h1>
            <span className="badge bg-primary fs-6">{taskCounts.total} tasks</span>
          </div>
          
          {/* Stats Cards */}
          <div className="row mb-4">
            <div className="col-md-4 mb-3">
              <div className="card border-primary border-2 h-100">
                <div className="card-body text-center">
                  <h5 className="card-title text-primary">To Do</h5>
                  <h2 className="fw-bold text-primary">{taskCounts.todo}</h2>
                </div>
              </div>
            </div>
            <div className="col-md-4 mb-3">
              <div className="card border-warning border-2 h-100">
                <div className="card-body text-center">
                  <h5 className="card-title text-warning">In Progress</h5>
                  <h2 className="fw-bold text-warning">{taskCounts.inProgress}</h2>
                </div>
              </div>
            </div>
            <div className="col-md-4 mb-3">
              <div className="card border-success border-2 h-100">
                <div className="card-body text-center">
                  <h5 className="card-title text-success">Done</h5>
                  <h2 className="fw-bold text-success">{taskCounts.done}</h2>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="row">
        <div className="col-lg-4 mb-4">
          <div className="card shadow-sm h-100">
            <div className="card-header bg-primary text-white">
              <h5 className="mb-0">
                <i className="bi bi-plus-circle me-2"></i>
                Create New Task
              </h5>
            </div>
            <div className="card-body">
              <form onSubmit={handleCreateTask}>
                <div className="mb-3">
                  <label className="form-label">Title</label>
                  <input
                    type="text"
                    className="form-control"
                    placeholder="Task Title"
                    name="title"
                    value={newTask.title}
                    onChange={handleNewTaskChange}
                    required
                  />
                </div>
                <div className="mb-3">
                  <label className="form-label">Description</label>
                  <textarea
                    className="form-control"
                    placeholder="Task Description"
                    name="description"
                    value={newTask.description}
                    onChange={handleNewTaskChange}
                    rows={3}
                  ></textarea>
                </div>
                <div className="row mb-3">
                  <div className="col-md-6">
                    <label className="form-label">Status</label>
                    <select
                      className="form-select"
                      name="status"
                      value={newTask.status}
                      onChange={handleNewTaskChange}
                    >
                      <option value="To Do">To Do</option>
                      <option value="In Progress">In Progress</option>
                      <option value="Done">Done</option>
                    </select>
                  </div>
                  <div className="col-md-6">
                    <label className="form-label">Priority</label>
                    <select
                      className="form-select"
                      name="priority"
                      value={newTask.priority}
                      onChange={handleNewTaskChange}
                    >
                      <option value="Low">Low</option>
                      <option value="Medium">Medium</option>
                      <option value="High">High</option>
                    </select>
                  </div>
                </div>
                <button type="submit" className="btn btn-primary w-100">
                  <i className="bi bi-plus-circle me-1"></i>
                  Create Task
                </button>
              </form>
            </div>
          </div>
        </div>

        <div className="col-lg-8">
          {/* Filters and Sorting */}
          <div className="card shadow-sm mb-4">
            <div className="card-body">
              <div className="row">
                <div className="col-md-5 mb-3 mb-md-0">
                  <label className="form-label">Filter by Status</label>
                  <select 
                    className="form-select"
                    value={filterStatus}
                    onChange={(e) => setFilterStatus(e.target.value)}
                  >
                    <option value="All">All Statuses</option>
                    <option value="To Do">To Do</option>
                    <option value="In Progress">In Progress</option>
                    <option value="Done">Done</option>
                  </select>
                </div>
                <div className="col-md-5 mb-3 mb-md-0">
                  <label className="form-label">Filter by Priority</label>
                  <select 
                    className="form-select"
                    value={filterPriority}
                    onChange={(e) => setFilterPriority(e.target.value)}
                  >
                    <option value="All">All Priorities</option>
                    <option value="Low">Low</option>
                    <option value="Medium">Medium</option>
                    <option value="High">High</option>
                  </select>
                </div>
                <div className="col-md-2">
                  <label className="form-label">Sort By</label>
                  <select 
                    className="form-select"
                    value={sortOrder}
                    onChange={(e) => setSortOrder(e.target.value)}
                  >
                    <option value="newest">Newest</option>
                    <option value="oldest">Oldest</option>
                    <option value="priority">Priority</option>
                  </select>
                </div>
              </div>
            </div>
          </div>

          {/* Task List */}
          <div className="card shadow-sm">
            <div className="card-header bg-white">
              <h5 className="mb-0">
                <i className="bi bi-card-list me-2"></i>
                Task List
              </h5>
            </div>
            <div className="card-body p-0">
              {filteredAndSortedTasks.length === 0 ? (
                <div className="text-center py-5">
                  <i className="bi bi-clipboard-check text-muted" style={{ fontSize: '3rem' }}></i>
                  <p className="mt-3 mb-0 text-muted">No tasks found. Create a new task to get started!</p>
                </div>
              ) : (
                <div className="list-group list-group-flush">
                  {filteredAndSortedTasks.map((task: Task) => (
                    <div key={task._id} className="list-group-item">
                      {editingTask && editingTask._id === task._id ? (
                        // Edit mode
                        <div className="row g-3">
                          <div className="col-12">
                            <label className="form-label">Title</label>
                            <input
                              type="text"
                              className="form-control"
                              name="title"
                              value={editingTask.title}
                              onChange={handleEditChange}
                              required
                            />
                          </div>
                          <div className="col-12">
                            <label className="form-label">Description</label>
                            <textarea
                              className="form-control"
                              name="description"
                              value={editingTask.description}
                              onChange={handleEditChange}
                              rows={2}
                            ></textarea>
                          </div>
                          <div className="col-md-6">
                            <label className="form-label">Status</label>
                            <select
                              className="form-select"
                              name="status"
                              value={editingTask.status}
                              onChange={handleEditChange}
                            >
                              <option value="To Do">To Do</option>
                              <option value="In Progress">In Progress</option>
                              <option value="Done">Done</option>
                            </select>
                          </div>
                          <div className="col-md-6">
                            <label className="form-label">Priority</label>
                            <select
                              className="form-select"
                              name="priority"
                              value={editingTask.priority}
                              onChange={handleEditChange}
                            >
                              <option value="Low">Low</option>
                              <option value="Medium">Medium</option>
                              <option value="High">High</option>
                            </select>
                          </div>
                          <div className="col-12">
                            <button 
                              className="btn btn-primary me-2"
                              onClick={handleUpdateTask}
                            >
                              Save
                            </button>
                            <button 
                              className="btn btn-secondary"
                              onClick={() => setEditingTask(null)}
                            >
                              Cancel
                            </button>
                          </div>
                        </div>
                      ) : (
                        // Display mode
                        <div className="row align-items-center">
                          <div className="col-md-7">
                            <h6 className="mb-1">{task.title}</h6>
                            <p className="mb-2 text-muted small">{task.description}</p>
                            <div className="d-flex flex-wrap gap-2">
                              <span className={`badge ${
                                task.status === 'To Do' ? 'bg-secondary' : 
                                task.status === 'In Progress' ? 'bg-warning' : 'bg-success'
                              }`}>
                                {task.status}
                              </span>
                              <span className={`badge ${
                                task.priority === 'Low' ? 'bg-info' : 
                                task.priority === 'Medium' ? 'bg-primary' : 'bg-danger'
                              }`}>
                                {task.priority}
                              </span>
                            </div>
                            
                            {/* Comments section */}
                            {task.comments && task.comments.length > 0 && (
                              <div className="mt-2">
                                <small className="text-muted">Comments:</small>
                                <div className="small">
                                  {task.comments.slice(-2).map((comment, index) => (
                                    <div key={index} className="text-muted fst-italic">
                                      "{comment.text}"
                                    </div>
                                  ))}
                                  {task.comments.length > 2 && (
                                    <div className="text-muted fst-italic">
                                      + {task.comments.length - 2} more comments
                                    </div>
                                  )}
                                </div>
                              </div>
                            )}
                          </div>
                          <div className="col-md-3">
                            <div className="d-flex flex-column gap-2">
                              <select 
                                className="form-select form-select-sm"
                                value={task.status}
                                onChange={(e) => handleStatusChange(task._id, e.target.value as any)}
                              >
                                <option value="To Do">To Do</option>
                                <option value="In Progress">In Progress</option>
                                <option value="Done">Done</option>
                              </select>
                              
                              {/* Comment input */}
                              <div className="input-group input-group-sm">
                                <input
                                  type="text"
                                  className="form-control"
                                  placeholder="Add comment..."
                                  value={commentText[task._id] || ''}
                                  onChange={(e) => handleCommentChange(task._id, e.target.value)}
                                />
                                <button 
                                  className="btn btn-outline-primary btn-sm"
                                  type="button"
                                  onClick={() => handleAddComment(task._id)}
                                >
                                  <i className="bi bi-chat"></i>
                                </button>
                              </div>
                            </div>
                          </div>
                          <div className="col-md-2 text-end">
                            <button 
                              className="btn btn-outline-primary btn-sm me-1"
                              onClick={() => handleEditTask(task)}
                            >
                              <i className="bi bi-pencil"></i>
                            </button>
                            <button 
                              className="btn btn-outline-danger btn-sm"
                              onClick={() => handleDelete(task._id)}
                            >
                              <i className="bi bi-trash"></i>
                            </button>
                          </div>
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Edit Task Modal - Fixed the issue with modal display */}
      {editingTask && (
        <div className="modal fade show d-block" tabIndex={-1} style={{ backgroundColor: 'rgba(0,0,0,0.5)', zIndex: 1050 }}>
          <div className="modal-dialog">
            <div className="modal-content">
              <div className="modal-header">
                <h5 className="modal-title">Edit Task</h5>
                <button 
                  type="button" 
                  className="btn-close" 
                  onClick={() => setEditingTask(null)}
                ></button>
              </div>
              <div className="modal-body">
                <div className="row g-3">
                  <div className="col-12">
                    <label className="form-label">Title</label>
                    <input
                      type="text"
                      className="form-control"
                      name="title"
                      value={editingTask.title}
                      onChange={handleEditChange}
                      required
                    />
                  </div>
                  <div className="col-12">
                    <label className="form-label">Description</label>
                    <textarea
                      className="form-control"
                      name="description"
                      value={editingTask.description}
                      onChange={handleEditChange}
                      rows={3}
                    ></textarea>
                  </div>
                  <div className="col-md-6">
                    <label className="form-label">Status</label>
                    <select
                      className="form-select"
                      name="status"
                      value={editingTask.status}
                      onChange={handleEditChange}
                    >
                      <option value="To Do">To Do</option>
                      <option value="In Progress">In Progress</option>
                      <option value="Done">Done</option>
                    </select>
                  </div>
                  <div className="col-md-6">
                    <label className="form-label">Priority</label>
                    <select
                      className="form-select"
                      name="priority"
                      value={editingTask.priority}
                      onChange={handleEditChange}
                    >
                      <option value="Low">Low</option>
                      <option value="Medium">Medium</option>
                      <option value="High">High</option>
                    </select>
                  </div>
                  <div className="col-12">
                    <div className="d-flex justify-content-between">
                      <button 
                        className="btn btn-primary"
                        onClick={handleUpdateTask}
                      >
                        Save Changes
                      </button>
                      <button 
                        className="btn btn-secondary"
                        onClick={() => setEditingTask(null)}
                      >
                        Cancel
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default DashboardPage;
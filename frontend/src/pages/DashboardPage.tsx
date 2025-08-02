import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { fetchTasks, deleteTask, createTask } from '../features/tasks/taskSlice';
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

  const [filterStatus, setFilterStatus] = useState<string>('All');
  const [filterPriority, setFilterPriority] = useState<string>('All');
  const [sortOrder, setSortOrder] = useState<string>('newest');

  useEffect(() => {
    dispatch(fetchTasks());
  }, [dispatch]);

  const handleDelete = (id: string) => {
    dispatch(deleteTask(id));
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
                <div className="mb-3">
                  <label className="form-label">Assignee</label>
                  <input
                    type="text"
                    className="form-control"
                    placeholder="Assignee (Optional)"
                    name="assignee"
                    value={newTask.assignee}
                    onChange={handleNewTaskChange}
                  />
                </div>
                <div className="d-grid">
                  <button type="submit" className="btn btn-primary">
                    <i className="bi bi-plus-circle me-2"></i>
                    Add Task
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>

        <div className="col-lg-8">
          <div className="card shadow-sm">
            <div className="card-header bg-white">
              <div className="row align-items-center">
                <div className="col-md-4 mb-2 mb-md-0">
                  <h5 className="mb-0">Task List</h5>
                </div>
                <div className="col-md-8">
                  <div className="row">
                    <div className="col-md-4 mb-2">
                      <select
                        className="form-select form-select-sm"
                        value={filterStatus}
                        onChange={(e) => setFilterStatus(e.target.value)}
                      >
                        <option value="All">All Statuses</option>
                        <option value="To Do">To Do</option>
                        <option value="In Progress">In Progress</option>
                        <option value="Done">Done</option>
                      </select>
                    </div>
                    <div className="col-md-4 mb-2">
                      <select
                        className="form-select form-select-sm"
                        value={filterPriority}
                        onChange={(e) => setFilterPriority(e.target.value)}
                      >
                        <option value="All">All Priorities</option>
                        <option value="Low">Low</option>
                        <option value="Medium">Medium</option>
                        <option value="High">High</option>
                      </select>
                    </div>
                    <div className="col-md-4">
                      <select
                        className="form-select form-select-sm"
                        value={sortOrder}
                        onChange={(e) => setSortOrder(e.target.value)}
                      >
                        <option value="newest">Newest First</option>
                        <option value="oldest">Oldest First</option>
                        <option value="priority">By Priority</option>
                      </select>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div className="card-body">
              {isLoading && tasks.length > 0 && (
                <div className="text-center py-3">
                  <div className="spinner-border text-primary" role="status">
                    <span className="visually-hidden">Updating tasks...</span>
                  </div>
                </div>
              )}
              
              {filteredAndSortedTasks.length === 0 ? (
                <div className="text-center py-5">
                  <i className="bi bi-clipboard-check text-muted fs-1 mb-3"></i>
                  <h5 className="text-muted">No tasks found</h5>
                  <p className="text-muted">Try changing your filters or create a new task</p>
                </div>
              ) : (
                <div className="row">
                  {filteredAndSortedTasks.map((task) => (
                    <div key={task._id} className="col-md-6 mb-3">
                      <div className="card h-100 border">
                        <div className="card-body d-flex flex-column">
                          <div className="d-flex justify-content-between align-items-start mb-2">
                            <h5 className="card-title mb-0">{task.title}</h5>
                            <span className={`badge ${
                              task.priority === 'High' ? 'bg-danger' : 
                              task.priority === 'Medium' ? 'bg-warning text-dark' : 'bg-success'
                            }`}>
                              {task.priority}
                            </span>
                          </div>
                          
                          {task.description && (
                            <p className="card-text flex-grow-1 text-muted">
                              {task.description}
                            </p>
                          )}
                          
                          <div className="mb-3">
                            <span className={`badge ${
                              task.status === 'To Do' ? 'bg-primary' : 
                              task.status === 'In Progress' ? 'bg-warning text-dark' : 'bg-success'
                            }`}>
                              {task.status}
                            </span>
                            {task.assignee && (
                              <span className="badge bg-secondary ms-2">
                                <i className="bi bi-person-fill me-1"></i>
                                {task.assignee}
                              </span>
                            )}
                          </div>
                          
                          <div className="mt-auto">
                            <small className="text-muted">
                              Created: {new Date(task.createdAt).toLocaleDateString()}
                            </small>
                          </div>
                          
                          <div className="mt-3">
                            <button
                              className="btn btn-outline-danger btn-sm"
                              onClick={() => handleDelete(task._id)}
                            >
                              <i className="bi bi-trash me-1"></i>
                              Delete
                            </button>
                          </div>
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
    </div>
  );
}; 

export default DashboardPage;
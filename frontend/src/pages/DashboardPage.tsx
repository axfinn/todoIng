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

  if (isLoading) return <div className="text-center mt-5">Loading tasks...</div>;
  if (error) return <div className="alert alert-danger mt-5">Error: {error}</div>;

  return (
    <div className="container mt-5">
      <h1 className="mb-4">My Tasks</h1>

      <div className="card mb-4">
        <div className="card-header">Create New Task</div>
        <div className="card-body">
          <form onSubmit={handleCreateTask}>
            <div className="mb-3">
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
              <textarea
                className="form-control"
                placeholder="Task Description"
                name="description"
                value={newTask.description}
                onChange={handleNewTaskChange}
              ></textarea>
            </div>
            <div className="row mb-3">
              <div className="col">
                <select
                  className="form-select"
                  name="status"
                  value={newTask.status}
                  onChange={handleNewTaskChange}
                >
                  <option value="To Do">To Do</option>
                  <option value="In Progress">In Progress</option>
                  <option value="Done">Done</option>
                  <option value="created">Created</option>
                  <option value="in-progress">In Progress (Legacy)</option>
                  <option value="paused">Paused</option>
                  <option value="completed">Completed</option>
                  <option value="cancelled">Cancelled</option>
                </select>
              </div>
              <div className="col">
                <select
                  className="form-select"
                  name="priority"
                  value={newTask.priority}
                  onChange={handleNewTaskChange}
                >
                  <option value="Low">Low</option>
                  <option value="Medium">Medium</option>
                  <option value="High">High</option>
                  <option value="low">Low (Legacy)</option>
                  <option value="medium">Medium (Legacy)</option>
                  <option value="high">High (Legacy)</option>
                </select>
              </div>
            </div>
            <div className="mb-3">
              <input
                type="text"
                className="form-control"
                placeholder="Assignee (Optional)"
                name="assignee"
                value={newTask.assignee}
                onChange={handleNewTaskChange}
              />
            </div>
            <button type="submit" className="btn btn-primary">
              Add Task
            </button>
          </form>
        </div>
      </div>

      <div className="row mb-3">
        <div className="col">
          <select
            className="form-select"
            value={filterStatus}
            onChange={(e) => setFilterStatus(e.target.value)}
          >
            <option value="All">All Statuses</option>
            <option value="To Do">To Do</option>
            <option value="In Progress">In Progress</option>
            <option value="Done">Done</option>
            <option value="created">Created</option>
            <option value="in-progress">In Progress (Legacy)</option>
            <option value="paused">Paused</option>
            <option value="completed">Completed</option>
            <option value="cancelled">Cancelled</option>
          </select>
        </div>
        <div className="col">
          <select
            className="form-select"
            value={filterPriority}
            onChange={(e) => setFilterPriority(e.target.value)}
          >
            <option value="All">All Priorities</option>
            <option value="Low">Low</option>
            <option value="Medium">Medium</option>
            <option value="High">High</option>
            <option value="low">Low (Legacy)</option>
            <option value="medium">Medium (Legacy)</option>
            <option value="high">High (Legacy)</option>
          </select>
        </div>
        <div className="col">
          <select
            className="form-select"
            value={sortOrder}
            onChange={(e) => setSortOrder(e.target.value)}
          >
            <option value="newest">Sort by Newest</option>
            <option value="oldest">Sort by Oldest</option>
            <option value="priority">Sort by Priority</option>
          </select>
        </div>
      </div>

      {filteredAndSortedTasks.length === 0 ? (
        <p className="text-center">No tasks found. Create one!</p>
      ) : (
        <div className="list-group">
          {filteredAndSortedTasks.map((task) => (
            <div key={task._id} className="list-group-item list-group-item-action mb-3">
              <div className="d-flex w-100 justify-content-between">
                <h5 className="mb-1">{task.title}</h5>
                <small className="text-muted">Status: {task.status}</small>
              </div>
              <p className="mb-1">{task.description}</p>
              <small className="text-muted">Priority: {task.priority}</small>
              {task.assignee && <small className="text-muted ms-2">Assignee: {task.assignee}</small>}
              <div className="mt-2">
                <button
                  className="btn btn-danger btn-sm"
                  onClick={() => handleDelete(task._id)}
                >
                  Delete
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}; 

export default DashboardPage;
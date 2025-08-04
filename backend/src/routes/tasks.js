const express = require('express');
const router = express.Router();
const auth = require('../middleware/auth');
const { check, validationResult } = require('express-validator');
const Task = require('../models/Task');

// Create Task
router.post(
  '/',
  [auth, [check('title', 'Title is required').not().isEmpty()]],
  async (req, res) => {
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(400).json({ errors: errors.array() });
    }

    const { title, description, status, priority, assignee, comments, deadline, scheduledDate } = req.body;

    try {
      const newTask = new Task({
        title,
        description,
        status: status || 'To Do',
        priority: priority || 'Medium',
        assignee,
        comments: comments || [],
        deadline: deadline || null,
        scheduledDate: scheduledDate || null,
        createdBy: req.user.id,
      });

      const task = await newTask.save();
      res.json(task);
    } catch (err) {
      console.error(err.message);
      res.status(500).send('Server Error');
    }
  }
);

// Get All Tasks
router.get('/', auth, async (req, res) => {
  try {
    const tasks = await Task.find({ createdBy: req.user.id }).sort({
      createdAt: -1,
    });
    res.json(tasks);
  } catch (err) {
    console.error(err.message);
    res.status(500).send('Server Error');
  }
});

// Get Task by ID
router.get('/:id', auth, async (req, res) => {
  try {
    const task = await Task.findById(req.params.id);

    if (!task) {
      return res.status(404).json({ msg: 'Task not found' });
    }

    // Make sure user owns task
    if (task.createdBy.toString() !== req.user.id) {
      return res.status(401).json({ msg: 'User not authorized' });
    }

    res.json(task);
  } catch (err) {
    console.error(err.message);
    res.status(500).send('Server Error');
  }
});

// Update Task
router.put(
  '/:id',
  [
    auth,
    [
      check('title', 'Title is required').optional().not().isEmpty(),
      check('status', 'Invalid status')
        .optional()
        .isIn(['To Do', 'In Progress', 'Done']),
      check('priority', 'Invalid priority')
        .optional()
        .isIn(['Low', 'Medium', 'High']),
    ],
  ],
  async (req, res) => {
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(400).json({ errors: errors.array() });
    }

    const { title, description, status, priority, assignee, comments, deadline, scheduledDate } = req.body;

    // Build task object
    const taskFields = {};
    if (title) taskFields.title = title;
    if (description) taskFields.description = description;
    if (status) taskFields.status = status;
    if (priority) taskFields.priority = priority;
    if (assignee !== undefined) taskFields.assignee = assignee;
    if (deadline !== undefined) taskFields.deadline = deadline;
    if (scheduledDate !== undefined) taskFields.scheduledDate = scheduledDate;

    // 处理评论，确保添加创建者信息
    if (comments) {
      taskFields.comments = comments.map(comment => {
        // 如果评论没有创建者信息，则添加当前用户作为创建者
        if (!comment.createdBy) {
          return {
            ...comment,
            createdBy: req.user.id
          };
        }
        return comment;
      });
    }

    try {
      let task = await Task.findById(req.params.id);

      if (!task) return res.status(404).json({ msg: 'Task not found' });

      // Make sure user owns task
      if (task.createdBy.toString() !== req.user.id) {
        return res.status(401).json({ msg: 'User not authorized' });
      }

      task = await Task.findByIdAndUpdate(
        req.params.id,
        { $set: taskFields },
        { new: true }
      );

      res.json(task);
    } catch (err) {
      console.error(err.message);
      res.status(500).send('Server Error');
    }
  }
);

// Delete Task
router.delete('/:id', auth, async (req, res) => {
  try {
    const task = await Task.findById(req.params.id);

    if (!task) {
      return res.status(404).json({ msg: 'Task not found' });
    }

    // Make sure user owns task
    if (task.createdBy.toString() !== req.user.id) {
      return res.status(401).json({ msg: 'User not authorized' });
    }

    await Task.findByIdAndDelete(req.params.id);

    res.json({ msg: 'Task removed' });
  } catch (err) {
    console.error(err.message);
    res.status(500).send('Server Error');
  }
});

// Helper function to format date for filename
function formatDateForFilename(date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0'); // Months are zero-based
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}

// Export all tasks
router.get('/export/all', auth, async (req, res) => {
  try {
    const tasks = await Task.find({ createdBy: req.user.id });
    
    // Set response headers
    res.setHeader('Content-Type', 'application/json');
    res.setHeader('Content-Disposition', `attachment; filename="todoing-backup-${formatDateForFilename(new Date())}.json"`);
    
    // Send tasks data
    res.json(tasks);
  } catch (err) {
    console.error(err.message);
    res.status(500).send('Server Error');
  }
});

// Import tasks
router.post('/import', auth, async (req, res) => {
  try {
    const tasksData = req.body.tasks;
    
    if (!Array.isArray(tasksData)) {
      return res.status(400).json({ msg: 'Invalid data format. Expected an array of tasks.' });
    }
    
    // Validate and import tasks
    const importedTasks = [];
    const errors = [];
    
    for (let i = 0; i < tasksData.length; i++) {
      const taskData = tasksData[i];
      
      try {
        // Remove _id field to avoid conflicts
        delete taskData._id;
        
        // Set creator to current user
        taskData.createdBy = req.user.id;
        
        // Create new task
        const newTask = new Task(taskData);
        const savedTask = await newTask.save();
        importedTasks.push(savedTask);
      } catch (error) {
        errors.push({ index: i, error: error.message });
      }
    }
    
    res.json({
      msg: `Imported ${importedTasks.length} tasks successfully`,
      imported: importedTasks.length,
      errors
    });
  } catch (err) {
    console.error(err.message);
    res.status(500).send('Server Error');
  }
});

module.exports = router;
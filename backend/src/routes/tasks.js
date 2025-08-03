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
    if (comments) taskFields.comments = comments;
    if (deadline !== undefined) taskFields.deadline = deadline;
    if (scheduledDate !== undefined) taskFields.scheduledDate = scheduledDate;

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

    await Task.findByIdAndRemove(req.params.id);

    res.json({ msg: 'Task removed' });
  } catch (err) {
    console.error(err.message);
    res.status(500).send('Server Error');
  }
});

module.exports = router;
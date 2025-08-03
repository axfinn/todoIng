const express = require('express');
const router = express.Router();
const auth = require('../middleware/auth');
const Report = require('../models/Report');
const Task = require('../models/Task');

// 获取用户的所有报告
router.get('/', auth, async (req, res) => {
  try {
    const reports = await Report.find({ userId: req.user.id })
      .sort({ createdAt: -1 });
    res.json(reports);
  } catch (err) {
    console.error(err.message);
    res.status(500).send('Server Error');
  }
});

// 获取特定报告
router.get('/:id', auth, async (req, res) => {
  try {
    const report = await Report.findById(req.params.id)
      .populate('tasks', ['title', 'description', 'status', 'priority', 'deadline', 'createdAt', 'updatedAt', 'comments']);

    if (!report) {
      return res.status(404).json({ msg: 'Report not found' });
    }

    if (report.userId.toString() !== req.user.id) {
      return res.status(401).json({ msg: 'User not authorized' });
    }

    res.json(report);
  } catch (err) {
    console.error(err.message);
    if (err.kind === 'ObjectId') {
      return res.status(404).json({ msg: 'Report not found' });
    }
    res.status(500).send('Server Error');
  }
});

// 生成报告
router.post('/generate', auth, async (req, res) => {
  try {
    const { type, period, startDate, endDate } = req.body;
    
    // 验证参数
    if (!type || !period || !startDate || !endDate) {
      return res.status(400).json({ msg: 'Please provide type, period, startDate, and endDate' });
    }
    
    // 获取该时间段内的任务
    const tasks = await Task.find({
      createdBy: req.user.id,
      createdAt: {
        $gte: new Date(startDate),
        $lte: new Date(endDate)
      }
    });
    
    // 填充任务的评论创建者信息
    await Task.populate(tasks, {
      path: 'comments.createdBy',
      select: 'username'
    });
    
    // 统计数据
    const totalTasks = tasks.length;
    const completedTasks = tasks.filter(task => task.status === 'Done').length;
    const inProgressTasks = tasks.filter(task => task.status === 'In Progress').length;
    const overdueTasks = tasks.filter(task => 
      task.deadline && new Date(task.deadline) < new Date() && task.status !== 'Done'
    ).length;
    const completionRate = totalTasks > 0 ? Math.round((completedTasks / totalTasks) * 100) : 0;
    
    // 生成报告标题
    const titles = {
      daily: `日报 - ${period}`,
      weekly: `周报 - ${period}`,
      monthly: `月报 - ${period}`
    };
    
    // 生成报告内容
    let content = `# ${titles[type]}\n\n`;
    content += `报告周期: ${new Date(startDate).toLocaleDateString('zh-CN')} - ${new Date(endDate).toLocaleDateString('zh-CN')}\n\n`;
    
    content += `## 统计信息\n`;
    content += `- 总任务数: ${totalTasks}\n`;
    content += `- 已完成任务: ${completedTasks}\n`;
    content += `- 进行中任务: ${inProgressTasks}\n`;
    content += `- 过期任务: ${overdueTasks}\n`;
    content += `- 完成率: ${completionRate}%\n\n`;
    
    content += `## 任务详情\n`;
    if (tasks.length > 0) {
      // 按任务聚合显示任务详情
      // 填充任务的评论创建者信息
      await Task.populate(tasks, {
        path: 'comments.createdBy',
        select: 'username'
      });
      
      tasks.forEach(task => {
        content += `### 任务: ${task.title}\n`;
        content += `- **任务状态**: ${task.status}\n`;
        content += `- **任务优先级**: ${task.priority}\n`;
        content += `- **创建时间**: ${task.createdAt ? task.createdAt.toLocaleString('zh-CN') : '未知'}\n`;
        content += `- **更新时间**: ${task.updatedAt ? task.updatedAt.toLocaleString('zh-CN') : '未知'}\n`;
        if (task.deadline) {
          content += `- **截止日期**: ${task.deadline.toLocaleString('zh-CN')}\n`;
        }
        if (task.scheduledDate) {
          content += `- **计划日期**: ${task.scheduledDate.toLocaleString('zh-CN')}\n`;
        }
        content += `- **任务描述**: ${task.description || '无'}\n`;
        
        // 添加任务相关事件时间线
        content += `\n#### 任务活动时间线\n`;
        
        // 收集该任务的所有事件
        const taskEvents = [];
        
        // 添加任务创建事件
        taskEvents.push({
          date: task.createdAt,
          message: '任务已创建'
        });
        
        // 添加任务更新事件
        if (task.updatedAt && task.updatedAt.getTime() !== task.createdAt.getTime()) {
          taskEvents.push({
            date: task.updatedAt,
            message: '任务已更新'
          });
        }
        
        // 添加任务评论事件
        if (task.comments && task.comments.length > 0) {
          task.comments.forEach(comment => {
            taskEvents.push({
              date: comment.createdAt,
              message: `添加了新评论: ${comment.text.substring(0, 30)}${comment.text.length > 30 ? '...' : ''}`
            });
          });
        }
        
        // 按时间排序事件
        taskEvents.sort((a, b) => a.date - b.date);
        
        // 显示任务事件时间线
        if (taskEvents.length > 0) {
          taskEvents.forEach(event => {
            content += `- ${event.date.toLocaleString('zh-CN')}: ${event.message}\n`;
          });
        } else {
          content += `- 暂无活动记录\n`;
        }
        
        // 显示详细评论信息
        if (task.comments && task.comments.length > 0) {
          content += `\n#### 评论详情\n`;
          task.comments.forEach((comment, index) => {
            content += `${index + 1}. **评论内容**: ${comment.text}\n`;
            content += `   **评论人**: ${comment.createdBy ? comment.createdBy.username : '未知用户'}\n`;
            content += `   **评论时间**: ${comment.createdAt.toLocaleString('zh-CN')}\n\n`;
          });
        }
        
        content += '\n---\n\n';
      });
    } else {
      content += `此周期内未找到任务。\n`;
    }
    
    // 创建报告
    const report = new Report({
      userId: req.user.id,
      type,
      period,
      title: titles[type],
      content,
      tasks: tasks.map(task => task.id),
      statistics: {
        totalTasks,
        completedTasks,
        inProgressTasks,
        overdueTasks,
        completionRate
      }
    });
    
    const savedReport = await report.save();
    res.json(savedReport);
  } catch (err) {
    console.error(err.message);
    res.status(500).send('Server Error');
  }
});

// 删除报告
router.delete('/:id', auth, async (req, res) => {
  try {
    const report = await Report.findById(req.params.id);

    if (!report) {
      return res.status(404).json({ msg: 'Report not found' });
    }

    // 修复用户授权检查逻辑
    if (report.userId.toString() !== req.user.id.toString()) {
      return res.status(401).json({ msg: 'User not authorized' });
    }

    // 使用findByIdAndDelete替代已弃用的findByIdAndRemove
    await Report.findByIdAndDelete(req.params.id);

    res.json({ msg: 'Report removed' });
  } catch (err) {
    console.error(err.message);
    if (err.kind === 'ObjectId') {
      return res.status(404).json({ msg: 'Report not found' });
    }
    res.status(500).send('Server Error');
  }
});

// 导出报告
router.get('/:id/export/:format', auth, async (req, res) => {
  try {
    const { id, format } = req.params;
    
    // 获取报告
    const report = await Report.findById(id);
    
    if (!report) {
      return res.status(404).json({ msg: 'Report not found' });
    }
    
    if (report.userId.toString() !== req.user.id) {
      return res.status(401).json({ msg: 'User not authorized' });
    }
    
    // 设置响应头
    let contentType = 'text/plain';
    let fileExtension = 'txt';
    
    switch (format.toLowerCase()) {
      case 'md':
      case 'markdown':
        contentType = 'text/markdown';
        fileExtension = 'md';
        break;
      case 'txt':
        contentType = 'text/plain';
        fileExtension = 'txt';
        break;
    }
    
    // 设置响应头
    res.setHeader('Content-Type', contentType);
    res.setHeader('Content-Disposition', `attachment; filename="report-${report.period}.${fileExtension}"`);
    
    // 发送报告内容
    const content = report.polishedContent || report.content;
    res.send(content);
  } catch (err) {
    console.error(err.message);
    if (err.kind === 'ObjectId') {
      return res.status(404).json({ msg: 'Report not found' });
    }
    res.status(500).send('Server Error');
  }
});

module.exports = router;
## [1.8.18] - 2025-08-07
### 修复
- 修复 Go 后端的响应格式，使其与 backend 保持一致
  - GET /api/tasks - 现在直接返回任务数组，而不是 {"tasks": [...], "message": "..."}
  - GET /api/tasks/:id - 现在直接返回任务对象，而不是 {"task": {...}, "message": "..."}
  - POST /api/tasks - 现在直接返回创建的任务对象
  - PUT /api/tasks/:id - 现在直接返回更新后的任务对象
  - DELETE /api/tasks/:id - 现在返回 {"msg": "Task removed"} 而不是 {"message": "Task deleted successfully"}
  - POST /api/tasks/:id/assign - 现在直接返回更新后的任务对象
  - POST /api/tasks/:id/comments - 现在直接返回更新后的任务对象
- 更新了相关的 Swagger 文档类型定义
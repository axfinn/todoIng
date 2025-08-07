## [1.8.19] - 2025-08-07
### 新增
- 完善报告接口实现
  - 为所有报告接口添加了认证中间件
  - 实现了从 JWT token 中正确提取用户 ID（路径：claims["user"].(map[string]interface{})["id"]）
  - 实现了真实的数据库操作，替换了之前的模拟数据
  - 实现了报告生成的业务逻辑，包括：
    - 根据时间范围查询用户任务
    - 计算任务统计信息（总数、完成数、进行中、逾期、完成率）
    - 生成格式化的报告内容（Markdown 格式）
    - 支持日报、周报、月报三种类型
  - 实现了报告导出功能，支持 Markdown 和纯文本格式
  - 确保所有响应格式与 backend 一致（直接返回数组或对象）
  - 修复了所有错误响应格式（使用 msg 而不是 error）

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
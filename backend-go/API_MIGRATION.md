# API 接口迁移进度文档

此文档用于跟踪从 Node.js 后端到 Go 后端的 API 接口迁移进度。

## 已完成的接口

### 认证相关接口

| 接口路径 | 方法 | 功能 | Protobuf 定义 | 实现状态 |
|---------|------|------|---------------|----------|
| `/api/auth/register` | POST | 用户注册 | `RegisterRequest/RegisterResponse` | ✅ 已定义 |
| `/api/auth/login` | POST | 用户登录 | `LoginRequest/LoginResponse` | ✅ 已定义 |
| `/api/auth/me` | GET | 获取当前用户信息 | `GetUserRequest/GetUserResponse` | ✅ 已定义 |
| `/api/auth/captcha` | GET | 生成验证码 | `GenerateCaptchaResponse` | ✅ 已定义 |
| `/api/auth/verify-captcha` | POST | 验证验证码 | `CaptchaRequest/CaptchaResponse` | ✅ 已定义 |
| `/api/auth/send-email-code` | POST | 发送注册邮箱验证码 | `SendEmailCodeRequest/SendEmailCodeResponse` | ✅ 已定义 |
| `/api/auth/send-login-email-code` | POST | 发送登录邮箱验证码 | `SendLoginEmailCodeRequest/SendLoginEmailCodeResponse` | ✅ 已定义 |

### 任务相关接口

| 接口路径 | 方法 | 功能 | Protobuf 定义 | 实现状态 |
|---------|------|------|---------------|----------|
| `/api/tasks` | POST | 创建任务 | `CreateTaskRequest/CreateTaskResponse` | ✅ 已定义 |
| `/api/tasks` | GET | 获取所有任务 | `GetTasksRequest/GetTasksResponse` | ✅ 已定义 |
| `/api/tasks/:id` | GET | 获取特定任务 | `GetTaskRequest/GetTaskResponse` | ✅ 已定义 |
| `/api/tasks/:id` | PUT | 更新任务 | `UpdateTaskRequest/UpdateTaskResponse` | ✅ 已定义 |
| `/api/tasks/:id` | DELETE | 删除任务 | `DeleteTaskRequest/DeleteTaskResponse` | ✅ 已定义 |
| `/api/tasks/:id/assign` | POST | 分配任务 | `AssignTaskRequest/AssignTaskResponse` | ✅ 已定义 |
| `/api/tasks/:id/comments` | POST | 添加评论 | `AddCommentRequest/AddCommentResponse` | ✅ 已定义 |
| `/api/tasks/export/all` | GET | 导出所有任务 | `ExportTasksRequest/ExportTasksResponse` | ✅ 已定义 |
| `/api/tasks/import` | POST | 导入任务 | `ImportTasksRequest/ImportTasksResponse` | ✅ 已定义 |

### 报告相关接口

| 接口路径 | 方法 | 功能 | Protobuf 定义 | 实现状态 |
|---------|------|------|---------------|----------|
| `/api/reports` | GET | 获取所有报告 | `GetReportsRequest/GetReportsResponse` | ✅ 已定义 |
| `/api/reports/:id` | GET | 获取特定报告 | `GetReportRequest/GetReportResponse` | ✅ 已定义 |
| `/api/reports/generate` | POST | 生成报告 | `GenerateReportRequest/GenerateReportResponse` | ✅ 已定义 |
| `/api/reports/:id/polish` | POST | AI润色报告 | `PolishReportRequest/PolishReportResponse` | ✅ 已定义 |
| `/api/reports/:id/export/:format` | GET | 导出报告 | `ExportReportRequest/ExportReportResponse` | ✅ 已定义 |
| `/api/reports/:id` | DELETE | 删除报告 | `DeleteReportRequest/DeleteReportResponse` | ✅ 已定义 |

## 待实现的接口

### 认证相关接口实现

- [x] 实现 `/api/auth/register` 接口
- [x] 实现 `/api/auth/login` 接口
- [x] 实现 `/api/auth/me` 接口
- [x] 实现 `/api/auth/captcha` 接口
- [x] 实现 `/api/auth/verify-captcha` 接口
- [x] 实现 `/api/auth/send-email-code` 接口
- [x] 实现 `/api/auth/send-login-email-code` 接口

### 任务相关接口实现

- [x] 实现 `/api/tasks` POST 接口
- [x] 实现 `/api/tasks` GET 接口
- [x] 实现 `/api/tasks/:id` GET 接口
- [x] 实现 `/api/tasks/:id` PUT 接口
- [x] 实现 `/api/tasks/:id` DELETE 接口
- [x] 实现 `/api/tasks/:id/assign` POST 接口
- [x] 实现 `/api/tasks/:id/comments` POST 接口
- [x] 实现 `/api/tasks/export/all` GET 接口
- [x] 实现 `/api/tasks/import` POST 接口

### 报告相关接口实现

- [x] 实现 `/api/reports` GET 接口
- [x] 实现 `/api/reports/:id` GET 接口
- [x] 实现 `/api/reports/generate` POST 接口
- [x] 实现 `/api/reports/:id/polish` POST 接口
- [x] 实现 `/api/reports/:id/export/:format` GET 接口
- [x] 实现 `/api/reports/:id` DELETE 接口

## Protobuf 服务定义

### UserService 服务

```protobuf
service UserService {
  // 注册
  rpc Register (RegisterRequest) returns (RegisterResponse);
  // 登录
  rpc Login (LoginRequest) returns (LoginResponse);
  // 获取用户信息
  rpc GetUser (GetUserRequest) returns (GetUserResponse);
  // 发送注册邮箱验证码
  rpc SendEmailCode (SendEmailCodeRequest) returns (SendEmailCodeResponse);
  // 发送登录邮箱验证码
  rpc SendLoginEmailCode (SendLoginEmailCodeRequest) returns (SendLoginEmailCodeResponse);
  // 生成验证码
  rpc GenerateCaptcha (common.Empty) returns (GenerateCaptchaResponse);
  // 验证验证码
  rpc VerifyCaptcha (CaptchaRequest) returns (CaptchaResponse);
}
```

### TaskService 服务

```protobuf
service TaskService {
  // 创建任务
  rpc CreateTask (CreateTaskRequest) returns (CreateTaskResponse);
  // 获取任务列表
  rpc GetTasks (GetTasksRequest) returns (GetTasksResponse);
  // 获取任务
  rpc GetTask (GetTaskRequest) returns (GetTaskResponse);
  // 更新任务
  rpc UpdateTask (UpdateTaskRequest) returns (UpdateTaskResponse);
  // 删除任务
  rpc DeleteTask (DeleteTaskRequest) returns (DeleteTaskResponse);
  // 添加评论
  rpc AddComment (AddCommentRequest) returns (AddCommentResponse);
  // 分配任务
  rpc AssignTask (AssignTaskRequest) returns (AssignTaskResponse);
  // 导入任务
  rpc ImportTasks (ImportTasksRequest) returns (ImportTasksResponse);
  // 导出任务
  rpc ExportTasks (ExportTasksRequest) returns (ExportTasksResponse);
}
```

### ReportService 服务

```protobuf
service ReportService {
  // 获取报告列表
  rpc GetReports (GetReportsRequest) returns (GetReportsResponse);
  // 获取报告
  rpc GetReport (GetReportRequest) returns (GetReportResponse);
  // 生成报告
  rpc GenerateReport (GenerateReportRequest) returns (GenerateReportResponse);
  // 润色报告
  rpc PolishReport (PolishReportRequest) returns (PolishReportResponse);
  // 导出报告
  rpc ExportReport (ExportReportRequest) returns (ExportReportResponse);
  // 删除报告
  rpc DeleteReport (DeleteReportRequest) returns (DeleteReportResponse);
}
```

## 下一步计划

1. ✅ 为每个接口添加 Swagger 注解，便于生成 API 文档
2. 进行全面测试，确保接口功能正确性
3. 部署和上线新版本
4. 逐步将前端请求切换到新的 Go 后端
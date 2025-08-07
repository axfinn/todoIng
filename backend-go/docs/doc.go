// Package docs TodoIng API.
//
// The purpose of this service is to provide a RESTful API for the TodoIng application.
//
//	Schemes: http
//	Host: localhost:5001
//	BasePath: /api
//	Version: 1.0.0
//	License: MIT http://opensource.org/licenses/MIT
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Security:
//	- BearerAuth:
//
//	SecurityDefinitions:
//	BearerAuth:
//	     type: apiKey
//	     name: Authorization
//	     in: header
//	     description: 'Type "Bearer" followed by a space and JWT token.'
//
// swagger:meta
package docs

// GenerateCaptchaResponse 生成验证码响应
// swagger:model
type GenerateCaptchaResponse struct {
	// 验证码
	Captcha int32 `json:"captcha"`
	// 消息
	Message string `json:"message"`
}

// SendEmailCodeRequest 发送邮箱验证码请求
// swagger:model
type SendEmailCodeRequest struct {
	// 邮箱地址
	Email string `json:"email"`
}

// SendEmailCodeResponse 发送邮箱验证码响应
// swagger:model
type SendEmailCodeResponse struct {
	// 消息
	Message string `json:"message"`
	// ID
	Id string `json:"id"`
}

// SendLoginEmailCodeRequest 发送登录邮箱验证码请求
// swagger:model
type SendLoginEmailCodeRequest struct {
	// 邮箱地址
	Email string `json:"email"`
}

// SendLoginEmailCodeResponse 发送登录邮箱验证码响应
// swagger:model
type SendLoginEmailCodeResponse struct {
	// 消息
	Message string `json:"message"`
	// ID
	Id string `json:"id"`
}

// RegisterRequest 用户注册请求
// swagger:model
type RegisterRequest struct {
	// 用户名
	Username string `json:"username"`
	// 邮箱地址
	Email string `json:"email"`
	// 密码
	Password string `json:"password"`
}

// RegisterResponse 用户注册响应
// swagger:model
type RegisterResponse struct {
	// JWT token
	Token string `json:"token"`
	// 用户信息
	User User `json:"user"`
	// 消息
	Message string `json:"message"`
}

// LoginRequest 用户登录请求
// swagger:model
type LoginRequest struct {
	// 邮箱地址
	Email string `json:"email"`
	// 密码
	Password string `json:"password"`
}

// LoginResponse 用户登录响应
// swagger:model
type LoginResponse struct {
	// JWT token
	Token string `json:"token"`
	// 用户信息
	User User `json:"user"`
	// 消息
	Message string `json:"message"`
}

// GetUserResponse 获取用户信息响应
// swagger:model
type GetUserResponse struct {
	// 用户信息
	User User `json:"user"`
	// 消息
	Message string `json:"message"`
}

// User 用户信息
// swagger:model
type User struct {
	// 用户ID
	Id string `json:"id"`
	// 用户名
	Username string `json:"username"`
	// 邮箱地址
	Email string `json:"email"`
}

// ErrorResponse 通用错误响应
// swagger:model
type ErrorResponse struct {
	// 错误信息
	Error string `json:"error"`
}

// GetTasksResponse 获取任务列表响应
// swagger:model
type GetTasksResponse struct {
	// 任务列表
	Tasks []Task `json:"tasks"`
	// 消息
	Message string `json:"message"`
}

// GetTaskResponse 获取任务响应
// swagger:model
type GetTaskResponse struct {
	// 任务信息
	Task Task `json:"task"`
	// 消息
	Message string `json:"message"`
}

// CreateTaskRequest 创建任务请求
// swagger:model
type CreateTaskRequest struct {
	// 任务标题
	Title string `json:"title"`
	// 任务描述
	Description string `json:"description"`
	// 任务状态
	Status int32 `json:"status"`
	// 任务优先级
	Priority int32 `json:"priority"`
	// 分配给
	Assignee string `json:"assignee"`
	// 截止日期
	Deadline string `json:"deadline"`
	// 计划日期
	ScheduledDate string `json:"scheduledDate"`
	// Token
	Token string `json:"token"`
}

// CreateTaskResponse 创建任务响应
// swagger:model
type CreateTaskResponse struct {
	// 任务信息
	Task Task `json:"task"`
	// 消息
	Message string `json:"message"`
}

// UpdateTaskRequest 更新任务请求
// swagger:model
type UpdateTaskRequest struct {
	// 任务ID
	Id string `json:"id"`
	// 任务标题
	Title string `json:"title"`
	// 任务描述
	Description string `json:"description"`
	// 任务状态
	Status int32 `json:"status"`
	// 任务优先级
	Priority int32 `json:"priority"`
	// 分配给
	Assignee string `json:"assignee"`
	// Token
	Token string `json:"token"`
}

// UpdateTaskResponse 更新任务响应
// swagger:model
type UpdateTaskResponse struct {
	// 任务信息
	Task Task `json:"task"`
	// 消息
	Message string `json:"message"`
}

// DeleteTaskResponse 删除任务响应
// swagger:model
type DeleteTaskResponse struct {
	// 消息
	Message string `json:"message"`
}

// AddCommentRequest 添加评论请求
// swagger:model
type AddCommentRequest struct {
	// 任务ID
	TaskId string `json:"task_id"`
	// 评论内容
	Text string `json:"text"`
	// Token
	Token string `json:"token"`
}

// AddCommentResponse 添加评论响应
// swagger:model
type AddCommentResponse struct {
	// 评论信息
	Comment Comment `json:"comment"`
	// 消息
	Message string `json:"message"`
}

// Task 任务信息
// swagger:model
type Task struct {
	// 任务ID
	Id string `json:"id"`
	// 任务标题
	Title string `json:"title"`
	// 任务描述
	Description string `json:"description"`
	// 任务状态
	Status int32 `json:"status"`
	// 任务优先级
	Priority int32 `json:"priority"`
	// 分配给
	Assignee string `json:"assignee"`
	// 创建者
	CreatedBy string `json:"created_by"`
}

// Comment 评论信息
// swagger:model
type Comment struct {
	// 评论ID
	Id string `json:"id"`
	// 评论内容
	Text string `json:"text"`
	// 创建者
	CreatedBy string `json:"created_by"`
}

// GetReportsResponse 获取报告列表响应
// swagger:model
type GetReportsResponse struct {
	// 报告列表
	Reports []Report `json:"reports"`
	// 消息
	Message string `json:"message"`
}

// GetReportResponse 获取报告响应
// swagger:model
type GetReportResponse struct {
	// 报告信息
	Report Report `json:"report"`
	// 消息
	Message string `json:"message"`
}

// GenerateReportRequest 生成报告请求
// swagger:model
type GenerateReportRequest struct {
	// 报告类型
	Type int32 `json:"type"`
	// 周期
	Period string `json:"period"`
	// 开始日期
	StartDate string `json:"start_date"`
	// 结束日期
	EndDate string `json:"end_date"`
	// Token
	Token string `json:"token"`
}

// GenerateReportResponse 生成报告响应
// swagger:model
type GenerateReportResponse struct {
	// 报告信息
	Report Report `json:"report"`
	// 消息
	Message string `json:"message"`
}

// PolishReportRequest 润色报告请求
// swagger:model
type PolishReportRequest struct {
	// 报告ID
	ReportId string `json:"report_id"`
	// API密钥
	ApiKey string `json:"api_key"`
	// 模型
	Model string `json:"model"`
	// API URL
	ApiUrl string `json:"api_url"`
	// 提供商
	Provider string `json:"provider"`
	// Token
	Token string `json:"token"`
}

// PolishReportResponse 润色报告响应
// swagger:model
type PolishReportResponse struct {
	// 报告信息
	Report Report `json:"report"`
	// 消息
	Message string `json:"message"`
}

// ExportReportRequest 导出报告请求
// swagger:model
type ExportReportRequest struct {
	// 报告ID
	Id string `json:"id"`
	// 格式
	Format string `json:"format"`
	// Token
	Token string `json:"token"`
}

// ExportReportResponse 导出报告响应
// swagger:model
type ExportReportResponse struct {
	// 内容
	Content string `json:"content"`
	// 文件名
	Filename string `json:"filename"`
	// 消息
	Message string `json:"message"`
}

// DeleteReportResponse 删除报告响应
// swagger:model
type DeleteReportResponse struct {
	// 消息
	Message string `json:"message"`
}

// Report 报告信息
// swagger:model
type Report struct {
	// 报告ID
	Id string `json:"id"`
	// 用户ID
	UserId string `json:"user_id"`
	// 报告类型
	Type int32 `json:"type"`
	// 周期
	Period string `json:"period"`
	// 标题
	Title string `json:"title"`
	// 内容
	Content string `json:"content"`
	// 润色后的内容
	PolishedContent string `json:"polished_content"`
	// 任务列表
	Tasks []string `json:"tasks"`
	// 统计信息
	Statistics ReportStatistics `json:"statistics"`
}

// ReportStatistics 报告统计信息
// swagger:model
type ReportStatistics struct {
	// 总任务数
	TotalTasks int32 `json:"total_tasks"`
	// 完成任务数
	CompletedTasks int32 `json:"completed_tasks"`
	// 进行中任务数
	InProgressTasks int32 `json:"in_progress_tasks"`
	// 逾期任务数
	OverdueTasks int32 `json:"overdue_tasks"`
	// 完成率
	CompletionRate int32 `json:"completion_rate"`
}


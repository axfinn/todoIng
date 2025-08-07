package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "todoing-backend/docs" // Swagger docs
	pb "todoing-backend/proto/gen/proto"
)

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

// RegisterReportRoutes 注册报告相关路由
func RegisterReportRoutes(router gin.IRouter) {
	reports := router.Group("/api/reports")
	{
		reports.GET("", getReports)
		reports.GET("/:id", getReport)
		// reports.POST("", createReport) // 暂时注释掉未实现的路由
		reports.DELETE("/:id", deleteReport)
	}
}

// getReports 获取报告列表
// @Summary 获取报告列表
// @Description 获取用户的所有报告
// @Tags reports
// @Accept json
// @Produce json
// @Param userId query string true "用户ID"
// @Success 200 {object} GetReportsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reports [get]
func getReports(c *gin.Context) {
	userId := c.Query("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	// 实际实现中，这里会从数据库获取报告列表
	reports := []*pb.Report{
		{
			Id:      "report-1",
			UserId:  userId,
			Type:    pb.ReportType_DAILY,
			Period:  "2023-01-01",
			Title:   "Sample Report 1",
			Content: "Sample report content",
			Statistics: &pb.ReportStatistics{
				TotalTasks:       10,
				CompletedTasks:   5,
				InProgressTasks:  3,
				OverdueTasks:     2,
				CompletionRate:   50,
			},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"reports": reports,
		"message": "Reports retrieved successfully",
	})
}

// getReport 获取单个报告
// @Summary 获取单个报告
// @Description 根据ID获取报告详情
// @Tags reports
// @Accept json
// @Produce json
// @Param id path string true "报告ID"
// @Success 200 {object} GetReportResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reports/{id} [get]
func getReport(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	// 实际实现中，这里会从数据库获取报告
	report := &pb.Report{
		Id:      id,
		UserId:  "user-id",
		Type:    pb.ReportType_DAILY,
		Period:  "2023-01-01",
		Title:   "Sample Report",
		Content: "Sample report content",
		Statistics: &pb.ReportStatistics{
			TotalTasks:       0,
			CompletedTasks:   0,
			InProgressTasks:  0,
			OverdueTasks:     0,
			CompletionRate:   0,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"report":  report,
		"message": "Report retrieved successfully",
	})
}

// generateReport 生成报告
// @Summary 生成报告
// @Description 根据类型和周期生成报告
// @Tags reports
// @Accept json
// @Produce json
// @Param request body GenerateReportRequest true "生成报告请求"
// @Success 200 {object} GenerateReportResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reports/generate [post]
func generateReport(c *gin.Context) {
	var req pb.GenerateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 实际实现中，这里会处理生成报告逻辑
	report := &pb.Report{
		Id:      "generated-report-id",
		UserId:  "user-id", // 这里应该是从上下文获取的用户ID
		Type:    req.Type,
		Period:  req.Period,
		Title:   "Generated Report",
		Content: "Generated report content",
		Statistics: &pb.ReportStatistics{
			TotalTasks:       0,
			CompletedTasks:   0,
			InProgressTasks:  0,
			OverdueTasks:     0,
			CompletionRate:   0,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"report":  report,
		"message": "Report generated successfully",
	})
}

// polishReport 润色报告
// @Summary 润色报告
// @Description 使用AI润色报告内容
// @Tags reports
// @Accept json
// @Produce json
// @Param request body PolishReportRequest true "润色报告请求"
// @Success 200 {object} PolishReportResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reports/polish [post]
func polishReport(c *gin.Context) {
	var req pb.PolishReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 实际实现中，这里会处理润色报告逻辑
	report := &pb.Report{
		Id:              req.ReportId,
		UserId:          "user-id",
		Type:            pb.ReportType_DAILY,
		Period:          "2023-01-01",
		Title:           "Polished Report",
		Content:         "Original report content",
		PolishedContent: "Polished report content",
		Statistics: &pb.ReportStatistics{
			TotalTasks:       0,
			CompletedTasks:   0,
			InProgressTasks:  0,
			OverdueTasks:     0,
			CompletionRate:   0,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"report":  report,
		"message": "Report polished successfully",
	})
}

// exportReport 导出报告
// @Summary 导出报告
// @Description 导出报告为指定格式
// @Tags reports
// @Accept json
// @Produce json
// @Param request body ExportReportRequest true "导出报告请求"
// @Success 200 {object} ExportReportResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reports/export [post]
func exportReport(c *gin.Context) {
	var req pb.ExportReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 实际实现中，这里会处理导出报告逻辑
	c.JSON(http.StatusOK, gin.H{
		"content":  "Report content",
		"filename": "report.txt",
		"message":  "Report exported successfully",
	})
}

// deleteReport 删除报告
// @Summary 删除报告
// @Description 根据ID删除报告
// @Tags reports
// @Accept json
// @Produce json
// @Param id path string true "报告ID"
// @Success 200 {object} DeleteReportResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reports/{id} [delete]
func deleteReport(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	// 实际实现中，这里会从数据库删除报告
	c.JSON(http.StatusOK, gin.H{
		"message": "Report deleted successfully",
	})
}
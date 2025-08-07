package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"todoing-backend/src/config"
	"todoing-backend/src/middleware"
	"todoing-backend/src/models"
)

// RegisterReportRoutes 注册报告相关路由
func RegisterReportRoutes(router gin.IRouter) {
	reports := router.Group("/api/reports")
	reports.Use(middleware.AuthMiddleware()) // 添加认证中间件
	{
		reports.GET("", getReports)
		reports.GET("/:id", getReport)
		reports.POST("/generate", generateReport)
		reports.POST("/:id/polish", polishReport)
		reports.GET("/:id/export/:format", exportReport)
		reports.DELETE("/:id", deleteReport)
	}
}

// getReports 获取报告列表
func getReports(c *gin.Context) {
	// 从JWT token中获取用户ID
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "User not authenticated"})
		return
	}

	// 提取用户信息
	claims := userClaims.(map[string]interface{})
	userMap, ok := claims["user"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid token structure"})
		return
	}

	userIDStr, ok := userMap["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "User ID not found in token"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid user ID"})
		return
	}

	// 从数据库获取报告列表
	var reports []models.Report
	filter := bson.M{"userId": userID}
	cursor, err := config.DB.Collection("reports").Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to fetch reports"})
		return
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &reports); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to decode reports"})
		return
	}

	// 转换为响应格式
	var responseReports []map[string]interface{}
	for _, report := range reports {
		responseReport := map[string]interface{}{
			"_id":             report.ID.Hex(),
			"userId":          report.UserID.Hex(),
			"type":            string(report.Type),
			"period":          report.Period,
			"title":           report.Title,
			"content":         report.Content,
			"polishedContent": report.PolishedContent,
			"tasks":           convertTaskIDs(report.Tasks),
			"statistics":      report.Statistics,
			"createdAt":       report.CreatedAt,
			"updatedAt":       report.UpdatedAt,
		}
		responseReports = append(responseReports, responseReport)
	}

	// 前端期望直接返回报告数组
	c.JSON(http.StatusOK, responseReports)
}

// getReport 获取单个报告
func getReport(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "ID is required"})
		return
	}

	reportID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid report ID"})
		return
	}

	// 从JWT token中获取用户ID
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "User not authenticated"})
		return
	}

	claims := userClaims.(map[string]interface{})
	userMap, ok := claims["user"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid token structure"})
		return
	}

	userIDStr, ok := userMap["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "User ID not found in token"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid user ID"})
		return
	}

	// 从数据库获取报告
	var report models.Report
	filter := bson.M{"_id": reportID, "userId": userID}
	err = config.DB.Collection("reports").FindOne(context.TODO(), filter).Decode(&report)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"msg": "Report not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to fetch report"})
		return
	}

	// 转换为响应格式
	responseReport := map[string]interface{}{
		"_id":             report.ID.Hex(),
		"userId":          report.UserID.Hex(),
		"type":            string(report.Type),
		"period":          report.Period,
		"title":           report.Title,
		"content":         report.Content,
		"polishedContent": report.PolishedContent,
		"tasks":           convertTaskIDs(report.Tasks),
		"statistics":      report.Statistics,
		"createdAt":       report.CreatedAt,
		"updatedAt":       report.UpdatedAt,
	}

	// 前端期望直接返回报告对象
	c.JSON(http.StatusOK, responseReport)
}

// generateReport 生成报告
func generateReport(c *gin.Context) {
	var req struct {
		Type      string `json:"type" binding:"required"`
		Period    string `json:"period"`
		StartDate string `json:"startDate" binding:"required"`
		EndDate   string `json:"endDate" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	// 验证报告类型
	var reportType models.ReportType
	switch req.Type {
	case "daily":
		reportType = models.DailyReport
	case "weekly":
		reportType = models.WeeklyReport
	case "monthly":
		reportType = models.MonthlyReport
	default:
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid report type"})
		return
	}

	// 从JWT token中获取用户ID
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "User not authenticated"})
		return
	}

	claims := userClaims.(map[string]interface{})
	userMap, ok := claims["user"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid token structure"})
		return
	}

	userIDStr, ok := userMap["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "User ID not found in token"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid user ID"})
		return
	}

	// 解析日期
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid start date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid end date format"})
		return
	}

	// 获取时间范围内的任务
	tasks, statistics, err := getTasksAndStatistics(userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to fetch tasks"})
		return
	}

	// 生成报告内容
	content := generateReportContent(reportType, startDate, endDate, tasks, statistics)

	// 生成报告标题
	title := generateReportTitle(reportType, startDate, endDate)

	// 设置period
	period := req.Period
	if period == "" {
		period = fmt.Sprintf("%s - %s", req.StartDate, req.EndDate)
	}

	// 创建报告
	report := models.Report{
		ID:         primitive.NewObjectID(),
		UserID:     userID,
		Type:       reportType,
		Period:     period,
		Title:      title,
		Content:    content,
		Tasks:      getTaskIDs(tasks),
		Statistics: statistics,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 保存到数据库
	_, err = config.DB.Collection("reports").InsertOne(context.TODO(), report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to save report"})
		return
	}

	// 转换为响应格式
	responseReport := map[string]interface{}{
		"_id":             report.ID.Hex(),
		"userId":          report.UserID.Hex(),
		"type":            string(report.Type),
		"period":          report.Period,
		"title":           report.Title,
		"content":         report.Content,
		"polishedContent": report.PolishedContent,
		"tasks":           convertTaskIDs(report.Tasks),
		"statistics":      report.Statistics,
		"createdAt":       report.CreatedAt,
		"updatedAt":       report.UpdatedAt,
	}

	// 前端期望直接返回报告对象
	c.JSON(http.StatusOK, responseReport)
}

// polishReport 润色报告
func polishReport(c *gin.Context) {
	reportID := c.Param("id")
	if reportID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Report ID is required"})
		return
	}

	var req struct {
		ApiKey   string `json:"api_key"`
		Model    string `json:"model"`
		ApiUrl   string `json:"api_url"`
		Provider string `json:"provider"`
		Token    string `json:"token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	// 从JWT token中获取用户ID
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "User not authenticated"})
		return
	}

	claims := userClaims.(map[string]interface{})
	userMap, ok := claims["user"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid token structure"})
		return
	}

	userIDStr, ok := userMap["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "User ID not found in token"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid user ID"})
		return
	}

	reportObjID, err := primitive.ObjectIDFromHex(reportID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid report ID"})
		return
	}

	// 获取报告
	var report models.Report
	filter := bson.M{"_id": reportObjID, "userId": userID}
	err = config.DB.Collection("reports").FindOne(context.TODO(), filter).Decode(&report)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"msg": "Report not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to fetch report"})
		return
	}

	// TODO: 实现AI润色逻辑
	// 这里暂时模拟润色后的内容
	polishedContent := fmt.Sprintf("【AI润色版】\n\n%s\n\n---\n此内容已经过AI优化，使其更加专业和易读。", report.Content)

	// 更新报告
	update := bson.M{
		"$set": bson.M{
			"polishedContent": polishedContent,
			"updatedAt":       time.Now(),
		},
	}

	_, err = config.DB.Collection("reports").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to update report"})
		return
	}

	// 返回更新后的报告
	report.PolishedContent = &polishedContent
	report.UpdatedAt = time.Now()

	// 转换为响应格式
	responseReport := map[string]interface{}{
		"_id":             report.ID.Hex(),
		"userId":          report.UserID.Hex(),
		"type":            string(report.Type),
		"period":          report.Period,
		"title":           report.Title,
		"content":         report.Content,
		"polishedContent": report.PolishedContent,
		"tasks":           convertTaskIDs(report.Tasks),
		"statistics":      report.Statistics,
		"createdAt":       report.CreatedAt,
		"updatedAt":       report.UpdatedAt,
	}

	// 前端期望直接返回报告对象
	c.JSON(http.StatusOK, responseReport)
}

// exportReport 导出报告
func exportReport(c *gin.Context) {
	reportID := c.Param("id")
	format := c.Param("format")

	if reportID == "" || format == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid parameters"})
		return
	}

	// 从JWT token中获取用户ID
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "User not authenticated"})
		return
	}

	claims := userClaims.(map[string]interface{})
	userMap, ok := claims["user"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid token structure"})
		return
	}

	userIDStr, ok := userMap["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "User ID not found in token"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid user ID"})
		return
	}

	reportObjID, err := primitive.ObjectIDFromHex(reportID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid report ID"})
		return
	}

	// 获取报告
	var report models.Report
	filter := bson.M{"_id": reportObjID, "userId": userID}
	err = config.DB.Collection("reports").FindOne(context.TODO(), filter).Decode(&report)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"msg": "Report not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to fetch report"})
		return
	}

	// 根据格式生成内容
	var content string
	var contentType string
	var fileExtension string

	// 优先使用润色后的内容
	reportContent := report.Content
	if report.PolishedContent != nil && *report.PolishedContent != "" {
		reportContent = *report.PolishedContent
	}

	switch format {
	case "md", "markdown":
		content = formatReportAsMarkdown(report, reportContent)
		contentType = "text/markdown"
		fileExtension = "md"
	case "txt":
		content = formatReportAsText(report, reportContent)
		contentType = "text/plain"
		fileExtension = "txt"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Unsupported format"})
		return
	}

	filename := fmt.Sprintf("report-%s-%s.%s", report.Type, report.Period, fileExtension)
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	// 返回文件内容
	c.String(http.StatusOK, content)
}

// deleteReport 删除报告
func deleteReport(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "ID is required"})
		return
	}

	reportID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid report ID"})
		return
	}

	// 从JWT token中获取用户ID
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "User not authenticated"})
		return
	}

	claims := userClaims.(map[string]interface{})
	userMap, ok := claims["user"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid token structure"})
		return
	}

	userIDStr, ok := userMap["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "User ID not found in token"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid user ID"})
		return
	}

	// 删除报告
	filter := bson.M{"_id": reportID, "userId": userID}
	result, err := config.DB.Collection("reports").DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to delete report"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Report not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Report deleted successfully",
	})
}

// 辅助函数

// getTasksAndStatistics 获取任务和统计信息
func getTasksAndStatistics(userID primitive.ObjectID, startDate, endDate time.Time) ([]models.Task, models.ReportStatistics, error) {
	// 设置结束日期为当天的23:59:59
	endDate = endDate.Add(24*time.Hour - 1*time.Second)

	// 查询条件
	filter := bson.M{
		"createdBy": userID,
		"createdAt": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	// 获取任务
	cursor, err := config.DB.Collection("tasks").Find(context.TODO(), filter)
	if err != nil {
		return nil, models.ReportStatistics{}, err
	}
	defer cursor.Close(context.TODO())

	var tasks []models.Task
	if err = cursor.All(context.TODO(), &tasks); err != nil {
		return nil, models.ReportStatistics{}, err
	}

	// 计算统计信息
	statistics := models.ReportStatistics{
		TotalTasks: len(tasks),
	}

	now := time.Now()
	for _, task := range tasks {
		switch task.Status {
		case models.StatusDone:
			statistics.CompletedTasks++
		case models.StatusInProgress:
			statistics.InProgressTasks++
		}

		// 检查是否逾期
		if task.Deadline != nil && task.Deadline.Before(now) && task.Status != models.StatusDone {
			statistics.OverdueTasks++
		}
	}

	// 计算完成率
	if statistics.TotalTasks > 0 {
		statistics.CompletionRate = (statistics.CompletedTasks * 100) / statistics.TotalTasks
	}

	return tasks, statistics, nil
}

// generateReportContent 生成报告内容
func generateReportContent(reportType models.ReportType, startDate, endDate time.Time, tasks []models.Task, statistics models.ReportStatistics) string {
	content := fmt.Sprintf("# %s 报告\n\n", reportType)
	content += fmt.Sprintf("报告周期：%s 至 %s\n\n", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	content += "## 任务概览\n\n"
	content += fmt.Sprintf("- 总任务数：%d\n", statistics.TotalTasks)
	content += fmt.Sprintf("- 已完成：%d\n", statistics.CompletedTasks)
	content += fmt.Sprintf("- 进行中：%d\n", statistics.InProgressTasks)
	content += fmt.Sprintf("- 已逾期：%d\n", statistics.OverdueTasks)
	content += fmt.Sprintf("- 完成率：%d%%\n\n", statistics.CompletionRate)

	// 添加任务详情
	if len(tasks) > 0 {
		content += "## 任务详情\n\n"

		// 按状态分组
		completedTasks := []models.Task{}
		inProgressTasks := []models.Task{}
		todoTasks := []models.Task{}

		for _, task := range tasks {
			switch task.Status {
			case models.StatusDone:
				completedTasks = append(completedTasks, task)
			case models.StatusInProgress:
				inProgressTasks = append(inProgressTasks, task)
			case models.StatusToDo:
				todoTasks = append(todoTasks, task)
			}
		}

		// 已完成任务
		if len(completedTasks) > 0 {
			content += "### 已完成任务\n\n"
			for i, task := range completedTasks {
				content += fmt.Sprintf("%d. %s\n", i+1, task.Title)
				if task.Description != "" {
					content += fmt.Sprintf("   - %s\n", task.Description)
				}
			}
			content += "\n"
		}

		// 进行中任务
		if len(inProgressTasks) > 0 {
			content += "### 进行中任务\n\n"
			for i, task := range inProgressTasks {
				content += fmt.Sprintf("%d. %s", i+1, task.Title)
				if task.Deadline != nil {
					content += fmt.Sprintf(" (截止日期：%s)", task.Deadline.Format("2006-01-02"))
				}
				content += "\n"
				if task.Description != "" {
					content += fmt.Sprintf("   - %s\n", task.Description)
				}
			}
			content += "\n"
		}

		// 待办任务
		if len(todoTasks) > 0 {
			content += "### 待办任务\n\n"
			for i, task := range todoTasks {
				content += fmt.Sprintf("%d. %s", i+1, task.Title)
				if task.Deadline != nil {
					content += fmt.Sprintf(" (截止日期：%s)", task.Deadline.Format("2006-01-02"))
				}
				content += "\n"
				if task.Description != "" {
					content += fmt.Sprintf("   - %s\n", task.Description)
				}
			}
			content += "\n"
		}
	}

	// 添加总结
	content += "## 总结\n\n"
	if statistics.CompletionRate >= 80 {
		content += "本期任务完成情况良好，继续保持！"
	} else if statistics.CompletionRate >= 60 {
		content += "本期任务完成率尚可，仍有提升空间。"
	} else {
		content += "本期任务完成率较低，建议调整任务安排或工作计划。"
	}

	return content
}

// generateReportTitle 生成报告标题
func generateReportTitle(reportType models.ReportType, startDate, endDate time.Time) string {
	switch reportType {
	case models.DailyReport:
		return fmt.Sprintf("%s 日报", startDate.Format("2006年01月02日"))
	case models.WeeklyReport:
		return fmt.Sprintf("%s 至 %s 周报", startDate.Format("01月02日"), endDate.Format("01月02日"))
	case models.MonthlyReport:
		return fmt.Sprintf("%s 月报", startDate.Format("2006年01月"))
	default:
		return fmt.Sprintf("%s 至 %s 报告", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	}
}

// getTaskIDs 获取任务ID列表
func getTaskIDs(tasks []models.Task) []primitive.ObjectID {
	ids := make([]primitive.ObjectID, len(tasks))
	for i, task := range tasks {
		ids[i] = task.ID
	}
	return ids
}

// convertTaskIDs 转换任务ID为字符串数组
func convertTaskIDs(ids []primitive.ObjectID) []string {
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = id.Hex()
	}
	return result
}

// formatReportAsMarkdown 格式化报告为Markdown
func formatReportAsMarkdown(report models.Report, content string) string {
	markdown := fmt.Sprintf("# %s\n\n", report.Title)
	markdown += fmt.Sprintf("**类型**: %s  \n", report.Type)
	markdown += fmt.Sprintf("**周期**: %s  \n", report.Period)
	markdown += fmt.Sprintf("**创建时间**: %s  \n\n", report.CreatedAt.Format("2006-01-02 15:04:05"))
	markdown += "---\n\n"
	markdown += content
	markdown += "\n\n---\n\n"
	markdown += fmt.Sprintf("*报告生成时间: %s*", time.Now().Format("2006-01-02 15:04:05"))
	return markdown
}

// formatReportAsText 格式化报告为纯文本
func formatReportAsText(report models.Report, content string) string {
	text := fmt.Sprintf("%s\n", report.Title)
	text += fmt.Sprintf("类型: %s\n", report.Type)
	text += fmt.Sprintf("周期: %s\n", report.Period)
	text += fmt.Sprintf("创建时间: %s\n\n", report.CreatedAt.Format("2006-01-02 15:04:05"))
	text += "=====================================\n\n"
	text += content
	text += "\n\n=====================================\n\n"
	text += fmt.Sprintf("报告生成时间: %s", time.Now().Format("2006-01-02 15:04:05"))
	return text
}

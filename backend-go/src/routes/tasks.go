package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	pb "todoing-backend/proto/gen/proto"
	"todoing-backend/src/config"
	"todoing-backend/src/models"
)

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
	Status string `json:"status"`
	// 任务优先级
	Priority string `json:"priority"`
	// 分配给
	Assignee *string `json:"assignee"`
	// 截止日期
	Deadline *string `json:"deadline"`
	// 计划日期
	ScheduledDate *string `json:"scheduledDate"`
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
	Status string `json:"status"`
	// 任务优先级
	Priority string `json:"priority"`
	// 分配给
	Assignee *string `json:"assignee"`
	// 截止日期
	Deadline *string `json:"deadline"`
	// 计划日期
	ScheduledDate *string `json:"scheduledDate"`
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
	Status string `json:"status"`
	// 任务优先级
	Priority string `json:"priority"`
	// 分配给
	Assignee string `json:"assignee"`
	// 创建者
	CreatedBy string `json:"created_by"`
	// 创建时间
	CreatedAt string `json:"created_at"`
	// 更新时间
	UpdatedAt string `json:"updated_at"`
	// 截止日期
	Deadline *string `json:"deadline,omitempty"`
	// 计划日期
	ScheduledDate *string `json:"scheduled_date,omitempty"`
	// 评论
	Comments []Comment `json:"comments,omitempty"`
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
	// 创建时间
	CreatedAt string `json:"created_at"`
}

// RegisterTaskRoutes 注册任务相关路由
func RegisterTaskRoutes(router gin.IRouter) {
	tasks := router.Group("/api/tasks")
	{
		tasks.GET("", getTasks)
		tasks.GET("/export/all", exportTasks)
		tasks.POST("/import", importTasks)
		tasks.GET("/:id", getTask)
		tasks.POST("", createTask)
		tasks.PUT("/:id", updateTask)
		tasks.DELETE("/:id", deleteTask)
		tasks.POST("/:id/assign", assignTask)
		tasks.POST("/:id/comments", addComment)
	}
}

// getTasks 获取任务列表
// @Summary 获取任务列表
// @Description 获取用户的所有任务
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {object} GetTasksResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks [get]
func getTasks(c *gin.Context) {
	// 从JWT token中提取用户信息
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token, authorization denied"})
		return
	}

	claims, ok := userClaims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userId, ok := claims["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}
	
	// 将用户ID字符串转换为ObjectID
	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// 从数据库获取任务列表
	tasks, err := models.GetTasksByUserID(userObjId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	// 转换为前端可以直接使用的格式
	var respTasks []Task
	for _, task := range tasks {
		var assignee string
		if task.Assignee != nil {
			assignee = *task.Assignee
		}
		
		respTask := Task{
			Id:          task.ID.Hex(),
			Title:       task.Title,
			Description: task.Description,
			Status:      convertTaskStatusToString(task.Status),
			Priority:    convertTaskPriorityToString(task.Priority),
			Assignee:    assignee,
			CreatedBy:   task.CreatedBy.Hex(),
			CreatedAt:   task.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
		}
		
		// 处理可选的Deadline和ScheduledDate字段
		if task.Deadline != nil {
			formattedDeadline := task.Deadline.Format(time.RFC3339)
			respTask.Deadline = &formattedDeadline
		}
		if task.ScheduledDate != nil {
			formattedScheduledDate := task.ScheduledDate.Format(time.RFC3339)
			respTask.ScheduledDate = &formattedScheduledDate
		}
		
		// 转换评论
		for _, comment := range task.Comments {
			respComment := Comment{
				Id:        comment.ID.Hex(),
				Text:      comment.Text,
				CreatedBy: comment.CreatedBy.Hex(),
				CreatedAt: comment.CreatedAt.Format(time.RFC3339),
			}
			respTask.Comments = append(respTask.Comments, respComment)
		}
		
		respTasks = append(respTasks, respTask)
	}

	// 返回JSON格式的任务列表
	c.JSON(http.StatusOK, gin.H{
		"tasks":   respTasks,
		"message": "Tasks retrieved successfully",
	})
}

// convertTaskStatusToProto 将TaskStatus转换为pb.TaskStatus
func convertTaskStatusToProto(status models.TaskStatus) pb.TaskStatus {
	switch status {
	case models.StatusToDo:
		return pb.TaskStatus_TO_DO
	case models.StatusInProgress:
		return pb.TaskStatus_IN_PROGRESS
	case models.StatusDone:
		return pb.TaskStatus_DONE
	default:
		return pb.TaskStatus_TO_DO
	}
}

// convertTaskPriorityToProto 将TaskPriority转换为pb.TaskPriority
func convertTaskPriorityToProto(priority models.TaskPriority) pb.TaskPriority {
	switch priority {
	case models.PriorityLow:
		return pb.TaskPriority_LOW
	case models.PriorityMedium:
		return pb.TaskPriority_MEDIUM
	case models.PriorityHigh:
		return pb.TaskPriority_HIGH
	default:
		return pb.TaskPriority_MEDIUM
	}
}

// convertTaskStatusToString 将TaskStatus转换为字符串
func convertTaskStatusToString(status models.TaskStatus) string {
	switch status {
	case models.StatusToDo:
		return "TO_DO"
	case models.StatusInProgress:
		return "IN_PROGRESS"
	case models.StatusDone:
		return "DONE"
	default:
		return "TO_DO"
	}
}

// convertTaskPriorityToString 将TaskPriority转换为字符串
func convertTaskPriorityToString(priority models.TaskPriority) string {
	switch priority {
	case models.PriorityLow:
		return "LOW"
	case models.PriorityMedium:
		return "MEDIUM"
	case models.PriorityHigh:
		return "HIGH"
	default:
		return "MEDIUM"
	}
}

// getTask 获取单个任务
// @Summary 获取单个任务
// @Description 根据ID获取任务详情
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Success 200 {object} GetTaskResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks/{id} [get]
func getTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	// 将ID字符串转换为ObjectID
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// 从数据库获取任务
	var task models.Task
	err = config.DB.Collection("tasks").FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		return
	}

	// 转换为前端可以直接使用的格式
	var assignee string
	if task.Assignee != nil {
		assignee = *task.Assignee
	}
	
	respTask := Task{
		Id:          task.ID.Hex(),
		Title:       task.Title,
		Description: task.Description,
		Status:      convertTaskStatusToString(task.Status),
		Priority:    convertTaskPriorityToString(task.Priority),
		Assignee:    assignee,
		CreatedBy:   task.CreatedBy.Hex(),
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
	}
	
	// 处理可选的Deadline和ScheduledDate字段
	if task.Deadline != nil {
		formattedDeadline := task.Deadline.Format(time.RFC3339)
		respTask.Deadline = &formattedDeadline
	}
	if task.ScheduledDate != nil {
		formattedScheduledDate := task.ScheduledDate.Format(time.RFC3339)
		respTask.ScheduledDate = &formattedScheduledDate
	}
	
	// 转换评论
	for _, comment := range task.Comments {
		respComment := Comment{
			Id:        comment.ID.Hex(),
			Text:      comment.Text,
			CreatedBy: comment.CreatedBy.Hex(),
			CreatedAt: comment.CreatedAt.Format(time.RFC3339),
		}
		respTask.Comments = append(respTask.Comments, respComment)
	}

	c.JSON(http.StatusOK, gin.H{
		"task":    respTask,
		"message": "Task retrieved successfully",
	})
}

// createTask 创建任务
// @Summary 创建任务
// @Description 创建一个新的任务
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body CreateTaskRequest true "Create task"
// @Success 201 {object} CreateTaskResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks [post]
func createTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT token中提取用户信息
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token, authorization denied"})
		return
	}

	claims, ok := userClaims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userId, ok := claims["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	// 将用户ID字符串转换为ObjectID
	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// 创建任务模型
	now := time.Now()
	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      convertStringToTaskStatus(req.Status),
		Priority:    convertStringToTaskPriority(req.Priority),
		CreatedBy:   userObjId,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// 处理可选字段
	if req.Assignee != nil {
		task.Assignee = req.Assignee
	}
	
	// 处理可选的时间字段
	if req.Deadline != nil && *req.Deadline != "" {
		deadlineTime, err := time.Parse(time.RFC3339, *req.Deadline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deadline format, must be RFC3339"})
			return
		}
		task.Deadline = &deadlineTime
	}
	
	if req.ScheduledDate != nil && *req.ScheduledDate != "" {
		scheduledDateTime, err := time.Parse(time.RFC3339, *req.ScheduledDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scheduled date format, must be RFC3339"})
			return
		}
		task.ScheduledDate = &scheduledDateTime
	}

	// 保存到数据库
	if err := models.CreateTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	// 转换为响应格式并返回
	var assignee string
	if task.Assignee != nil {
		assignee = *task.Assignee
	}

	// 构建返回的Task对象
	respTask := Task{
		Id:          task.ID.Hex(),
		Title:       task.Title,
		Description: task.Description,
		Status:      convertTaskStatusToString(task.Status),
		Priority:    convertTaskPriorityToString(task.Priority),
		Assignee:    assignee,
		CreatedBy:   task.CreatedBy.Hex(),
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
	}
	
	// 处理可选的Deadline和ScheduledDate字段
	if task.Deadline != nil {
		formattedDeadline := task.Deadline.Format(time.RFC3339)
		respTask.Deadline = &formattedDeadline
	}
	if task.ScheduledDate != nil {
		formattedScheduledDate := task.ScheduledDate.Format(time.RFC3339)
		respTask.ScheduledDate = &formattedScheduledDate
	}

	c.JSON(http.StatusCreated, gin.H{
		"task":    respTask,
		"message": "Task created successfully",
	})
}

// convertProtoToTaskStatus 将pb.TaskStatus转换为TaskStatus
func convertProtoToTaskStatus(status pb.TaskStatus) models.TaskStatus {
	switch status {
	case pb.TaskStatus_TO_DO:
		return models.StatusToDo
	case pb.TaskStatus_IN_PROGRESS:
		return models.StatusInProgress
	case pb.TaskStatus_DONE:
		return models.StatusDone
	default:
		return models.StatusToDo
	}
}

// convertProtoToTaskPriority 将pb.TaskPriority转换为TaskPriority
func convertProtoToTaskPriority(priority pb.TaskPriority) models.TaskPriority {
	switch priority {
	case pb.TaskPriority_LOW:
		return models.PriorityLow
	case pb.TaskPriority_MEDIUM:
		return models.PriorityMedium
	case pb.TaskPriority_HIGH:
		return models.PriorityHigh
	default:
		return models.PriorityMedium
	}
}

// convertStringToTaskStatus 将字符串转换为TaskStatus
func convertStringToTaskStatus(status string) models.TaskStatus {
	switch status {
	case "TO_DO", "to_do", "ToDo", "todo":
		return models.StatusToDo
	case "IN_PROGRESS", "in_progress", "InProgress", "progress":
		return models.StatusInProgress
	case "DONE", "done", "Done":
		return models.StatusDone
	default:
		return models.StatusToDo
	}
}

// convertStringToTaskPriority 将字符串转换为TaskPriority
func convertStringToTaskPriority(priority string) models.TaskPriority {
	switch priority {
	case "LOW", "low", "Low":
		return models.PriorityLow
	case "MEDIUM", "medium", "Medium":
		return models.PriorityMedium
	case "HIGH", "high", "High":
		return models.PriorityHigh
	default:
		return models.PriorityMedium
	}
}

// updateTask 更新任务
// @Summary 更新任务
// @Description 更新任务信息
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Param request body UpdateTaskRequest true "更新任务请求"
// @Success 200 {object} UpdateTaskResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks/{id} [put]
func updateTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	// 将ID字符串转换为ObjectID
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 构建更新字段
	update := bson.M{
		"$set": bson.M{
			"title":       req.Title,
			"description": req.Description,
			"status":      models.TaskStatus(req.Status),
			"priority":    models.TaskPriority(req.Priority),
			"assignee":    req.Assignee,
			"updatedAt":   time.Now(),
		},
	}

	// 处理可选的时间字段
	if req.Deadline != nil && *req.Deadline != "" {
		deadlineTime, err := time.Parse(time.RFC3339, *req.Deadline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deadline format, must be RFC3339"})
			return
		}
		update["$set"].(bson.M)["deadline"] = deadlineTime
	}

	if req.ScheduledDate != nil && *req.ScheduledDate != "" {
		scheduledDateTime, err := time.Parse(time.RFC3339, *req.ScheduledDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scheduled date format, must be RFC3339"})
			return
		}
		update["$set"].(bson.M)["scheduled_date"] = scheduledDateTime
	}

	// 更新数据库中的任务
	result, err := config.DB.Collection("tasks").UpdateOne(context.TODO(), bson.M{"_id": objId}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// 获取更新后的任务
	var updatedTask models.Task
	err = config.DB.Collection("tasks").FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&updatedTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated task"})
		return
	}

	// 转换为响应格式
	var assignee string
	if updatedTask.Assignee != nil {
		assignee = *updatedTask.Assignee
	}
	
	respTask := Task{
		Id:          updatedTask.ID.Hex(),
		Title:       updatedTask.Title,
		Description: updatedTask.Description,
		Status:      convertTaskStatusToString(updatedTask.Status),
		Priority:    convertTaskPriorityToString(updatedTask.Priority),
		Assignee:    assignee,
		CreatedBy:   updatedTask.CreatedBy.Hex(),
		CreatedAt:   updatedTask.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   updatedTask.UpdatedAt.Format(time.RFC3339),
	}
	
	// 处理可选的Deadline和ScheduledDate字段
	if updatedTask.Deadline != nil {
		formattedDeadline := updatedTask.Deadline.Format(time.RFC3339)
		respTask.Deadline = &formattedDeadline
	}
	if updatedTask.ScheduledDate != nil {
		formattedScheduledDate := updatedTask.ScheduledDate.Format(time.RFC3339)
		respTask.ScheduledDate = &formattedScheduledDate
	}
	
	// 转换评论
	for _, comment := range updatedTask.Comments {
		respComment := Comment{
			Id:        comment.ID.Hex(),
			Text:      comment.Text,
			CreatedBy: comment.CreatedBy.Hex(),
			CreatedAt: comment.CreatedAt.Format(time.RFC3339),
		}
		respTask.Comments = append(respTask.Comments, respComment)
	}

	c.JSON(http.StatusOK, gin.H{
		"task":    respTask,
		"message": "Task updated successfully",
	})
}

// deleteTask 删除任务
// @Summary 删除任务
// @Description 根据ID删除任务
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Success 200 {object} DeleteTaskResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks/{id} [delete]
func deleteTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	// 将ID字符串转换为ObjectID
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// 从数据库删除任务
	result, err := config.DB.Collection("tasks").DeleteOne(context.TODO(), bson.M{"_id": objId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})
}

// assignTask 分配任务
// @Summary 分配任务
// @Description 将任务分配给用户
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Param request body UpdateTaskRequest true "分配任务请求"
// @Success 200 {object} UpdateTaskResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks/{id}/assign [post]
func assignTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	// 将ID字符串转换为ObjectID
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 准备更新数据
	update := bson.M{
		"$set": bson.M{
			"assignee":  req.Assignee,
			"updatedAt": time.Now(),
		},
	}

	// 更新数据库中的任务分配信息
	result, err := config.DB.Collection("tasks").UpdateOne(context.TODO(), bson.M{"_id": objId}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign task"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// 获取更新后的任务
	var updatedTask models.Task
	err = config.DB.Collection("tasks").FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&updatedTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated task"})
		return
	}

	// 转换为响应格式
	var assignee string
	if updatedTask.Assignee != nil {
		assignee = *updatedTask.Assignee
	}
	
	respTask := Task{
		Id:          updatedTask.ID.Hex(),
		Title:       updatedTask.Title,
		Description: updatedTask.Description,
		Status:      convertTaskStatusToString(updatedTask.Status),
		Priority:    convertTaskPriorityToString(updatedTask.Priority),
		Assignee:    assignee,
		CreatedBy:   updatedTask.CreatedBy.Hex(),
		CreatedAt:   updatedTask.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   updatedTask.UpdatedAt.Format(time.RFC3339),
	}
	
	// 处理可选的Deadline和ScheduledDate字段
	if updatedTask.Deadline != nil {
		formattedDeadline := updatedTask.Deadline.Format(time.RFC3339)
		respTask.Deadline = &formattedDeadline
	}
	if updatedTask.ScheduledDate != nil {
		formattedScheduledDate := updatedTask.ScheduledDate.Format(time.RFC3339)
		respTask.ScheduledDate = &formattedScheduledDate
	}
	
	// 转换评论
	for _, comment := range updatedTask.Comments {
		respComment := Comment{
			Id:        comment.ID.Hex(),
			Text:      comment.Text,
			CreatedBy: comment.CreatedBy.Hex(),
			CreatedAt: comment.CreatedAt.Format(time.RFC3339),
		}
		respTask.Comments = append(respTask.Comments, respComment)
	}

	c.JSON(http.StatusOK, gin.H{
		"task":    respTask,
		"message": "Task assigned successfully",
	})
}

// addComment 添加评论
// @Summary 添加评论
// @Description 为任务添加评论
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Param request body AddCommentRequest true "添加评论请求"
// @Success 200 {object} AddCommentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks/{id}/comments [post]
func addComment(c *gin.Context) {
	taskId := c.Param("id")
	if taskId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task id is required"})
		return
	}

	// 将任务ID字符串转换为ObjectID
	taskObjId, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
		return
	}

	var req AddCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT token中提取用户信息
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token, authorization denied"})
		return
	}

	claims, ok := userClaims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userId, ok := claims["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	// 将用户ID字符串转换为ObjectID
	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// 创建评论对象
	now := time.Now()
	comment := models.Comment{
		ID:        primitive.NewObjectID(),
		Text:      req.Text,
		CreatedBy: userObjId,
		CreatedAt: now,
	}

	// 将评论添加到任务中
	filter := bson.M{"_id": taskObjId}
	update := bson.M{
		"$push": bson.M{"comments": comment},
		"$set":  bson.M{"updatedAt": now},
	}

	result, err := config.DB.Collection("tasks").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// 获取更新后的任务
	var updatedTask models.Task
	err = config.DB.Collection("tasks").FindOne(context.TODO(), bson.M{"_id": taskObjId}).Decode(&updatedTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated task"})
		return
	}

	// 转换为前端可以直接使用的格式
	var assignee string
	if updatedTask.Assignee != nil {
		assignee = *updatedTask.Assignee
	}
	
	respTask := Task{
		Id:          updatedTask.ID.Hex(),
		Title:       updatedTask.Title,
		Description: updatedTask.Description,
		Status:      convertTaskStatusToString(updatedTask.Status),
		Priority:    convertTaskPriorityToString(updatedTask.Priority),
		Assignee:    assignee,
		CreatedBy:   updatedTask.CreatedBy.Hex(),
		CreatedAt:   updatedTask.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   updatedTask.UpdatedAt.Format(time.RFC3339),
	}
	
	// 处理可选的Deadline和ScheduledDate字段
	if updatedTask.Deadline != nil {
		formattedDeadline := updatedTask.Deadline.Format(time.RFC3339)
		respTask.Deadline = &formattedDeadline
	}
	if updatedTask.ScheduledDate != nil {
		formattedScheduledDate := updatedTask.ScheduledDate.Format(time.RFC3339)
		respTask.ScheduledDate = &formattedScheduledDate
	}
	
	// 转换评论
	for _, comment := range updatedTask.Comments {
		respComment := Comment{
			Id:        comment.ID.Hex(),
			Text:      comment.Text,
			CreatedBy: comment.CreatedBy.Hex(),
			CreatedAt: comment.CreatedAt.Format(time.RFC3339),
		}
		respTask.Comments = append(respTask.Comments, respComment)
	}

	c.JSON(http.StatusOK, gin.H{
		"task":    respTask,
		"message": "Comment added successfully",
	})
}

// ExportTasksResponse 导出任务响应
// swagger:model
type ExportTasksResponse struct {
	Tasks []Task `json:"tasks"`
}

// ImportTasksRequest 导入任务请求
// swagger:model
type ImportTasksRequest struct {
	Tasks []Task `json:"tasks"`
}

// ImportTasksResponse 导入任务响应
// swagger:model
type ImportTasksResponse struct {
	Message  string `json:"msg"`
	Imported int    `json:"imported"`
	Errors   []ImportError `json:"errors"`
}

// ImportError 导入错误信息
// swagger:model
type ImportError struct {
	Index int    `json:"index"`
	Error string `json:"error"`
}

// formatDateForFilename 格式化日期为文件名
func formatDateForFilename(date time.Time) string {
	return date.Format("2006-01-02")
}

// exportTasks 导出所有任务
// @Summary 导出所有任务
// @Description 导出当前用户的所有任务
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {object} ExportTasksResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks/export/all [get]
func exportTasks(c *gin.Context) {
	// 从上下文获取用户信息
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token, authorization denied"})
		return
	}

	claims, ok := userClaims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	// 从token中提取用户ID
	userID, err := primitive.ObjectIDFromHex(claims["id"].(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// 从数据库获取所有任务
	var tasks []models.Task
	cursor, err := config.DB.Collection("tasks").Find(context.TODO(), bson.M{"created_by": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode tasks"})
		return
	}

	// 转换为响应格式
	var respTasks []Task
	for _, task := range tasks {
		var assignee string
		if task.Assignee != nil {
			assignee = *task.Assignee
		}
		
		respTask := Task{
			Id:          task.ID.Hex(),
			Title:       task.Title,
			Description: task.Description,
			Status:      convertTaskStatusToString(task.Status),
			Priority:    convertTaskPriorityToString(task.Priority),
			Assignee:    assignee,
			CreatedBy:   task.CreatedBy.Hex(),
			CreatedAt:   task.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
		}
		
		// 处理可选的Deadline和ScheduledDate字段
		if task.Deadline != nil {
			formattedDeadline := task.Deadline.Format(time.RFC3339)
			respTask.Deadline = &formattedDeadline
		}
		if task.ScheduledDate != nil {
			formattedScheduledDate := task.ScheduledDate.Format(time.RFC3339)
			respTask.ScheduledDate = &formattedScheduledDate
		}
		
		// 转换评论
		for _, comment := range task.Comments {
			respComment := Comment{
				Id:        comment.ID.Hex(),
				Text:      comment.Text,
				CreatedBy: comment.CreatedBy.Hex(),
				CreatedAt: comment.CreatedAt.Format(time.RFC3339),
			}
			respTask.Comments = append(respTask.Comments, respComment)
		}
		
		respTasks = append(respTasks, respTask)
	}

	// 设置响应头
	filename := fmt.Sprintf("todoing-backup-%s.json", formatDateForFilename(time.Now()))
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	// 返回任务数据
	c.JSON(http.StatusOK, respTasks)
}

// importTasks 导入任务
// @Summary 导入任务
// @Description 批量导入任务
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body ImportTasksRequest true "导入任务请求"
// @Success 200 {object} ImportTasksResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks/import [post]
func importTasks(c *gin.Context) {
	// 从上下文获取用户信息
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token, authorization denied"})
		return
	}

	claims, ok := userClaims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	// 从token中提取用户ID
	userID, err := primitive.ObjectIDFromHex(claims["id"].(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var req ImportTasksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if len(req.Tasks) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No tasks to import"})
		return
	}

	// 导入任务
	importedCount := 0
	var errors []ImportError

	for i, taskData := range req.Tasks {
		// 创建新任务
		task := models.Task{
			ID:          primitive.NewObjectID(),
			Title:       taskData.Title,
			Description: taskData.Description,
			Status:      models.StringToTaskStatus(taskData.Status),
			Priority:    models.StringToTaskPriority(taskData.Priority),
			CreatedBy:   userID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// 设置可选字段
		if taskData.Assignee != "" {
			task.Assignee = &taskData.Assignee
		}
		if taskData.Deadline != nil {
			deadline, err := time.Parse(time.RFC3339, *taskData.Deadline)
			if err == nil {
				task.Deadline = &deadline
			}
		}
		if taskData.ScheduledDate != nil {
			scheduledDate, err := time.Parse(time.RFC3339, *taskData.ScheduledDate)
			if err == nil {
				task.ScheduledDate = &scheduledDate
			}
		}

		// 处理评论
		for _, commentData := range taskData.Comments {
			comment := models.Comment{
				ID:        primitive.NewObjectID(),
				Text:      commentData.Text,
				CreatedBy: userID,
				CreatedAt: time.Now(),
			}
			task.Comments = append(task.Comments, comment)
		}

		// 插入到数据库
		if _, err := config.DB.Collection("tasks").InsertOne(context.TODO(), task); err != nil {
			errors = append(errors, ImportError{
				Index: i,
				Error: err.Error(),
			})
		} else {
			importedCount++
		}
	}

	c.JSON(http.StatusOK, ImportTasksResponse{
		Message:  fmt.Sprintf("Imported %d tasks successfully", importedCount),
		Imported: importedCount,
		Errors:   errors,
	})
}
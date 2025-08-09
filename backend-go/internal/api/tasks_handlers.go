package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/axfinn/todoIng/backend-go/internal/observability"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskDeps struct{ DB *mongo.Database }

type taskRequest struct {
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	Status        string  `json:"status"`
	Priority      string  `json:"priority"`
	Assignee      *string `json:"assignee"`
	Deadline      *string `json:"deadline"`      // 改为 string 类型以兼容前端
	ScheduledDate *string `json:"scheduledDate"` // 改为 string 类型以兼容前端
	Comments      []struct {
		Text      string `json:"text"`
		CreatedBy string `json:"createdBy,omitempty"`
		CreatedAt string `json:"createdAt,omitempty"`
	} `json:"comments"`
}

var allowedStatus = map[string]bool{"To Do": true, "In Progress": true, "Done": true}
var allowedPriority = map[string]bool{"Low": true, "Medium": true, "High": true}

// CreateTask 创建新任务
// @Summary 创建新任务
// @Description 创建一个新的任务项
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param task body taskRequest true "任务信息"
// @Success 200 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/tasks [post]
func (d *TaskDeps) CreateTask(w http.ResponseWriter, r *http.Request) {
	uid := GetUserID(r)
	if uid == "" {
		JSON(w, 401, map[string]string{"msg": "Unauthorized"})
		return
	}
	var req taskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		observability.CtxLog(r.Context(), "CreateTask decode error: %v", err)
		JSON(w, 400, map[string]string{"msg": "Invalid body"})
		return
	}
	observability.CtxLog(r.Context(), "CreateTask received: title=%q, description=%q, status=%q, priority=%q", req.Title, req.Description, req.Status, req.Priority)
	if strings.TrimSpace(req.Title) == "" {
		JSON(w, 400, map[string]string{"msg": "Title is required"})
		return
	}
	if req.Status == "" {
		req.Status = "To Do"
	}
	if !allowedStatus[req.Status] {
		JSON(w, 400, map[string]string{"msg": "Invalid status"})
		return
	}
	if req.Priority == "" {
		req.Priority = "Medium"
	}
	if !allowedPriority[req.Priority] {
		JSON(w, 400, map[string]string{"msg": "Invalid priority"})
		return
	}

	now := time.Now()
	doc := bson.M{
		"title":         req.Title,
		"description":   req.Description,
		"status":        req.Status,
		"priority":      req.Priority,
		"assignee":      req.Assignee,
		"deadline":      parseDate(req.Deadline),
		"scheduledDate": parseDate(req.ScheduledDate),
		"comments":      []bson.M{},
		"createdBy":     uid,
		"createdAt":     now,
		"updatedAt":     now,
	}

	// 处理评论数据，确保兼容原有格式
	for _, c := range req.Comments {
		if strings.TrimSpace(c.Text) != "" {
			comment := bson.M{
				"text":      c.Text,
				"createdBy": uid,
				"createdAt": now,
			}
			doc["comments"] = append(doc["comments"].([]bson.M), comment)
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	res, err := d.DB.Collection("tasks").InsertOne(ctx, doc)
	if err != nil {
		JSON(w, 500, map[string]string{"msg": "DB error"})
		return
	}
	doc["_id"] = res.InsertedID.(primitive.ObjectID).Hex()
	JSON(w, 200, doc)
}

// 辅助函数：解析日期字符串
func parseDate(dateStr *string) *time.Time {
	if dateStr == nil || *dateStr == "" {
		return nil
	}
	if t, err := time.Parse("2006-01-02", *dateStr); err == nil {
		return &t
	}
	if t, err := time.Parse(time.RFC3339, *dateStr); err == nil {
		return &t
	}
	return nil
}

// ListTasks 获取任务列表
// @Summary 获取用户的所有任务
// @Description 获取当前用户创建的所有任务列表，按创建时间倒序排列
// @Tags 任务管理
// @Accept json
// @Produce json
// @Success 200 {object} []map[string]interface{} "任务列表"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/tasks [get]
func (d *TaskDeps) ListTasks(w http.ResponseWriter, r *http.Request) {
	uid := GetUserID(r)
	if uid == "" {
		JSON(w, 401, map[string]string{"msg": "Unauthorized"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	cur, err := d.DB.Collection("tasks").Find(ctx, bson.M{"createdBy": uid}, optionsFindSortCreatedAtDesc())
	if err != nil {
		JSON(w, 500, map[string]string{"msg": "DB error"})
		return
	}
	defer cur.Close(ctx)
	tasks := []bson.M{}
	for cur.Next(ctx) {
		var m bson.M
		if err := cur.Decode(&m); err == nil {
			if id, ok := m["_id"].(primitive.ObjectID); ok {
				m["_id"] = id.Hex()
			}
			tasks = append(tasks, m)
		}
	}
	JSON(w, 200, tasks)
}

// GetTask 获取单个任务详情
// @Summary 获取任务详情
// @Description 根据任务ID获取任务的详细信息
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Success 200 {object} map[string]interface{} "任务详情"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 404 {object} map[string]string "任务不存在"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/tasks/{id} [get]
func (d *TaskDeps) GetTask(w http.ResponseWriter, r *http.Request) {
	uid := GetUserID(r)
	if uid == "" {
		JSON(w, 401, map[string]string{"msg": "Unauthorized"})
		return
	}
	id := muxVar(r, "id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		JSON(w, 404, map[string]string{"msg": "Task not found"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var m bson.M
	if err := d.DB.Collection("tasks").FindOne(ctx, bson.M{"_id": objID, "createdBy": uid}).Decode(&m); err != nil {
		JSON(w, 404, map[string]string{"msg": "Task not found"})
		return
	}
	m["_id"] = id
	JSON(w, 200, m)
}

// UpdateTask 更新任务
// @Summary 更新任务信息
// @Description 根据任务ID更新任务的详细信息
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Param task body taskRequest true "更新的任务信息"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 404 {object} map[string]string "任务不存在"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/tasks/{id} [put]
func (d *TaskDeps) UpdateTask(w http.ResponseWriter, r *http.Request) {
	uid := GetUserID(r)
	if uid == "" {
		JSON(w, 401, map[string]string{"msg": "Unauthorized"})
		return
	}
	id := muxVar(r, "id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		JSON(w, 404, map[string]string{"msg": "Task not found"})
		return
	}
	var req taskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSON(w, 400, map[string]string{"msg": "Invalid body"})
		return
	}
	update := bson.M{}
	if req.Title != "" {
		update["title"] = req.Title
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if req.Status != "" {
		if !allowedStatus[req.Status] {
			JSON(w, 400, map[string]string{"msg": "Invalid status"})
			return
		}
		update["status"] = req.Status
	}
	if req.Priority != "" {
		if !allowedPriority[req.Priority] {
			JSON(w, 400, map[string]string{"msg": "Invalid priority"})
			return
		}
		update["priority"] = req.Priority
	}
	if req.Assignee != nil {
		update["assignee"] = req.Assignee
	}
	if req.Deadline != nil {
		update["deadline"] = req.Deadline
	}
	if req.ScheduledDate != nil {
		update["scheduledDate"] = req.ScheduledDate
	}
	if len(req.Comments) > 0 { // replace comments
		now := time.Now()
		comments := make([]bson.M, 0, len(req.Comments))
		for _, c := range req.Comments {
			if strings.TrimSpace(c.Text) != "" {
				comments = append(comments, bson.M{"text": c.Text, "createdBy": uid, "createdAt": now})
			}
		}
		update["comments"] = comments
	}
	if len(update) == 0 {
		JSON(w, 400, map[string]string{"msg": "No fields to update"})
		return
	}
	update["updatedAt"] = time.Now()
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	res := d.DB.Collection("tasks").FindOneAndUpdate(ctx, bson.M{"_id": objID, "createdBy": uid}, bson.M{"$set": update}, optionsFindOneAndUpdateReturnAfter())
	var m bson.M
	if err := res.Decode(&m); err != nil {
		JSON(w, 404, map[string]string{"msg": "Task not found"})
		return
	}
	if idObj, ok := m["_id"].(primitive.ObjectID); ok {
		m["_id"] = idObj.Hex()
	}
	JSON(w, 200, m)
}

// DeleteTask 删除任务
// @Summary 删除任务
// @Description 根据任务ID删除指定的任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Success 200 {object} map[string]string "删除成功"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 404 {object} map[string]string "任务不存在"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/tasks/{id} [delete]
func (d *TaskDeps) DeleteTask(w http.ResponseWriter, r *http.Request) {
	uid := GetUserID(r)
	if uid == "" {
		JSON(w, 401, map[string]string{"msg": "Unauthorized"})
		return
	}
	id := muxVar(r, "id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		JSON(w, 404, map[string]string{"msg": "Task not found"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	res, err := d.DB.Collection("tasks").DeleteOne(ctx, bson.M{"_id": objID, "createdBy": uid})
	if err != nil || res.DeletedCount == 0 {
		JSON(w, 404, map[string]string{"msg": "Task not found"})
		return
	}
	JSON(w, 200, map[string]string{"msg": "Task removed"})
}

// GET /api/tasks/export/all
func (d *TaskDeps) ExportAll(w http.ResponseWriter, r *http.Request) {
	uid := GetUserID(r)
	if uid == "" {
		JSON(w, 401, map[string]string{"msg": "Unauthorized"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	cur, err := d.DB.Collection("tasks").Find(ctx, bson.M{"createdBy": uid})
	if err != nil {
		JSON(w, 500, map[string]string{"msg": "DB error"})
		return
	}
	defer cur.Close(ctx)
	tasks := []bson.M{}
	for cur.Next(ctx) {
		var m bson.M
		if cur.Decode(&m) == nil {
			if idObj, ok := m["_id"].(primitive.ObjectID); ok {
				m["_id"] = idObj.Hex()
			}
			tasks = append(tasks, m)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=todoing-backup-"+time.Now().Format("2006-01-02")+".json")
	_ = json.NewEncoder(w).Encode(tasks)
}

// POST /api/tasks/import
func (d *TaskDeps) ImportTasks(w http.ResponseWriter, r *http.Request) {
	uid := GetUserID(r)
	if uid == "" {
		JSON(w, 401, map[string]string{"msg": "Unauthorized"})
		return
	}
	var body struct {
		Tasks []map[string]any `json:"tasks"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		JSON(w, 400, map[string]string{"msg": "Invalid body"})
		return
	}
	if len(body.Tasks) == 0 {
		JSON(w, 400, map[string]string{"msg": "No tasks"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()
	col := d.DB.Collection("tasks")
	imported := 0
	var errorsArr []map[string]any
	for i, t := range body.Tasks {
		title, _ := t["title"].(string)
		if strings.TrimSpace(title) == "" {
			continue
		}
		now := time.Now()
		status, _ := t["status"].(string)
		priority, _ := t["priority"].(string)
		if status == "" || !allowedStatus[status] {
			status = "To Do"
		}
		if priority == "" || !allowedPriority[priority] {
			priority = "Medium"
		}
		doc := bson.M{
			"title":         title,
			"description":   t["description"],
			"status":        status,
			"priority":      priority,
			"assignee":      t["assignee"],
			"deadline":      t["deadline"],
			"scheduledDate": t["scheduledDate"],
			"comments":      []bson.M{},
			"createdBy":     uid,
			"createdAt":     now,
			"updatedAt":     now,
		}
		if _, err := col.InsertOne(ctx, doc); err != nil {
			errorsArr = append(errorsArr, map[string]any{"index": i, "error": err.Error()})
			continue
		}
		imported++
	}
	JSON(w, 200, map[string]any{"msg": "Imported tasks", "imported": imported, "errors": errorsArr})
}

// Helper utilities
func muxVar(r *http.Request, key string) string { return mux.Vars(r)[key] }

func optionsFindOneAndUpdateReturnAfter() *options.FindOneAndUpdateOptions {
	return &options.FindOneAndUpdateOptions{ReturnDocument: func(rd options.ReturnDocument) *options.ReturnDocument { v := options.After; return &v }(options.After)}
}

func optionsFindSortCreatedAtDesc() *options.FindOptions {
	return &options.FindOptions{Sort: bson.D{{Key: "createdAt", Value: -1}}}
}

func SetupTaskRoutes(r *mux.Router, deps *TaskDeps) {
	s := r.PathPrefix("/api/tasks").Subrouter()
	s.Handle("", Auth(http.HandlerFunc(deps.ListTasks))).Methods(http.MethodGet)
	s.Handle("", Auth(http.HandlerFunc(deps.CreateTask))).Methods(http.MethodPost)
	s.Handle("/export/all", Auth(http.HandlerFunc(deps.ExportAll))).Methods(http.MethodGet)
	s.Handle("/import", Auth(http.HandlerFunc(deps.ImportTasks))).Methods(http.MethodPost)
	s.Handle("/{id}", Auth(http.HandlerFunc(deps.GetTask))).Methods(http.MethodGet)
	s.Handle("/{id}", Auth(http.HandlerFunc(deps.UpdateTask))).Methods(http.MethodPut)
	s.Handle("/{id}", Auth(http.HandlerFunc(deps.DeleteTask))).Methods(http.MethodDelete)
}

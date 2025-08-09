package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReportDeps struct{ DB *mongo.Database }

// ListReports 获取报表列表
// @Summary 获取用户的所有报表
// @Description 获取当前用户创建的所有报表列表，按创建时间倒序排列
// @Tags 报表管理
// @Accept json
// @Produce json
// @Success 200 {object} []map[string]interface{} "报表列表"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/reports [get]
func (d *ReportDeps) ListReports(w http.ResponseWriter, r *http.Request) {
	uid := GetUserID(r)
	if uid == "" {
		JSON(w, 401, map[string]string{"msg": "Unauthorized"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	cur, err := d.DB.Collection("reports").Find(ctx, bson.M{"userId": uid})
	if err != nil {
		JSON(w, 500, map[string]string{"msg": "DB error"})
		return
	}
	defer cur.Close(ctx)
	var reports []bson.M
	for cur.Next(ctx) {
		var m bson.M
		if cur.Decode(&m) == nil {
			if id, ok := m["_id"].(primitive.ObjectID); ok {
				m["_id"] = id.Hex()
			}
			reports = append(reports, m)
		}
	}
	sort.SliceStable(reports, func(i, j int) bool {
		ti, _ := reports[i]["createdAt"].(time.Time)
		tj, _ := reports[j]["createdAt"].(time.Time)
		return ti.After(tj)
	})
	JSON(w, 200, reports)
}

// GetReport 获取报表详情
// @Summary 获取单个报表详情
// @Description 根据报表ID获取报表的详细信息，包含关联的任务数据
// @Tags 报表管理
// @Accept json
// @Produce json
// @Param id path string true "报表ID"
// @Success 200 {object} map[string]interface{} "报表详情"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 404 {object} map[string]string "报表不存在"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/reports/{id} [get]
func (d *ReportDeps) GetReport(w http.ResponseWriter, r *http.Request) {
	uid := GetUserID(r)
	if uid == "" {
		JSON(w, 401, map[string]string{"msg": "Unauthorized"})
		return
	}
	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		JSON(w, 404, map[string]string{"msg": "Report not found"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	var rep bson.M
	if err := d.DB.Collection("reports").FindOne(ctx, bson.M{"_id": objID, "userId": uid}).Decode(&rep); err != nil {
		JSON(w, 404, map[string]string{"msg": "Report not found"})
		return
	}
	rep["_id"] = id

	// 兼容 Node.js 版本：populate 任务详情
	if taskIDs, ok := rep["tasks"].([]interface{}); ok && len(taskIDs) > 0 {
		var tasks []bson.M
		for _, taskID := range taskIDs {
			if taskIDStr, ok := taskID.(string); ok {
				if taskObjID, err := primitive.ObjectIDFromHex(taskIDStr); err == nil {
					var task bson.M
					if err := d.DB.Collection("tasks").FindOne(ctx, bson.M{"_id": taskObjID}).Decode(&task); err == nil {
						if id, ok := task["_id"].(primitive.ObjectID); ok {
							task["_id"] = id.Hex()
						}
						tasks = append(tasks, task)
					}
				}
			}
		}
		rep["tasks"] = tasks
	}

	JSON(w, 200, rep)
}

type generateReportRequest struct {
	Type      string `json:"type"`
	Period    string `json:"period"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

// POST /api/reports/generate
func (d *ReportDeps) GenerateReport(w http.ResponseWriter, r *http.Request) {
	uid := GetUserID(r)
	if uid == "" {
		JSON(w, 401, map[string]string{"msg": "Unauthorized"})
		return
	}
	var req generateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSON(w, 400, map[string]string{"msg": "Invalid body"})
		return
	}
	if req.Type == "" || req.Period == "" || req.StartDate == "" || req.EndDate == "" {
		JSON(w, 400, map[string]string{"msg": "Please provide type, period, startDate, and endDate"})
		return
	}
	if req.Type != "daily" && req.Type != "weekly" && req.Type != "monthly" {
		JSON(w, 400, map[string]string{"msg": "Invalid type"})
		return
	}

	// 支持多种日期格式，兼容前端
	start, err1 := parseFlexibleDate(req.StartDate)
	end, err2 := parseFlexibleDate(req.EndDate)
	if err1 != nil || err2 != nil {
		JSON(w, 400, map[string]string{"msg": "Invalid date format"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	tasksCur, err := d.DB.Collection("tasks").Find(ctx, bson.M{"createdBy": uid, "createdAt": bson.M{"$gte": start, "$lte": end}})
	if err != nil {
		JSON(w, 500, map[string]string{"msg": "DB error"})
		return
	}
	defer tasksCur.Close(ctx)
	var tasks []bson.M
	for tasksCur.Next(ctx) {
		var t bson.M
		if tasksCur.Decode(&t) == nil {
			if id, ok := t["_id"].(primitive.ObjectID); ok {
				t["_id"] = id.Hex()
			}
			tasks = append(tasks, t)
		}
	}
	total := len(tasks)
	completed := 0
	inProgress := 0
	overdue := 0
	now := time.Now()
	for _, t := range tasks {
		status, _ := t["status"].(string)
		if status == "Done" {
			completed++
		} else if status == "In Progress" {
			inProgress++
		}
		if ddl, ok := t["deadline"].(time.Time); ok {
			if ddl.Before(now) && status != "Done" {
				overdue++
			}
		}
	}
	completionRate := 0
	if total > 0 {
		completionRate = int(float64(completed) / float64(total) * 100)
	}
	titles := map[string]string{"daily": "日报 - " + req.Period, "weekly": "周报 - " + req.Period, "monthly": "月报 - " + req.Period}
	var sb strings.Builder
	sb.WriteString("# " + titles[req.Type] + "\n\n")
	sb.WriteString("报告周期: " + start.Format("2006/01/02") + " - " + end.Format("2006/01/02") + "\n\n")
	sb.WriteString("## 统计信息\n")
	sb.WriteString("- 总任务数: " + itoa(total) + "\n")
	sb.WriteString("- 已完成任务: " + itoa(completed) + "\n")
	sb.WriteString("- 进行中任务: " + itoa(inProgress) + "\n")
	sb.WriteString("- 过期任务: " + itoa(overdue) + "\n")
	sb.WriteString("- 完成率: " + itoa(completionRate) + "%\n\n")
	sb.WriteString("## 任务详情\n")
	if len(tasks) == 0 {
		sb.WriteString("此周期内未找到任务。\n")
	} else {
		for _, t := range tasks {
			title, _ := t["title"].(string)
			status, _ := t["status"].(string)
			priority, _ := t["priority"].(string)
			createdAt, _ := t["createdAt"].(time.Time)
			updatedAt, _ := t["updatedAt"].(time.Time)
			sb.WriteString("### 任务: " + title + "\n")
			sb.WriteString("- **任务状态**: " + status + "\n")
			sb.WriteString("- **任务优先级**: " + priority + "\n")
			if !createdAt.IsZero() {
				sb.WriteString("- **创建时间**: " + createdAt.Format("2006-01-02 15:04:05") + "\n")
			}
			if !updatedAt.IsZero() {
				sb.WriteString("- **更新时间**: " + updatedAt.Format("2006-01-02 15:04:05") + "\n")
			}
			if ddl, ok := t["deadline"].(time.Time); ok && !ddl.IsZero() {
				sb.WriteString("- **截止日期**: " + ddl.Format("2006-01-02 15:04:05") + "\n")
			}
			if sch, ok := t["scheduledDate"].(time.Time); ok && !sch.IsZero() {
				sb.WriteString("- **计划日期**: " + sch.Format("2006-01-02 15:04:05") + "\n")
			}
			desc, _ := t["description"].(string)
			if desc == "" {
				desc = "无"
			}
			sb.WriteString("- **任务描述**: " + desc + "\n\n")
			sb.WriteString("#### 任务活动时间线\n")
			sb.WriteString("- " + createdAt.Format("2006-01-02 15:04:05") + ": 任务已创建\n")
			if !updatedAt.IsZero() && updatedAt.After(createdAt.Add(5*time.Second)) {
				sb.WriteString("- " + updatedAt.Format("2006-01-02 15:04:05") + ": 任务已更新\n")
			}
			sb.WriteString("\n---\n\n")
		}
	}
	content := sb.String()
	reportDoc := bson.M{
		"userId":          uid,
		"type":            req.Type,
		"period":          req.Period,
		"title":           titles[req.Type],
		"content":         content,
		"polishedContent": nil,
		"tasks":           extractTaskIDs(tasks),
		"statistics": bson.M{
			"totalTasks":      total,
			"completedTasks":  completed,
			"inProgressTasks": inProgress,
			"overdueTasks":    overdue,
			"completionRate":  completionRate,
		},
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
	}
	res, err := d.DB.Collection("reports").InsertOne(ctx, reportDoc)
	if err != nil {
		JSON(w, 500, map[string]string{"msg": "DB error"})
		return
	}
	reportDoc["_id"] = res.InsertedID.(primitive.ObjectID).Hex()
	JSON(w, 200, reportDoc)
}

func extractTaskIDs(tasks []bson.M) []string {
	ids := make([]string, 0, len(tasks))
	for _, t := range tasks {
		if id, ok := t["_id"].(string); ok {
			ids = append(ids, id)
		}
	}
	return ids
}

// POST /api/reports/{id}/polish
func (d *ReportDeps) PolishReport(w http.ResponseWriter, r *http.Request) {
	uid := GetUserID(r)
	if uid == "" {
		JSON(w, 401, map[string]string{"msg": "Unauthorized"})
		return
	}
	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		JSON(w, 404, map[string]string{"msg": "Report not found"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	var rep bson.M
	if err := d.DB.Collection("reports").FindOne(ctx, bson.M{"_id": objID, "userId": uid}).Decode(&rep); err != nil {
		JSON(w, 404, map[string]string{"msg": "Report not found"})
		return
	}
	original := rep["content"].(string)
	polished := original
	if pc, ok := rep["polishedContent"].(string); ok && pc != "" {
		polished = pc
	} else {
		polished = original + "\n\n[Polished Placeholder]"
	}
	_, err = d.DB.Collection("reports").UpdateByID(ctx, objID, bson.M{"$set": bson.M{"polishedContent": polished, "updatedAt": time.Now()}})
	if err != nil {
		JSON(w, 500, map[string]string{"msg": "DB error"})
		return
	}
	rep["polishedContent"] = polished
	rep["_id"] = id
	JSON(w, 200, rep)
}

// DELETE /api/reports/{id}
func (d *ReportDeps) DeleteReport(w http.ResponseWriter, r *http.Request) {
	uid := GetUserID(r)
	if uid == "" {
		JSON(w, 401, map[string]string{"msg": "Unauthorized"})
		return
	}
	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		JSON(w, 404, map[string]string{"msg": "Report not found"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	res, err := d.DB.Collection("reports").DeleteOne(ctx, bson.M{"_id": objID, "userId": uid})
	if err != nil || res.DeletedCount == 0 {
		JSON(w, 404, map[string]string{"msg": "Report not found"})
		return
	}
	JSON(w, 200, map[string]string{"msg": "Report removed"})
}

// GET /api/reports/{id}/export/{format}
func (d *ReportDeps) ExportReport(w http.ResponseWriter, r *http.Request) {
	uid := GetUserID(r)
	if uid == "" {
		JSON(w, 401, map[string]string{"msg": "Unauthorized"})
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	format := strings.ToLower(vars["format"])
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		JSON(w, 404, map[string]string{"msg": "Report not found"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var rep bson.M
	if err := d.DB.Collection("reports").FindOne(ctx, bson.M{"_id": objID, "userId": uid}).Decode(&rep); err != nil {
		JSON(w, 404, map[string]string{"msg": "Report not found"})
		return
	}
	content := ""
	if pc, ok := rep["polishedContent"].(string); ok && pc != "" {
		content = pc
	} else {
		content = rep["content"].(string)
	}
	fileExt := "txt"
	contentType := "text/plain"
	if format == "md" || format == "markdown" {
		fileExt = "md"
		contentType = "text/markdown"
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "attachment; filename=report-"+rep["period"].(string)+"."+fileExt)
	_, _ = w.Write([]byte(content))
}

func SetupReportRoutes(r *mux.Router, deps *ReportDeps) {
	s := r.PathPrefix("/api/reports").Subrouter()
	s.Handle("", Auth(http.HandlerFunc(deps.ListReports))).Methods(http.MethodGet)
	s.Handle("/generate", Auth(http.HandlerFunc(deps.GenerateReport))).Methods(http.MethodPost)
	s.Handle("/{id}", Auth(http.HandlerFunc(deps.GetReport))).Methods(http.MethodGet)
	s.Handle("/{id}", Auth(http.HandlerFunc(deps.DeleteReport))).Methods(http.MethodDelete)
	s.Handle("/{id}/polish", Auth(http.HandlerFunc(deps.PolishReport))).Methods(http.MethodPost)
	s.Handle("/{id}/export/{format}", Auth(http.HandlerFunc(deps.ExportReport))).Methods(http.MethodGet)
}

// small helpers without importing strconv to keep file self-contained
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	b := make([]byte, 0, 6)
	for i > 0 {
		b = append([]byte{byte('0' + i%10)}, b...)
		i /= 10
	}
	return string(b)
}

// 灵活的日期解析，兼容前端多种格式
func parseFlexibleDate(dateStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, errors.New("unsupported date format")
}

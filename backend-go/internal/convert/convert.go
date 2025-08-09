package convert

import (
	"github.com/axfinn/todoIng/backend-go/internal/models"
	pb "github.com/axfinn/todoIng/backend-go/pkg/api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// 这样可以保持现有代码正常工作

// UserToProto 将内部用户模型转换为 protobuf 模型
func UserToProto(user *models.User) *pb.User {
	if user == nil {
		return nil
	}

	return &pb.User{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.CreatedAt), // 使用 CreatedAt 因为模型中没有 UpdatedAt
	}
}

// ProtoToUser 将 protobuf 用户模型转换为内部模型
func ProtoToUser(pbUser *pb.User) *models.User {
	if pbUser == nil {
		return nil
	}

	user := &models.User{
		Username: pbUser.Username,
		Email:    pbUser.Email,
	}

	if pbUser.CreatedAt != nil {
		user.CreatedAt = pbUser.CreatedAt.AsTime()
	}

	return user
}

// TaskStatusToProto 将内部任务状态转换为 protobuf 状态
func TaskStatusToProto(status string) pb.TaskStatus {
	switch status {
	case "todo":
		return pb.TaskStatus_TASK_STATUS_TODO
	case "in_progress":
		return pb.TaskStatus_TASK_STATUS_IN_PROGRESS
	case "done":
		return pb.TaskStatus_TASK_STATUS_DONE
	default:
		return pb.TaskStatus_TASK_STATUS_UNSPECIFIED
	}
}

// ProtoToTaskStatus 将 protobuf 状态转换为内部任务状态
func ProtoToTaskStatus(status pb.TaskStatus) string {
	switch status {
	case pb.TaskStatus_TASK_STATUS_TODO:
		return "todo"
	case pb.TaskStatus_TASK_STATUS_IN_PROGRESS:
		return "in_progress"
	case pb.TaskStatus_TASK_STATUS_DONE:
		return "done"
	default:
		return "todo"
	}
}

// TaskPriorityToProto 将内部任务优先级转换为 protobuf 优先级
func TaskPriorityToProto(priority string) pb.TaskPriority {
	switch priority {
	case "low":
		return pb.TaskPriority_TASK_PRIORITY_LOW
	case "medium":
		return pb.TaskPriority_TASK_PRIORITY_MEDIUM
	case "high":
		return pb.TaskPriority_TASK_PRIORITY_HIGH
	default:
		return pb.TaskPriority_TASK_PRIORITY_UNSPECIFIED
	}
}

// ProtoToTaskPriority 将 protobuf 优先级转换为内部任务优先级
func ProtoToTaskPriority(priority pb.TaskPriority) string {
	switch priority {
	case pb.TaskPriority_TASK_PRIORITY_LOW:
		return "low"
	case pb.TaskPriority_TASK_PRIORITY_MEDIUM:
		return "medium"
	case pb.TaskPriority_TASK_PRIORITY_HIGH:
		return "high"
	default:
		return "medium"
	}
}

// TaskToProto 将内部任务模型转换为 protobuf 模型
func TaskToProto(task *models.Task) *pb.Task {
	if task == nil {
		return nil
	}

	var dueDate *timestamppb.Timestamp
	if task.Deadline != nil && !task.Deadline.IsZero() {
		dueDate = timestamppb.New(*task.Deadline)
	}

	return &pb.Task{
		Id:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      TaskStatusToProto(task.Status),
		Priority:    TaskPriorityToProto(task.Priority),
		DueDate:     dueDate,
		UserId:      task.CreatedBy, // 使用 CreatedBy 作为 UserId
		CreatedAt:   timestamppb.New(task.CreatedAt),
		UpdatedAt:   timestamppb.New(task.UpdatedAt),
	}
}

// ProtoToTask 将 protobuf 任务模型转换为内部模型
func ProtoToTask(pbTask *pb.Task) *models.Task {
	if pbTask == nil {
		return nil
	}

	task := &models.Task{
		Title:       pbTask.Title,
		Description: pbTask.Description,
		Status:      ProtoToTaskStatus(pbTask.Status),
		Priority:    ProtoToTaskPriority(pbTask.Priority),
	}

	if pbTask.DueDate != nil {
		deadline := pbTask.DueDate.AsTime()
		task.Deadline = &deadline
	}
	if pbTask.CreatedAt != nil {
		task.CreatedAt = pbTask.CreatedAt.AsTime()
	}
	if pbTask.UpdatedAt != nil {
		task.UpdatedAt = pbTask.UpdatedAt.AsTime()
	}

	return task
}

// ReportTypeToProto 将内部报表类型转换为 protobuf 类型
func ReportTypeToProto(reportType string) pb.ReportType {
	switch reportType {
	case "daily":
		return pb.ReportType_REPORT_TYPE_DAILY
	case "weekly":
		return pb.ReportType_REPORT_TYPE_WEEKLY
	case "monthly":
		return pb.ReportType_REPORT_TYPE_MONTHLY
	default:
		return pb.ReportType_REPORT_TYPE_UNSPECIFIED
	}
}

// ProtoToReportType 将 protobuf 类型转换为内部报表类型
func ProtoToReportType(reportType pb.ReportType) string {
	switch reportType {
	case pb.ReportType_REPORT_TYPE_DAILY:
		return "daily"
	case pb.ReportType_REPORT_TYPE_WEEKLY:
		return "weekly"
	case pb.ReportType_REPORT_TYPE_MONTHLY:
		return "monthly"
	default:
		return "daily"
	}
}

// ReportToProto 将内部报表模型转换为 protobuf 模型
func ReportToProto(report *models.Report) *pb.Report {
	if report == nil {
		return nil
	}

	// 注意：这里的 Tasks 字段在内部模型中是 []string，但 proto 中期望的是 []*pb.Task
	// 实际使用时可能需要额外的查询来获取完整的任务信息
	pbTasks := make([]*pb.Task, 0) // 暂时为空数组

	stats := &pb.ReportStats{
		TotalTasks:      int32(report.Statistics.TotalTasks),
		CompletedTasks:  int32(report.Statistics.CompletedTasks),
		PendingTasks:    int32(report.Statistics.TotalTasks - report.Statistics.CompletedTasks - report.Statistics.InProgressTasks),
		InProgressTasks: int32(report.Statistics.InProgressTasks),
		CompletionRate:  float64(report.Statistics.CompletionRate),
	}

	return &pb.Report{
		Id:        report.ID,
		Title:     report.Title,
		Type:      ReportTypeToProto(report.Type),
		StartDate: timestamppb.New(report.CreatedAt), // 使用 CreatedAt 作为开始时间
		EndDate:   timestamppb.New(report.UpdatedAt), // 使用 UpdatedAt 作为结束时间
		UserId:    report.UserID,
		CreatedAt: timestamppb.New(report.CreatedAt),
		Tasks:     pbTasks,
		Stats:     stats,
	}
}

// TimeToString 将时间转换为字符串（保持与现有API兼容）
func TimeToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"todoing-backend/src/config"
)

type TaskStatus string
type TaskPriority string

const (
	StatusToDo       TaskStatus = "To Do"
	StatusInProgress TaskStatus = "In Progress"
	StatusDone       TaskStatus = "Done"
)

const (
	PriorityLow    TaskPriority = "Low"
	PriorityMedium TaskPriority = "Medium"
	PriorityHigh   TaskPriority = "High"
)

type Comment struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Text      string             `json:"text" bson:"text"`
	CreatedBy primitive.ObjectID `json:"createdBy" bson:"createdBy"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

type Task struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title         string             `json:"title" bson:"title" validate:"required"`
	Description   string             `json:"description" bson:"description"`
	Status        TaskStatus         `json:"status" bson:"status"`
	Priority      TaskPriority       `json:"priority" bson:"priority"`
	Assignee      *string            `json:"assignee" bson:"assignee"`
	CreatedBy     primitive.ObjectID `json:"createdBy" bson:"createdBy" validate:"required"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt"`
	Deadline      *time.Time         `json:"deadline" bson:"deadline"`
	ScheduledDate *time.Time         `json:"scheduledDate" bson:"scheduledDate"`
	Comments      []Comment          `json:"comments" bson:"comments"`
}

// GetTasksByUserID 根据用户ID获取任务列表
func GetTasksByUserID(userID primitive.ObjectID) ([]*Task, error) {
	var tasks []*Task

	// 查询条件：createdBy等于用户ID
	filter := bson.M{"createdBy": userID}

	// 按创建时间倒序排序
	opts := options.Find()
	opts.SetSort(bson.D{{"createdAt", -1}})

	cursor, err := config.DB.Collection("tasks").Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// 遍历结果
	for cursor.Next(context.TODO()) {
		var task Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}

		tasks = append(tasks, &task)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// StringToTaskStatus 将字符串转换为TaskStatus
func StringToTaskStatus(s string) TaskStatus {
	switch s {
	case "To Do", "TO_DO":
		return StatusToDo
	case "In Progress", "IN_PROGRESS":
		return StatusInProgress
	case "Done", "DONE":
		return StatusDone
	default:
		return StatusToDo
	}
}

// StringToTaskPriority 将字符串转换为TaskPriority
func StringToTaskPriority(s string) TaskPriority {
	switch s {
	case "Low", "LOW":
		return PriorityLow
	case "Medium", "MEDIUM":
		return PriorityMedium
	case "High", "HIGH":
		return PriorityHigh
	default:
		return PriorityMedium
	}
}

// CreateTask 创建新任务
func CreateTask(task *Task) error {
	// 设置创建时间和更新时间
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	// 注意：不要为nil的时间字段设置默认值，保持它们为nil
	// 这样在后续处理中可以正确地将它们视为可选字段

	// 插入数据库
	_, err := config.DB.Collection("tasks").InsertOne(context.TODO(), task)
	return err
}

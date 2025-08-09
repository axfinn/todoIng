package convert

import (
	"testing"
	"time"

	"github.com/axfinn/todoIng/backend-go/internal/models"
	pb "github.com/axfinn/todoIng/backend-go/pkg/api/v1"
)

func TestUserToProto(t *testing.T) {
	// 测试数据
	userID := "60d5ecb74eb3b8001f8b4567"
	createdAt := time.Now()

	user := &models.User{
		ID:        userID,
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: createdAt,
	}

	// 转换为 protobuf
	pbUser := UserToProto(user)

	// 验证转换结果
	if pbUser == nil {
		t.Fatal("Expected non-nil pbUser")
	}

	if pbUser.Id != userID {
		t.Errorf("Expected ID %s, got %s", userID, pbUser.Id)
	}

	if pbUser.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, pbUser.Username)
	}

	if pbUser.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, pbUser.Email)
	}

	if pbUser.CreatedAt.AsTime().Unix() != createdAt.Unix() {
		t.Errorf("Expected createdAt %v, got %v", createdAt, pbUser.CreatedAt.AsTime())
	}
}

func TestProtoToUser(t *testing.T) {
	// 测试数据
	pbUser := &pb.User{
		Id:       "60d5ecb74eb3b8001f8b4567",
		Username: "testuser",
		Email:    "test@example.com",
	}

	// 转换为内部模型
	user := ProtoToUser(pbUser)

	// 验证转换结果
	if user == nil {
		t.Fatal("Expected non-nil user")
	}

	if user.Username != pbUser.Username {
		t.Errorf("Expected username %s, got %s", pbUser.Username, user.Username)
	}

	if user.Email != pbUser.Email {
		t.Errorf("Expected email %s, got %s", pbUser.Email, user.Email)
	}
}

func TestTaskStatusConversion(t *testing.T) {
	tests := []struct {
		internal string
		proto    pb.TaskStatus
	}{
		{"todo", pb.TaskStatus_TASK_STATUS_TODO},
		{"in_progress", pb.TaskStatus_TASK_STATUS_IN_PROGRESS},
		{"done", pb.TaskStatus_TASK_STATUS_DONE},
		{"unknown", pb.TaskStatus_TASK_STATUS_UNSPECIFIED},
	}

	for _, tt := range tests {
		t.Run(tt.internal, func(t *testing.T) {
			// 测试内部到 proto 的转换
			proto := TaskStatusToProto(tt.internal)
			if proto != tt.proto {
				t.Errorf("Expected proto status %v, got %v", tt.proto, proto)
			}

			// 测试 proto 到内部的转换（除了 unknown 情况）
			if tt.internal != "unknown" {
				internal := ProtoToTaskStatus(tt.proto)
				if internal != tt.internal {
					t.Errorf("Expected internal status %s, got %s", tt.internal, internal)
				}
			}
		})
	}
}

func TestTaskPriorityConversion(t *testing.T) {
	tests := []struct {
		internal string
		proto    pb.TaskPriority
	}{
		{"low", pb.TaskPriority_TASK_PRIORITY_LOW},
		{"medium", pb.TaskPriority_TASK_PRIORITY_MEDIUM},
		{"high", pb.TaskPriority_TASK_PRIORITY_HIGH},
		{"unknown", pb.TaskPriority_TASK_PRIORITY_UNSPECIFIED},
	}

	for _, tt := range tests {
		t.Run(tt.internal, func(t *testing.T) {
			// 测试内部到 proto 的转换
			proto := TaskPriorityToProto(tt.internal)
			if proto != tt.proto {
				t.Errorf("Expected proto priority %v, got %v", tt.proto, proto)
			}

			// 测试 proto 到内部的转换（除了 unknown 情况）
			if tt.internal != "unknown" {
				internal := ProtoToTaskPriority(tt.proto)
				if internal != tt.internal {
					t.Errorf("Expected internal priority %s, got %s", tt.internal, internal)
				}
			}
		})
	}
}

func TestNilInputs(t *testing.T) {
	// 测试 nil 输入的处理
	if UserToProto(nil) != nil {
		t.Error("Expected nil result for nil input")
	}

	if ProtoToUser(nil) != nil {
		t.Error("Expected nil result for nil input")
	}

	if TaskToProto(nil) != nil {
		t.Error("Expected nil result for nil input")
	}

	if ProtoToTask(nil) != nil {
		t.Error("Expected nil result for nil input")
	}

	if ReportToProto(nil) != nil {
		t.Error("Expected nil result for nil input")
	}
}

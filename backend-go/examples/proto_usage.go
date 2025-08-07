package main

import (
	"fmt"
	"log"
	"time"

	"todoing-backend/proto/gen/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// 示例：展示如何使用protobuf生成的结构
func main() {
	// 创建一个用户对象
	user := &proto.User{
		Id:       "12345",
		Username: "john_doe",
		Email:    "john@example.com",
		CreatedAt: timestamppb.New(time.Now()),
	}

	// 创建注册请求
	registerReq := &proto.RegisterRequest{
		Username: "jane_doe",
		Email:    "jane@example.com",
		Password: "secret123",
	}

	// 创建注册响应
	registerResp := &proto.RegisterResponse{
		Token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		User:    user,
		Message: "User registered successfully",
	}

	// 创建任务对象
	task := &proto.Task{
		Id:          "task123",
		Title:       "Complete project",
		Description: "Finish the protobuf integration",
		Status:      proto.TaskStatus_TO_DO,
		Priority:    proto.TaskPriority_HIGH,
		CreatedBy:   user.Id,
		CreatedAt:   timestamppb.New(time.Now()),
		UpdatedAt:   timestamppb.New(time.Now()),
	}

	// 创建任务列表响应
	tasksResp := &proto.GetTasksResponse{
		Tasks:   []*proto.Task{task},
		Message: "Tasks retrieved successfully",
	}

	// 创建报告对象
	report := &proto.Report{
		Id:      "report123",
		UserId:  user.Id,
		Type:    proto.ReportType_WEEKLY,
		Period:  "2023-W20",
		Title:   "Weekly Report",
		Content: "# Weekly Report\n\nThis is a sample report content.",
		Statistics: &proto.ReportStatistics{
			TotalTasks:       10,
			CompletedTasks:   7,
			InProgressTasks:  2,
			OverdueTasks:     1,
			CompletionRate:   70,
		},
		CreatedAt: timestamppb.New(time.Now()),
		UpdatedAt: timestamppb.New(time.Now()),
	}

	// 打印对象信息
	fmt.Println("User:", user.Username, "(", user.Email, ")")
	fmt.Println("Registration Token:", registerResp.Token)
	fmt.Println("Task:", task.Title, " - Status:", task.Status)
	fmt.Println("Report:", report.Title, " - Type:", report.Type)

	// 展示任务列表
	fmt.Println("\nTasks:")
	for _, t := range tasksResp.Tasks {
		fmt.Printf("- %s (Priority: %s, Status: %s)\n", t.Title, t.Priority, t.Status)
	}

	// 展示报告统计
	fmt.Printf("\nReport Statistics:\n")
	fmt.Printf("Total Tasks: %d\n", report.Statistics.TotalTasks)
	fmt.Printf("Completed Tasks: %d\n", report.Statistics.CompletedTasks)
	fmt.Printf("Completion Rate: %d%%\n", report.Statistics.CompletionRate)
}

// 示例：展示如何在API服务中使用protobuf结构
type UserService struct{}

func (s *UserService) Register(req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	// 实际实现中，这里会处理用户注册逻辑
	log.Printf("Registering user: %s (%s)", req.Username, req.Email)
	
	// 创建用户对象
	user := &proto.User{
		Id:       "generated-user-id",
		Username: req.Username,
		Email:    req.Email,
		CreatedAt: timestamppb.New(time.Now()),
	}
	
	// 返回响应
	return &proto.RegisterResponse{
		Token:   "generated-jwt-token",
		User:    user,
		Message: "User registered successfully",
	}, nil
}

func (s *UserService) Login(req *proto.LoginRequest) (*proto.LoginResponse, error) {
	// 实际实现中，这里会处理用户登录逻辑
	log.Printf("Logging in user: %s", req.Email)
	
	// 创建用户对象
	user := &proto.User{
		Id:       "found-user-id",
		Username: "found-username",
		Email:    req.Email,
		CreatedAt: timestamppb.New(time.Now()),
	}
	
	// 返回响应
	return &proto.LoginResponse{
		Token:   "generated-jwt-token",
		User:    user,
		Message: "User logged in successfully",
	}, nil
}

// 示例：展示如何在任务服务中使用protobuf结构
type TaskService struct{}

func (s *TaskService) CreateTask(req *proto.CreateTaskRequest) (*proto.CreateTaskResponse, error) {
	// 实际实现中，这里会处理创建任务逻辑
	log.Printf("Creating task: %s", req.Title)
	
	// 创建任务对象
	task := &proto.Task{
		Id:          "generated-task-id",
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		Assignee:    req.Assignee,
		CreatedBy:   "user-id",
		CreatedAt:   timestamppb.New(time.Now()),
		UpdatedAt:   timestamppb.New(time.Now()),
		Deadline:    req.Deadline,
		ScheduledDate: req.ScheduledDate,
	}
	
	// 返回响应
	return &proto.CreateTaskResponse{
		Task:    task,
		Message: "Task created successfully",
	}, nil
}

func (s *TaskService) GetTasks(req *proto.GetTasksRequest) (*proto.GetTasksResponse, error) {
	// 实际实现中，这里会从数据库获取任务列表
	log.Printf("Getting tasks for user with token: %s", req.Token)
	
	// 返回响应
	return &proto.GetTasksResponse{
		Tasks:   []*proto.Task{}, // 实际实现中会包含任务列表
		Message: "Tasks retrieved successfully",
	}, nil
}

// 示例：展示如何在报告服务中使用protobuf结构
type ReportService struct{}

func (s *ReportService) GenerateReport(req *proto.GenerateReportRequest) (*proto.GenerateReportResponse, error) {
	// 实际实现中，这里会处理生成报告逻辑
	log.Printf("Generating %s report for period: %s", req.Type, req.Period)
	
	// 创建报告对象
	report := &proto.Report{
		Id:      "generated-report-id",
		UserId:  "user-id",
		Type:    req.Type,
		Period:  req.Period,
		Title:   fmt.Sprintf("%s Report", req.Type),
		Content: "Generated report content",
		Statistics: &proto.ReportStatistics{
			TotalTasks:       0,
			CompletedTasks:   0,
			InProgressTasks:  0,
			OverdueTasks:     0,
			CompletionRate:   0,
		},
		CreatedAt: timestamppb.New(time.Now()),
		UpdatedAt: timestamppb.New(time.Now()),
	}
	
	// 返回响应
	return &proto.GenerateReportResponse{
		Report:  report,
		Message: "Report generated successfully",
	}, nil
}
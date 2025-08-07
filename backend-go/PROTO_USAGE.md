# 在项目中使用 Protobuf 结构

本文档说明了如何在项目中使用 Protobuf 生成的结构体。

## 目录结构

```
backend-go/
├── proto/
│   ├── user.proto          # 用户相关消息和服务定义
│   ├── task.proto          # 任务相关消息和服务定义
│   ├── report.proto        # 报告相关消息和服务定义
│   └── gen/                # 生成的Go代码
│       └── proto/
│           ├── user.pb.go
│           ├── task.pb.go
│           ├── report.pb.go
│           ├── user_grpc.pb.go
│           ├── task_grpc.pb.go
│           └── report_grpc.pb.go
```

## Protobuf 结构使用示例

### 1. 导入生成的 Protobuf 包

```go
import (
    pb "todoing-backend/proto/gen/proto"
)
```

### 2. 使用消息结构体

```go
// 创建用户对象
user := &pb.User{
    Id:       "user123",
    Username: "john_doe",
    Email:    "john@example.com",
}

// 创建任务对象
task := &pb.Task{
    Id:          "task123",
    Title:       "完成项目",
    Description: "使用 protobuf 重构后端",
    Status:      pb.TaskStatus_TO_DO,
    Priority:    pb.TaskPriority_HIGH,
}
```

### 3. 在 gRPC 服务中使用

```go
type userServiceServer struct {
    pb.UnimplementedUserServiceServer
}

func (s *userServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
    // 处理注册逻辑
    user := &pb.User{
        Id:       "generated-id",
        Username: req.Username,
        Email:    req.Email,
    }
    
    return &pb.RegisterResponse{
        Token: "jwt-token",
        User:  user,
    }, nil
}
```

## 生成 Protobuf 代码

### 1. 安装必要的工具

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 2. 生成 Go 代码

```bash
protoc --go_out=./proto/gen --go-grpc_out=./proto/gen proto/*.proto
```

## 主要消息类型

### 用户相关
- `User`: 用户信息
- `RegisterRequest/RegisterResponse`: 注册请求/响应
- `LoginRequest/LoginResponse`: 登录请求/响应

### 任务相关
- `Task`: 任务信息
- `Comment`: 评论信息
- `CreateTaskRequest/CreateTaskResponse`: 创建任务请求/响应
- `UpdateTaskRequest/UpdateTaskResponse`: 更新任务请求/响应

### 报告相关
- `Report`: 报告信息
- `ReportStatistics`: 报告统计信息
- `GenerateReportRequest/GenerateReportResponse`: 生成报告请求/响应

## 枚举类型

### 任务状态
```go
TaskStatus_TO_DO = 0
TaskStatus_IN_PROGRESS = 1
TaskStatus_DONE = 2
```

### 任务优先级
```go
TaskPriority_LOW = 0
TaskPriority_MEDIUM = 1
TaskPriority_HIGH = 2
```

### 报告类型
```go
ReportType_DAILY = 0
ReportType_WEEKLY = 1
ReportType_MONTHLY = 2
```

## 在 REST API 中使用 Protobuf 结构

虽然项目当前使用 REST API，但可以在处理请求和响应时使用 Protobuf 结构：

```go
func createUser(c *gin.Context) {
    var req pb.CreateTaskRequest
    // 将 JSON 请求体解析为 Protobuf 结构
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 使用 req 中的数据创建任务
    task := &pb.Task{
        Title:       req.Title,
        Description: req.Description,
        Status:      req.Status,
        Priority:    req.Priority,
    }
    
    // 返回 Protobuf 结构作为 JSON 响应
    c.JSON(200, task)
}
```

## 服务监听配置

服务默认监听在所有网络接口上（0.0.0.0），支持局域网访问：

- HTTP API 服务: `0.0.0.0:5001`
- gRPC 服务: `0.0.0.0:5002`
- Swagger UI: `http://0.0.0.0:5001/swagger/index.html`

## Docker 部署

### 构建 Docker 镜像

```bash
cd backend-go
docker build -t todoing-backend .
```

### 运行容器

```bash
docker run -d \
  -p 5001:5001 \
  -p 5002:5002 \
  --name todoing-backend \
  todoing-backend
```

### 环境变量配置

可以通过以下环境变量配置应用：

- `PORT`: HTTP服务端口，默认为5001
- `MONGODB_URI`: MongoDB连接字符串
- `JWT_SECRET`: JWT密钥

示例：

```bash
docker run -d \
  -p 5001:5001 \
  -p 5002:5002 \
  -e PORT=5001 \
  -e MONGODB_URI=mongodb://mongo:27017/todoing \
  -e JWT_SECRET=mysecretkey \
  --name todoing-backend \
  todoing-backend
```

## 优势

1. **强类型**: Protobuf 提供强类型检查，减少运行时错误
2. **向前兼容**: 可以安全地添加新字段而不破坏现有代码
3. **跨语言支持**: 可以在不同语言间共享数据结构
4. **高效序列化**: Protobuf 序列化比 JSON 更高效
5. **文档化**: .proto 文件本身就是 API 文档

## 最佳实践

1. 在处理 API 请求时，将 JSON 数据解析为 Protobuf 结构
2. 在数据库操作中，使用 Protobuf 结构作为数据传输对象
3. 在服务间通信中，使用 Protobuf 进行高效的数据传输
4. 保持 .proto 文件的清晰和文档化
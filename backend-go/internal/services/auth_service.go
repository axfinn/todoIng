package services

import (
	"context"
	"fmt"

	"github.com/axfinn/todoIng/backend-go/internal/convert"
	"github.com/axfinn/todoIng/backend-go/internal/models"
	pb "github.com/axfinn/todoIng/backend-go/pkg/api/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthService gRPC 认证服务实现
type AuthService struct {
	pb.UnimplementedAuthServiceServer
	db *mongo.Database
}

// NewAuthService 创建新的认证服务
func NewAuthService(db *mongo.Database) *AuthService {
	return &AuthService{
		db: db,
	}
}

// Register 实现用户注册
func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// 验证请求参数
	if req.Username == "" || req.Email == "" || len(req.Password) < 6 {
		return nil, status.Error(codes.InvalidArgument, "Invalid username, email, or password")
	}

	// 检查用户是否已存在
	users := s.db.Collection("users")
	err := users.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"email": req.Email},
			{"username": req.Username},
		},
	}).Err()

	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "User already exists")
	}

	if err != mongo.ErrNoDocuments {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Database error: %v", err))
	}

	// 这里应该添加密码哈希逻辑
	// 为了示例，我们简化处理
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // 实际应用中需要哈希
	}

	// 插入用户到数据库
	result, err := users.InsertOne(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to create user: %v", err))
	}

	// 设置用户ID
	user.ID = result.InsertedID.(string)

	// 转换为 protobuf 响应
	pbUser := convert.UserToProto(user)

	return &pb.RegisterResponse{
		Response: &pb.Response{
			Code:    201,
			Message: "User created successfully",
		},
		User: pbUser,
	}, nil
}

// Login 实现用户登录
func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "Email and password are required")
	}

	// 查找用户
	users := s.db.Collection("users")
	var user models.User
	err := users.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.Unauthenticated, "Invalid email or password")
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("Database error: %v", err))
	}

	// 这里应该验证密码哈希
	// 为了示例，我们简化处理
	if user.Password != req.Password {
		return nil, status.Error(codes.Unauthenticated, "Invalid email or password")
	}

	// 生成 JWT token（这里简化处理）
	token := "jwt_token_here"

	// 转换为 protobuf 响应
	pbUser := convert.UserToProto(&user)

	return &pb.LoginResponse{
		Response: &pb.Response{
			Code:    200,
			Message: "Login successful",
		},
		Token: token,
		User:  pbUser,
	}, nil
}

// EmailCodeLogin 实现邮箱验证码登录
func (s *AuthService) EmailCodeLogin(ctx context.Context, req *pb.EmailCodeLoginRequest) (*pb.LoginResponse, error) {
	// 实现邮箱验证码登录逻辑
	return nil, status.Error(codes.Unimplemented, "Email code login not implemented yet")
}

// SendLoginEmailCode 发送登录邮箱验证码
func (s *AuthService) SendLoginEmailCode(ctx context.Context, req *pb.SendLoginEmailCodeRequest) (*pb.Response, error) {
	// 实现发送邮箱验证码逻辑
	return nil, status.Error(codes.Unimplemented, "Send login email code not implemented yet")
}

// VerifyToken 验证 JWT token
func (s *AuthService) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	// 实现 token 验证逻辑
	return nil, status.Error(codes.Unimplemented, "Token verification not implemented yet")
}

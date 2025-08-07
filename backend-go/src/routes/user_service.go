package routes

import (
	"context"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "todoing-backend/proto/gen/proto"
	"todoing-backend/src/config"
	"todoing-backend/src/models"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
}

// Register 实现用户注册
func (s *UserServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// 检查用户是否已存在
	var existingUser models.User
	err := config.DB.Collection("users").FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"username": req.Username},
			{"email": req.Email},
		},
	}).Decode(&existingUser)

	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "User already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Server error")
	}

	// Create user
	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	result, err := config.DB.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Server error")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  result.InsertedID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // 1 week
	})

	tokenString, err := token.SignedString([]byte("your-jwt-secret")) // 应该从环境变量获取
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Server error")
	}

	return &pb.RegisterResponse{
		Token: tokenString,
		User: &pb.User{
			Id:       result.InsertedID.(primitive.ObjectID).Hex(),
			Username: user.Username,
			Email:    user.Email,
		},
		Message: "User registered successfully",
	}, nil
}

// Login 实现用户登录
func (s *UserServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Find user
	var user models.User
	err := config.DB.Collection("users").FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid credentials")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid credentials")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // 1 week
	})

	tokenString, err := token.SignedString([]byte("your-jwt-secret")) // 应该从环境变量获取
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Server error")
	}

	return &pb.LoginResponse{
		Token: tokenString,
		User: &pb.User{
			Id:       user.ID.Hex(),
			Username: user.Username,
			Email:    user.Email,
		},
		Message: "User logged in successfully",
	}, nil
}

// GetUser 实现获取用户信息
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// 这里应该从 token 中解析用户信息
	// 为简化实现，我们暂时返回一个示例用户
	return &pb.GetUserResponse{
		User: &pb.User{
			Id:       "user-id",
			Username: "example_user",
			Email:    "user@example.com",
		},
		Message: "User information retrieved successfully",
	}, nil
}

// SendEmailCode 实现发送注册邮箱验证码
func (s *UserServiceServer) SendEmailCode(ctx context.Context, req *pb.SendEmailCodeRequest) (*pb.SendEmailCodeResponse, error) {
	// 检查用户是否已存在
	var existingUser models.User
	err := config.DB.Collection("users").FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)

	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "User already exists")
	}

	// 实际项目中应该发送邮件验证码，这里仅作演示
	return &pb.SendEmailCodeResponse{
		Message: "Registration email code sent successfully",
		Id:      "email-code-id",
	}, nil
}

// SendLoginEmailCode 实现发送登录邮箱验证码
func (s *UserServiceServer) SendLoginEmailCode(ctx context.Context, req *pb.SendLoginEmailCodeRequest) (*pb.SendLoginEmailCodeResponse, error) {
	// 检查用户是否存在
	var existingUser models.User
	err := config.DB.Collection("users").FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	// 实际项目中应该发送邮件验证码，这里仅作演示
	return &pb.SendLoginEmailCodeResponse{
		Message: "Login email code sent successfully",
		Id:      "login-email-code-id",
	}, nil
}

// GenerateCaptcha 实现生成验证码
func (s *UserServiceServer) GenerateCaptcha(ctx context.Context, req *pb.Empty) (*pb.GenerateCaptchaResponse, error) {
	// 生成4位数字验证码
	captcha := int32(rand.Intn(9000) + 1000)

	return &pb.GenerateCaptchaResponse{
		Captcha: captcha,
		Message: "Captcha generated successfully",
	}, nil
}

// VerifyCaptcha 实现验证验证码
func (s *UserServiceServer) VerifyCaptcha(ctx context.Context, req *pb.CaptchaRequest) (*pb.CaptchaResponse, error) {
	// 实际项目中应该验证验证码，这里仅作演示
	return &pb.CaptchaResponse{
		Message: "Captcha verified successfully",
	}, nil
}


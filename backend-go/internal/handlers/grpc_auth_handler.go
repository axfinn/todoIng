package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/axfinn/todoIng/backend-go/internal/api"
	pb "github.com/axfinn/todoIng/backend-go/pkg/api/v1"
)

// GrpcAuthHandler HTTP 到 gRPC 的认证处理器
type GrpcAuthHandler struct {
	authClient pb.AuthServiceClient
}

// NewGrpcAuthHandler 创建新的 gRPC 认证处理器
func NewGrpcAuthHandler(authClient pb.AuthServiceClient) *GrpcAuthHandler {
	return &GrpcAuthHandler{
		authClient: authClient,
	}
}

// RegisterRequest HTTP 注册请求结构
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register 处理用户注册请求（HTTP -> gRPC）
// @Summary 用户注册 (gRPC)
// @Description 通过 gRPC 服务处理用户注册
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册信息"
// @Success 201 {object} map[string]interface{} "注册成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 409 {object} map[string]string "用户已存在"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/v2/auth/register [post]
func (h *GrpcAuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.JSON(w, http.StatusBadRequest, map[string]string{"msg": "Invalid request body"})
		return
	}

	// 创建 gRPC 请求
	grpcReq := &pb.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	// 调用 gRPC 服务
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.authClient.Register(ctx, grpcReq)
	if err != nil {
		// 处理 gRPC 错误
		api.JSON(w, http.StatusInternalServerError, map[string]string{"msg": "Registration failed"})
		return
	}

	// 转换响应
	response := map[string]interface{}{
		"user": map[string]interface{}{
			"id":       resp.User.Id,
			"username": resp.User.Username,
			"email":    resp.User.Email,
		},
	}

	api.JSON(w, int(resp.Response.Code), response)
}

// LoginRequest HTTP 登录请求结构
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login 处理用户登录请求（HTTP -> gRPC）
// @Summary 用户登录 (gRPC)
// @Description 通过 gRPC 服务处理用户登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} map[string]interface{} "登录成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 401 {object} map[string]string "认证失败"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/v2/auth/login [post]
func (h *GrpcAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.JSON(w, http.StatusBadRequest, map[string]string{"msg": "Invalid request body"})
		return
	}

	// 创建 gRPC 请求
	grpcReq := &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	// 调用 gRPC 服务
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.authClient.Login(ctx, grpcReq)
	if err != nil {
		// 处理 gRPC 错误
		api.JSON(w, http.StatusUnauthorized, map[string]string{"msg": "Login failed"})
		return
	}

	// 转换响应
	response := map[string]interface{}{
		"token": resp.Token,
		"user": map[string]interface{}{
			"id":       resp.User.Id,
			"username": resp.User.Username,
			"email":    resp.User.Email,
		},
	}

	api.JSON(w, int(resp.Response.Code), response)
}

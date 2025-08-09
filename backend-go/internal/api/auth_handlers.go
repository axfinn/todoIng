package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/axfinn/todoIng/backend-go/internal/auth"
	"github.com/axfinn/todoIng/backend-go/internal/captcha"
	"github.com/axfinn/todoIng/backend-go/internal/email"
	"github.com/axfinn/todoIng/backend-go/internal/models"
	"github.com/axfinn/todoIng/backend-go/internal/observability"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthDeps struct {
	DB         *mongo.Database
	EmailCodes *email.Store
}

// RegisterRequest 用户注册请求结构
type registerRequest struct {
	Username    string `json:"username" example:"johndoe" validate:"required"`
	Email       string `json:"email" example:"john@example.com" validate:"required,email"`
	Password    string `json:"password" example:"password123" validate:"required,min=6"`
	EmailCode   string `json:"emailCode" example:"123456"`
	EmailCodeId string `json:"emailCodeId" example:"code-id-123"`
}

// LoginRequest 用户登录请求结构
type loginRequest struct {
	Email       string `json:"email" example:"john@example.com" validate:"required,email"`
	Password    string `json:"password" example:"password123" validate:"required"`
	EmailCode   string `json:"emailCode" example:"123456"`
	EmailCodeId string `json:"emailCodeId" example:"code-id-123"`
}

// UserResponse 用户响应结构
type UserResponse struct {
	ID        string    `json:"id" example:"60d5ecb74eb3b8001f8b4567"`
	Username  string    `json:"username" example:"johndoe"`
	Email     string    `json:"email" example:"john@example.com"`
	CreatedAt time.Time `json:"createdAt" example:"2023-12-01T10:00:00Z"`
	UpdatedAt time.Time `json:"updatedAt" example:"2023-12-01T10:00:00Z"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserResponse `json:"user"`
}

// Register 用户注册
// @Summary 用户注册
// @Description 注册新用户账户
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body registerRequest true "注册信息"
// @Success 201 {object} map[string]interface{} "注册成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 409 {object} map[string]string "用户已存在"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/auth/register [post]

func (d *AuthDeps) Register(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("DISABLE_REGISTRATION") == "true" {
		JSON(w, http.StatusForbidden, map[string]string{"msg": "Registration is disabled"})
		return
	}
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSON(w, 400, map[string]string{"msg": "Invalid body"})
		return
	}
	if req.Username == "" || req.Email == "" || len(req.Password) < 6 {
		JSON(w, 400, map[string]string{"msg": "Invalid fields"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	users := d.DB.Collection("users")
	// email verify
	if os.Getenv("ENABLE_EMAIL_VERIFICATION") == "true" {
		if err := d.EmailCodes.Verify(req.EmailCodeId, strings.ToLower(req.Email), req.EmailCode); err != nil {
			JSON(w, 400, map[string]string{"msg": err.Error()})
			return
		}
	}
	// uniqueness
	if err := users.FindOne(ctx, bson.M{"$or": []bson.M{{"email": req.Email}, {"username": req.Username}}}).Err(); err == nil {
		JSON(w, 409, map[string]string{"msg": "User already exists"})
		return
	}
	pwHash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	user := models.User{Username: req.Username, Email: strings.ToLower(req.Email), Password: string(pwHash), CreatedAt: time.Now()}
	res, err := users.InsertOne(ctx, user)
	if err != nil {
		JSON(w, 500, map[string]string{"msg": "DB error"})
		return
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()
	token, _ := auth.Generate(id, time.Hour)

	userResponse := UserResponse{
		ID:        id,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.CreatedAt, // 使用CreatedAt作为初始UpdatedAt
	}

	JSON(w, 201, map[string]interface{}{
		"token": token,
		"user":  userResponse,
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户通过邮箱和密码登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body loginRequest true "登录信息"
// @Success 200 {object} LoginResponse "登录成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 401 {object} map[string]string "认证失败"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/auth/login [post]
func (d *AuthDeps) Login(w http.ResponseWriter, r *http.Request) {
	observability.CtxLog(r.Context(), "Login attempt from IP: %s", r.RemoteAddr)

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		observability.LogWarn("Invalid login request body from IP: %s", r.RemoteAddr)
		JSON(w, 400, map[string]string{"msg": "Invalid body"})
		return
	}
	if req.Email == "" {
		observability.LogWarn("Login attempt without email from IP: %s", r.RemoteAddr)
		JSON(w, 400, map[string]string{"msg": "Email required"})
		return
	}

	normalizedEmail := strings.ToLower(req.Email)
	observability.CtxLog(r.Context(), "Login attempt for email: %s", normalizedEmail)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	users := d.DB.Collection("users")
	var user models.User
	err := users.FindOne(ctx, bson.M{"email": normalizedEmail}).Decode(&user)
	if err != nil {
		observability.LogWarn("Login failed - user not found for email: %s", normalizedEmail)
		// 如果是邮箱验证码登录但用户不存在，返回用户不存在的错误
		if req.EmailCode != "" && req.EmailCodeId != "" && os.Getenv("ENABLE_EMAIL_VERIFICATION") == "true" {
			JSON(w, 401, map[string]string{"msg": "User does not exist"})
			return
		}
		JSON(w, 401, map[string]string{"msg": "Invalid credentials"})
		return
	}
	// email code login
	if req.EmailCode != "" && req.EmailCodeId != "" && os.Getenv("ENABLE_EMAIL_VERIFICATION") == "true" {
		observability.CtxLog(r.Context(), "Attempting email code verification for user: %s", normalizedEmail)
		// 使用小写邮箱地址进行验证，确保与存储时一致
		if err := d.EmailCodes.Verify(req.EmailCodeId, normalizedEmail, req.EmailCode); err != nil {
			observability.LogWarn("Email code verification failed for user: %s, error: %v", normalizedEmail, err)
			// 统一验证码相关的错误信息，避免暴露具体的验证码错误
			JSON(w, 400, map[string]string{"msg": "Invalid verification code"})
			return
		}
		observability.LogInfo("Email code verification successful for user: %s", normalizedEmail)
	} else {
		observability.CtxLog(r.Context(), "Attempting password verification for user: %s", normalizedEmail)
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
			observability.LogWarn("Password verification failed for user: %s", normalizedEmail)
			JSON(w, 401, map[string]string{"msg": "Invalid credentials"})
			return
		}
		observability.LogInfo("Password verification successful for user: %s", normalizedEmail)
	}

	token, _ := auth.Generate(user.ID, time.Hour)
	observability.LogInfo("Login successful for user: %s (ID: %s)", normalizedEmail, user.ID)
	JSON(w, 200, map[string]string{"token": token})
}

func (d *AuthDeps) Me(w http.ResponseWriter, r *http.Request) {
	uid := GetUserID(r)
	if uid == "" {
		JSON(w, 401, map[string]string{"msg": "Unauthorized"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	users := d.DB.Collection("users")
	var user models.User
	objID, _ := primitive.ObjectIDFromHex(uid)
	err := users.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		JSON(w, 500, map[string]string{"msg": "Not found"})
		return
	}
	user.Password = ""
	JSON(w, 200, user)
}

func (d *AuthDeps) SendRegisterEmailCode(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("ENABLE_EMAIL_VERIFICATION") != "true" {
		JSON(w, 400, map[string]string{"msg": "Email verification disabled"})
		return
	}
	var body struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		JSON(w, 400, map[string]string{"msg": "Invalid body"})
		return
	}
	if body.Email == "" {
		JSON(w, 400, map[string]string{"msg": "Email required"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	users := d.DB.Collection("users")
	normalizedEmail := strings.ToLower(body.Email)

	// 注册时检查用户是否已存在
	if err := users.FindOne(ctx, bson.M{"email": normalizedEmail}).Err(); err == nil {
		JSON(w, 400, map[string]string{"msg": "User already exists"})
		return
	}

	id, code := d.EmailCodes.Generate(normalizedEmail, 6)
	_ = email.Send(body.Email, code) // 发送邮件使用原始邮箱格式
	JSON(w, 200, map[string]string{"id": id, "msg": "Verification code sent"})
}

func (d *AuthDeps) SendLoginEmailCode(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("ENABLE_EMAIL_VERIFICATION") != "true" {
		JSON(w, 400, map[string]string{"msg": "Email verification disabled"})
		return
	}
	var body struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		JSON(w, 400, map[string]string{"msg": "Invalid body"})
		return
	}
	if body.Email == "" {
		JSON(w, 400, map[string]string{"msg": "Email required"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	users := d.DB.Collection("users")
	normalizedEmail := strings.ToLower(body.Email)

	// 登录时检查用户是否存在
	var user models.User
	if err := users.FindOne(ctx, bson.M{"email": normalizedEmail}).Decode(&user); err != nil {
		JSON(w, 400, map[string]string{"msg": "User does not exist"})
		return
	}

	id, code := d.EmailCodes.Generate(normalizedEmail, 6)
	_ = email.Send(body.Email, code) // 发送邮件使用原始邮箱格式
	JSON(w, 200, map[string]string{"id": id, "msg": "Login verification code sent"})
}

func SetupAuthRoutes(r *mux.Router, deps *AuthDeps) {
	r.HandleFunc("/api/auth/register", deps.Register).Methods(http.MethodPost)
	r.HandleFunc("/api/auth/login", deps.Login).Methods(http.MethodPost)
	r.Handle("/api/auth/me", Auth(http.HandlerFunc(deps.Me))).Methods(http.MethodGet)
	r.HandleFunc("/api/auth/send-email-code", deps.SendRegisterEmailCode).Methods(http.MethodPost)
	// 登录邮箱验证码使用专门的函数，检查用户是否存在
	r.HandleFunc("/api/auth/send-login-email-code", deps.SendLoginEmailCode).Methods(http.MethodPost)
}

// 包装函数，用于兼容测试代码

func RegisterHandler(db *mongo.Database, emailStore *email.Store, captchaStore *captcha.Store) http.HandlerFunc {
	deps := &AuthDeps{DB: db, EmailCodes: emailStore}
	return deps.Register
}

func LoginHandler(db *mongo.Database) http.HandlerFunc {
	deps := &AuthDeps{DB: db}
	return deps.Login
}

func VerifyTokenHandler(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			JSON(w, 401, map[string]string{"msg": "Unauthorized"})
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			JSON(w, 401, map[string]string{"msg": "Unauthorized"})
			return
		}
		_, err := auth.Parse(parts[1])
		if err != nil {
			JSON(w, 401, map[string]string{"msg": "Unauthorized"})
			return
		}
		JSON(w, 200, map[string]string{"msg": "Token valid"})
	}
}

func EmailCodeLoginHandler(db *mongo.Database, emailStore *email.Store) http.HandlerFunc {
	deps := &AuthDeps{DB: db, EmailCodes: emailStore}
	return deps.Login // 邮箱验证码登录使用相同的登录逻辑
}

func SendLoginEmailCodeHandler(db *mongo.Database, emailStore *email.Store) http.HandlerFunc {
	deps := &AuthDeps{DB: db, EmailCodes: emailStore}
	return deps.SendLoginEmailCode
}

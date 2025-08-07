package routes

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"todoing-backend/src/config"
)

// GenerateCaptchaResponse 生成验证码响应
// swagger:model
type GenerateCaptchaResponse struct {
	// 验证码图片 (base64 encoded)
	Image string `json:"image"`
	// 验证码ID
	Id string `json:"id"`
	// 消息
	Message string `json:"message"`
}

// SendEmailCodeRequest 发送邮箱验证码请求
// swagger:model
type SendEmailCodeRequest struct {
	// 邮箱地址
	Email string `json:"email"`
}

// SendEmailCodeResponse 发送邮箱验证码响应
// swagger:model
type SendEmailCodeResponse struct {
	// 消息
	Message string `json:"message"`
	// ID
	Id string `json:"id"`
}

// SendLoginEmailCodeRequest 发送登录邮箱验证码请求
// swagger:model
type SendLoginEmailCodeRequest struct {
	// 邮箱地址
	Email string `json:"email"`
}

// SendLoginEmailCodeResponse 发送登录邮箱验证码响应
// swagger:model
type SendLoginEmailCodeResponse struct {
	// 消息
	Message string `json:"message"`
	// ID
	Id string `json:"id"`
}

// RegisterRequest 用户注册请求
// swagger:model
type RegisterRequest struct {
	// 用户名
	Username string `json:"username"`
	// 邮箱地址
	Email string `json:"email"`
	// 密码
	Password string `json:"password"`
}

// RegisterResponse 用户注册响应
// swagger:model
type RegisterResponse struct {
	// JWT token
	Token string `json:"token"`
	// 用户信息
	User User `json:"user"`
	// 消息
	Message string `json:"message"`
}

// LoginRequest 用户登录请求
// swagger:model
type LoginRequest struct {
	// 邮箱地址
	Email string `json:"email"`
	// 密码
	Password string `json:"password"`
	// 验证码ID
	CaptchaId string `json:"captchaId"`
	// 验证码值
	Captcha string `json:"captcha"`
}

// LoginResponse 用户登录响应
// swagger:model
type LoginResponse struct {
	// JWT token
	Token string `json:"token"`
	// 用户信息
	User User `json:"user"`
	// 消息
	Message string `json:"message"`
}

// GetUserResponse 获取用户信息响应
// swagger:model
type GetUserResponse struct {
	// 用户信息
	User User `json:"user"`
	// 消息
	Message string `json:"message"`
}

// User 用户信息
// swagger:model
type User struct {
	// 用户ID
	Id string `json:"id"`
	// 用户名
	Username string `json:"username"`
	// 邮箱地址
	Email string `json:"email"`
}

// ErrorResponse 通用错误响应
// swagger:model
type ErrorResponse struct {
	// 错误信息
	Msg string `json:"msg"`
}

// VerifyCaptchaRequest 验证码验证请求
// swagger:model
type VerifyCaptchaRequest struct {
	// 验证码
	Captcha string `json:"captcha"`
	// 验证码ID
	CaptchaId string `json:"captchaId"`
}

// VerifyCaptchaResponse 验证码验证响应
// swagger:model
type VerifyCaptchaResponse struct {
	// 消息
	Message string `json:"message"`
}

// RegisterAuthRoutes 注册认证相关路由
func RegisterAuthRoutes(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", register)
		auth.POST("/login", login)
		auth.GET("/me", getUser)
		auth.GET("/captcha", getCaptcha)
		auth.POST("/verify-captcha", verifyCaptcha)
		auth.POST("/send-email-code", sendEmailCode)
		auth.POST("/send-login-email-code", sendLoginEmailCode)
	}
}

// 定义验证码存储结构
type captchaData struct {
	Text      string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// 定义一个简单的内存存储来保存验证码
var captchaStore = make(map[string]*captchaData)

// GenerateCaptcha 生成验证码
// @Summary 生成验证码
// @Description 生成图形验证码用于注册/登录验证
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} GenerateCaptchaResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/captcha [get]
func getCaptcha(c *gin.Context) {
	rand.Seed(time.Now().UnixNano())

	// 生成4位数字验证码
	captcha := rand.Intn(9000) + 1000
	
	// 创建图片
	width := 150
	height := 50
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 设置背景色
	bgColor := color.RGBA{255, 255, 255, 255} // 白色背景
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, bgColor)
		}
	}

	// 添加干扰线
	for i := 0; i < 5; i++ {
		x1 := rand.Intn(width)
		y1 := rand.Intn(height)
		x2 := rand.Intn(width)
		y2 := rand.Intn(height)
		lineColor := color.RGBA{uint8(rand.Intn(200)), uint8(rand.Intn(200)), uint8(rand.Intn(200)), 255}
		
		// Bresenham直线算法
		dx := abs(x2 - x1)
		dy := abs(y2 - y1)
		var sx, sy int
		if x1 < x2 {
			sx = 1
		} else {
			sx = -1
		}
		if y1 < y2 {
			sy = 1
		} else {
			sy = -1
		}
		err := dx - dy
		
		x, y := x1, y1
		for {
			if x >= 0 && x < width && y >= 0 && y < height {
				img.Set(x, y, lineColor)
			}
			if x == x2 && y == y2 {
				break
			}
			e2 := 2 * err
			if e2 > -dy {
				err -= dy
				x += sx
			}
			if e2 < dx {
				err += dx
				y += sy
			}
		}
	}

	// 添加干扰点
	for i := 0; i < 30; i++ {
		x := rand.Intn(width)
		y := rand.Intn(height)
		img.Set(x, y, color.RGBA{uint8(rand.Intn(200)), uint8(rand.Intn(200)), uint8(rand.Intn(200)), 255})
	}

	// 将验证码转换为字符串并绘制
	captchaStr := fmt.Sprintf("%04d", captcha)
	
	// 绘制数字，每个数字用不同的字体大小和位置
	for i, digit := range captchaStr {
		// 计算数字位置
		x := 20 + i*30 + rand.Intn(5) - 2
		y := 25 + rand.Intn(10) - 5
		
		// 绘制数字（简单但清晰的方式）
		drawSimpleDigit(img, digit, x, y)
	}

	// 将图片编码为base64
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &GenerateCaptchaResponse{
			Message: "Failed to generate captcha image",
		})
		return
	}

	// 生成验证码ID
	captchaID := fmt.Sprintf("%d", time.Now().UnixNano())

	// 存储验证码到内存（设置5分钟过期时间）
	captchaStore[captchaID] = &captchaData{
		Text:      captchaStr,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	c.JSON(http.StatusOK, &GenerateCaptchaResponse{
		Image:   "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()),
		Id:      captchaID,
		Message: "Captcha generated successfully",
	})
}

// abs 计算绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// drawSimpleDigit 简单但清晰地绘制数字
func drawSimpleDigit(img *image.RGBA, digit rune, x, y int) {
	// 定义数字的模式（7段显示）
	patterns := map[rune][7]bool{
		'0': {true, true, true, true, true, true, false},
		'1': {false, true, true, false, false, false, false},
		'2': {true, true, false, true, true, false, true},
		'3': {true, true, true, true, false, false, true},
		'4': {false, true, true, false, false, true, true},
		'5': {true, false, true, true, false, true, true},
		'6': {true, false, true, true, true, true, true},
		'7': {true, true, true, false, false, false, false},
		'8': {true, true, true, true, true, true, true},
		'9': {true, true, true, true, false, true, true},
	}
	
	// 定义7段的位置
	segments := [7][4]int{
		// a段 (顶部)
		{x + 5, y + 2, x + 20, y + 2},
		// b段 (右上)
		{x + 23, y + 5, x + 23, y + 15},
		// c段 (右下)
		{x + 23, y + 18, x + 23, y + 28},
		// d段 (底部)
		{x + 5, y + 30, x + 20, y + 30},
		// e段 (左下)
		{x + 2, y + 18, x + 2, y + 28},
		// f段 (左上)
		{x + 2, y + 5, x + 2, y + 15},
		// g段 (中间)
		{x + 5, y + 16, x + 20, y + 16},
	}
	
	// 绘制每个激活的段
	color := color.RGBA{0, 0, 0, 255} // 黑色
	pattern := patterns[digit]
	
	for i, active := range pattern {
		if active {
			drawSegment(img, segments[i], color)
		}
	}
}

// drawSegment 绘制一个段
func drawSegment(img *image.RGBA, segment [4]int, color color.RGBA) {
	x1, y1, x2, y2 := segment[0], segment[1], segment[2], segment[3]
	
	// 如果是水平线
	if y1 == y2 {
		for x := x1; x <= x2; x++ {
			// 绘制加粗效果
			if x >= 0 && x < img.Bounds().Max.X && y1 >= 0 && y1 < img.Bounds().Max.Y {
				img.Set(x, y1, color)
			}
			if x >= 0 && x < img.Bounds().Max.X && y1-1 >= 0 && y1-1 < img.Bounds().Max.Y {
				img.Set(x, y1-1, color)
			}
			if x >= 0 && x < img.Bounds().Max.X && y1+1 >= 0 && y1+1 < img.Bounds().Max.Y {
				img.Set(x, y1+1, color)
			}
		}
	} else { // 如果是垂直线
		for y := y1; y <= y2; y++ {
			// 绘制加粗效果
			if x1 >= 0 && x1 < img.Bounds().Max.X && y >= 0 && y < img.Bounds().Max.Y {
				img.Set(x1, y, color)
			}
			if x1-1 >= 0 && x1-1 < img.Bounds().Max.X && y >= 0 && y < img.Bounds().Max.Y {
				img.Set(x1-1, y, color)
			}
			if x1+1 >= 0 && x1+1 < img.Bounds().Max.X && y >= 0 && y < img.Bounds().Max.Y {
				img.Set(x1+1, y, color)
			}
		}
	}
}

// VerifyCaptcha 验证码验证
// @Summary 验证码验证
// @Description 验证图形验证码是否正确
// @Tags auth
// @Accept json
// @Produce json
// @Param request body VerifyCaptchaRequest true "验证码验证请求"
// @Success 200 {object} VerifyCaptchaResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/verify-captcha [post]
func verifyCaptcha(c *gin.Context) {
	var req VerifyCaptchaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Invalid request"})
		return
	}

	// 检查是否启用了验证码功能
	if os.Getenv("ENABLE_CAPTCHA") != "true" {
		c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Captcha is not enabled"})
		return
	}

	// 验证验证码逻辑
	if req.Captcha == "" {
		c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Captcha is required"})
		return
	}
	
	// 检查是否提供了验证码ID
	if req.CaptchaId == "" {
		c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Captcha ID is required"})
		return
	}
	
	// 验证验证码是否正确
	storedCaptcha, exists := captchaStore[req.CaptchaId]
	if !exists {
		c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Invalid or expired captcha"})
		return
	}
	
	// 检查验证码是否过期
	if time.Now().After(storedCaptcha.ExpiresAt) {
		delete(captchaStore, req.CaptchaId)
		c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Captcha has expired"})
		return
	}
	
	// 比较验证码（不区分大小写）
	if storedCaptcha.Text != req.Captcha {
		c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Invalid captcha"})
		return
	}
	
	// 验证成功后删除验证码，防止重复使用
	delete(captchaStore, req.CaptchaId)
	
	c.JSON(http.StatusOK, &VerifyCaptchaResponse{
		Message: "Captcha verified successfully",
	})
}

// SendEmailCode 发送注册邮箱验证码
// @Summary 发送注册邮箱验证码
// @Description 发送邮箱验证码用于用户注册验证
// @Tags auth
// @Accept json
// @Produce json
// @Param request body SendEmailCodeRequest true "发送邮箱验证码请求"
// @Success 200 {object} SendEmailCodeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/send-email-code [post]
func sendEmailCode(c *gin.Context) {
	var req SendEmailCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: err.Error()})
		return
	}

	// 检查用户是否已存在
	var existingUser bson.M
	err := config.DB.Collection("users").FindOne(c, bson.M{"email": req.Email}).Decode(&existingUser)

	if err == nil {
		c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "User already exists"})
		return
	}

	// TODO: 实现真实的邮件发送功能
	// 生成验证码ID (模拟)
	codeId := fmt.Sprintf("email_%d", time.Now().UnixNano())
	
	c.JSON(http.StatusOK, &SendEmailCodeResponse{
		Message: "Registration email code sent successfully",
		Id: codeId,
	})
}

// SendLoginEmailCode 发送登录邮箱验证码
// @Summary 发送登录邮箱验证码
// @Description 发送邮箱验证码用于用户登录验证
// @Tags auth
// @Accept json
// @Produce json
// @Param request body SendLoginEmailCodeRequest true "发送登录邮箱验证码请求"
// @Success 200 {object} SendLoginEmailCodeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/send-login-email-code [post]
func sendLoginEmailCode(c *gin.Context) {
	var req SendLoginEmailCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: err.Error()})
		return
	}

	// 检查用户是否存在
	var existingUser bson.M
	err := config.DB.Collection("users").FindOne(c, bson.M{"email": req.Email}).Decode(&existingUser)

	if err != nil {
		c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "User not found"})
		return
	}

	// TODO: 实现真实的邮件发送功能
	// 生成验证码ID (模拟)
	codeId := fmt.Sprintf("login_email_%d", time.Now().UnixNano())
	
	c.JSON(http.StatusOK, &SendLoginEmailCodeResponse{
		Message: "Login email code sent successfully",
		Id: codeId,
	})
}

// Register 用户注册
// @Summary 用户注册
// @Description 注册新用户
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "用户注册请求"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: err.Error()})
		return
	}

	// 检查是否禁用注册功能
	if os.Getenv("DISABLE_REGISTRATION") == "true" {
		c.JSON(http.StatusForbidden, &ErrorResponse{Msg: "Registration is disabled"})
		return
	}

	// Check if user already exists
	var existingUser bson.M

	// Check username
	err := config.DB.Collection("users").FindOne(c, bson.M{"username": req.Username}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Username already exists"})
		return
	}

	// Check email
	err = config.DB.Collection("users").FindOne(c, bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Email already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorResponse{Msg: "Server error"})
		return
	}

	// Create user document
	userID := primitive.NewObjectID()
	userDoc := bson.M{
		"_id":        userID,
		"username":   req.Username,
		"email":      req.Email,
		"password":   string(hashedPassword),
		"created_at": time.Now(),
	}

	_, err = config.DB.Collection("users").InsertOne(c, userDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorResponse{Msg: "Server error"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    userID.Hex(),
		"email": req.Email,
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorResponse{Msg: "Server error"})
		return
	}

	c.JSON(http.StatusCreated, &RegisterResponse{
		Token: tokenString,
		User: User{
			Id:       userID.Hex(),
			Username: req.Username,
			Email:    req.Email,
		},
		Message: "User registered successfully",
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取访问令牌
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "用户登录请求"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: err.Error()})
		return
	}

	// 检查是否启用了验证码功能
	if os.Getenv("ENABLE_CAPTCHA") == "true" {
		// 验证验证码逻辑
		if req.Captcha == "" {
			c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Captcha is required"})
			return
		}
		
		// 检查是否提供了验证码ID
		if req.CaptchaId == "" {
			c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Captcha ID is required"})
			return
		}
		
		// 验证验证码是否正确
		storedCaptcha, exists := captchaStore[req.CaptchaId]
		if !exists {
			c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Invalid or expired captcha"})
			return
		}
		
		// 检查验证码是否过期
		if time.Now().After(storedCaptcha.ExpiresAt) {
			delete(captchaStore, req.CaptchaId)
			c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Captcha has expired"})
			return
		}
		
		// 比较验证码（不区分大小写）
		if storedCaptcha.Text != req.Captcha {
			c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Invalid captcha"})
			return
		}
		
		// 验证成功后删除验证码，防止重复使用
		delete(captchaStore, req.CaptchaId)
	}

	// 验证密码
	// 注意：由于User对象中不包含密码字段，我们需要从数据库中获取完整信息
	var userDoc bson.M
	err := config.DB.Collection("users").FindOne(c, bson.M{"email": req.Email}).Decode(&userDoc)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &ErrorResponse{Msg: "Invalid credentials"})
		return
	}

	hashedPassword, ok := userDoc["password"].(string)
	if !ok || bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, &ErrorResponse{Msg: "Invalid credentials"})
		return
	}

	// 从用户文档中获取用户ID
	userID := userDoc["_id"].(primitive.ObjectID).Hex()

	// 更新最后登录时间
	_, err = config.DB.Collection("users").UpdateOne(
		c,
		bson.M{"email": req.Email},
		bson.M{"$set": bson.M{"last_login": time.Now()}},
	)
	if err != nil {
		// 这里我们只记录错误，但不中断登录流程
		fmt.Printf("Failed to update last login time: %v\n", err)
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // 1 week
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorResponse{Msg: "Server error"})
		return
	}

	c.JSON(http.StatusOK, &LoginResponse{
		Token: tokenString,
		User: User{
			Id:       userDoc["_id"].(primitive.ObjectID).Hex(),
			Username: userDoc["username"].(string),
			Email:    userDoc["email"].(string),
		},
		Message: "User logged in successfully",
	})
}

// GetUser 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取当前认证用户的信息
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} GetUserResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/me [get]
func getUser(c *gin.Context) {
	// Get user from context (set by auth middleware)
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, &ErrorResponse{Msg: "No token, authorization denied"})
		return
	}

	claims, ok := userClaims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, &ErrorResponse{Msg: "Invalid token claims"})
		return
	}

	// 从token中提取用户ID
	userID, err := primitive.ObjectIDFromHex(claims["id"].(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, &ErrorResponse{Msg: "Invalid token"})
		return
	}

	// 获取用户信息
	var userDoc bson.M
	err = config.DB.Collection("users").FindOne(c, bson.M{"_id": userID}).Decode(&userDoc)
	if err != nil {
		c.JSON(http.StatusNotFound, &ErrorResponse{Msg: "User not found"})
		return
	}

	// 构造返回的用户信息
	user := User{
		Id:       userDoc["_id"].(primitive.ObjectID).Hex(),
		Username: userDoc["username"].(string),
		Email:    userDoc["email"].(string),
	}

	c.JSON(http.StatusOK, &GetUserResponse{
		User:    user,
		Message: "User profile retrieved successfully",
	})
}

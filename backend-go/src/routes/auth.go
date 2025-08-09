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
	"strings"
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
	Msg string `json:"msg"`
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
	Msg string `json:"msg"`
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
	// 邮箱验证码
	EmailCode string `json:"emailCode"`
	// 邮箱验证码ID
	EmailCodeId string `json:"emailCodeId"`
}

// RegisterResponse 用户注册响应
// swagger:model
type RegisterResponse struct {
	// JWT token
	Token string `json:"token"`
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
	// 邮箱验证码
	EmailCode string `json:"emailCode"`
	// 邮箱验证码ID
	EmailCodeId string `json:"emailCodeId"`
}

// LoginResponse 用户登录响应
// swagger:model
type LoginResponse struct {
	// JWT token
	Token string `json:"token"`
}

// GetUserResponse 获取用户信息响应 - 直接返回用户对象
// 注意：backend 直接返回用户对象，不包装在响应中
type GetUserResponse = User

// User 用户信息
// swagger:model
type User struct {
	// 用户ID
	Id string `json:"_id"`
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
	Msg string `json:"msg"`
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

// cleanupExpiredCaptchas 清理过期的验证码
func cleanupExpiredCaptchas() {
	now := time.Now()
	for id, captcha := range captchaStore {
		if now.After(captcha.ExpiresAt) {
			delete(captchaStore, id)
		}
	}
}

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
	// 检查是否启用了验证码功能
	if os.Getenv("ENABLE_CAPTCHA") != "true" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Captcha is not enabled"})
		return
	}

	rand.Seed(time.Now().UnixNano())

	// 生成6位字母数字混合验证码（排除易混淆字符）
	const charset = "ABCDEFGHJKMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789"
	captchaText := make([]byte, 6)
	for i := range captchaText {
		captchaText[i] = charset[rand.Intn(len(charset))]
	}
	captcha := string(captchaText)

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
	captchaStr := captcha

	// 绘制验证码字符
	// 由于Go标准库不支持字体渲染，我们使用简单的方式绘制
	// 实际项目中应该使用 github.com/golang/freetype 等库
	for i, char := range captchaStr {
		// 计算字符位置
		x := 10 + i*22 + rand.Intn(5) - 2
		y := 25 + rand.Intn(10) - 5

		// 使用随机颜色
		charColor := color.RGBA{
			uint8(rand.Intn(100)),
			uint8(rand.Intn(100)),
			uint8(rand.Intn(100)),
			255,
		}

		// 简单地绘制字符占位符（实际应使用字体库）
		// 检查字符是否为数字，如果是则使用7段显示方式绘制
		if char >= '0' && char <= '9' {
			drawSimpleDigit(img, char, x, y)
		} else {
			drawCharPlaceholder(img, x, y, charColor, char)
		}
	}

	// 将图片编码为base64
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Msg: "Failed to generate captcha image",
		})
		return
	}

	// 生成验证码ID（使用更安全的随机字符串）
	b := make([]byte, 16)
	_, err = rand.Read(b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Msg: "Failed to generate captcha ID",
		})
		return
	}
	captchaID := fmt.Sprintf("%x", b)

	// 存储验证码到内存（转换为大写，设置5分钟过期时间）
	captchaStore[captchaID] = &captchaData{
		Text:      strings.ToUpper(captchaStr),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	// 定期清理过期验证码（10%概率触发）
	if rand.Float32() < 0.1 {
		go cleanupExpiredCaptchas()
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

// drawCharPlaceholder 改进版字符绘制
func drawCharPlaceholder(img *image.RGBA, x, y int, color color.RGBA, char rune) {
	// 绘制更清晰的字符框架
	width := 20
	height := 30
	
	// 绘制字符框架的四个边
	for i := 0; i < width; i++ {
		// 顶部边
		px := x - 10 + i
		py := y - 15
		if px >= 0 && px < img.Bounds().Max.X && py >= 0 && py < img.Bounds().Max.Y {
			img.Set(px, py, color)
		}
		
		// 底部边
		px = x - 10 + i
		py = y + 15
		if px >= 0 && px < img.Bounds().Max.X && py >= 0 && py < img.Bounds().Max.Y {
			img.Set(px, py, color)
		}
	}
	
	for i := 0; i < height; i++ {
		// 左侧边
		px := x - 10
		py := y - 15 + i
		if px >= 0 && px < img.Bounds().Max.X && py >= 0 && py < img.Bounds().Max.Y {
			img.Set(px, py, color)
		}
		
		// 右侧边
		px = x + 10
		py = y - 15 + i
		if px >= 0 && px < img.Bounds().Max.X && py >= 0 && py < img.Bounds().Max.Y {
			img.Set(px, py, color)
		}
	}
}

// drawSimpleDigit 改进版数字绘制
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

	// 定义7段的位置（调整位置以适应验证码图片大小）
	segments := [7][4]int{
		// a段 (顶部)
		{x + 3, y + 1, x + 12, y + 1},
		// b段 (右上)
		{x + 14, y + 3, x + 14, y + 10},
		// c段 (右下)
		{x + 14, y + 13, x + 14, y + 20},
		// d段 (底部)
		{x + 3, y + 22, x + 12, y + 22},
		// e段 (左下)
		{x + 1, y + 13, x + 1, y + 20},
		// f段 (左上)
		{x + 1, y + 3, x + 1, y + 10},
		// g段 (中间)
		{x + 3, y + 12, x + 12, y + 12},
	}

	// 绘制每个激活的段
	clr := color.RGBA{0, 0, 0, 255} // 黑色
	pattern := patterns[digit]

	for i, active := range pattern {
		if active {
			drawSegment(img, segments[i], clr)
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
	if storedCaptcha.Text != strings.ToUpper(req.Captcha) {
		c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Invalid captcha"})
		return
	}

	// 验证成功后删除验证码，防止重复使用
	delete(captchaStore, req.CaptchaId)

	c.JSON(http.StatusOK, &VerifyCaptchaResponse{
		Msg: "Captcha verified successfully",
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

	// 生成邮箱验证码
	emailCode := config.GenerateEmailCode()
	emailCodeId := config.GenerateEmailCodeId()

	// 存储验证码并设置过期时间
	config.EmailCodeStore.Set(emailCodeId, &config.EmailCodeData{
		Email:     req.Email,
		Code:      emailCode,
		CreatedAt: time.Now().UnixMilli(),
		ExpiresAt: time.Now().UnixMilli() + config.EMAIL_CODE_CONFIG.Lifetime,
		Attempts:  0,
	})

	// 发送验证码邮件
	err = config.SendVerificationEmail(req.Email, emailCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorResponse{Msg: "Failed to send email"})
		return
	}

	// 定期清理过期验证码
	if rand.Float32() < 0.1 { // 10%的概率触发清理
		go config.CleanupExpiredEmailCodes()
	}

	c.JSON(http.StatusOK, &SendEmailCodeResponse{
		Msg: "Verification code sent successfully",
		Id:  emailCodeId,
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
		c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "User does not exist"})
		return
	}

	// 生成邮箱验证码
	emailCode := config.GenerateEmailCode()
	emailCodeId := config.GenerateEmailCodeId()

	// 存储验证码并设置过期时间
	config.EmailCodeStore.Set(emailCodeId, &config.EmailCodeData{
		Email:     req.Email,
		Code:      emailCode,
		CreatedAt: time.Now().UnixMilli(),
		ExpiresAt: time.Now().UnixMilli() + config.EMAIL_CODE_CONFIG.Lifetime,
		Attempts:  0,
	})

	// 发送验证码邮件
	err = config.SendVerificationEmail(req.Email, emailCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorResponse{Msg: "Failed to send email"})
		return
	}

	// 定期清理过期验证码
	if rand.Float32() < 0.1 { // 10%的概率触发清理
		go config.CleanupExpiredEmailCodes()
	}

	c.JSON(http.StatusOK, &SendLoginEmailCodeResponse{
		Msg: "Login verification code sent successfully",
		Id:  emailCodeId,
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

	// 检查是否启用了邮箱验证码功能
	if os.Getenv("ENABLE_EMAIL_VERIFICATION") == "true" {
		// 验证邮箱验证码逻辑
		if req.EmailCode == "" {
			c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Email verification code is required"})
			return
		}

		// 检查是否提供了邮箱验证码ID
		if req.EmailCodeId == "" {
			c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Email verification code ID is required"})
			return
		}

		// 验证验证码是否正确
		storedEmailCode := config.EmailCodeStore.Get(req.EmailCodeId)
		if storedEmailCode == nil {
			c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Invalid or expired email verification code"})
			return
		}

		// 检查邮箱是否匹配
		if storedEmailCode.Email != req.Email {
			c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Email does not match the verification code"})
			return
		}

		// 检查尝试次数
		if storedEmailCode.Attempts >= config.EMAIL_CODE_CONFIG.MaxAttempts {
			config.EmailCodeStore.Delete(req.EmailCodeId)
			c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Too many attempts. Please request a new verification code."})
			return
		}

		// 增加尝试次数
		storedEmailCode.Attempts++
		config.EmailCodeStore.Set(req.EmailCodeId, storedEmailCode)

		// 验证验证码是否正确
		if storedEmailCode.Code != strings.ToUpper(req.EmailCode) {
			c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Invalid email verification code"})
			return
		}

		// 验证成功后删除验证码，防止重复使用
		config.EmailCodeStore.Delete(req.EmailCodeId)
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
	payload := jwt.MapClaims{
		"user": map[string]interface{}{
			"id": userID.Hex(),
		},
		"exp": time.Now().Add(time.Hour * 1).Unix(), // 1 hour
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorResponse{Msg: "Server error"})
		return
	}

	c.JSON(http.StatusOK, &RegisterResponse{
		Token: tokenString,
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

	// 检查用户是否存在
	var userDoc bson.M
	err := config.DB.Collection("users").FindOne(c, bson.M{"email": req.Email}).Decode(&userDoc)
	if err != nil {
		c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Invalid Credentials"})
		return
	}

	// 检查是否启用了邮箱验证码登录
	isEmailCodeLogin := os.Getenv("ENABLE_EMAIL_VERIFICATION") == "true" && req.EmailCode != "" && req.EmailCodeId != ""

	// 如果提供了邮箱验证码，则验证邮箱验证码登录
	if isEmailCodeLogin {
		// 验证邮箱验证码逻辑
		storedEmailCode := config.EmailCodeStore.Get(req.EmailCodeId)
		if storedEmailCode == nil {
			c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Invalid or expired email verification code"})
			return
		}

		// 检查邮箱是否匹配
		if storedEmailCode.Email != req.Email {
			c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Email does not match the verification code"})
			return
		}

		// 检查尝试次数
		if storedEmailCode.Attempts >= config.EMAIL_CODE_CONFIG.MaxAttempts {
			config.EmailCodeStore.Delete(req.EmailCodeId)
			c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Too many attempts. Please request a new verification code."})
			return
		}

		// 增加尝试次数
		storedEmailCode.Attempts++
		config.EmailCodeStore.Set(req.EmailCodeId, storedEmailCode)

		// 验证验证码是否正确
		if storedEmailCode.Code != strings.ToUpper(req.EmailCode) {
			c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Invalid email verification code"})
			return
		}

		// 验证成功后删除验证码，防止重复使用
		config.EmailCodeStore.Delete(req.EmailCodeId)
	} else {
		// 否则验证密码登录
		// 检查是否提供了密码
		if req.Password == "" {
			c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Password is required"})
			return
		}

		hashedPassword, ok := userDoc["password"].(string)
		if !ok || bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)) != nil {
			c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Invalid Credentials"})
			return
		}
	}

	// 检查是否启用了验证码功能
	// 但邮箱验证码登录时不需要图片验证码
	if os.Getenv("ENABLE_CAPTCHA") == "true" && !isEmailCodeLogin {
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
		if storedCaptcha.Text != strings.ToUpper(req.Captcha) {
			c.JSON(http.StatusBadRequest, &ErrorResponse{Msg: "Invalid captcha"})
			return
		}

		// 验证成功后删除验证码，防止重复使用
		delete(captchaStore, req.CaptchaId)
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
	payload := jwt.MapClaims{
		"user": map[string]interface{}{
			"id": userID,
		},
		"exp": time.Now().Add(time.Hour * 1).Unix(), // 1 hour
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorResponse{Msg: "Server error"})
		return
	}

	c.JSON(http.StatusOK, &LoginResponse{
		Token: tokenString,
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
	// JWT payload 格式: { "user": { "id": "xxx" } }
	userData, ok := claims["user"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, &ErrorResponse{Msg: "Invalid token structure"})
		return
	}

	userIDStr, ok := userData["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, &ErrorResponse{Msg: "Invalid user ID in token"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
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

	// 构造返回的用户信息 - 直接返回用户对象
	user := User{
		Id:       userDoc["_id"].(primitive.ObjectID).Hex(),
		Username: userDoc["username"].(string),
		Email:    userDoc["email"].(string),
	}

	// backend 直接返回用户对象，不包装在响应对象中
	c.JSON(http.StatusOK, user)
}

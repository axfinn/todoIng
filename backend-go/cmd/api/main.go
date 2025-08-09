// @title TodoIng Backend API
// @version 1.0
// @description 这是 TodoIng 项目的后端API服务，提供任务管理、用户认证、报表生成等功能
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:5004
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

// Tasks API
// @Summary 创建新任务
// @Description 创建一个新的任务项
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param task body object true "任务信息"
// @Success 200 {object} map[string]interface{} "创建成功"
// @Router /api/tasks [post]

// @Summary 获取任务列表
// @Description 获取当前用户的所有任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} map[string]interface{} "任务列表"
// @Router /api/tasks [get]

// @Summary 获取任务详情
// @Description 根据ID获取任务详情
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "任务ID"
// @Success 200 {object} map[string]interface{} "任务详情"
// @Router /api/tasks/{id} [get]

// @Summary 更新任务
// @Description 更新任务信息
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "任务ID"
// @Param task body object true "任务信息"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Router /api/tasks/{id} [put]

// @Summary 删除任务
// @Description 删除指定任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "任务ID"
// @Success 200 {object} map[string]string "删除成功"
// @Router /api/tasks/{id} [delete]

// Reports API
// @Summary 获取报表列表
// @Description 获取用户的所有报表
// @Tags 报表管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} map[string]interface{} "报表列表"
// @Router /api/reports [get]

// @Summary 获取报表详情
// @Description 根据ID获取报表详情
// @Tags 报表管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "报表ID"
// @Success 200 {object} map[string]interface{} "报表详情"
// @Router /api/reports/{id} [get]

// Captcha API
// @Summary 生成验证码
// @Description 生成验证码图片
// @Tags 验证码
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "验证码图片和ID"
// @Router /api/auth/captcha [get]

// @Summary 验证验证码
// @Description 验证用户输入的验证码
// @Tags 验证码
// @Accept json
// @Produce json
// @Param body body object true "验证码信息"
// @Success 200 {object} map[string]string "验证成功"
// @Router /api/auth/verify-captcha [post]

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/swaggo/http-swagger"

	"github.com/axfinn/todoIng/backend-go/internal/api"
	"github.com/axfinn/todoIng/backend-go/internal/email"
	"github.com/axfinn/todoIng/backend-go/internal/captcha"
	"github.com/axfinn/todoIng/backend-go/internal/observability"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/axfinn/todoIng/backend-go/docs" // 导入生成的文档
)

// 确保swag能够扫描到所有handler类型和swagger注释
func init() {
	// 引用所有handler依赖类型，确保swag扫描时能发现它们
	_ = api.TaskDeps{}
	_ = api.ReportDeps{}
	_ = api.CaptchaDeps{}
	_ = api.AuthDeps{}
}

var client *mongo.Client

func main() {
	_ = godotenv.Load()

	// 初始化日志系统
	observability.InitLogger()
	observability.LogInfo("Application starting up...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// --- tracing init ---
	shutdown, errTrace := observability.InitTracer(context.Background(), "todoing-api", os.Getenv("ENVIRONMENT"), "1.0")
	if errTrace != nil { log.Printf("tracing init error: %v", errTrace) } else { defer func(){ _ = shutdown(context.Background()) }() }

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		observability.LogError("MONGO_URI environment variable not set")
		log.Fatal("MONGO_URI not set")
	}
	observability.LogInfo("Connecting to MongoDB at %s", mongoURI)
	
	var err error
	clientOpts := options.Client().ApplyURI(mongoURI)
	// 暂时移除 mongo tracing 监控
	client, err = mongo.Connect(ctx, clientOpts)
	if err != nil {
		observability.LogError("Failed to connect to MongoDB: %v", err)
		log.Fatal(err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		observability.LogError("Failed to ping MongoDB: %v", err)
		log.Fatal(err)
	}
	observability.LogInfo("MongoDB connected successfully")

	db := client.Database("todoing")
	emailStore := email.NewStore(10*time.Minute, 3)
	captchaStore := captcha.NewStore(5 * time.Minute)
	observability.LogInfo("Email store and captcha store initialized")

	r := api.NewRouter()
	// 暂时直接使用普通的 router，不使用 otelhttp
	handler := r
	observability.LogInfo("Router initialized")
	
	// Swagger 文档路由
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	
	// 静态文件服务 - API 文档
	docsHandler := http.StripPrefix("/docs/", http.FileServer(http.Dir("docs/")))
	r.PathPrefix("/docs/").Handler(docsHandler)
	
	// 完整 API 文档路由
	r.HandleFunc("/api-docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "docs/api_complete.json")
	}).Methods(http.MethodGet)
	
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	api.SetupAuthRoutes(r, &api.AuthDeps{DB: db, EmailCodes: emailStore})
	api.SetupCaptchaRoutes(r, &api.CaptchaDeps{Store: captchaStore})
	api.SetupTaskRoutes(r, &api.TaskDeps{DB: db})
	api.SetupReportRoutes(r, &api.ReportDeps{DB: db})
	observability.LogInfo("All API routes configured")

	port := os.Getenv("PORT")
	if port == "" { port = "5001" }
	server := &http.Server{Addr: ":"+port, Handler: handler}
	observability.LogInfo("HTTP server configured on port %s", port)

	// create default user if not exists
	go func() {
		time.Sleep(500 * time.Millisecond)
		observability.LogInfo("Starting default user creation check...")
		username := os.Getenv("DEFAULT_USERNAME")
		password := os.Getenv("DEFAULT_PASSWORD")
		emailAddr := os.Getenv("DEFAULT_EMAIL")
		if username == "" || password == "" || emailAddr == "" { 
			observability.LogWarn("Default user environment variables not set, skipping default user creation")
			return 
		}
		ctxDef, cancelDef := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelDef()
		usersCol := db.Collection("users")
		err := usersCol.FindOne(ctxDef, bson.M{"$or": []bson.M{{"username": username},{"email": emailAddr}}}).Err()
		if err == mongo.ErrNoDocuments {
			hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
			_, errIns := usersCol.InsertOne(ctxDef, bson.M{"username": username, "email": emailAddr, "password": string(hash), "createdAt": time.Now()})
			if errIns != nil { 
				observability.LogError("Failed to create default user: %v", errIns)
			} else { 
				observability.LogInfo("Default user created successfully: %s (%s)", username, emailAddr)
			}
		} else if err == nil {
			observability.LogInfo("Default user already exists: %s", username)
		} else {
			observability.LogError("Error checking for default user: %v", err)
		}
	}()

	go func() {
		observability.LogInfo("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			observability.LogError("Server error: %s", err)
			log.Fatalf("listen: %s", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	observability.LogInfo("Server is ready and listening for requests")
	<-quit
	observability.LogInfo("Shutdown signal received, starting graceful shutdown...")
	
	ctxShut, cancelShut := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShut()
	
	if err := server.Shutdown(ctxShut); err != nil { 
		observability.LogError("Server shutdown error: %v", err)
	} else {
		observability.LogInfo("HTTP server shutdown successfully")
	}
	
	if err := client.Disconnect(ctxShut); err != nil { 
		observability.LogError("MongoDB disconnect error: %v", err)
	} else {
		observability.LogInfo("MongoDB disconnected successfully")
	}
	
	observability.LogInfo("Application shutdown complete")
	fmt.Println("Server exiting")
}

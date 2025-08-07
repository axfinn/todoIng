package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "todoing-backend/docs" // Swagger docs
	"todoing-backend/src/config"
	"todoing-backend/src/middleware"
	"todoing-backend/src/routes"
)

var (
	// 默认环境变量值
	defaultPort = "5001"
	// 修改MongoDB URI指向docker-compose中的mongodb服务
	defaultMongoURI  = "mongodb://todoing_mongodb_dev:27017/todoing"
	defaultJWTSecret = "todoing_secret_key"
)

// @title Todoing API
// @version 1.0
// @description This is a todoing server.
// @host localhost:5001
// @BasePath /
func main() {
	// 加载 .env 文件（如果存在）
	godotenv.Load()

	// 获取环境变量，如果不存在则使用默认值
	port := getEnv("PORT", defaultPort)
	mongoURI := getEnv("MONGO_URI", defaultMongoURI)
	jwtSecret := getEnv("JWT_SECRET", defaultJWTSecret)

	// 设置 JWT 密钥
	os.Setenv("JWT_SECRET", jwtSecret)

	// 连接 MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	// 设置数据库
	config.DB = client.Database("todoing")

	// 创建 Gin 路由
	router := gin.Default()

	// 注册路由
	routes.RegisterAuthRoutes(router)

	// 为任务和报告路由添加认证中间件
	authorized := router.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		routes.RegisterTaskRoutes(authorized)
		routes.RegisterReportRoutes(authorized)
	}

	// 注册 Swagger 路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动 HTTP 服务器
	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Server running on 0.0.0.0:%s", port)
	log.Printf("Swagger UI available at http://0.0.0.0:%s/swagger/index.html", port)

	// 在单独的 goroutine 中启动 HTTP 服务器
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen and serve: %v", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 优雅地关闭 HTTP 服务器
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}


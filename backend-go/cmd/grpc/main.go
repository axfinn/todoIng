package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/axfinn/todoIng/backend-go/internal/services"
	pb "github.com/axfinn/todoIng/backend-go/pkg/api/v1"
)

// @title TodoIng gRPC API
// @version 1.0
// @description 这是 TodoIng 项目的 gRPC API 服务
// @host localhost:9001
// @BasePath /

func main() {
	_ = godotenv.Load()

	// 连接数据库
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not set")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}
	log.Println("MongoDB connected for gRPC server")

	db := client.Database("todoing")

	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer()

	// 注册服务
	authService := services.NewAuthService(db)
	pb.RegisterAuthServiceServer(grpcServer, authService)

	// 启用反射（用于 grpcurl 等工具）
	reflection.Register(grpcServer)

	// 监听端口
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "9001"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// 启动服务器
	go func() {
		log.Printf("gRPC server running on port %s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()

	ctxShut, cancelShut := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShut()
	if err := client.Disconnect(ctxShut); err != nil {
		log.Printf("MongoDB disconnect error: %v", err)
	}

	fmt.Println("gRPC server exited")
}

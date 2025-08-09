package testutil

import (
	"context"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestDB 测试数据库结构
type TestDB struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// SetupTestDB 设置测试数据库
func SetupTestDB(t *testing.T) *TestDB {
	// 使用环境变量或默认的测试数据库URI
	mongoURI := os.Getenv("TEST_MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	// 使用唯一的测试数据库名
	dbName := "todoing_test_" + generateTestDBName()
	db := client.Database(dbName)

	return &TestDB{
		Client: client,
		DB:     db,
	}
}

// TeardownTestDB 清理测试数据库
func (tdb *TestDB) TeardownTestDB(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 删除测试数据库
	if err := tdb.DB.Drop(ctx); err != nil {
		t.Logf("Warning: Failed to drop test database: %v", err)
	}

	// 关闭连接
	if err := tdb.Client.Disconnect(ctx); err != nil {
		t.Logf("Warning: Failed to disconnect from test database: %v", err)
	}
}

// CleanCollection 清理指定集合
func (tdb *TestDB) CleanCollection(t *testing.T, collectionName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := tdb.DB.Collection(collectionName)
	if _, err := collection.DeleteMany(ctx, bson.M{}); err != nil {
		t.Logf("Warning: Failed to clean collection %s: %v", collectionName, err)
	}
}

// generateTestDBName 生成唯一的测试数据库名
func generateTestDBName() string {
	return time.Now().Format("20060102_150405_000")
}

// MockUser 创建测试用户数据
func MockUser() map[string]interface{} {
	return map[string]interface{}{
		"username":   "testuser",
		"email":      "test@example.com",
		"password":   "$2a$10$hash", // bcrypt hash
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
}

// MockTask 创建测试任务数据
func MockTask(userID string) map[string]interface{} {
	return map[string]interface{}{
		"title":       "Test Task",
		"description": "Test Description",
		"status":      "todo",
		"priority":    "medium",
		"due_date":    time.Now().Add(24 * time.Hour),
		"user_id":     userID,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}
}

// MockReport 创建测试报表数据
func MockReport(userID string) map[string]interface{} {
	return map[string]interface{}{
		"title":      "Test Report",
		"type":       "daily",
		"start_date": time.Now().Add(-24 * time.Hour),
		"end_date":   time.Now(),
		"user_id":    userID,
		"created_at": time.Now(),
	}
}

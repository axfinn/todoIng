package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/axfinn/todoIng/backend-go/internal/auth"
	"github.com/axfinn/todoIng/backend-go/internal/captcha"
	"github.com/axfinn/todoIng/backend-go/internal/email"
	"github.com/axfinn/todoIng/backend-go/test/testutil"
)

func TestRegister(t *testing.T) {
	// 设置测试数据库
	testDB := testutil.SetupTestDB(t)
	defer testDB.TeardownTestDB(t)

	// 设置测试依赖
	emailStore := email.NewStore(10*time.Minute, 3)
	captchaStore := captcha.NewStore(5 * time.Minute)

	// 创建路由
	r := mux.NewRouter()
	setupAuthRoutes(r, testDB.DB, emailStore, captchaStore)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedField  string
	}{
		{
			name: "成功注册",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusCreated,
			expectedField:  "user",
		},
		{
			name: "缺少用户名",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "缺少邮箱",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "缺少密码",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "重复邮箱",
			requestBody: map[string]interface{}{
				"username": "testuser2",
				"email":    "test@example.com", // 与第一个测试相同的邮箱
				"password": "password123",
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 对于重复邮箱测试，先清理数据库
			if tt.name == "重复邮箱" {
				// 先注册一个用户
				firstUser := map[string]interface{}{
					"username": "firstuser",
					"email":    "test@example.com",
					"password": "password123",
				}
				jsonData, _ := json.Marshal(firstUser)
				req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)
			}

			jsonData, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Response: %s",
					tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.expectedField != "" {
				var response map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if _, exists := response[tt.expectedField]; !exists {
					t.Errorf("Expected field '%s' not found in response", tt.expectedField)
				}
			}

			// 清理数据库以避免测试间干扰
			testDB.CleanCollection(t, "users")
		})
	}
}

func TestLogin(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	defer testDB.TeardownTestDB(t)

	emailStore := email.NewStore(10*time.Minute, 3)
	captchaStore := captcha.NewStore(5 * time.Minute)

	r := mux.NewRouter()
	setupAuthRoutes(r, testDB.DB, emailStore, captchaStore)

	// 先注册一个用户
	registerData := map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonData, _ := json.Marshal(registerData)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "成功登录",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "错误密码",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "不存在的用户",
			requestBody: map[string]interface{}{
				"email":    "nonexistent@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Response: %s",
					tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if _, exists := response["token"]; !exists {
					t.Error("Expected 'token' field not found in response")
				}
			}
		})
	}
}

func TestVerifyToken(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	defer testDB.TeardownTestDB(t)

	// 创建一个有效的JWT令牌
	userID := primitive.NewObjectID()
	token, err := auth.GenerateJWT(userID.Hex())
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	emailStore := email.NewStore(10*time.Minute, 3)
	captchaStore := captcha.NewStore(5 * time.Minute)

	r := mux.NewRouter()
	setupAuthRoutes(r, testDB.DB, emailStore, captchaStore)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "有效令牌",
			authHeader:     "Bearer " + token,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "无效令牌",
			authHeader:     "Bearer invalid_token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "缺少令牌",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/auth/verify", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Response: %s",
					tt.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

// setupAuthRoutes 设置认证路由用于测试
func setupAuthRoutes(r *mux.Router, db *mongo.Database, emailStore *email.Store, captchaStore *captcha.Store) {
	authRouter := r.PathPrefix("/auth").Subrouter()

	authRouter.HandleFunc("/register", RegisterHandler(db, emailStore, captchaStore)).Methods("POST")
	authRouter.HandleFunc("/login", LoginHandler(db)).Methods("POST")
	authRouter.HandleFunc("/verify", VerifyTokenHandler(db)).Methods("GET")
	authRouter.HandleFunc("/login/email-code", EmailCodeLoginHandler(db, emailStore)).Methods("POST")
	authRouter.HandleFunc("/send-login-email-code", SendLoginEmailCodeHandler(db, emailStore)).Methods("POST")
}

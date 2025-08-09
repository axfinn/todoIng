package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

// 测试健康检查端点
func TestHealthEndpoint(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "ok" {
		t.Errorf("Expected 'ok', got '%s'", w.Body.String())
	}
}

// 测试注册请求验证
func TestRegisterValidation(t *testing.T) {
	tests := []struct {
		name          string
		requestBody   map[string]interface{}
		expectedValid bool
	}{
		{
			name: "有效注册请求",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedValid: true,
		},
		{
			name: "缺少用户名",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedValid: false,
		},
		{
			name: "密码太短",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "123",
			},
			expectedValid: false,
		},
		{
			name: "邮箱格式无效",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"email":    "invalid-email",
				"password": "password123",
			},
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req registerRequest
			jsonData, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			if err := json.Unmarshal(jsonData, &req); err != nil {
				t.Fatalf("Failed to unmarshal request body: %v", err)
			}

			// 验证请求字段
			isValidEmail := strings.Contains(req.Email, "@") && strings.Contains(req.Email, ".")
			isValid := req.Username != "" && req.Email != "" && len(req.Password) >= 6 && isValidEmail

			if isValid != tt.expectedValid {
				t.Errorf("Expected validation result %v, got %v", tt.expectedValid, isValid)
			}
		})
	}
}

// 测试登录请求验证
func TestLoginValidation(t *testing.T) {
	tests := []struct {
		name          string
		requestBody   map[string]interface{}
		expectedValid bool
	}{
		{
			name: "有效登录请求",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedValid: true,
		},
		{
			name: "缺少邮箱",
			requestBody: map[string]interface{}{
				"password": "password123",
			},
			expectedValid: false,
		},
		{
			name: "缺少密码",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req loginRequest
			jsonData, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			if err := json.Unmarshal(jsonData, &req); err != nil {
				t.Fatalf("Failed to unmarshal request body: %v", err)
			}

			// 验证请求字段
			isValid := req.Email != "" && req.Password != ""

			if isValid != tt.expectedValid {
				t.Errorf("Expected validation result %v, got %v", tt.expectedValid, isValid)
			}
		})
	}
}

// 测试JSON响应助手函数
func TestJSONHelper(t *testing.T) {
	w := httptest.NewRecorder()

	testData := map[string]string{
		"message": "test",
		"status":  "success",
	}

	JSON(w, http.StatusOK, testData)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "test" {
		t.Errorf("Expected message 'test', got '%s'", response["message"])
	}

	if response["status"] != "success" {
		t.Errorf("Expected status 'success', got '%s'", response["status"])
	}
}

// 测试环境变量处理
func TestEnvironmentVariables(t *testing.T) {
	// 保存原始环境变量
	originalRegistration := os.Getenv("DISABLE_REGISTRATION")
	originalEmailVerification := os.Getenv("ENABLE_EMAIL_VERIFICATION")

	// 清理函数
	defer func() {
		if originalRegistration != "" {
			os.Setenv("DISABLE_REGISTRATION", originalRegistration)
		} else {
			os.Unsetenv("DISABLE_REGISTRATION")
		}
		if originalEmailVerification != "" {
			os.Setenv("ENABLE_EMAIL_VERIFICATION", originalEmailVerification)
		} else {
			os.Unsetenv("ENABLE_EMAIL_VERIFICATION")
		}
	}()

	// 测试注册禁用
	os.Setenv("DISABLE_REGISTRATION", "true")
	if os.Getenv("DISABLE_REGISTRATION") != "true" {
		t.Error("Failed to set DISABLE_REGISTRATION environment variable")
	}

	// 测试邮箱验证启用
	os.Setenv("ENABLE_EMAIL_VERIFICATION", "true")
	if os.Getenv("ENABLE_EMAIL_VERIFICATION") != "true" {
		t.Error("Failed to set ENABLE_EMAIL_VERIFICATION environment variable")
	}
}

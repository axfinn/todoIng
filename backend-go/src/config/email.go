package config

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"
)

// EmailCodeStore 存储邮箱验证码（生产环境建议使用Redis）
var EmailCodeStore = &emailCodeStore{
	store: make(map[string]*EmailCodeData),
}

// emailCodeStore 邮箱验证码存储结构
type emailCodeStore struct {
	store map[string]*EmailCodeData
	mu    sync.RWMutex
}

// EmailCodeData 邮箱验证码数据
type EmailCodeData struct {
	Email     string
	Code      string
	CreatedAt int64
	ExpiresAt int64
	Attempts  int
}

// EMAIL_CODE_CONFIG 邮箱验证码配置
var EMAIL_CODE_CONFIG = struct {
	Length      int
	Lifetime    int64
	MaxAttempts int
}{
	Length:      6,             // 验证码长度
	Lifetime:    5 * 60 * 1000, // 5分钟有效期（毫秒）
	MaxAttempts: 3,             // 最大尝试次数
}

// GenerateEmailCode 生成邮箱验证码
func GenerateEmailCode() string {
	const digits = "0123456789"
	code := make([]byte, EMAIL_CODE_CONFIG.Length)
	for i := range code {
		code[i] = digits[randInt(len(digits))]
	}
	return string(code)
}

// GenerateEmailCodeId 生成邮箱验证码ID
func GenerateEmailCodeId() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// SendVerificationEmail 发送验证邮件（模拟）
func SendVerificationEmail(email string, code string) error {
	// TODO: 实现真实的邮件发送功能
	// 这里只是模拟，实际应该调用邮件服务
	fmt.Printf("Sending verification email to %s with code: %s\n", email, code)
	return nil
}

// CleanupExpiredEmailCodes 清理过期的邮箱验证码
func CleanupExpiredEmailCodes() {
	EmailCodeStore.mu.Lock()
	defer EmailCodeStore.mu.Unlock()

	now := time.Now().UnixMilli()
	for id, codeData := range EmailCodeStore.store {
		if codeData.ExpiresAt < now {
			delete(EmailCodeStore.store, id)
		}
	}
}

// Set 设置邮箱验证码
func (s *emailCodeStore) Set(id string, data *EmailCodeData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[id] = data
}

// Get 获取邮箱验证码
func (s *emailCodeStore) Get(id string) *EmailCodeData {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.store[id]
}

// Delete 删除邮箱验证码
func (s *emailCodeStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, id)
}

// randInt 生成随机整数
func randInt(max int) int {
	b := make([]byte, 1)
	rand.Read(b)
	return int(b[0]) % max
}

// 启动定时清理任务
func init() {
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			CleanupExpiredEmailCodes()
		}
	}()
}


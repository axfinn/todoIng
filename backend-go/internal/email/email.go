package email

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	mail "github.com/go-mail/mail/v2"
)

type Code struct {
	Email     string
	Code      string
	Attempts  int
	ExpiresAt time.Time
}

type Store struct {
	m           sync.Map
	TTL         time.Duration
	MaxAttempts int
}

func NewStore(ttl time.Duration, maxAttempts int) *Store {
	return &Store{TTL: ttl, MaxAttempts: maxAttempts}
}

func randomID() string { b := make([]byte, 16); rand.Read(b); return hex.EncodeToString(b) }

func randomCode(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return strings.ToUpper(hex.EncodeToString(b)[:n])
}

func (s *Store) Generate(email string, length int) (id, code string) {
	id = randomID()
	code = randomCode(length)
	s.m.Store(id, &Code{Email: email, Code: code, Attempts: 0, ExpiresAt: time.Now().Add(s.TTL)})
	return
}

func (s *Store) Verify(id, email, input string) error {
	v, ok := s.m.Load(id)
	if !ok {
		return errors.New("invalid or expired code")
	}
	c := v.(*Code)
	if time.Now().After(c.ExpiresAt) {
		s.m.Delete(id)
		return errors.New("code expired")
	}
	if c.Email != email {
		return errors.New("email mismatch")
	}
	if c.Attempts >= s.MaxAttempts {
		s.m.Delete(id)
		return errors.New("too many attempts")
	}
	c.Attempts++
	if strings.ToUpper(input) != c.Code {
		return errors.New("invalid code")
	}
	s.m.Delete(id)
	return nil
}

func (s *Store) Cleanup() {
	now := time.Now()
	s.m.Range(func(key, value any) bool {
		c := value.(*Code)
		if now.After(c.ExpiresAt) {
			s.m.Delete(key)
		}
		return true
	})
}

func Send(to, code string) error {
	host := os.Getenv("EMAIL_HOST")
	user := os.Getenv("EMAIL_USER")
	pass := os.Getenv("EMAIL_PASS")
	port := os.Getenv("EMAIL_PORT")
	if host == "" || user == "" || pass == "" {
		return fmt.Errorf("email config incomplete")
	}
	p := 587
	fmt.Sscanf(port, "%d", &p)
	d := mail.NewDialer(host, p, user, pass)
	m := mail.NewMessage()
	from := os.Getenv("EMAIL_FROM")
	if from == "" {
		from = user
	}
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "TodoIng 邮箱验证码")
	m.SetBody("text/html", fmt.Sprintf("<p>您的验证码: <b>%s</b> (10分钟内有效)</p>", code))
	return d.DialAndSend(m)
}

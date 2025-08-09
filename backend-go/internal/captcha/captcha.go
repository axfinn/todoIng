package captcha

import (
	crand "crypto/rand"
	"encoding/hex"
	mrand "math/rand"
	"strings"
	"sync"
	"time"
)

type Item struct {
	Text      string
	ExpiresAt time.Time
}

type Store struct {
	m   sync.Map // key -> Item
	TTL time.Duration
}

func NewStore(ttl time.Duration) *Store {
	return &Store{TTL: ttl}
}

const letters = "ABCDEFGHJKMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789"

func RandomText(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(letters[mrand.Intn(len(letters))])
	}
	return sb.String()
}

func randomID() string {
	b := make([]byte, 16)
	crand.Read(b)
	return hex.EncodeToString(b)
}

func init() {
	mrand.Seed(time.Now().UnixNano())
}

func (s *Store) Generate(n int) (id, text string) {
	text = strings.ToUpper(RandomText(n))
	id = randomID()
	s.m.Store(id, Item{Text: text, ExpiresAt: time.Now().Add(s.TTL)})
	return
}

func (s *Store) Verify(id, value string) bool {
	v, ok := s.m.Load(id)
	if !ok {
		return false
	}
	item := v.(Item)
	if time.Now().After(item.ExpiresAt) {
		s.m.Delete(id)
		return false
	}
	if strings.ToUpper(value) != item.Text {
		return false
	}
	s.m.Delete(id)
	return true
}

func (s *Store) Cleanup() {
	now := time.Now()
	s.m.Range(func(key, value any) bool {
		item := value.(Item)
		if now.After(item.ExpiresAt) {
			s.m.Delete(key)
		}
		return true
	})
}

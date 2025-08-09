package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func Secret() []byte { return []byte(os.Getenv("JWT_SECRET")) }

func Generate(userID string, ttl time.Duration) (string, error) {
	claims := Claims{UserID: userID, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl))}}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(Secret())
}

func Parse(token string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) { return Secret(), nil })
	if err != nil {
		return nil, err
	}
	if claims.UserID == "" {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// GenerateJWT 兼容测试的函数名
func GenerateJWT(userID string) (string, error) {
	return Generate(userID, time.Hour)
}

package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/axfinn/todoIng/backend-go/internal/auth"
	"github.com/axfinn/todoIng/backend-go/internal/observability"
	"github.com/gorilla/mux"
)

type contextKey string

const userKey contextKey = "userId"

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &responseWriter{ResponseWriter: w, status: 200}

		next.ServeHTTP(ww, r)
		dur := time.Since(start)

		observability.CtxLog(r.Context(), "%s %s %d %s", r.Method, r.URL.Path, ww.status, dur)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				ctx := r.Context()
				observability.CtxLog(ctx, "panic: %v", err)

				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "server error"})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "No token", http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid auth header", http.StatusUnauthorized)
			return
		}
		claims, err := auth.Parse(parts[1])
		if err != nil {
			http.Error(w, "Token invalid", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(r *http.Request) string {
	v := r.Context().Value(userKey)
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func JSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func decodeJSON(r *http.Request, v interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(Logging, Recover)
	return r
}

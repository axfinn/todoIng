package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/axfinn/todoIng/backend-go/internal/auth"
	"github.com/axfinn/todoIng/backend-go/internal/email"
	"github.com/axfinn/todoIng/backend-go/internal/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthDeps struct {
	DB *mongo.Database
	EmailCodes *email.Store
}

type registerRequest struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	EmailCode string `json:"emailCode"`
	EmailCodeId string `json:"emailCodeId"`
}

type loginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
	EmailCode string `json:"emailCode"`
	EmailCodeId string `json:"emailCodeId"`
}

func (d *AuthDeps) Register(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("DISABLE_REGISTRATION") == "true" {
		JSON(w, http.StatusForbidden, map[string]string{"msg":"Registration is disabled"}); return
	}
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { JSON(w, 400, map[string]string{"msg":"Invalid body"}); return }
	if req.Username=="" || req.Email=="" || len(req.Password)<6 { JSON(w,400,map[string]string{"msg":"Invalid fields"}); return }
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second); defer cancel()
	users := d.DB.Collection("users")
	// email verify
	if os.Getenv("ENABLE_EMAIL_VERIFICATION") == "true" {
		if err := d.EmailCodes.Verify(req.EmailCodeId, strings.ToLower(req.Email), req.EmailCode); err != nil { JSON(w,400,map[string]string{"msg":err.Error()}); return }
	}
	// uniqueness
	if err := users.FindOne(ctx, bson.M{"$or": []bson.M{{"email": req.Email},{"username": req.Username}}}).Err(); err == nil { JSON(w,400,map[string]string{"msg":"User already exists"}); return }
	pwHash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	user := models.User{Username: req.Username, Email: strings.ToLower(req.Email), Password: string(pwHash), CreatedAt: time.Now()}
	res, err := users.InsertOne(ctx, user)
	if err != nil { JSON(w,500,map[string]string{"msg":"DB error"}); return }
	id := res.InsertedID.(primitive.ObjectID).Hex()
	token, _ := auth.Generate(id, time.Hour)
	JSON(w,200,map[string]string{"token": token})
}

func (d *AuthDeps) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { JSON(w,400,map[string]string{"msg":"Invalid body"}); return }
	if req.Email=="" { JSON(w,400,map[string]string{"msg":"Email required"}); return }
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second); defer cancel()
	users := d.DB.Collection("users")
	var user models.User
	err := users.FindOne(ctx, bson.M{"email": strings.ToLower(req.Email)}).Decode(&user)
	if err != nil { JSON(w,400,map[string]string{"msg":"Invalid credentials"}); return }
	// email code login
	if req.EmailCode != "" && req.EmailCodeId != "" && os.Getenv("ENABLE_EMAIL_VERIFICATION") == "true" {
		// 使用小写邮箱地址进行验证，确保与存储时一致
		if err := d.EmailCodes.Verify(req.EmailCodeId, strings.ToLower(req.Email), req.EmailCode); err != nil { JSON(w,400,map[string]string{"msg":err.Error()}); return }
	} else {
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil { JSON(w,400,map[string]string{"msg":"Invalid credentials"}); return }
	}
	token, _ := auth.Generate(user.ID, time.Hour)
	JSON(w,200,map[string]string{"token": token})
}

func (d *AuthDeps) Me(w http.ResponseWriter, r *http.Request) {
	uid := GetUserID(r)
	if uid == "" { JSON(w,401,map[string]string{"msg":"Unauthorized"}); return }
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second); defer cancel()
	users := d.DB.Collection("users")
	var user models.User
	objID, _ := primitive.ObjectIDFromHex(uid)
	err := users.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil { JSON(w,500,map[string]string{"msg":"Not found"}); return }
	user.Password = ""
	JSON(w,200,user)
}

func (d *AuthDeps) SendRegisterEmailCode(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("ENABLE_EMAIL_VERIFICATION") != "true" { JSON(w,400,map[string]string{"msg":"Email verification disabled"}); return }
	var body struct { Email string `json:"email"` }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil { JSON(w,400,map[string]string{"msg":"Invalid body"}); return }
	if body.Email == "" { JSON(w,400,map[string]string{"msg":"Email required"}); return }
	
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second); defer cancel()
	users := d.DB.Collection("users")
	normalizedEmail := strings.ToLower(body.Email)
	
	// 注册时检查用户是否已存在
	if err := users.FindOne(ctx, bson.M{"email": normalizedEmail}).Err(); err == nil { 
		JSON(w,400,map[string]string{"msg":"User already exists"}); return 
	}
	
	id, code := d.EmailCodes.Generate(normalizedEmail, 6)
	_ = email.Send(body.Email, code) // 发送邮件使用原始邮箱格式
	JSON(w,200,map[string]string{"id": id, "msg":"Verification code sent"})
}

func (d *AuthDeps) SendLoginEmailCode(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("ENABLE_EMAIL_VERIFICATION") != "true" { JSON(w,400,map[string]string{"msg":"Email verification disabled"}); return }
	var body struct { Email string `json:"email"` }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil { JSON(w,400,map[string]string{"msg":"Invalid body"}); return }
	if body.Email == "" { JSON(w,400,map[string]string{"msg":"Email required"}); return }
	
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second); defer cancel()
	users := d.DB.Collection("users")
	normalizedEmail := strings.ToLower(body.Email)
	
	// 登录时检查用户是否存在
	var user models.User
	if err := users.FindOne(ctx, bson.M{"email": normalizedEmail}).Decode(&user); err != nil { 
		JSON(w,400,map[string]string{"msg":"User does not exist"}); return 
	}
	
	id, code := d.EmailCodes.Generate(normalizedEmail, 6)
	_ = email.Send(body.Email, code) // 发送邮件使用原始邮箱格式
	JSON(w,200,map[string]string{"id": id, "msg":"Login verification code sent"})
}

func SetupAuthRoutes(r *mux.Router, deps *AuthDeps) {
	r.HandleFunc("/api/auth/register", deps.Register).Methods(http.MethodPost)
	r.HandleFunc("/api/auth/login", deps.Login).Methods(http.MethodPost)
	r.Handle("/api/auth/me", Auth(http.HandlerFunc(deps.Me))).Methods(http.MethodGet)
	r.HandleFunc("/api/auth/send-email-code", deps.SendRegisterEmailCode).Methods(http.MethodPost)
	// 登录邮箱验证码使用专门的函数，检查用户是否存在
	r.HandleFunc("/api/auth/send-login-email-code", deps.SendLoginEmailCode).Methods(http.MethodPost)
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/axfinn/todoIng/backend-go/internal/api"
	"github.com/axfinn/todoIng/backend-go/internal/email"
	"github.com/axfinn/todoIng/backend-go/internal/captcha"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var client *mongo.Client

func main() {
	_ = godotenv.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not set")
	}
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}
	log.Println("MongoDB connected")

	db := client.Database("todoing")
	emailStore := email.NewStore(10*time.Minute, 3)
	captchaStore := captcha.NewStore(5 * time.Minute)

	r := api.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	api.SetupAuthRoutes(r, &api.AuthDeps{DB: db, EmailCodes: emailStore})
	api.SetupCaptchaRoutes(r, &api.CaptchaDeps{Store: captchaStore})
	api.SetupTaskRoutes(r, &api.TaskDeps{DB: db})
	api.SetupReportRoutes(r, &api.ReportDeps{DB: db})

	port := os.Getenv("PORT")
	if port == "" { port = "5001" }
	server := &http.Server{Addr: ":"+port, Handler: r}

	// create default user if not exists
	go func() {
		time.Sleep(500 * time.Millisecond)
		username := os.Getenv("DEFAULT_USERNAME")
		password := os.Getenv("DEFAULT_PASSWORD")
		emailAddr := os.Getenv("DEFAULT_EMAIL")
		if username == "" || password == "" || emailAddr == "" { return }
		ctxDef, cancelDef := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelDef()
		usersCol := db.Collection("users")
		err := usersCol.FindOne(ctxDef, bson.M{"$or": []bson.M{{"username": username},{"email": emailAddr}}}).Err()
		if err == mongo.ErrNoDocuments {
			hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
			_, errIns := usersCol.InsertOne(ctxDef, bson.M{"username": username, "email": emailAddr, "password": string(hash), "createdAt": time.Now()})
			if errIns != nil { log.Printf("create default user failed: %v", errIns) } else { log.Println("Default user created") }
		} else if err == nil {
			log.Println("Default user already exists")
		}
	}()

	go func() {
		log.Printf("Server running on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctxShut, cancelShut := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShut()
	if err := server.Shutdown(ctxShut); err != nil { log.Printf("Server Shutdown: %v", err) }
	if err := client.Disconnect(ctxShut); err != nil { log.Printf("Mongo disconnect: %v", err) }
	fmt.Println("Server exiting")
}

package main

import (
	"context"
	"log"
	"net/http"
	"user-feed/db"
	"user-feed/handlers"
	"user-feed/redis"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found: %v", err)
	}
}

func main() {
	var ctx = context.Background()
	pgBase, err := db.NewDataBase()

	if err != nil {
		log.Fatalf("Failed to connect postgres: %v", err)
	}
	err = pgBase.AutoMigrate()
	if err != nil {
		log.Fatalf("Failed to migrate data postgres: %v", err)
	}

	redisClient, err := redis.InitRedis(ctx)
	if err != nil {
		log.Fatalf("Failed to connect postgres: %v", err)
	}
	customHandler := handlers.NewCustomHandler(pgBase)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /user", customHandler.CreateUserHandler)
	//mux.HandleFunc("POST /post", handlers.CreatePostHandler)
	//mux.HandleFunc("GET /feed/{userID}", handlers.GetUserHandler)

	http.ListenAndServe(":8080", mux)

	log.Println("Сервер запущен")
	_ = redisClient
}

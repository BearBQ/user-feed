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

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HelloHandler)

	http.ListenAndServe(":8080", mux)
	_ = redisClient
}

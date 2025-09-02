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

var ctx = context.Background()
var err error

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found: %v", err)
	}
}

func main() {

	err = db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect postgres: %v", err)
	}
	err = redis.InitRedis(ctx)
	if err != nil {
		log.Fatalf("Failed to connect postgres: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HelloHandler)

	http.ListenAndServe(":8080", mux)
}

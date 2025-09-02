package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis(ctx context.Context) error {
	host := os.Getenv("RS_HOST")
	pswd := os.Getenv("RS_PASS")
	rDB, err := strconv.Atoi(os.Getenv("DB"))
	if err != nil {
		rDB = 0
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pswd,
		DB:       rDB,
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		fmt.Println(redisClient.Options().Addr, redisClient.Options().Password, redisClient.Options().DB)
		log.Fatal("failed to connect redis")
	}
	log.Printf("подключение к базе Redis %v успешно", rDB)
	return nil
}

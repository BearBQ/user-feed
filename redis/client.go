package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

// Close() закрывает соединение с редис
func (rc *RedisClient) Close() error {
	if rc.Client != nil {
		return rc.Client.Close()
	}
	return nil
}

// InitRedis активирует подключение к базе редис
func InitRedis(ctx context.Context) (*RedisClient, error) {
	host := os.Getenv("RS_HOST")
	pswd := os.Getenv("RS_PASS")
	rDB, err := strconv.Atoi(os.Getenv("DB"))
	if err != nil {
		rDB = 0
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pswd,
		DB:       rDB,
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		fmt.Println(redisClient.Options().Addr, redisClient.Options().Password, redisClient.Options().DB)
		return nil, fmt.Errorf("failed to connect redis: %v", err)
	}
	log.Printf("подключение к базе Redis %v успешно", rDB)
	return &RedisClient{Client: redisClient}, nil
}

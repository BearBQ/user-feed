package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisBase interface {
	Close() error
	GetDataFromRedis(ctx context.Context, userID uint) ([]byte, error)
	PushData(ctx context.Context, userID uint, data []byte) error
	ClearData(ctx context.Context, userID uint) error
}

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

func (rc *RedisClient) GetDataFromRedis(ctx context.Context, userID uint) ([]byte, error) {
	user := fmt.Sprintf("feed:%d", userID)
	json, err := rc.Client.Get(ctx, user).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("redis error: %w", err)
	}
	return []byte(json), nil
}

func (rc *RedisClient) PushData(ctx context.Context, userID uint, data []byte) error {
	user := fmt.Sprintf("feed:%d", userID)
	err := rc.Client.Set(ctx, user, data, 5*time.Minute)
	if err != nil {
		return fmt.Errorf("redis set error: %v", err)
	}
	return nil
}

func (rc *RedisClient) ClearData(ctx context.Context, userID uint) error {
	user := fmt.Sprintf("feed:%d", userID)
	err := rc.Client.Del(ctx, user).Err()
	if err != nil {
		if err == redis.Nil {
			// Ключ не существует - это не ошибка для операции удаления
			return nil
		}
		return fmt.Errorf("failed to delete key %d: %w", userID, err)
	}
	return nil
}

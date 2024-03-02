package connect

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
)

func NewRedisConnection() (*redis.Client, error) {

	redis_addr := os.Getenv("REDIS_SOURCE")
	redis_password := os.Getenv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     redis_addr,
		Password: redis_password,
		DB:       0,
	})

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Error connecting to the redis:", err)
		return nil, err
	}
	fmt.Println("Connected to Redis: ", pong)

	return client, nil
}

package redisconn

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func NewRedisConnection() (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
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

package db

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var Context = context.Background()

func CreateClient (num int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: num,
	})

	return client
}
package redis

import (
	"github.com/redis/go-redis/v9"
)

func RedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})
}

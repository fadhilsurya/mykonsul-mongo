package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func RedisClient(addr string, port string) *redis.Client {
	address := fmt.Sprintf("%s:%v", addr, port)
	return redis.NewClient(&redis.Options{
		Addr:     address, // Redis address
		Password: "",      // No password set
		DB:       0,       // Use default DB
	})
}

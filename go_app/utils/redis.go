package utils

import (
	"chat_with_go/config"

	"github.com/redis/go-redis/v9"
)

func GetRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: config.RedisAddress, // Redis server address
		Password: "", // No password set
		DB: 0, // Use default DB
	})
	return rdb
}

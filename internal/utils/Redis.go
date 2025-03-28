package utils

import (
	"yunshu_Forum/config"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func RedisInit() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
		PoolSize: config.Redis.Poolsize,
	})
	_, err := RedisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func GetRedis() *redis.Client {
	return RedisClient
}

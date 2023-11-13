package database

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func CreateRedisClient(dbNo int) (*redis.Client, error) {
	var redisClient = redis.NewClient(&redis.Options{
		Addr:     "db:6379",
		Password: "",
		DB:       dbNo,
	})
	_, error := redisClient.Ping(Ctx).Result()
	if error != nil {
		return nil, error
	} else {
		return redisClient, nil
	}
}

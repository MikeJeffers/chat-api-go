package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func redisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", REDIS_HOST, REDIS_PORT),
		Password: REDIS_PASSWORD,
		DB:       0,
		Protocol: 3,
	})
}

func storeToken(token string, userId int64, redisClient *redis.Client, c *gin.Context) *redis.StatusCmd {
	return redisClient.SetEx(c, fmt.Sprintf("jwt:%v", userId), token, time.Hour*1)
}

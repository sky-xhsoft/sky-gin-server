// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: token.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/11
// Project Description:
// ----------------------------------------------------------------------------

package token

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"time"
)

const tokenPrefix = "TOKEN:"
const tokenTTL = 24 * time.Hour * 7

func GenerateToken() string {
	return uuid.NewString()
}

func Save(redisClient *redis.Client, token string, info interface{}) error {
	return redisClient.Set(context.Background(), tokenPrefix+token, info, tokenTTL).Err()
}

func Get(redisClient *redis.Client, token string) (interface{}, error) {
	return redisClient.Get(context.Background(), tokenPrefix+token).Result()
}

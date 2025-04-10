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
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sky-xhsoft/sky-gin-server/models"
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

func SaveUser(redis *redis.Client, token string, user *models.SysUser) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return redis.Set(context.Background(), tokenPrefix+token, data, tokenTTL).Err()
}

func GetUser(redis *redis.Client, token string) (*models.SysUser, error) {
	val, err := redis.Get(context.Background(), tokenPrefix+token).Result()
	if err != nil {
		return nil, err
	}
	var user models.SysUser
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func Get(redisClient *redis.Client, token string) (interface{}, error) {
	return redisClient.Get(context.Background(), tokenPrefix+token).Result()
}

func DeleteToken(redis *redis.Client, token string) error {
	return redis.Del(context.Background(), tokenPrefix+token).Err()
}

// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: redis.go
// Author: xhsoftware-skyzhou
// Created On: 2025/1/24
// Project Description:
// ----------------------------------------------------------------------------

package store

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/sky-xhsoft/sky-gin-server/config"
)

func NewRedis(c *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr + ":" + c.Redis.Port,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})

	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		return nil, err
	}
	return client, nil
}

// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: getTableNewId.go
// Author: xhsoftware-skyzhou
// Created On: 2025/5/15
// Project Description:
// ----------------------------------------------------------------------------

package utils

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// 动态生成表单ID
func GetTableId(ctx context.Context, db *gorm.DB, redis *redis.Client, table string) (int64, error) {
	table = strings.ToUpper(table)
	seqKey := fmt.Sprintf("SEQ:%s", table)
	lockKey := fmt.Sprintf("LOCK:%s", table)

	// 1. 尝试从 Redis 获取
	exists, err := redis.Exists(ctx, seqKey).Result()
	if err != nil {
		return 0, err
	}

	if exists == 1 {
		return redis.Incr(ctx, seqKey).Result()
	}

	// 2. Redis 无值，加锁防止并发初始化
	lockValue := fmt.Sprintf("%d", time.Now().UnixNano())
	acquired, err := redis.SetNX(ctx, lockKey, lockValue, 5*time.Second).Result()
	if err != nil {
		return 0, err
	}

	if !acquired {
		// 未获得锁，稍后重试
		time.Sleep(100 * time.Millisecond)
		return GetTableId(ctx, db, redis, table)
	}
	defer redis.Del(ctx, lockKey) // 解锁

	// 3. 读取数据库最大 ID
	var maxID int64
	if err := db.Table(table).Select("MAX(ID)").Scan(&maxID).Error; err != nil {
		return 0, err
	}

	newID := maxID + 1
	if err := redis.Set(ctx, seqKey, newID, 0).Err(); err != nil {
		return 0, err
	}

	return newID, nil
}

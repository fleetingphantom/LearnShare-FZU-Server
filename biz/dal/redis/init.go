package redis

import (
	"LearnShare/config"
	"LearnShare/pkg/errno"
	"LearnShare/pkg/logger"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

// Init 初始化Redis连接
func Init() error {
	if config.Redis == nil {
		logger.Error("Redis配置为空")
		return errno.NewErrNo(errno.InternalServiceErrorCode, "Redis配置为空")
	}
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	_, err := RDB.Ping(context.TODO()).Result()
	if err != nil {
		logger.Errorf("Redis连接测试失败: %v", err)
		return errno.NewErrNo(errno.InternalRedisErrorCode, fmt.Sprintf("Redis连接测试失败: %v", err))
	}
	logger.Info("Redis初始化成功")
	return nil
}

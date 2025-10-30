package redis

import (
	"LearnShare/config"
	"LearnShare/pkg/errno"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func Init() error {
	if config.Redis == nil {
		return errno.NewErrNo(errno.InternalServiceErrorCode, "redis config is nil")
	}
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	_, err := RDB.Ping(context.TODO()).Result()
	if err != nil {
		return errno.NewErrNo(errno.InternalRedisErrorCode, fmt.Sprintf("client.NewRedisClient: ping redis failed: %v", err))
	}
	return nil
}

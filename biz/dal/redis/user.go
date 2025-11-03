package redis

import (
	"LearnShare/pkg/errno"
	"context"
	"fmt"
	"strings"
	"time"
)

func IsKeyExist(ctx context.Context, key string) bool {
	return RDB.Exists(ctx, key).Val() == 1
}

func GetCodeCache(ctx context.Context, key string) (code string, err error) {
	value, err := RDB.Get(ctx, key).Result()
	if err != nil {
		return "", errno.NewErrNo(errno.InternalRedisErrorCode, "write code to cache error:"+err.Error())
	}
	var storedCode string
	parts := strings.Split(value, "_")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid code format, expected 2 parts, got %d", len(parts))
	}
	storedCode = parts[0]
	return storedCode, nil
}
func PutCodeToCache(ctx context.Context, key, code string) error {
	timeNow := time.Now().Unix()
	value := fmt.Sprintf("%s_%d", code, timeNow)
	expiration := 10 * time.Minute
	err := RDB.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return errno.NewErrNo(errno.InternalRedisErrorCode, "write code to cache error:"+err.Error())
	}
	return nil
}

func DeleteCodeCache(ctx context.Context, key string) error {
	err := RDB.Del(ctx, key).Err()
	if err != nil {
		return errno.NewErrNo(errno.InternalRedisErrorCode, "delete code from cache error:"+err.Error())
	}
	return nil
}

func SetBlacklistToken(ctx context.Context, token string) error {
	err := RDB.Set(ctx, token, "blacklisted", 12*time.Hour).Err()
	if err != nil {
		return errno.NewErrNo(errno.InternalRedisErrorCode, "set blacklist token error:"+err.Error())
	}
	return nil
}
func IsBlacklistToken(ctx context.Context, token string) (bool, error) {
	result, err := RDB.Get(ctx, token).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return false, nil
		}
		return false, errno.NewErrNo(errno.InternalRedisErrorCode, "get blacklist token error:"+err.Error())
	}
	if result == "blacklisted" {
		return true, nil
	}
	return false, nil
}

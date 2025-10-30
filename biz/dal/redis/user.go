package redis

import (
	"LearnShare/pkg/errno"
	"context"
	"fmt"
	"time"
)

func IsKeyExist(ctx context.Context, key string) bool {
	return RDB.Exists(ctx, key).Val() == 1
}

func GetCodeCache(ctx context.Context, key string) (string, error) {
	value, err := RDB.Get(ctx, key).Result()
	if err != nil {
		return "", errno.NewErrNo(errno.InternalRedisErrorCode, "write code to cache error:"+err.Error())
	}
	var storedCode, timestampStr string
	_, err = fmt.Sscanf(value, "%s_%s", &storedCode, &timestampStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse code: %v", err)
	}
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

func SetBlacklistToken(ctx context.Context, token string) error {
	err := RDB.Set(ctx, token, "blacklisted", 12*time.Hour).Err()
	if err != nil {
		return errno.NewErrNo(errno.InternalRedisErrorCode, "set blacklist token error:"+err.Error())
	}
	return nil
}

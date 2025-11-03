package redis

import (
	"LearnShare/pkg/errno"
	"context"
	"fmt"
	"strings"
	"time"
)

// IsKeyExist 检查键是否存在
func IsKeyExist(ctx context.Context, key string) bool {
	return RDB.Exists(ctx, key).Val() == 1
}

// GetCodeCache 获取验证码缓存
func GetCodeCache(ctx context.Context, key string) (code string, err error) {
	value, err := RDB.Get(ctx, key).Result()
	if err != nil {
		return "", errno.NewErrNo(errno.InternalRedisErrorCode, "获取验证码缓存失败")
	}
	var storedCode string
	parts := strings.Split(value, "_")
	if len(parts) != 2 {
		return "", errno.NewErrNo(errno.InternalRedisErrorCode, "验证码格式错误")
	}
	storedCode = parts[0]
	return storedCode, nil
}

// PutCodeToCache 将验证码写入缓存
func PutCodeToCache(ctx context.Context, key, code string) error {
	timeNow := time.Now().Unix()
	value := fmt.Sprintf("%s_%d", code, timeNow)
	expiration := 10 * time.Minute
	err := RDB.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return errno.NewErrNo(errno.InternalRedisErrorCode, "写入验证码缓存失败")
	}
	return nil
}

// DeleteCodeCache 删除验证码缓存
func DeleteCodeCache(ctx context.Context, key string) error {
	err := RDB.Del(ctx, key).Err()
	if err != nil {
		return errno.NewErrNo(errno.InternalRedisErrorCode, "删除验证码缓存失败: "+err.Error())
	}
	return nil
}

// SetBlacklistToken 将令牌加入黑名单
func SetBlacklistToken(ctx context.Context, token string) error {
	err := RDB.Set(ctx, token, "blacklisted", 12*time.Hour).Err()
	if err != nil {
		return errno.NewErrNo(errno.InternalRedisErrorCode, "设置令牌黑名单失败: "+err.Error())
	}
	return nil
}

// IsBlacklistToken 检查令牌是否在黑名单中
func IsBlacklistToken(ctx context.Context, token string) (bool, error) {
	result, err := RDB.Get(ctx, token).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return false, nil
		}
		return false, errno.NewErrNo(errno.InternalRedisErrorCode, "获取令牌黑名单状态失败: "+err.Error())
	}
	if result == "blacklisted" {
		return true, nil
	}
	return false, nil
}

package redis

import (
	"LearnShare/pkg/errno"
	"context"
)

func SetPermissionCache(ctx context.Context, key string, value string) error {
	err := RDB.Set(ctx, key, value, 0).Err()
	if err != nil {
		return errno.NewErrNo(errno.InternalRedisErrorCode, "设置权限缓存失败: "+err.Error())
	}
	return nil
}

func GetPermissionCache(ctx context.Context, key string) (string, error) {
	value, err := RDB.Get(ctx, key).Result()
	if err != nil {
		return "", errno.NewErrNo(errno.InternalRedisErrorCode, "获取权限缓存失败: "+err.Error())
	}
	return value, nil
}

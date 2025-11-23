package redis

import "context"

// IsKeyExist 检查键是否存在
func IsKeyExist(ctx context.Context, key string) bool {
	return RDB.Exists(ctx, key).Val() == 1
}

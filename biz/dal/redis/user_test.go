package redis

import (
	"context"
	"testing"

	miniredis "github.com/alicebob/miniredis/v2"
	goRedis "github.com/redis/go-redis/v9"
)

// initTestRedis 初始化一个内存版Redis用于测试
func initTestRedis(t *testing.T) (*miniredis.Miniredis, func()) {
	t.Helper()
	server, err := miniredis.Run()
	if err != nil {
		t.Fatalf("启动 MiniRedis 失败: %v", err)
	}

	RDB = goRedis.NewClient(&goRedis.Options{Addr: server.Addr()})

	cleanup := func() {
		if err := RDB.Close(); err != nil {
			t.Fatalf("关闭 Redis 客户端失败: %v", err)
		}
		server.Close()
	}
	return server, cleanup
}

func TestPutAndGetCodeCache(t *testing.T) {
	_, cleanup := initTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	if err := PutCodeToCache(ctx, "user@example.com", "123456"); err != nil {
		t.Fatalf("写入验证码缓存失败: %v", err)
	}

	code, err := GetCodeCache(ctx, "user@example.com")
	if err != nil {
		t.Fatalf("读取验证码缓存失败: %v", err)
	}
	if code != "123456" {
		t.Fatalf("期望验证码为 123456, 实际为 %s", code)
	}

	if !IsKeyExist(ctx, "user@example.com") {
		t.Fatalf("键应当存在")
	}

	if err := DeleteCodeCache(ctx, "user@example.com"); err != nil {
		t.Fatalf("删除验证码缓存失败: %v", err)
	}
	if IsKeyExist(ctx, "user@example.com") {
		t.Fatalf("删除后键仍然存在")
	}
}

func TestBlacklistToken(t *testing.T) {
	_, cleanup := initTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	if err := SetBlacklistToken(ctx, "token123"); err != nil {
		t.Fatalf("写入黑名单失败: %v", err)
	}

	blacklisted, err := IsBlacklistToken(ctx, "token123")
	if err != nil {
		t.Fatalf("查询黑名单失败: %v", err)
	}
	if !blacklisted {
		t.Fatalf("预期 token 在黑名单中")
	}

	inBlacklist, err := IsBlacklistToken(ctx, "token456")
	if err != nil {
		t.Fatalf("查询黑名单失败: %v", err)
	}
	if inBlacklist {
		t.Fatalf("未写入的 token 被误判为在黑名单中")
	}
}

package service

import (
	"testing"

	"LearnShare/pkg/constants"

	"github.com/cloudwego/hertz/pkg/app"
)

func TestGetUidFormContext(t *testing.T) {
	ctx := app.NewContext(0)
	ctx.Set(constants.ContextUid, int64(42))

	if uid := GetUidFormContext(ctx); uid != 42 {
		t.Fatalf("预期返回 42, 实际为 %d", uid)
	}
}

func TestGetUidFormContextInvalid(t *testing.T) {
	ctx := app.NewContext(0)
	ctx.Set(constants.ContextUid, "invalid")

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("预期出现 panic")
		}
	}()

	_ = GetUidFormContext(ctx)
}

func TestGetUuidFormContext(t *testing.T) {
	ctx := app.NewContext(0)
	ctx.Set(constants.UUID, "uuid-123")

	if val := GetUuidFormContext(ctx); val != "uuid-123" {
		t.Fatalf("预期 UUID 为 uuid-123, 实际为 %s", val)
	}
}

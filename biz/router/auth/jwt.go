package auth

import (
	"LearnShare/biz/dal/redis"
	"LearnShare/biz/middleware"
	"LearnShare/biz/service"
	"LearnShare/pkg/errno"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func AccessTokenAuth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 1. 验证 access-token 是否有效
		if err := middleware.IsAccessTokenAvailable(ctx, c); err != nil {
			fail(c, err)
			return
		}

		// 2. 取出 UUID
		Uuid := service.GetUuidFormContext(c)

		// 3. 判断是否已登出（黑名单）
		ok, err := redis.IsBlacklistToken(ctx, Uuid)
		if err != nil {
			fail(c, err)
			return
		}
		if ok {
			fail(c, errno.NewErrNo(errno.AuthInvalidCode, "令牌已被注销"))
			return
		}

		// 4. 放行
		c.Next(ctx)
	}
}

package auth

import (
	"LearnShare/biz/dal/redis"
	"LearnShare/biz/middleware"
	"LearnShare/biz/pack"
	"LearnShare/biz/service"
	"LearnShare/pkg/errno"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func AccessTokenAuth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 1. 验证 access-token 是否有效
		if !middleware.IsAccessTokenAvailable(ctx, c) {
			fail(c, errno.NewErrNo(errno.AuthInvalidCode, "访问令牌无效"))
			return
		}

		// 2. 取出 UUID
		Uuid := service.GetUuidFormContext(c)
		if Uuid == "" {
			fail(c, errno.NewErrNo(errno.AuthInvalidCode, "上下文中未找到 UUID"))
			return
		}

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

// fail 统一返回错误并终止后续中间件。
func fail(c *app.RequestContext, err error) {
	pack.BuildFailResponse(c, err)
	c.Abort()
}

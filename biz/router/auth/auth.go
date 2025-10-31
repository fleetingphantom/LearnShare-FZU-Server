package auth

import (
	"LearnShare/biz/middleware"
	"LearnShare/biz/pack"
	"LearnShare/pkg/errno"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func Auth() []app.HandlerFunc {
	return append(make([]app.HandlerFunc, 0),
		DoubleTokenAuthFunc(),
	)
}

func DoubleTokenAuthFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if !middleware.IsAccessTokenAvailable(ctx, c) {
			pack.BuildFailResponse(c, errno.AuthInvalid)
			c.Abort()
			return
		}
		c.Next(ctx)
	}
}

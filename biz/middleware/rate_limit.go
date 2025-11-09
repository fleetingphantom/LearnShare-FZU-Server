package middleware

import (
	"LearnShare/biz/dal/redis"
	"LearnShare/biz/pack"
	"LearnShare/pkg/errno"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// EmailRateLimitMiddleware 邮件发送频率限制中间件（同一IP 1分钟间隔）
func EmailRateLimitMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 获取客户端 IP
		clientIP := c.ClientIP()

		// 检查是否存在限制
		if redis.CheckEmailRateLimit(ctx, clientIP) {
			pack.BuildFailResponse(c, errno.NewErrNo(42901, "发送邮件过于频繁，请稍后再试"))
			c.Abort()
			return
		}

		// 设置限制标记，1分钟过期
		err := redis.SetEmailRateLimit(ctx, clientIP)
		if err != nil {
			// Redis 错误不应该阻止请求，记录日志后继续
			// 这里可以添加日志记录
		}

		c.Next(ctx)
	}
}

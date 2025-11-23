package auth

import (
	"LearnShare/biz/dal/redis"
	"LearnShare/biz/middleware"
	"LearnShare/config"
	"LearnShare/pkg/errno"
	"LearnShare/pkg/logger"
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
			fail(c, errno.NewErrNo(errno.ServiceEmailSendLimit, "发送邮件过于频繁，请稍后再试"))
			return
		}

		// 设置限制标记，1分钟过期
		err := redis.SetEmailRateLimit(ctx, clientIP)
		if err != nil {
			logger.Errorf("设置邮件频率限制失败: client_ip=%s error=%v", clientIP, err)
		}

		c.Next(ctx)
	}
}

// TurnstileMiddleware Turnstile 验证中间件
func TurnstileMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 检查是否启用 Turnstile 验证
		if !config.Turnstile.Enabled {
			c.Next(ctx)
			return
		}

		// 从请求头中获取 Turnstile token
		token := c.GetHeader("CF-Turnstile-Token")
		if len(token) == 0 {
			// 尝试从表单或 JSON 中获取
			tokenStr := c.PostForm("cf_turnstile_token")
			if tokenStr == "" {
				// 尝试从 JSON body 中获取
				type TurnstileReq struct {
					CfTurnstileToken string `json:"cf_turnstile_token"`
				}
				var req TurnstileReq
				_ = c.BindJSON(&req)
				tokenStr = req.CfTurnstileToken
			}
			token = []byte(tokenStr)
		}

		if len(token) == 0 {
			fail(c, errno.NewErrNo(errno.TurnstileMissingTokenCode, "缺少 Turnstile 验证 token"))
			return
		}

		// 验证 Turnstile token
		if !middleware.VerifyTurnstile(string(token), c.ClientIP()) {
			fail(c, errno.NewErrNo(errno.TurnstileInvalidTokenCode, "Turnstile 验证失败"))
			return
		}

		c.Next(ctx)
	}
}

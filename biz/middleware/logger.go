package middleware

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// RequestLogger HTTP请求日志中间件
func RequestLogger() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		startTime := time.Now()

		// 处理请求
		c.Next(ctx)

		// 记录请求日志
		latency := time.Since(startTime)
		statusCode := c.Response.StatusCode()
		method := string(c.Method())
		path := string(c.Path())
		clientIP := c.ClientIP()

		hlog.CtxInfof(ctx, "HTTP Request | method=%s path=%s client_ip=%s status=%d latency=%dms",
			method, path, clientIP, statusCode, latency.Milliseconds())
	}
}

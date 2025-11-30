package middleware

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/google/uuid"
)

// RequestLogger HTTP请求日志中间件
func RequestLogger() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 生成或获取请求 ID
		requestID := string(c.GetHeader("X-Request-ID"))
		if requestID == "" {
			requestID = uuid.New().String()
		}
		// 设置请求 ID 到上下文和响应头
		c.Set("request_id", requestID)
		c.Response.Header.Set("X-Request-ID", requestID)

		startTime := time.Now()

		// 处理请求
		c.Next(ctx)

		// 记录请求日志
		latency := time.Since(startTime)
		statusCode := c.Response.StatusCode()
		method := string(c.Method())
		path := string(c.Path())
		clientIP := c.ClientIP()

		hlog.CtxInfof(ctx, "HTTP Request | request_id=%s method=%s path=%s client_ip=%s status=%d latency=%dms",
			requestID, method, path, clientIP, statusCode, latency.Milliseconds())
	}
}

// SlowQueryLogger 慢请求监控中间件
// threshold 参数指定慢请求的阈值（毫秒），默认 1000ms
func SlowQueryLogger(threshold ...int) app.HandlerFunc {
	thresholdMs := 1000
	if len(threshold) > 0 && threshold[0] > 0 {
		thresholdMs = threshold[0]
	}

	return func(ctx context.Context, c *app.RequestContext) {
		startTime := time.Now()

		// 处理请求
		c.Next(ctx)

		// 检查是否为慢请求
		duration := time.Since(startTime)
		if duration.Milliseconds() > int64(thresholdMs) {
			method := string(c.Method())
			path := string(c.Path())
			clientIP := c.ClientIP()
			statusCode := c.Response.StatusCode()
			requestID, _ := c.Get("request_id")

			hlog.CtxWarnf(ctx, "Slow Request | request_id=%v method=%s path=%s client_ip=%s status=%d duration=%dms threshold=%dms",
				requestID, method, path, clientIP, statusCode, duration.Milliseconds(), thresholdMs)
		}
	}
}

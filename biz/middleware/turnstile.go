package middleware

import (
	"LearnShare/biz/pack"
	"LearnShare/config"
	"LearnShare/pkg/errno"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/cloudwego/hertz/pkg/app"
)

// TurnstileResponse Cloudflare Turnstile API 响应结构
type TurnstileResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
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
			pack.BuildFailResponse(c, errno.NewErrNo(40001, "缺少 Turnstile 验证 token"))
			c.Abort()
			return
		}

		// 验证 Turnstile token
		if !verifyTurnstile(string(token), c.ClientIP()) {
			pack.BuildFailResponse(c, errno.NewErrNo(40002, "Turnstile 验证失败"))
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}

// verifyTurnstile 验证 Cloudflare Turnstile token
func verifyTurnstile(token, remoteIP string) bool {
	secretKey := config.Turnstile.SecretKey
	if secretKey == "" {
		return false
	}

	// 构建验证请求
	data := url.Values{}
	data.Set("secret", secretKey)
	data.Set("response", token)
	data.Set("remoteip", remoteIP)

	// 发送验证请求到 Cloudflare
	resp, err := http.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify", data)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	// 解析响应
	var result TurnstileResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return false
	}

	return result.Success
}

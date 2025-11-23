package middleware

import (
	"LearnShare/config"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// TurnstileResponse Cloudflare Turnstile API 响应结构
type TurnstileResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

// VerifyTurnstile 验证 Cloudflare Turnstile token
func VerifyTurnstile(token, remoteIP string) bool {
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

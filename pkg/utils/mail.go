package utils

import (
	"LearnShare/pkg/errno"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"math/big"
	"net/smtp"
	"strconv"

	"LearnShare/config"

	"github.com/jordan-wright/email"
)

// MailSendCode 发送验证码邮件到指定地址。优先使用 config.Smtp 配置，若未初始化则回退到环境变量。
func MailSendCode(to string, code string) error {
	if to == "" {
		return errno.NewErrNo(errno.InternalServiceErrorCode, "收件人邮箱为空")
	}

	host := config.Smtp.Host
	port := strconv.Itoa(config.Smtp.Port)
	user := config.Smtp.User
	pass := config.Smtp.Password
	from := config.Smtp.From
	fromName := config.Smtp.FromName

	if host == "" || port == "" || user == "" || pass == "" || from == "" {
		return errno.NewErrNo(errno.InternalServiceErrorCode, "SMTP 配置不完整，请在 config.yaml 或环境变量中设置 SMTP 配置")
	}

	addr := host + ":" + port

	e := email.NewEmail()
	if fromName != "" {
		e.From = fmt.Sprintf("%s <%s>", fromName, from)
	} else {
		e.From = from
	}
	e.To = []string{to}
	e.Subject = "验证码"
	e.HTML = []byte(fmt.Sprintf("你的验证码为：<h1>%s</h1><p>有效期请以系统设置为准。</p>", code))

	auth := smtp.PlainAuth("", user, pass, host)

	tlsCfg := &tls.Config{ServerName: host}
	if err := e.SendWithTLS(addr, auth, tlsCfg); err == nil {
		return nil
	}

	if err := e.Send(addr, auth); err != nil {
		return errno.NewErrNo(errno.InternalServiceErrorCode, "MailSendCode: 无法发送邮件，请检查 SMTP 配置"+err.Error())
	}
	return nil
}

// GenerateCode 生成指定长度的数字验证码，使用 crypto/rand 保证安全性。
func GenerateCode() (string, error) {
	digits := "0123456789"
	maxInt := big.NewInt(int64(len(digits)))
	b := make([]byte, config.Verify.CodeLength)
	for i := 0; i < config.Verify.CodeLength; i++ {
		rn, err := rand.Int(rand.Reader, maxInt)
		if err != nil {
			return "", errno.NewErrNo(errno.InternalServiceErrorCode, err.Error())
		}
		b[i] = digits[rn.Int64()]
	}
	return string(b), nil
}

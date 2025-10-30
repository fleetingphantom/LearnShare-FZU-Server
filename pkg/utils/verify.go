package utils

import (
	"LearnShare/pkg/errno"
	"regexp"
)

func VerifyUsername(username string) (bool, error) {
	usernameRe := regexp.MustCompile(`^[A-Za-z0-9_]{4,16}$`)
	if !usernameRe.MatchString(username) {
		return false, errno.NewErrNo(errno.ServiceInvalidUsername, "用户名格式不正确，应为4-16位字母、数字或下划线组成")
	}
	return true, nil
}

func VerifyPassword(password string) (bool, error) {
	if !regexp.MustCompile(`^[A-Za-z0-9]{8,20}$`).MatchString(password) {
		return false, errno.NewErrNo(errno.ServiceInvalidPassword, "密码格式不正确，应为8-20位字母和数字组成")
	}

	if !regexp.MustCompile(`[A-Za-z]`).MatchString(password) || !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return false, errno.NewErrNo(errno.ServiceInvalidPassword, "密码必须同时包含字母和数字")
	}
	return true, nil
}

func VerifyEmail(email string) (bool, error) {
	emailRe := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRe.MatchString(email) {
		return false, errno.NewErrNo(errno.ServiceInvalidEmail, "邮箱格式不正确")
	}
	return true, nil
}

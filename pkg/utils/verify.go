package utils

import (
	"LearnShare/pkg/errno"
	"regexp"
)

func VerifyUsername(username string) (bool, error) {
	usernameRe := regexp.MustCompile(`^[A-Za-z0-9_]{4,16}$`)
	if !usernameRe.MatchString(username) {
		return false, errno.UserUsernameFormatInvalidError
	}
	return true, nil
}

func VerifyPassword(password string) (bool, error) {
	if !regexp.MustCompile(`^[A-Za-z0-9]{8,20}$`).MatchString(password) {
		return false, errno.UserPasswordFormatInvalidError
	}

	if !regexp.MustCompile(`[A-Za-z]`).MatchString(password) || !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return false, errno.UserPasswordFormatInvalidError
	}
	return true, nil
}

func VerifyEmail(email string) (bool, error) {
	emailRe := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRe.MatchString(email) {
		return false, errno.UserEmailFormatInvalidError
	}
	return true, nil
}

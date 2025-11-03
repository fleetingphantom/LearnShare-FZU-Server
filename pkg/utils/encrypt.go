package utils

import (
	"LearnShare/pkg/constants"
	"LearnShare/pkg/errno"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(pwd string) (string, error) {
	passwordDigest, err := bcrypt.GenerateFromPassword([]byte(pwd), constants.UserDefaultEncryptPasswordCost)
	if err != nil {
		return "", errno.NewErrNo(errno.InternalServiceErrorCode, fmt.Sprintf("密码加密失败, pwd: %s, err: %v", pwd, err))
	}
	return string(passwordDigest), nil
}

func ComparePassword(passwordDigest, password string) error {
	if bcrypt.CompareHashAndPassword([]byte(passwordDigest), []byte(password)) != nil {
		return errno.UserPasswordIncorrectError
	}
	return nil
}

package errno

import (
	"errors"
	"fmt"
)

type ErrNo struct {
	ErrorCode int64
	ErrorMsg  string
}

func (e ErrNo) Error() string {
	return fmt.Sprintf("[%d] %s", e.ErrorCode, e.ErrorMsg) // 错误信息
}

func NewErrNo(code int64, msg string) ErrNo {
	return ErrNo{
		ErrorCode: code,
		ErrorMsg:  msg, // 错误消息
	}
}

// WithMessage 将替换默认消息为新消息
func (e ErrNo) WithMessage(msg string) ErrNo {
	e.ErrorMsg = msg // 替换为中文消息
	return e
}

// WithError 将在消息后添加错误信息
func (e ErrNo) WithError(err error) ErrNo {
	e.ErrorMsg = e.ErrorMsg + ", " + err.Error() // 添加中文错误信息
	return e
}

// ConvertErr 将错误转换为ErrNo类型
// 默认使用用户服务错误码
func ConvertErr(err error) ErrNo {
	if err == nil {
		return Success
	}
	errno := ErrNo{}
	if errors.As(err, &errno) {
		return errno
	}

	s := InternalServiceError
	s.ErrorMsg = err.Error() // 转换为中文错误消息
	return s
}

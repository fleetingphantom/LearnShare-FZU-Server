package service

import (
	"LearnShare/pkg/constants"
	"LearnShare/pkg/errno"

	"github.com/cloudwego/hertz/pkg/app"
)

func GetUidFormContext(c *app.RequestContext) int64 {
	uid, _ := c.Get(constants.ContextUid)
	userid, err := convertToInt64(uid)
	if err != nil {
		panic(err)
	}

	return userid
}

func convertToInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	default:
		return 0, errno.NewErrNo(errno.InternalServiceErrorCode, "无法转换为int64类型")
	}
}

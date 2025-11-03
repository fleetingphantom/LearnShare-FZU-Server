package pack

import (
	"LearnShare/biz/model/module"
	"LearnShare/pkg/errno"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// SendResponse 发送成功响应
func SendResponse(c *app.RequestContext, data interface{}) {
	c.JSON(consts.StatusOK, data)
}

// SendFailResponse 发送失败响应
func SendFailResponse(c *app.RequestContext, data *module.BaseResp) {
	c.JSON(consts.StatusBadRequest, utils.H{
		"baseResponse": data,
	})
}

// BuildBaseResp 构建基础响应
func BuildBaseResp(err errno.ErrNo) *module.BaseResp {
	return &module.BaseResp{
		Code:    int32(err.ErrorCode),
		Message: err.ErrorMsg,
	}
}

// BuildFailResponse 构建失败响应
func BuildFailResponse(c *app.RequestContext, err error) {
	if err == nil {
		SendFailResponse(c, BuildBaseResp(errno.Success))
		return
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		SendFailResponse(c, BuildBaseResp(e))
		return
	}

	e = errno.NewErrNo(errno.InternalServiceErrorCode, err.Error())
	SendFailResponse(c, BuildBaseResp(e))
}

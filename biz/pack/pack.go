package pack

import (
	"LearnShare/biz/model/module"
	"LearnShare/pkg/errno"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func SendResponse(c *app.RequestContext, data interface{}) {
	c.JSON(consts.StatusOK, data)
}

func SendFailResponse(c *app.RequestContext, data *module.BaseResp) {
	c.JSON(consts.StatusBadRequest, utils.H{
		"base": data,
	})
}

func BuildBaseResp(err errno.ErrNo) *module.BaseResp {
	return &module.BaseResp{
		Code:    int32(err.ErrorCode),
		Message: err.ErrorMsg,
	}
}

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

package auth

import (
	"LearnShare/biz/pack"

	"github.com/cloudwego/hertz/pkg/app"
)

// fail 统一返回错误并终止后续中间件。
func fail(c *app.RequestContext, err error) {
	pack.BuildFailResponse(c, err)
	c.Abort()
}

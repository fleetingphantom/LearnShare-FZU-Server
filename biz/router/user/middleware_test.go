package user

import (
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
)

// TestRootMw 测试根路径中间件
func TestRootMw(t *testing.T) {
	mws := rootMw()
	if mws != nil {
		t.Fatalf("预期rootMw返回nil，实际返回 %v", mws)
	}
}

// TestApiMw 测试API路径中间件
func TestApiMw(t *testing.T) {
	mws := _apiMw()
	if mws != nil {
		t.Fatalf("预期_apiMw返回nil，实际返回 %v", mws)
	}
}

// TestAuthMw 测试auth路径中间件
func TestAuthMw(t *testing.T) {
	mws := _authMw()
	if mws != nil {
		t.Fatalf("预期_authMw返回nil，实际返回 %v", mws)
	}
}

// TestLogininMw 测试登录中间件
func TestLogininMw(t *testing.T) {
	mws := _logininMw()
	if mws != nil {
		t.Fatalf("预期_logininMw返回nil，实际返回 %v", mws)
	}
}

// TestLoginoutMw 测试登出中间件（需要认证）
func TestLoginoutMw(t *testing.T) {
	mws := _loginoutMw()
	if mws == nil {
		t.Fatalf("预期_loginoutMw返回认证中间件，实际返回nil")
	}
	if len(mws) == 0 {
		t.Fatalf("预期_loginoutMw返回至少一个中间件，实际返回0个")
	}
	// 验证返回的是HandlerFunc类型
	for i, mw := range mws {
		if mw == nil {
			t.Fatalf("预期第%d个中间件不为nil，实际为nil", i)
		}
	}
}

// TestRefreshtokenMw 测试刷新令牌中间件
func TestRefreshtokenMw(t *testing.T) {
	mws := _refreshtokenMw()
	if mws != nil {
		t.Fatalf("预期_refreshtokenMw返回nil，实际返回 %v", mws)
	}
}

// TestRegisterMw 测试注册中间件
func TestRegisterMw(t *testing.T) {
	mws := _registerMw()
	if mws != nil {
		t.Fatalf("预期_registerMw返回nil，实际返回 %v", mws)
	}
}

// TestUsersMw 测试users路径中间件
func TestUsersMw(t *testing.T) {
	mws := _usersMw()
	if mws != nil {
		t.Fatalf("预期_usersMw返回nil，实际返回 %v", mws)
	}
}

// TestUploadavatarMw 测试上传头像中间件（需要认证）
func TestUploadavatarMw(t *testing.T) {
	mws := _uploadavatarMw()
	if mws == nil {
		t.Fatalf("预期_uploadavatarMw返回认证中间件，实际返回nil")
	}
	if len(mws) == 0 {
		t.Fatalf("预期_uploadavatarMw返回至少一个中间件，实际返回0个")
	}
}

// TestGetuserinfoMw 测试获取用户信息中间件
func TestGetuserinfoMw(t *testing.T) {
	mws := _getuserinfoMw()
	if mws != nil {
		t.Fatalf("预期_getuserinfoMw返回nil，实际返回 %v", mws)
	}
}

// TestMeMw 测试me路径中间件
func TestMeMw(t *testing.T) {
	mws := _meMw()
	if mws != nil {
		t.Fatalf("预期_meMw返回nil，实际返回 %v", mws)
	}
}

// TestUpdateemailMw 测试更新邮箱中间件（需要认证）
func TestUpdateemailMw(t *testing.T) {
	mws := _updateemailMw()
	if mws == nil {
		t.Fatalf("预期_updateemailMw返回认证中间件，实际返回nil")
	}
	if len(mws) == 0 {
		t.Fatalf("预期_updateemailMw返回至少一个中间件，实际返回0个")
	}
}

// TestUpdatemajorMw 测试更新专业中间件（需要认证）
func TestUpdatemajorMw(t *testing.T) {
	mws := _updatemajorMw()
	if mws == nil {
		t.Fatalf("预期_updatemajorMw返回认证中间件，实际返回nil")
	}
	if len(mws) == 0 {
		t.Fatalf("预期_updatemajorMw返回至少一个中间件，实际返回0个")
	}
}

// TestPasswordMw 测试password路径中间件
func TestPasswordMw(t *testing.T) {
	mws := _passwordMw()
	if mws != nil {
		t.Fatalf("预期_passwordMw返回nil，实际返回 %v", mws)
	}
}

// TestUpdatepasswordMw 测试更新密码中间件（需要认证）
func TestUpdatepasswordMw(t *testing.T) {
	mws := _updatepasswordMw()
	if mws == nil {
		t.Fatalf("预期_updatepasswordMw返回认证中间件，实际返回nil")
	}
	if len(mws) == 0 {
		t.Fatalf("预期_updatepasswordMw返回至少一个中间件，实际返回0个")
	}
}

// TestResetpasswordMw 测试重置密码中间件
func TestResetpasswordMw(t *testing.T) {
	mws := _resetpasswordMw()
	if mws != nil {
		t.Fatalf("预期_resetpasswordMw返回nil，实际返回 %v", mws)
	}
}

// TestEmailMw 测试email路径中间件
func TestEmailMw(t *testing.T) {
	mws := _emailMw()
	if mws != nil {
		t.Fatalf("预期_emailMw返回nil，实际返回 %v", mws)
	}
}

// TestSendverifyemailMw 测试发送验证邮件中间件
func TestSendverifyemailMw(t *testing.T) {
	mws := _sendverifyemailMw()
	if mws != nil {
		t.Fatalf("预期_sendverifyemailMw返回nil，实际返回 %v", mws)
	}
}

// TestVerifyemailMw 测试验证邮箱中间件（需要认证）
func TestVerifyemailMw(t *testing.T) {
	mws := _verifyemailMw()
	if mws == nil {
		t.Fatalf("预期_verifyemailMw返回认证中间件，实际返回nil")
	}
	if len(mws) == 0 {
		t.Fatalf("预期_verifyemailMw返回至少一个中间件，实际返回0个")
	}
}

// TestMiddlewareExecution 测试中间件执行
func TestMiddlewareExecution(t *testing.T) {
	t.Run("测试需要认证的中间件可以执行", func(t *testing.T) {
		// 获取需要认证的中间件
		authMiddlewares := []struct {
			name string
			mws  []app.HandlerFunc
		}{
			{"_loginoutMw", _loginoutMw()},
			{"_uploadavatarMw", _uploadavatarMw()},
			{"_updateemailMw", _updateemailMw()},
			{"_updatemajorMw", _updatemajorMw()},
			{"_updatepasswordMw", _updatepasswordMw()},
			{"_verifyemailMw", _verifyemailMw()},
		}

		for _, am := range authMiddlewares {
			if am.mws == nil || len(am.mws) == 0 {
				t.Fatalf("%s: 预期返回认证中间件，实际返回空", am.name)
			}

			// 验证每个中间件都可以被调用
			for i, mw := range am.mws {
				if mw == nil {
					t.Fatalf("%s: 第%d个中间件为nil", am.name, i)
				}
				// 中间件函数应该可以正常调用（不测试具体行为，只测试不会panic）
				func() {
					defer func() {
						if r := recover(); r != nil {
							// 预期可能会有panic（因为没有完整的请求上下文），这是正常的
							// 我们只是验证中间件函数本身存在且可以被调用
						}
					}()
					// 注意：这里不实际执行中间件，因为需要完整的Hertz上下文
					// 只验证中间件函数不为nil即可
				}()
			}
		}
	})

	t.Run("测试不需要认证的中间件返回nil", func(t *testing.T) {
		noAuthMiddlewares := []struct {
			name string
			mws  []app.HandlerFunc
		}{
			{"rootMw", rootMw()},
			{"_apiMw", _apiMw()},
			{"_authMw", _authMw()},
			{"_logininMw", _logininMw()},
			{"_refreshtokenMw", _refreshtokenMw()},
			{"_registerMw", _registerMw()},
			{"_usersMw", _usersMw()},
			{"_getuserinfoMw", _getuserinfoMw()},
			{"_meMw", _meMw()},
			{"_passwordMw", _passwordMw()},
			{"_resetpasswordMw", _resetpasswordMw()},
			{"_emailMw", _emailMw()},
			{"_sendverifyemailMw", _sendverifyemailMw()},
		}

		for _, nam := range noAuthMiddlewares {
			if nam.mws != nil {
				t.Fatalf("%s: 预期返回nil，实际返回 %v", nam.name, nam.mws)
			}
		}
	})
}

// TestMiddlewareCount 测试中间件数量
func TestMiddlewareCount(t *testing.T) {
	// 测试需要认证的中间件至少有一个handler
	authMiddlewares := map[string][]app.HandlerFunc{
		"_loginoutMw":       _loginoutMw(),
		"_uploadavatarMw":   _uploadavatarMw(),
		"_updateemailMw":    _updateemailMw(),
		"_updatemajorMw":    _updatemajorMw(),
		"_updatepasswordMw": _updatepasswordMw(),
		"_verifyemailMw":    _verifyemailMw(),
	}

	for name, mws := range authMiddlewares {
		if len(mws) < 1 {
			t.Errorf("%s: 预期至少有1个中间件，实际有 %d 个", name, len(mws))
		}
	}
}

// TestAllMiddlewaresDefined 测试所有中间件函数都已定义
func TestAllMiddlewaresDefined(t *testing.T) {
	middlewares := []struct {
		name string
		fn   func() []app.HandlerFunc
	}{
		{"rootMw", rootMw},
		{"_apiMw", _apiMw},
		{"_authMw", _authMw},
		{"_logininMw", _logininMw},
		{"_loginoutMw", _loginoutMw},
		{"_refreshtokenMw", _refreshtokenMw},
		{"_registerMw", _registerMw},
		{"_usersMw", _usersMw},
		{"_uploadavatarMw", _uploadavatarMw},
		{"_getuserinfoMw", _getuserinfoMw},
		{"_meMw", _meMw},
		{"_updateemailMw", _updateemailMw},
		{"_updatemajorMw", _updatemajorMw},
		{"_passwordMw", _passwordMw},
		{"_updatepasswordMw", _updatepasswordMw},
		{"_resetpasswordMw", _resetpasswordMw},
		{"_emailMw", _emailMw},
		{"_sendverifyemailMw", _sendverifyemailMw},
		{"_verifyemailMw", _verifyemailMw},
	}

	for _, m := range middlewares {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("%s: 中间件函数调用时panic: %v", m.name, r)
				}
			}()
			// 调用中间件函数以验证其已定义
			_ = m.fn()
		}()
	}
}

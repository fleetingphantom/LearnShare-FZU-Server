package resource

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

// TestResourcesMw 测试resources路径中间件
func TestResourcesMw(t *testing.T) {
	mws := _resourcesMw()
	if mws != nil {
		t.Fatalf("预期_resourcesMw返回nil，实际返回 %v", mws)
	}
}

// TestUploadResourceMw 测试上传资源中间件
func TestUploadResourceMw(t *testing.T) {
	mws := _uploadresourceMw()
	if mws != nil {
		t.Fatalf("预期_uploadresourceMw返回nil，实际返回 %v", mws)
	}
}

// TestGetResourceMw 测试获取资源中间件
func TestGetResourceMw(t *testing.T) {
	mws := _getresourceMw()
	if mws != nil {
		t.Fatalf("预期_getresourceMw返回nil，实际返回 %v", mws)
	}
}

// TestDownloadResourceMw 测试下载资源中间件
func TestDownloadResourceMw(t *testing.T) {
	mws := _downloadresourceMw()
	if mws != nil {
		t.Fatalf("预期_downloadresourceMw返回nil，实际返回 %v", mws)
	}
}

// TestReportResourceMw 测试举报资源中间件（需要认证）
func TestReportResourceMw(t *testing.T) {
	mws := _reportresourceMw()
	if mws == nil {
		t.Fatalf("预期_reportresourceMw返回认证中间件，实际返回nil")
	}
	if len(mws) == 0 {
		t.Fatalf("预期_reportresourceMw返回至少一个中间件，实际返回0个")
	}
	// 验证返回的是HandlerFunc类型
	for i, mw := range mws {
		if mw == nil {
			t.Fatalf("预期第%d个中间件不为nil，实际为nil", i)
		}
	}
}

// TestResources0Mw 测试resources0路径中间件
func TestResources0Mw(t *testing.T) {
	mws := _resources0Mw()
	if mws != nil {
		t.Fatalf("预期_resources0Mw返回nil，实际返回 %v", mws)
	}
}

// TestSearchResourcesMw 测试搜索资源中间件
func TestSearchResourcesMw(t *testing.T) {
	mws := _searchresourcesMw()
	if mws != nil {
		t.Fatalf("预期_searchresourcesMw返回nil，实际返回 %v", mws)
	}
}

// TestResourceMw 测试resource路径中间件
func TestResourceMw(t *testing.T) {
	mws := _resourceMw()
	if mws != nil {
		t.Fatalf("预期_resourceMw返回nil，实际返回 %v", mws)
	}
}

// TestGetResourceCommentsMw 测试获取资源评论中间件
func TestGetResourceCommentsMw(t *testing.T) {
	mws := _getresourcecommentsMw()
	if mws != nil {
		t.Fatalf("预期_getresourcecommentsMw返回nil，实际返回 %v", mws)
	}
}

// TestResourceCommentsMw 测试resource_comments路径中间件
func TestResourceCommentsMw(t *testing.T) {
	mws := _resource_commentsMw()
	if mws != nil {
		t.Fatalf("预期_resource_commentsMw返回nil，实际返回 %v", mws)
	}
}

// TestSubmitResourceCommentMw 测试提交资源评论中间件（需要认证）
func TestSubmitResourceCommentMw(t *testing.T) {
	mws := _submitresourcecommentMw()
	if mws == nil {
		t.Fatalf("预期_submitresourcecommentMw返回认证中间件，实际返回nil")
	}
	if len(mws) == 0 {
		t.Fatalf("预期_submitresourcecommentMw返回至少一个中间件，实际返回0个")
	}
}

// TestResourceRatingsMw 测试resource_ratings路径中间件
func TestResourceRatingsMw(t *testing.T) {
	mws := _resource_ratingsMw()
	if mws != nil {
		t.Fatalf("预期_resource_ratingsMw返回nil，实际返回 %v", mws)
	}
}

// TestDeleteResourceRatingMw 测试删除资源评分中间件（需要认证）
func TestDeleteResourceRatingMw(t *testing.T) {
	mws := _deleteresourceratingMw()
	if mws == nil {
		t.Fatalf("预期_deleteresourceratingMw返回认证中间件，实际返回nil")
	}
	if len(mws) == 0 {
		t.Fatalf("预期_deleteresourceratingMw返回至少一个中间件，实际返回0个")
	}
}

// TestSubmitResourceRatingMw 测试提交资源评分中间件（需要认证）
func TestSubmitResourceRatingMw(t *testing.T) {
	mws := _submitresourceratingMw()
	if mws == nil {
		t.Fatalf("预期_submitresourceratingMw返回认证中间件，实际返回nil")
	}
	if len(mws) == 0 {
		t.Fatalf("预期_submitresourceratingMw返回至少一个中间件，实际返回0个")
	}
}

// TestDeleteResourceCommentMw 测试删除资源评论中间件（需要认证）
func TestDeleteResourceCommentMw(t *testing.T) {
	mws := _deleteresourcecommentMw()
	if mws == nil {
		t.Fatalf("预期_deleteresourcecommentMw返回认证中间件，实际返回nil")
	}
	if len(mws) == 0 {
		t.Fatalf("预期_deleteresourcecommentMw返回至少一个中间件，实际返回0个")
	}
}

// TestResourcesCommentsMw 测试resources_comments路径中间件
func TestResourcesCommentsMw(t *testing.T) {
	mws := _resources_commentsMw()
	if mws != nil {
		t.Fatalf("预期_resources_commentsMw返回nil，实际返回 %v", mws)
	}
}

// TestReportMw 测试report路径中间件
func TestReportMw(t *testing.T) {
	mws := _reportMw()
	if mws != nil {
		t.Fatalf("预期_reportMw返回nil，实际返回 %v", mws)
	}
}

// TestResources1Mw 测试resources1路径中间件
func TestResources1Mw(t *testing.T) {
	mws := _resources1Mw()
	if mws != nil {
		t.Fatalf("预期_resources1Mw返回nil，实际返回 %v", mws)
	}
}

// TestResourceIdMw 测试resource_id路径中间件
func TestResourceIdMw(t *testing.T) {
	mws := _resource_idMw()
	if mws != nil {
		t.Fatalf("预期_resource_idMw返回nil，实际返回 %v", mws)
	}
}

// TestResourceId0Mw 测试resource_id0路径中间件
func TestResourceId0Mw(t *testing.T) {
	mws := _resource_id0Mw()
	if mws != nil {
		t.Fatalf("预期_resource_id0Mw返回nil，实际返回 %v", mws)
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
			{"_reportresourceMw", _reportresourceMw()},
			{"_submitresourcecommentMw", _submitresourcecommentMw()},
			{"_deleteresourceratingMw", _deleteresourceratingMw()},
			{"_submitresourceratingMw", _submitresourceratingMw()},
			{"_deleteresourcecommentMw", _deleteresourcecommentMw()},
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
			{"_resourcesMw", _resourcesMw()},
			{"_uploadresourceMw", _uploadresourceMw()},
			{"_getresourceMw", _getresourceMw()},
			{"_downloadresourceMw", _downloadresourceMw()},
			{"_resources0Mw", _resources0Mw()},
			{"_searchresourcesMw", _searchresourcesMw()},
			{"_resourceMw", _resourceMw()},
			{"_getresourcecommentsMw", _getresourcecommentsMw()},
			{"_resource_commentsMw", _resource_commentsMw()},
			{"_resource_ratingsMw", _resource_ratingsMw()},
			{"_resources_commentsMw", _resources_commentsMw()},
			{"_reportMw", _reportMw()},
			{"_resources1Mw", _resources1Mw()},
			{"_resource_idMw", _resource_idMw()},
			{"_resource_id0Mw", _resource_id0Mw()},
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
		"_reportresourceMw":        _reportresourceMw(),
		"_submitresourcecommentMw": _submitresourcecommentMw(),
		"_deleteresourceratingMw":  _deleteresourceratingMw(),
		"_submitresourceratingMw":  _submitresourceratingMw(),
		"_deleteresourcecommentMw": _deleteresourcecommentMw(),
	}

	for name, mws := range authMiddlewares {
		if len(mws) < 1 {
			t.Errorf("%s: 预期至少有1个中间件，实际有 %d 个", name, len(mws))
		}
	}
}

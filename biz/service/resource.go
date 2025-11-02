package service

import (
	"LearnShare/biz/dal/db"
	model "LearnShare/biz/model/module"
	"LearnShare/biz/model/resource"
	"LearnShare/pkg/errno"

	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// ResourceService 封装了资源相关的服务
type ResourceService struct {
	ctx context.Context
	c   *app.RequestContext
}

// NewResourceService 创建一个新的 ResourceService
func NewResourceService(ctx context.Context, c *app.RequestContext) *ResourceService {
	return &ResourceService{ctx: ctx, c: c}
}

// SearchResources 执行资源搜索
func (s *ResourceService) SearchResources(req *resource.SearchResourceReq) ([]*model.Resource, int64, error) {
	// 验证搜索关键词长度
	if req.Keyword != nil && *req.Keyword != "" && len(*req.Keyword) > 100 {
		return nil, 0, errno.NewErrNo(errno.ServiceInvalidParameter, "搜索关键词过长")
	}

	// 验证分页参数
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	resources, total, err := db.SearchResources(s.ctx, req.Keyword, req.TagId, req.CourseId, req.SortBy, int(req.PageNum), int(req.PageSize))
	if err != nil {
		return nil, 0, err
	}

	var modelResources []*model.Resource
	for _, r := range resources {
		modelResources = append(modelResources, r.ToResourceModule())
	}

	return modelResources, total, nil
}

// GetResource 执行获取单个资源信息
func (s *ResourceService) GetResource(req *resource.GetResourceReq) (*model.Resource, error) {
	// 验证资源ID
	if req.ResourceId <= 0 {
		return nil, errno.NewErrNo(errno.ServiceInvalidParameter, "资源ID无效")
	}

	resource, err := db.GetResourceByID(s.ctx, req.ResourceId)
	if err != nil {
		return nil, err
	}

	return resource.ToResourceModule(), nil
}

// GetResourceComments 执行获取资源评论列表
func (s *ResourceService) GetResourceComments(req *resource.GetResourceCommentsReq) ([]*model.ResourceComment, int64, error) {
	// 验证资源ID
	if req.ResourceId <= 0 {
		return nil, 0, errno.NewErrNo(errno.ServiceInvalidParameter, "资源ID无效")
	}

	// 验证分页参数
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// 调用数据库层获取评论数据
	comments, total, err := db.GetResourceComments(s.ctx, req.ResourceId, req.SortBy, int(req.PageNum), int(req.PageSize))
	if err != nil {
		return nil, 0, err
	}

	var modelComments []*model.ResourceComment
	for _, comment := range comments {
		modelComments = append(modelComments, comment.ToResourceCommentModule())
	}

	return modelComments, total, nil
}

// SubmitResourceRating 执行提交资源评分
func (s *ResourceService) SubmitResourceRating(req *resource.SubmitResourceRatingReq) (*model.ResourceRating, error) {
	userID := GetUidFormContext(s.c)

	// 验证评分范围
	if req.Rating < 0 || req.Rating > 5 {
		return nil, errno.NewErrNo(errno.ServiceInvalidParameter, "评分必须在0-5之间")
	}

	// 调用数据库层提交评分，使用rating字段
	rating, err := db.SubmitResourceRating(s.ctx, userID, req.ResourceId, req.Rating)
	if err != nil {
		return nil, err
	}

	return rating.ToResourceRatingModule(), nil
}

// SubmitResourceComment 执行提交资源评论
func (s *ResourceService) SubmitResourceComment(req *resource.SubmitResourceCommentReq) (*model.ResourceComment, error) {
	userID := GetUidFormContext(s.c)

	// 验证评论内容
	if req.Content == "" {
		return nil, errno.NewErrNo(errno.ServiceInvalidParameter, "评论内容不能为空")
	}

	if len(req.Content) > 1000 {
		return nil, errno.NewErrNo(errno.ServiceInvalidParameter, "评论内容不能超过1000字符")
	}

	// 处理父评论ID
	var parentID *int64
	if req.IsSetParentId() && req.ParentId != nil && *req.ParentId != 0 {
		parentID = req.ParentId
	}

	// 调用数据库层提交评论
	comment, err := db.SubmitResourceComment(s.ctx, userID, req.ResourceId, req.Content, parentID)
	if err != nil {
		return nil, err
	}

	return comment.ToResourceCommentModule(), nil
}

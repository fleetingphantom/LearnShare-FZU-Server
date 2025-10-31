package service

import (
	"LearnShare/biz/dal/db"
	model "LearnShare/biz/model/module"
	"LearnShare/biz/model/resource"
	"context"
)

// SearchResourcesService 封装了搜索资源的服务
type SearchResourcesService struct {
	ctx context.Context
}

// NewSearchResourcesService 创建一个新的 SearchResourcesService
func NewSearchResourcesService(ctx context.Context) *SearchResourcesService {
	return &SearchResourcesService{ctx: ctx}
}

// SearchResources 执行资源搜索
func (s *SearchResourcesService) SearchResources(req *resource.SearchResourceReq) ([]*model.Resource, int64, error) {
	resources, total, err := db.SearchResources(s.ctx, req.Keyword, req.TagId, req.CourseId, req.SortBy, int(req.PageNum), int(req.PageSize))
	if err != nil {
		return nil, 0, err
	}

	// 将 db.Resource 转换为 model.Resource
	var modelResources []*model.Resource
	for _, r := range resources {
		var tags []*model.ResourceTag
		for _, t := range r.Tags {
			tags = append(tags, &model.ResourceTag{
				TagId:   t.TagID,
				TagName: t.TagName,
			})
		}

		modelResources = append(modelResources, &model.Resource{
			ResourceId:    r.ResourceID,
			Title:         r.Title,
			Description:   &r.Description,
			FilePath:      r.FilePath,
			FileType:      r.FileType,
			FileSize:      r.FileSize,
			UploaderId:    r.UploaderID,
			CourseId:      r.CourseID,
			DownloadCount: r.DownloadCount,
			AverageRating: r.AverageRating,
			RatingCount:   r.RatingCount,
			Status:        r.Status,
			CreatedAt:     r.CreatedAt.Unix(),
			Tags:          tags,
		})
	}

	return modelResources, total, nil
}

// GetResourceService 封装了获取单个资源信息的服务
type GetResourceService struct {
	ctx context.Context
}

// NewGetResourceService 创建一个新的 GetResourceService
func NewGetResourceService(ctx context.Context) *GetResourceService {
	return &GetResourceService{ctx: ctx}
}

// GetResource 执行获取单个资源信息
func (s *GetResourceService) GetResource(req *resource.GetResourceReq) (*model.Resource, error) {
	resource, err := db.GetResourceByID(s.ctx, req.ResourceId)
	if err != nil {
		return nil, err
	}

	// 将 db.Resource 转换为 model.Resource
	var tags []*model.ResourceTag
	for _, t := range resource.Tags {
		tags = append(tags, &model.ResourceTag{
			TagId:   t.TagID,
			TagName: t.TagName,
		})
	}

	return &model.Resource{
		ResourceId:    resource.ResourceID,
		Title:         resource.Title,
		Description:   &resource.Description,
		FilePath:      resource.FilePath,
		FileType:      resource.FileType,
		FileSize:      resource.FileSize,
		UploaderId:    resource.UploaderID,
		CourseId:      resource.CourseID,
		DownloadCount: resource.DownloadCount,
		AverageRating: resource.AverageRating,
		RatingCount:   resource.RatingCount,
		Status:        resource.Status,
		CreatedAt:     resource.CreatedAt.Unix(),
		Tags:          tags,
	}, nil
}
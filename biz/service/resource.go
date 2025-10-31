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
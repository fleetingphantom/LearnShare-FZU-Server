package db

import (
	"context"
)

func SearchResources(ctx context.Context, keyword *string, tagID, courseID *int64, sortBy *string, pageNum, pageSize int) ([]*Resource, int64, error) {
	var resources []*Resource
	var total int64

	// 参数验证
	if pageNum <= 0 || pageSize <= 0 {
		return []*Resource{}, 0, nil
	}

	db := DB.WithContext(ctx)

	if keyword != nil && *keyword != "" {
		db = db.Where("title LIKE ? OR description LIKE ?", "%"+*keyword+"%", "%"+*keyword+"%")
	}

	if courseID != nil {
		db = db.Where("course_id = ?", *courseID)
	}

	if tagID != nil {
		db = db.Joins("JOIN resource_tag_mapping ON resource_tag_mapping.resource_id = resource.resource_id").
			Where("resource_tag_mapping.tag_id = ?", *tagID)
	}

	switch {
	case sortBy != nil && *sortBy == "hot":
		db = db.Order("download_count desc")
	case sortBy != nil && *sortBy == "rating":
		db = db.Order("average_rating desc")
	default:
		db = db.Order("created_at desc")
	}

	err := db.Model(&Resource{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Preload("Tags").Find(&resources).Error
	if err != nil {
		return nil, 0, err
	}

	return resources, total, nil
}

// GetResourceByID 根据资源ID获取单个资源信息
func GetResourceByID(ctx context.Context, resourceID int64) (*Resource, error) {
	var resource Resource

	err := DB.WithContext(ctx).
		Preload("Tags").
		Where("resource_id = ?", resourceID).
		First(&resource).Error

	if err != nil {
		return nil, err
	}

	return &resource, nil
}

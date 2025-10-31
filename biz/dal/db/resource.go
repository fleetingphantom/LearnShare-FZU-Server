package db

import (
	"context"
)

func SearchResources(ctx context.Context, keyword *string, tagID, courseID *int64, sortBy *string, pageNum, pageSize int) ([]*Resource, int64, error) {
	var resources []*Resource
	var total int64

	db := DB.WithContext(ctx)

	if keyword != nil && *keyword != "" {
		db = db.Where("title LIKE ? OR description LIKE ?", "%"+*keyword+"%", "%"+*keyword+"%")
	}

	if courseID != nil {
		db = db.Where("course_id = ?", *courseID)
	}

	if tagID != nil {
		db = db.Joins("JOIN resource_tag_mappings ON resource_tag_mappings.resource_id = resources.resource_id").
			Where("resource_tag_mappings.tag_id = ?", *tagID)
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
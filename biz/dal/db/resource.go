package db

import (
	"LearnShare/pkg/constants"
	"LearnShare/pkg/errno"
	"context"
	"errors"

	"gorm.io/gorm"
)

func SearchResources(ctx context.Context, keyword *string, tagID, courseID *int64, sortBy *string, pageNum, pageSize int) ([]*Resource, int64, error) {
	var resources []*Resource
	var total int64

	db := DB.WithContext(ctx).Table(constants.ResourceTableName)

	if keyword != nil && *keyword != "" {
		db = db.Where("resource_name LIKE ? OR description LIKE ?", "%"+*keyword+"%", "%"+*keyword+"%")
	}

	if courseID != nil {
		db = db.Where("course_id = ?", *courseID)
	}

	if tagID != nil {
		db = db.Joins("JOIN "+constants.ResourceTagMappingTableName+" ON "+constants.ResourceTagMappingTableName+".resource_id = "+constants.ResourceTableName+".resource_id").
			Where(constants.ResourceTagMappingTableName+".tag_id = ?", *tagID)
	}

	switch {
	case sortBy != nil && *sortBy == "hot":
		db = db.Order("download_count desc")
	case sortBy != nil && *sortBy == "rating":
		db = db.Order("average_rating desc")
	default:
		db = db.Order("created_at desc")
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "统计资源数量失败: "+err.Error())
	}

	err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Preload("Tags").Find(&resources).Error
	if err != nil {
		return nil, 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询资源列表失败: "+err.Error())
	}

	return resources, total, nil
}

// GetResourceByID 根据资源ID获取单个资源信息
func GetResourceByID(ctx context.Context, resourceID int64) (*Resource, error) {
	var resource Resource

	err := DB.WithContext(ctx).Table(constants.ResourceTableName).
		Preload("Tags").
		Where("resource_id = ?", resourceID).
		First(&resource).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "记录未找到")
		}
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "根据ID查询资源失败: "+err.Error())
	}

	return &resource, nil
}

// GetResourceComments 获取资源评论列表
func GetResourceComments(ctx context.Context, resourceID int64, sortBy *string, pageNum, pageSize int) ([]*ResourceComment, int64, error) {
	var comments []*ResourceComment
	var total int64

	db := DB.WithContext(ctx).Table(constants.ResourceCommentTableName).
		Preload("User").
		Where("resource_id = ?", resourceID).
		Where("is_visible = ?", true).
		Where("status = ?", "normal")

	// 根据排序参数进行排序
	if sortBy != nil {
		switch *sortBy {
		case "latest":
			db = db.Order("created_at DESC")
		case "hottest":
			db = db.Order("likes DESC, created_at DESC")
		default:
			db = db.Order("created_at DESC")
		}
	} else {
		db = db.Order("created_at DESC")
	}

	// 获取总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "统计资源评论数量失败: "+err.Error())
	}

	// 获取分页数据
	err = db.Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&comments).Error
	if err != nil {
		return nil, 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询资源评论列表失败: "+err.Error())
	}

	return comments, total, nil
}

// SubmitResourceRating 提交资源评分
func SubmitResourceRating(ctx context.Context, userID, resourceID int64, recommendation float64) (*ResourceRating, error) {
	// 开始事务
	tx := DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查是否已经评分过
	var existingRating ResourceRating
	err := tx.Table(constants.ResourceRatingTableName).Where("user_id = ? AND resource_id = ?", userID, resourceID).First(&existingRating).Error

	var rating *ResourceRating

	if err == nil {
		// 更新现有评分
		existingRating.Recommendation = recommendation
		existingRating.IsVisible = true // 确保在重新评分时，记录是可见的
		err = tx.Table(constants.ResourceRatingTableName).Save(&existingRating).Error
		rating = &existingRating
	} else {
		// 创建新评分
		rating = &ResourceRating{
			UserID:         userID,
			ResourceID:     resourceID,
			Recommendation: recommendation,
			IsVisible:      true,
		}
		err = tx.Table(constants.ResourceRatingTableName).Create(rating).Error
	}

	if err != nil {
		tx.Rollback()
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交资源评分失败: "+err.Error())
	}

	// 重新计算资源的平均评分
	var avgResult struct {
		AverageRating float64 `gorm:"column:average_rating"`
		RatingCount   int64   `gorm:"column:rating_count"`
	}

	err = tx.Table(constants.ResourceRatingTableName).
		Select("AVG(recommendation) as average_rating, COUNT(*) as rating_count").
		Where("resource_id = ? AND is_visible = ?", resourceID, true).
		Scan(&avgResult).Error

	if err != nil {
		tx.Rollback()
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "计算资源平均评分失败: "+err.Error())
	}

	// 更新资源的评分信息
	err = tx.Table(constants.ResourceTableName).
		Where("resource_id = ?", resourceID).
		Updates(map[string]interface{}{
			"average_rating": avgResult.AverageRating,
			"rating_count":   avgResult.RatingCount,
		}).Error

	if err != nil {
		tx.Rollback()
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新资源评分信息失败: "+err.Error())
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交评分事务失败: "+err.Error())
	}

	return rating, nil
}

// SubmitResourceComment 提交资源评论
func SubmitResourceComment(ctx context.Context, userID, resourceID int64, content string, parentID *int64) (*ResourceComment, error) {
	// 开始事务
	tx := DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建评论
	comment := &ResourceComment{
		UserID:     userID,
		ResourceID: resourceID,
		Content:    content,
		ParentID:   parentID,
		Likes:      0,
		IsVisible:  true,
		Status:     "normal",
	}

	// 保存评论
	err := tx.Table(constants.ResourceCommentTableName).Create(comment).Error
	if err != nil {
		tx.Rollback()
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "保存资源评论失败: "+err.Error())
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交评论事务失败: "+err.Error())
	}

	// 预加载用户信息
	err2 := DB.WithContext(ctx).Table(constants.ResourceCommentTableName).Preload("User").First(comment, comment.CommentID).Error
	if err2 != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "预加载评论用户信息失败: "+err2.Error())
	}

	return comment, nil
}

// DeleteResourceRating 删除资源评分
func DeleteResourceRating(ctx context.Context, ratingID, userID int64) error {
	// 开始事务
	tx := DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查询评分记录，确保用户只能删除自己的评分
	var rating ResourceRating
	err := tx.Table(constants.ResourceRatingTableName).Where("rating_id = ? AND user_id = ?", ratingID, userID).First(&rating).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.NewErrNo(errno.InternalDatabaseErrorCode, "未找到评分记录或无权删除")
		}
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询评分记录失败: "+err.Error())
	}

	// 直接从数据库中删除评分记录
	err = tx.Table(constants.ResourceRatingTableName).Delete(&rating).Error
	if err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除评分失败: "+err.Error())
	}

	// 重新计算资源的平均评分
	var avgResult struct {
		AverageRating float64 `gorm:"column:average_rating"`
		RatingCount   int64   `gorm:"column:rating_count"`
	}

	err = tx.Table(constants.ResourceRatingTableName).
		Select("AVG(recommendation) as average_rating, COUNT(*) as rating_count").
		Where("resource_id = ?", rating.ResourceID).
		Scan(&avgResult).Error

	if err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "计算资源平均评分失败: "+err.Error())
	}

	// 更新资源的评分信息
	err = tx.Table(constants.ResourceTableName).
		Where("resource_id = ?", rating.ResourceID).
		Updates(map[string]interface{}{
			"average_rating": avgResult.AverageRating,
			"rating_count":   avgResult.RatingCount,
		}).Error

	if err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新资源评分信息失败: "+err.Error())
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交删除评分事务失败: "+err.Error())
	}

	return nil
}

// DeleteResourceComment 删除资源评论
func DeleteResourceComment(ctx context.Context, commentID, userID int64) error {
	// 开始事务
	tx := DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 直接删除评论，确保用户只能删除自己的评论
	result := tx.Table(constants.ResourceCommentTableName).Where("comment_id = ? AND user_id = ?", commentID, userID).Delete(&ResourceComment{})
	if result.Error != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除评论失败: "+result.Error.Error())
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "未找到评论或无权删除")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交删除评论事务失败: "+err.Error())
	}

	return nil
}

// CreateReview 创建一个新的举报（审核）
func CreateReview(ctx context.Context, creatorID int64, targetID int64, targetType, reason string) error {
	review := &Review{
		TargetID:   targetID,
		TargetType: targetType,
		Reason:     reason,
		ReviewerID: &creatorID, // 使用 creatorID
	}

	result := DB.WithContext(ctx).Table(constants.ReviewTableName).Create(review)
	if result.Error != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "创建举报失败: "+result.Error.Error())
	}

	return nil
}

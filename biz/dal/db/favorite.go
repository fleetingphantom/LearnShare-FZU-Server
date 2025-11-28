package db

import (
	"LearnShare/pkg/constants"
	"LearnShare/pkg/errno"
	"context"
	"errors"

	"gorm.io/gorm"
)

// GetFavoritesByUser 获取用户收藏列表
func GetFavoritesByUser(ctx context.Context, userID int64, targetType string) ([]*Favorite, error) {
	var favorites []*Favorite

	query := DB.WithContext(ctx).Table(constants.FavoriteTableName).Where("user_id = ?", userID)

	if targetType != "" {
		query = query.Where("target_type = ?", targetType)
	}

	err := query.Order("created_at DESC").Find(&favorites).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询收藏列表失败: "+err.Error())
	}

	return favorites, nil
}

// GetFavoritesWithDetails 获取用户收藏列表及其详细信息（优化：批量查询避免N+1）
func GetFavoritesWithDetails(ctx context.Context, userID int64, targetType string) ([]*FavoriteDetail, error) {
	var details []*FavoriteDetail

	// 先获取所有收藏记录
	favorites, err := GetFavoritesByUser(ctx, userID, targetType)
	if err != nil {
		return nil, err
	}

	if len(favorites) == 0 {
		return []*FavoriteDetail{}, nil
	}

	// 根据 target_type 批量获取关联数据
	switch targetType {
	case "course":
		// 批量获取课程信息
		courseIDs := make([]int64, 0, len(favorites))
		for _, f := range favorites {
			courseIDs = append(courseIDs, f.TargetID)
		}

		var courses []*Course
		err := DB.WithContext(ctx).Table(constants.CourseTableName).
			Where("course_id IN (?)", courseIDs).
			Find(&courses).Error
		if err != nil {
			return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "批量查询课程失败: "+err.Error())
		}

		// 建立课程ID到课程信息的映射
		courseMap := make(map[int64]*Course)
		for _, c := range courses {
			courseMap[c.CourseID] = c
		}

		// 组装结果
		for _, f := range favorites {
			if course, exists := courseMap[f.TargetID]; exists {
				details = append(details, &FavoriteDetail{
					Favorite: f,
					Target:   course,
				})
			}
		}

	case "resource":
		// 批量获取资源信息
		resourceIDs := make([]int64, 0, len(favorites))
		for _, f := range favorites {
			resourceIDs = append(resourceIDs, f.TargetID)
		}

		var resources []*Resource
		err := DB.WithContext(ctx).Table(constants.ResourceTableName).
			Preload("Tags").
			Where("resource_id IN (?)", resourceIDs).
			Find(&resources).Error
		if err != nil {
			return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "批量查询资源失败: "+err.Error())
		}

		// 建立资源ID到资源信息的映射
		resourceMap := make(map[int64]*Resource)
		for _, r := range resources {
			resourceMap[r.ResourceID] = r
		}

		// 组装结果
		for _, f := range favorites {
			if resource, exists := resourceMap[f.TargetID]; exists {
				details = append(details, &FavoriteDetail{
					Favorite: f,
					Target:   resource,
				})
			}
		}

	default:
		// 对于其他类型，只返回收藏信息
		for _, f := range favorites {
			details = append(details, &FavoriteDetail{
				Favorite: f,
				Target:   nil,
			})
		}
	}

	return details, nil
}

// FavoriteDetail 包含收藏信息和目标对象的详细信息
type FavoriteDetail struct {
	Favorite *Favorite
	Target   interface{} // 可以是 *Course, *Resource 等
}

// AddFavorite 添加收藏
func AddFavorite(ctx context.Context, userID, targetID int64, targetType string) (*Favorite, error) {
	// 检查是否已经收藏
	var count int64
	err := DB.WithContext(ctx).Table(constants.FavoriteTableName).
		Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).
		Count(&count).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "检查收藏状态失败: "+err.Error())
	}

	if count > 0 {
		return nil, errno.NewErrNo(errno.ParamVerifyErrorCode, "已经收藏过了")
	}

	// 插入收藏记录
	favorite := &Favorite{
		UserID:     userID,
		TargetID:   targetID,
		TargetType: targetType,
	}

	err = DB.WithContext(ctx).Table(constants.FavoriteTableName).Create(favorite).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "添加收藏失败: "+err.Error())
	}

	return favorite, nil
}

// AddFavoriteAsync 添加收藏（异步）
func AddFavoriteAsync(ctx context.Context, userID, targetID int64, targetType string) (chan *Favorite, chan error) {
	favChan := make(chan *Favorite, 1)
	errChan := make(chan error, 1)

	pool := GetAsyncPool()
	pool.Submit(func() error {
		favorite, err := AddFavorite(ctx, userID, targetID, targetType)
		if err != nil {
			errChan <- err
			close(favChan)
			close(errChan)
			return err
		}
		favChan <- favorite
		close(favChan)
		close(errChan)
		return nil
	})

	return favChan, errChan
}

// RemoveFavorite 删除收藏
func RemoveFavorite(ctx context.Context, favoriteID int64, userID int64) error {
	result := DB.WithContext(ctx).Table(constants.FavoriteTableName).
		Where("favorite_id = ? AND user_id = ?", favoriteID, userID).
		Delete(&Favorite{})

	if result.Error != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除收藏失败: "+result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return errno.NewErrNo(errno.ParamVerifyErrorCode, "收藏不存在或无权删除")
	}

	return nil
}

// RemoveFavoriteAsync 删除收藏（异步）
func RemoveFavoriteAsync(ctx context.Context, favoriteID int64, userID int64) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return RemoveFavorite(ctx, favoriteID, userID)
	})
}

// IsFavorited 检查是否已收藏
func IsFavorited(ctx context.Context, userID, targetID int64, targetType string) (bool, error) {
	var count int64
	err := DB.WithContext(ctx).Table(constants.FavoriteTableName).
		Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).
		Count(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, errno.NewErrNo(errno.InternalDatabaseErrorCode, "检查收藏状态失败: "+err.Error())
	}

	return count > 0, nil
}

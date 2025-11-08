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

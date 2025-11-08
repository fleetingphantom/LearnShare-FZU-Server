package service

import (
	"LearnShare/biz/dal/db"
	"LearnShare/biz/model/favorite"
	"LearnShare/biz/model/module"
	"LearnShare/pkg/errno"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

type FavoriteService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewFavoriteService(ctx context.Context, c *app.RequestContext) *FavoriteService {
	return &FavoriteService{ctx: ctx, c: c}
}

// GetFavorites 获取收藏列表
func (s *FavoriteService) GetFavorites(req *favorite.GetFavoriteReq) ([]*module.Favorite, error) {
	// 从上下文获取用户ID
	userID := GetUidFormContext(s.c)

	// 获取收藏列表
	favorites, err := db.GetFavoritesByUser(s.ctx, userID, req.Type)
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "获取收藏列表失败: "+err.Error())
	}

	// 转换为模块类型
	var favoriteModules []*module.Favorite
	for _, f := range favorites {
		favoriteModules = append(favoriteModules, f.ToFavoriteModule())
	}

	return favoriteModules, nil
}

// AddFavorite 添加收藏
func (s *FavoriteService) AddFavorite(req *favorite.AddFavoriteReq) (*module.Favorite, error) {
	// 从上下文获取用户ID
	userID := GetUidFormContext(s.c)

	// 添加收藏
	favChan, errChan := db.AddFavoriteAsync(s.ctx, userID, req.TargetID, req.TargetType)

	select {
	case err := <-errChan:
		if err != nil {
			return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "添加收藏失败: "+err.Error())
		}
	case fav := <-favChan:
		return fav.ToFavoriteModule(), nil
	}

	return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "添加收藏失败")
}

// RemoveFavorite 删除收藏
func (s *FavoriteService) RemoveFavorite(req *favorite.RemoveFavoriteReq) error {
	// 从上下文获取用户ID
	userID := GetUidFormContext(s.c)

	// 删除收藏
	errChan := db.RemoveFavoriteAsync(s.ctx, req.FavoriteID, userID)
	if err := <-errChan; err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除收藏失败: "+err.Error())
	}

	return nil
}

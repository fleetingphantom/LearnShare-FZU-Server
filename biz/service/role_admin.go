package service

import (
	"LearnShare/biz/dal/db"
	"LearnShare/biz/dal/redis"
	"LearnShare/biz/model/module"
	"LearnShare/biz/model/user"
	"context"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

type RoleAdminService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewRoleAdminService(ctx context.Context, c *app.RequestContext) *RoleAdminService {
	return &RoleAdminService{ctx: ctx, c: c}
}

func (s *RoleAdminService) GetPermissionList() ([]*module.Permission, error) {
	permissions, err := db.GetAllPermissions(s.ctx)
	if err != nil {
		return nil, err
	}

	var permissionList []*module.Permission
	for _, perm := range permissions {
		permissionList = append(permissionList, perm.ToPermissionModule())
	}

	return permissionList, nil
}

func (s *RoleAdminService) GetRoleList() ([]*module.Role, error) {
	roles, err := db.GetAllRoles(s.ctx)
	if err != nil {
		return nil, err
	}

	var roleList []*module.Role
	for _, role := range roles {
		roleList = append(roleList, role.ToRoleModule())
	}

	return roleList, nil
}

func (s *RoleAdminService) AddRole(req *user.AddRoleReq) (int64, error) {
	// 创建角色
	roleID, err := db.CreateRole(s.ctx, req.RoleName, req.PermissionIds)
	if err != nil {
		return 0, err
	}

	return roleID, nil
}

func (s *RoleAdminService) GetRolePermissions(roleID int64) ([]string, error) {
	var (
		permissions []string
		err         error
	)
	// 优先从缓存中获取
	if redis.IsKeyExist(s.ctx, "role_permissions_"+strconv.FormatInt(roleID, 10)) {
		data, err := redis.GetPermissionCache(s.ctx, "role_permissions_"+strconv.FormatInt(roleID, 10))
		if err != nil {
			return nil, err
		}
		permissions = deserializePermissions(data)
		return permissions, nil

	} else {
		permissions, err = db.GetRolePermissions(s.ctx, roleID)
		if err != nil {
			return nil, err
		}
		err = redis.SetPermissionCache(s.ctx, "role_permissions_"+strconv.FormatInt(roleID, 10), serializePermissions(permissions))
		if err != nil {
			return nil, err
		}
	}
	return permissions, nil
}

func serializePermissions(permissions []string) string {
	result := ""
	for i, perm := range permissions {
		result += perm
		if i < len(permissions)-1 {
			result += ","
		}
	}
	return result
}

func deserializePermissions(data string) []string {
	if data == "" {
		return []string{}
	}
	return strings.Split(data, ",")

}

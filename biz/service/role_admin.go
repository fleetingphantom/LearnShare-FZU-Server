package service

import (
	"LearnShare/biz/dal/db"
	"LearnShare/biz/model/module"
	"LearnShare/biz/model/user"
	"context"

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

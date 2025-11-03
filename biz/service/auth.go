package service

import (
	"context"

	"LearnShare/biz/dal/db"
)

const (
	// RoleSuperAdminID 表示超级管理员，对应 config/sql/init.sql 中 role_id = 1
	RoleSuperAdminID int64 = 1
	// RoleAuditorID 表示审核员，对应 config/sql/init.sql 中 role_id = 3
	RoleAuditorID int64 = 3
)

// HasAuditorPermissions 判断指定用户是否具备审核员角色要求的权限集合
func HasAuditorPermissions(ctx context.Context, userID int64, requiredPermissions ...string) (bool, error) {
	allowed, _, err := verifyAuditorPermissions(ctx, userID, requiredPermissions...)
	if err != nil {
		return false, err
	}
	return allowed, nil
}

// HasSuperAdminPermissions 判断指定用户是否具备超级管理员权限，内部先复用审核员权限校验
func HasSuperAdminPermissions(ctx context.Context, userID int64, requiredPermissions ...string) (bool, error) {
	allowed, user, err := verifyAuditorPermissions(ctx, userID, requiredPermissions...)
	if err != nil || !allowed {
		return allowed, err
	}
	return user.RoleID == RoleSuperAdminID, nil
}

// verifyAuditorPermissions 核心逻辑：加载用户角色与权限，并判断是否满足审核员要求
func verifyAuditorPermissions(ctx context.Context, userID int64, requiredPermissions ...string) (bool, *db.User, error) {
	user, err := db.GetUserByID(ctx, userID)
	if err != nil {
		return false, nil, err
	}

	if user.RoleID != RoleAuditorID && user.RoleID != RoleSuperAdminID {
		return false, user, nil
	}

	// 无必需权限时，只要身份满足即可通过
	if len(requiredPermissions) == 0 {
		return true, user, nil
	}

	perms, err := db.GetRolePermissions(ctx, user.RoleID)
	if err != nil {
		return false, user, err
	}

	permSet := make(map[string]struct{}, len(perms))
	for _, perm := range perms {
		if perm == "" {
			continue
		}
		permSet[perm] = struct{}{}
	}

	for _, need := range requiredPermissions {
		if need == "" {
			continue
		}
		if _, ok := permSet[need]; !ok {
			return false, user, nil
		}
	}

	return true, user, nil
}

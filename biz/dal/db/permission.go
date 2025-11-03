package db

import (
	"context"

	"LearnShare/pkg/constants"
	"LearnShare/pkg/errno"
)

// GetRolePermissions 从数据库查询指定角色拥有的全部权限标识
func GetRolePermissions(ctx context.Context, roleID int64) ([]string, error) {
	var permissions []string
	err := DB.WithContext(ctx).
		Table(constants.RolePermissionTableName+" AS rp").
		Joins("JOIN "+constants.PermissionTableName+" AS p ON rp.permission_id = p.permission_id").
		Where("rp.role_id = ?", roleID).
		Pluck("p.permission_name", &permissions).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询角色权限失败: "+err.Error())
	}
	return permissions, nil
}

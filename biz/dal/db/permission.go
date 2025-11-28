package db

import (
	"LearnShare/pkg/constants"
	"LearnShare/pkg/errno"
	"context"
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

// GetAllPermissions 获取所有权限列表
func GetAllPermissions(ctx context.Context) ([]*Permission, error) {
	var permissions []*Permission
	err := DB.WithContext(ctx).Table(constants.PermissionTableName).Find(&permissions).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询权限列表失败: "+err.Error())
	}

	return permissions, nil
}

// GetAllRoles 获取所有角色列表（包含权限信息）
func GetAllRoles(ctx context.Context) ([]*RoleWithPermissions, error) {
	var roles []Role
	err := DB.WithContext(ctx).Table(constants.RoleTableName).Find(&roles).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询角色列表失败: "+err.Error())
	}

	if len(roles) == 0 {
		return []*RoleWithPermissions{}, nil
	}

	// 收集所有角色ID（用于限制 JOIN 查询范围，保持结果顺序由 roles 决定）
	roleIDs := make([]int64, len(roles))
	for i, r := range roles {
		roleIDs[i] = r.RoleID
	}

	var rows []RolePermRow
	err = DB.WithContext(ctx).
		Table(constants.RoleTableName+" AS r").
		Select("r.role_id, p.permission_id, p.permission_name, p.description AS permission_description").
		Joins("LEFT JOIN "+constants.RolePermissionTableName+" rp ON r.role_id = rp.role_id").
		Joins("LEFT JOIN "+constants.PermissionTableName+" p ON rp.permission_id = p.permission_id").
		Where("r.role_id IN ?", roleIDs).
		Scan(&rows).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询角色权限失败: "+err.Error())
	}

	// 将查询结果按角色聚合
	rolePermMap := make(map[int64][]Permission, len(roles))
	for _, r := range rows {
		if r.PermissionID == nil {
			// 角色可能没有任何权限关联，跳过
			continue
		}
		perm := Permission{
			PermissionID:   *r.PermissionID,
			PermissionName: "",
			Description:    "",
		}
		if r.PermissionName != nil {
			perm.PermissionName = *r.PermissionName
		}
		if r.PermissionDescription != nil {
			perm.Description = *r.PermissionDescription
		}
		rolePermMap[r.RoleID] = append(rolePermMap[r.RoleID], perm)
	}

	// 按原始 roles 顺序组装最终结果
	result := make([]*RoleWithPermissions, len(roles))
	for i, r := range roles {
		perms := rolePermMap[r.RoleID]
		if perms == nil {
			perms = []Permission{}
		}
		result[i] = &RoleWithPermissions{
			RoleID:      r.RoleID,
			RoleName:    r.RoleName,
			Description: r.Description,
			Permissions: perms,
		}
	}

	return result, nil
}

// CreateRole 创建新角色
func CreateRole(ctx context.Context, roleName string, permissionIds []int64) (int64, error) {
	// 开启事务
	tx := DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建角色
	role := &Role{
		RoleName:    roleName,
		Description: "",
	}
	err := tx.Table(constants.RoleTableName).Create(role).Error
	if err != nil {
		tx.Rollback()
		return 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "创建角色失败: "+err.Error())
	}

	// 关联权限
	for _, permissionID := range permissionIds {
		rolePermission := &RolePermission{
			RoleID:       role.RoleID,
			PermissionID: permissionID,
		}
		err = tx.Table(constants.RolePermissionTableName).Create(rolePermission).Error
		if err != nil {
			tx.Rollback()
			return 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "关联角色权限失败: "+err.Error())
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交事务失败: "+err.Error())
	}

	return role.RoleID, nil
}

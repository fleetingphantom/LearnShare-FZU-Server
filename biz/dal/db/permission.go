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

	// 收集所有角色ID
	roleIDs := make([]int64, len(roles))
	for i, r := range roles {
		roleIDs[i] = r.RoleID
	}

	// 批量查询所有角色的权限关联
	var rolePermissions []RolePermission
	err = DB.WithContext(ctx).
		Table(constants.RolePermissionTableName).
		Where("role_id IN ?", roleIDs).
		Find(&rolePermissions).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询角色权限失败: "+err.Error())
	}

	// 收集所有权限ID
	permissionIDSet := make(map[int64]bool)
	for _, rp := range rolePermissions {
		permissionIDSet[rp.PermissionID] = true
	}

	permissionIDs := make([]int64, 0, len(permissionIDSet))
	for id := range permissionIDSet {
		permissionIDs = append(permissionIDs, id)
	}

	// 批量查询所有权限详情
	var permissions []Permission
	if len(permissionIDs) > 0 {
		err = DB.WithContext(ctx).
			Table(constants.PermissionTableName).
			Where("permission_id IN ?", permissionIDs).
			Find(&permissions).Error
		if err != nil {
			return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询权限详情失败: "+err.Error())
		}
	}

	// 构建权限ID到权限对象的映射
	permissionMap := make(map[int64]Permission)
	for _, p := range permissions {
		permissionMap[p.PermissionID] = p
	}

	// 构建角色ID到权限列表的映射
	rolePermMap := make(map[int64][]Permission)
	for _, rp := range rolePermissions {
		if p, ok := permissionMap[rp.PermissionID]; ok {
			rolePermMap[rp.RoleID] = append(rolePermMap[rp.RoleID], Permission{
				PermissionID:   p.PermissionID,
				PermissionName: p.PermissionName,
				Description:    p.Description,
			})
		}
	}

	// 组装最终结果
	result := make([]*RoleWithPermissions, len(roles))
	for i, r := range roles {
		result[i] = &RoleWithPermissions{
			RoleID:      r.RoleID,
			RoleName:    r.RoleName,
			Description: r.Description,
			Permissions: rolePermMap[r.RoleID],
		}
		if result[i].Permissions == nil {
			result[i].Permissions = []Permission{}
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

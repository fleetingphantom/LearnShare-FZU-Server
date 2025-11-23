package auth

import (
	"LearnShare/biz/service"
	"LearnShare/pkg/errno"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// RequirePermission 返回需要特定权限的中间件（需要先经过 Auth 中间件）
func RequirePermission(permissionName string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {

		//1. 获取用户角色ID
		roleId := service.GetRoleIdFormContext(c)

		// 2. 获取该角色的所有权限

		permissions, err := service.NewRoleAdminService(ctx, c).GetRolePermissions(roleId)
		if err != nil {
			fail(c, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询角色权限失败"))
			return
		}

		// 3. 检查是否拥有所需权限
		hasPermission := false
		for _, perm := range permissions {
			if permissionName == perm {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			fail(c, errno.NewErrNo(errno.AuthNoOperatePermissionCode, "无权限访问"))
			return
		}

		// 5. 放行
		c.Next(ctx)
	}
}

// RequirePermissions 返回需要多个权限之一的中间件（OR 逻辑，需要先经过 Auth 中间件）
func RequirePermissions(permissionNames ...string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		//1. 获取用户角色ID
		roleId := service.GetRoleIdFormContext(c)

		// 2. 获取该角色的所有权限
		permissions, err := service.NewRoleAdminService(ctx, c).GetRolePermissions(roleId)
		if err != nil {
			fail(c, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询角色权限失败"))
			return
		}

		// 3. 检查是否拥有任一所需权限
		hasPermission := false
		for _, requiredPerm := range permissionNames {
			for _, perm := range permissions {
				if requiredPerm == perm {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			fail(c, errno.NewErrNo(errno.AuthNoOperatePermissionCode, "无权限访问"))
			return
		}

		// 5. 放行
		c.Next(ctx)
	}
}

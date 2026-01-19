package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Permission 权限中间件
type Permission struct {
	// 角色权限检查
	RoleRequired []string
	// 权限代码检查
	PermissionRequired []string
}

// RequireRole 角色权限检查
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查用户是否已认证
		userID, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "用户未认证",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// TODO: 这里需要从数据库查询用户角色信息
		// 暂时使用硬编码的角色检查逻辑
		userRoles := getUserRoles(userID.(uint))

		// 检查用户是否有所需角色
		hasRole := false
		for _, requiredRole := range roles {
			for _, userRole := range userRoles {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "权限不足",
				"data": gin.H{
					"required": roles,
					"current":  userRoles,
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermission 权限代码检查
func RequirePermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查用户是否已认证
		userID, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "用户未认证",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// TODO: 这里需要从数据库查询用户权限信息
		// 暂时使用硬编码的权限检查逻辑
		userPermissions := getUserPermissions(userID.(uint))

		// 检查用户是否有所需权限
		hasPermission := false
		for _, requiredPermission := range permissions {
			for _, userPermission := range userPermissions {
				if userPermission == requiredPermission {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "权限不足",
				"data": gin.H{
					"required": permissions,
					"current":  userPermissions,
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// getUserRoles 获取用户角色（临时实现，后续从数据库查询）
func getUserRoles(userID uint) []string {
	// 临时实现：管理员用户拥有所有角色
	if userID == 1 {
		return []string{"admin", "super_admin"}
	}

	// 默认用户角色
	return []string{"user"}
}

// getUserPermissions 获取用户权限（临时实现，后续从数据库查询）
func getUserPermissions(userID uint) []string {
	// 临时实现：管理员用户拥有所有权限
	if userID == 1 {
		return []string{"system:user:list", "system:user:create", "system:user:update", "system:user:delete", "system:role:list", "system:role:create", "system:role:update", "system:role:delete", "system:menu:list", "system:menu:create", "system:menu:update", "system:menu:delete"}
	}

	// 默认用户权限
	return []string{"system:user:list"}
}

// AdminOnly 仅管理员访问
func AdminOnly() gin.HandlerFunc {
	return RequireRole("admin", "super_admin")
}

// SuperAdminOnly 仅超级管理员访问
func SuperAdminOnly() gin.HandlerFunc {
	return RequireRole("super_admin")
}

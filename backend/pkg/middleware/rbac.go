package middleware

import (
	"webGL-720yun/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RBACMiddleware struct {
	db *gorm.DB
}

func NewRBACMiddleware(db *gorm.DB) *RBACMiddleware {
	return &RBACMiddleware{db: db}
}

// RequirePermission 检查用户是否拥有访问当前资源的权限
func (m *RBACMiddleware) RequirePermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 0. 检查是否为超级管理员
		if isSuper, exists := c.Get("is_super_admin"); exists {
			if super, ok := isSuper.(bool); ok && super {
				c.Next()
				return
			}
		}

		// 1. 获取用户角色 ID（从 AuthMiddleware 设置的上下文中获取）
		val, exists := c.Get("role_id")
		if !exists {
			utils.Forbidden(c.Writer, "未发现角色信息")
			c.Abort()
			return
		}

		roleID, ok := val.(uint)
		if !ok {
			utils.Forbidden(c.Writer, "角色 ID 格式错误")
			c.Abort()
			return
		}

		// 2. 获取当前请求的资源 (URI) 和动作 (Method)
		resource := c.Request.URL.Path
		action := c.Request.Method

		// 3. 查询数据库，检查该角色是否拥有对应权限
		// 注意：这里使用了 many2many 关联查询
		var count int64
		err := m.db.Table("sys_permissions").
			Joins("JOIN sys_role_permissions ON sys_role_permissions.permission_id = sys_permissions.id").
			Where("sys_role_permissions.role_id = ? AND sys_permissions.resource = ? AND sys_permissions.action = ?", roleID, resource, action).
			Count(&count).Error

		if err != nil {
			utils.InternalServerError(c.Writer, "权限检查失败")
			c.Abort()
			return
		}

		// 如果 count 为 0，说明没有权限
		if count == 0 {
			// 特殊处理：如果是管理员，可能默认拥有所有权限，或者在这里添加白名单逻辑
			// 这里简单处理，如果没有明确配置权限则拦截
			utils.Forbidden(c.Writer, "权限不足：无权访问当前资源")
			c.Abort()
			return
		}

		// 4. 验证通过，继续后续流程
		c.Next()
	}
}

// RequireRolePermission 检查是否拥有特定名称的权限
func (m *RBACMiddleware) RequireRolePermission(permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 0. 检查是否为超级管理员
		if isSuper, exists := c.Get("is_super_admin"); exists {
			if super, ok := isSuper.(bool); ok && super {
				c.Next()
				return
			}
		}

		val, exists := c.Get("role_id")
		if !exists {
			utils.Forbidden(c.Writer, "未发现角色信息")
			c.Abort()
			return
		}

		roleID := val.(uint)

		var count int64
		err := m.db.Table("sys_permissions").
			Joins("JOIN sys_role_permissions ON sys_role_permissions.permission_id = sys_permissions.id").
			Where("sys_role_permissions.role_id = ? AND sys_permissions.name = ?", roleID, permissionName).
			Count(&count).Error

		if err != nil || count == 0 {
			utils.Forbidden(c.Writer, "权限不足：缺少 " + permissionName + " 权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

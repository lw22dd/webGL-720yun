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

func (m *RBACMiddleware) RequirePermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isSuper, exists := c.Get("is_super_admin"); exists {
			if super, ok := isSuper.(bool); ok && super {
				c.Next()
				return
			}
		}

		val, exists := c.Get("role_id")
		if !exists {
			utils.Forbidden(c, "未发现角色信息")
			c.Abort()
			return
		}

		roleID, ok := val.(uint)
		if !ok {
			utils.Forbidden(c, "角色 ID 格式错误")
			c.Abort()
			return
		}

		resource := c.Request.URL.Path
		action := c.Request.Method

		var count int64
		err := m.db.Table("sys_permissions").
			Joins("JOIN sys_role_permissions ON sys_role_permissions.permission_id = sys_permissions.id").
			Where("sys_role_permissions.role_id = ? AND sys_permissions.resource = ? AND sys_permissions.action = ?", roleID, resource, action).
			Count(&count).Error

		if err != nil {
			utils.InternalServerError(c, "权限检查失败")
			c.Abort()
			return
		}

		if count == 0 {
			utils.Forbidden(c, "权限不足：无权访问当前资源")
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *RBACMiddleware) RequireRolePermission(permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isSuper, exists := c.Get("is_super_admin"); exists {
			if super, ok := isSuper.(bool); ok && super {
				c.Next()
				return
			}
		}

		val, exists := c.Get("role_id")
		if !exists {
			utils.Forbidden(c, "未发现角色信息")
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
			utils.Forbidden(c, "权限不足：缺少 "+permissionName+" 权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

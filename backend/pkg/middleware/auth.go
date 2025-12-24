package middleware

import (
	"strconv"
	"strings"
	"webGL-720yun/pkg/services/jwt"
	"webGL-720yun/pkg/services/redis"
	"webGL-720yun/pkg/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	jwtService   *jwt.JWTService
	redisService *redis.RedisService
	noAuthPaths  []string
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(jwtService *jwt.JWTService, redisService *redis.RedisService, noAuthPaths []string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService:   jwtService,
		redisService: redisService,
		noAuthPaths:  noAuthPaths,
	}
}

// RequireAuth 需要认证
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否在免认证路径中
		if m.isNoAuthPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c.Writer, "缺少认证令牌")
			c.Abort()
			return
		}

		// 解析Bearer令牌
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(c.Writer, "认证令牌格式错误")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 检查令牌是否在黑名单中
		if m.redisService.IsInBlacklist(tokenString) {
			utils.Unauthorized(c.Writer, "令牌已失效")
			c.Abort()
			return
		}

		// 验证令牌
		claims, err := m.jwtService.ParseToken(tokenString)
		if err != nil {
			utils.Unauthorized(c.Writer, "无效的认证令牌")
			c.Abort()
			return
		}

		// 设置用户信息到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}

// RequireRole 需要特定角色
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户角色
		userRole, exists := c.Get("role")
		if !exists {
			utils.Forbidden(c.Writer, "无法获取用户角色")
			c.Abort()
			return
		}

		// 检查角色权限
		roleStr, ok := userRole.(string)
		if !ok {
			utils.Forbidden(c.Writer, "用户角色格式错误")
			c.Abort()
			return
		}

		// 检查用户角色是否在允许的角色列表中
		hasPermission := false
		for _, role := range roles {
			if roleStr == role {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			utils.Forbidden(c.Writer, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin 需要管理员权限
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return m.RequireRole("admin")
}

// RequireStudent 需要学生权限
func (m *AuthMiddleware) RequireStudent() gin.HandlerFunc {
	return m.RequireRole("student")
}

// RequireAdminOrSelf 需要管理员权限或用户本人
func (m *AuthMiddleware) RequireAdminOrSelf() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID和角色
		userID, _ := c.Get("user_id")
		userRole, _ := c.Get("role")

		// 获取目标用户ID（从URL参数）
		targetUserID := c.Param("id")

		// 如果是管理员，直接通过
		if roleStr, ok := userRole.(string); ok {
			if roleStr == "admin" {
				c.Next()
				return
			}
		}

		// 如果是用户本人，也允许
		if userIDUint, ok := userID.(uint); ok {
			if targetID, err := strconv.ParseUint(targetUserID, 10, 32); err == nil && uint(targetID) == userIDUint {
				c.Next()
				return
			}
		}

		utils.Forbidden(c.Writer, "权限不足")
		c.Abort()
	}
}

// isNoAuthPath 检查路径是否在免认证列表中
func (m *AuthMiddleware) isNoAuthPath(path string) bool {
	for _, noAuthPath := range m.noAuthPaths {
		if strings.HasPrefix(path, noAuthPath) {
			return true
		}
	}
	return false
}

// GetCurrentUser 获取当前用户信息
func GetCurrentUser(c *gin.Context) (userID uint, username string, role string) {
	if id, exists := c.Get("user_id"); exists {
		if idUint, ok := id.(uint); ok {
			userID = idUint
		}
	}
	if name, exists := c.Get("username"); exists {
		if nameStr, ok := name.(string); ok {
			username = nameStr
		}
	}
	if r, exists := c.Get("role"); exists {
		if roleStr, ok := r.(string); ok {
			role = roleStr
		}
	}
	return
}

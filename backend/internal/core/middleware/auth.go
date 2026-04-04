package middleware

import (
	"strconv"
	"strings"

	"webGL-720yun/pkg/jwt"
	"webGL-720yun/pkg/redis"
	"webGL-720yun/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtService   *jwt.JWTService
	redisService *redis.RedisService
	noAuthPaths  []string
}

func NewAuthMiddleware(jwtService *jwt.JWTService, redisService *redis.RedisService, noAuthPaths []string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService:   jwtService,
		redisService: redisService,
		noAuthPaths:  noAuthPaths,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if m.isNoAuthPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c.Writer, "缺少认证令牌")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(c.Writer, "认证令牌格式错误")
			c.Abort()
			return
		}

		tokenString := parts[1]

		if m.redisService.IsInBlacklist(tokenString) {
			utils.Unauthorized(c.Writer, "令牌已失效")
			c.Abort()
			return
		}

		claims, err := m.jwtService.ParseToken(tokenString)
		if err != nil {
			utils.Unauthorized(c.Writer, "无效的认证令牌")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("role_id", claims.RoleID)
		c.Set("is_super_admin", claims.IsSuper)
		c.Set("claims", claims)

		c.Next()
	}
}

func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			utils.Forbidden(c.Writer, "无法获取用户角色")
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			utils.Forbidden(c.Writer, "用户角色格式错误")
			c.Abort()
			return
		}

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

func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return m.RequireRole("admin")
}

func (m *AuthMiddleware) RequireStudent() gin.HandlerFunc {
	return m.RequireRole("student")
}

func (m *AuthMiddleware) RequireAdminOrSelf() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		userRole, _ := c.Get("role")

		targetUserID := c.Param("id")

		if roleStr, ok := userRole.(string); ok {
			if roleStr == "admin" {
				c.Next()
				return
			}
		}

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

func (m *AuthMiddleware) isNoAuthPath(path string) bool {
	for _, noAuthPath := range m.noAuthPaths {
		if strings.HasPrefix(path, noAuthPath) {
			return true
		}
	}
	return false
}

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

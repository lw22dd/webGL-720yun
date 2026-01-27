package middleware

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"webGL-720yun/config"
	"webGL-720yun/pkg/services/redis"
	"webGL-720yun/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	RoleID   uint   `json:"role_id"`
	IsSuper  bool   `json:"is_super"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret         string
	accessTimeout  time.Duration
	refreshTimeout time.Duration
	issuer         string
}

func NewJWTService(cfg *config.JWTConfig) *JWTService {
	return &JWTService{
		secret:         cfg.Secret,
		accessTimeout:  time.Duration(cfg.AccessTimeout) * time.Minute,
		refreshTimeout: time.Duration(cfg.RefreshTimeout) * time.Hour,
		issuer:         cfg.Issuer,
	}
}

func (s *JWTService) GenerateAccessToken(userID uint, username string, role string, roleID uint, isSuper bool) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RoleID:   roleID,
		IsSuper:  isSuper,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTimeout)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *JWTService) GenerateRefreshToken(userID uint) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshTimeout)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *JWTService) ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *JWTService) ValidateToken(tokenString string) error {
	_, err := s.ParseToken(tokenString)
	return err
}

func (s *JWTService) GetTokenRemainingTime(tokenString string) (time.Duration, error) {
	claims, err := s.ParseToken(tokenString)
	if err != nil {
		return 0, err
	}

	if claims.ExpiresAt == nil {
		return 0, errors.New("token has no expiration time")
	}

	remainingTime := time.Until(claims.ExpiresAt.Time)
	if remainingTime < 0 {
		return 0, errors.New("token has expired")
	}

	return remainingTime, nil
}

type AuthMiddleware struct {
	jwtService   *JWTService
	redisService *redis.RedisService
	noAuthPaths  []string
}

func NewAuthMiddleware(jwtService *JWTService, redisService *redis.RedisService, noAuthPaths []string) *AuthMiddleware {
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

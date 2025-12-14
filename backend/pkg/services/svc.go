package services

import (
	"webGL-720yun/app/user"
	"webGL-720yun/pkg/middleware"
	"webGL-720yun/pkg/services/jwt"
	"webGL-720yun/pkg/services/redis"

	"gorm.io/gorm"
)

// 服务上下文，统一管理所有服务实例
type ServiceContext struct {
	// 数据库连接
	DB *gorm.DB

	// 缓存服务
	RedisService *redis.RedisService

	// JWT服务
	JWTService *jwt.JWTService

	// 用户服务
	UserService *user.UserService

	// 认证中间件
	AuthMiddleware *middleware.AuthMiddleware
}

// 创建服务上下文实例
func NewServiceContext(db *gorm.DB, redisService *redis.RedisService, jwtService *jwt.JWTService, userService *user.UserService, authMiddleware *middleware.AuthMiddleware) *ServiceContext {
	return &ServiceContext{
		DB:             db,
		RedisService:   redisService,
		JWTService:     jwtService,
		UserService:    userService,
		AuthMiddleware: authMiddleware,
	}
}

// 获取Redis服务实例
func (sc *ServiceContext) GetRedisService() *redis.RedisService {
	return sc.RedisService
}

// 获取JWT服务实例
func (sc *ServiceContext) GetJWTService() *jwt.JWTService {
	return sc.JWTService
}

// 获取用户服务实例
func (sc *ServiceContext) GetUserService() *user.UserService {
	return sc.UserService
}

// 获取认证中间件实例
func (sc *ServiceContext) GetAuthMiddleware() *middleware.AuthMiddleware {
	return sc.AuthMiddleware
}

// 获取数据库连接实例
func (sc *ServiceContext) GetDB() *gorm.DB {
	return sc.DB
}

// 关闭所有服务连接
func (sc *ServiceContext) Close() error {
	// 关闭Redis连接
	if sc.RedisService != nil {
		if err := sc.RedisService.Close(); err != nil {
			return err
		}
	}

	// 关闭数据库连接
	if sc.DB != nil {
		sqlDB, err := sc.DB.DB()
		if err != nil {
			return err
		}
		if err := sqlDB.Close(); err != nil {
			return err
		}
	}

	return nil
}

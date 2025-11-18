package services

import (
	"gorm.io/gorm"
	"webGL-720yun/app/user"
	"webGL-720yun/pkg/middleware"
	"webGL-720yun/pkg/services/jwt"
	"webGL-720yun/pkg/services/redis"
)

// ServiceManager 服务管理器，统一管理所有服务实例
type ServiceManager struct {
	// 数据库连接
	DB *gorm.DB
	
	// 缓存服务
	RedisService *redis.RedisService
	
	// JWT服务
	JWTService *jwt.JWTService
	
	// 用户服务
	UserService *user.UserService
	
	// 用户处理器
	UserHandler *user.UserHandler
	
	// 认证中间件
	AuthMiddleware *middleware.AuthMiddleware
}

// NewServiceManager 创建服务管理器实例
func NewServiceManager(db *gorm.DB, redisService *redis.RedisService, jwtService *jwt.JWTService, userService *user.UserService, userHandler *user.UserHandler, authMiddleware *middleware.AuthMiddleware) *ServiceManager {
	return &ServiceManager{
		DB:             db,
		RedisService:   redisService,
		JWTService:     jwtService,
		UserService:    userService,
		UserHandler:    userHandler,
		AuthMiddleware: authMiddleware,
	}
}

// GetRedisService 获取Redis服务实例
func (sm *ServiceManager) GetRedisService() *redis.RedisService {
	return sm.RedisService
}

// GetJWTService 获取JWT服务实例
func (sm *ServiceManager) GetJWTService() *jwt.JWTService {
	return sm.JWTService
}

// GetUserService 获取用户服务实例
func (sm *ServiceManager) GetUserService() *user.UserService {
	return sm.UserService
}

// GetUserHandler 获取用户处理器实例
func (sm *ServiceManager) GetUserHandler() *user.UserHandler {
	return sm.UserHandler
}

// GetAuthMiddleware 获取认证中间件实例
func (sm *ServiceManager) GetAuthMiddleware() *middleware.AuthMiddleware {
	return sm.AuthMiddleware
}

// GetDB 获取数据库连接实例
func (sm *ServiceManager) GetDB() *gorm.DB {
	return sm.DB
}

// Close 关闭所有服务连接
func (sm *ServiceManager) Close() error {
	// 关闭Redis连接
	if sm.RedisService != nil {
		if err := sm.RedisService.Close(); err != nil {
			return err
		}
	}
	
	// 关闭数据库连接
	if sm.DB != nil {
		sqlDB, err := sm.DB.DB()
		if err != nil {
			return err
		}
		if err := sqlDB.Close(); err != nil {
			return err
		}
	}
	
	return nil
}
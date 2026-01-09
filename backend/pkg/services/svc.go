package services

import (
	"webGL-720yun/app/user"
	"webGL-720yun/pkg/middleware"
	"webGL-720yun/pkg/services/redis"

	"gorm.io/gorm"
)

type ServiceContext struct {
	DB             *gorm.DB
	RedisService   *redis.RedisService
	JWTService     *middleware.JWTService
	UserService    *user.UserService
	AuthMiddleware *middleware.AuthMiddleware
}

func NewServiceContext(db *gorm.DB, redisService *redis.RedisService, jwtService *middleware.JWTService, userService *user.UserService, authMiddleware *middleware.AuthMiddleware) *ServiceContext {
	return &ServiceContext{
		DB:             db,
		RedisService:   redisService,
		JWTService:     jwtService,
		UserService:    userService,
		AuthMiddleware: authMiddleware,
	}
}

func (sc *ServiceContext) Close() error {
	if sc.RedisService != nil {
		if err := sc.RedisService.Close(); err != nil {
			return err
		}
	}

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

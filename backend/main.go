package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"webGL-720yun/app/user"
	"webGL-720yun/config"
	"webGL-720yun/pkg/database"
	"webGL-720yun/pkg/logger"
	"webGL-720yun/pkg/middleware"
	"webGL-720yun/pkg/services"
	"webGL-720yun/pkg/services/jwt"
	"webGL-720yun/pkg/services/redis"
	routes "webGL-720yun/route"
)

var (
	db             *gorm.DB
	redisClient    *redis.RedisService
	jwtService     *jwt.JWTService
	userService    *user.UserService
	userHandler    *user.UserHandler
	authMiddleware *middleware.AuthMiddleware
	serviceManager *services.ServiceManager
)

func init() {
	// 初始化配置
	if err := config.Init(); err != nil {
		fmt.Printf("初始化配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志系统
	logger.Setup(&config.Conf.Logger)

	// 初始化数据库
	if err := database.Init(&config.Conf.Database); err != nil {
		logger.Fatal("初始化数据库失败:", err)
	}
	db = database.GetDB()
	
	// 调试：打印数据库配置
	logger.Infof("数据库配置: host=%s, port=%d, database=%s", config.Conf.Database.Host, config.Conf.Database.Port, config.Conf.Database.Database)

	// 初始化Redis
	redisClient = redis.NewRedisService(&config.Conf.Redis)

	// 初始化JWT服务
	jwtService = jwt.NewJWTService(&config.Conf.JWT)

	// 初始化用户服务
	userService = user.NewUserService(jwtService, redisClient)

	// 初始化用户处理器
	userHandler = user.NewUserHandler(userService)

	// 初始化权限中间件
	noAuthPaths := config.Conf.NoAuth
	authMiddleware = middleware.NewAuthMiddleware(jwtService, redisClient, noAuthPaths)

	// 初始化服务管理器
	serviceManager = services.NewServiceManager(db, redisClient, jwtService, userService, userHandler, authMiddleware)
}

func main() {
	// 设置Gin模式
	if config.Conf.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin引擎
	r := gin.New()

	// 添加中间件
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(gin.Recovery())

	// 注册用户相关路由
	routes.UserAPIRoutes(r, serviceManager)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Conf.Server.Port),
		Handler: r,
	}

	// 启动服务器
	go func() {
		logger.Infof("服务器启动在端口 %d", config.Conf.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器启动失败:", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Infoln("正在关闭服务器...")

	// 设置关闭超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭HTTP服务器
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("服务器关闭失败:", err)
	}

	// 关闭服务管理器（包含所有服务连接）
	if err := serviceManager.Close(); err != nil {
		logger.Errorln("关闭服务管理器失败:", err)
	}

	logger.Infoln("服务器已关闭")
}

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

func main() {

	// 1. 初始化配置
	if err := config.Init(); err != nil {
		fmt.Printf("配置加载失败: %v\n", err)
		os.Exit(1)
	}

	// 2. 初始化日志系统
	logger.Setup(&config.Conf.Logger)

	// 3. 初始化数据库
	if err := database.Init(&config.Conf.Database); err != nil {
		logger.Fatal("初始化数据库失败:", err)
	}
	db := database.GetDB()

	// 调试：打印数据库配置
	logger.Infof("数据库配置: host=%s, port=%d, database=%s", config.Conf.Database.Host, config.Conf.Database.Port, config.Conf.Database.Database)

	// 4. 初始化Redis服务
	redisClient := redis.NewRedisService(&config.Conf.Redis)

	// 5. 初始化JWT服务
	jwtService := jwt.NewJWTService(&config.Conf.JWT)

	// 6. 初始化用户服务
	userService := user.NewUserService(jwtService, redisClient)

	// 7. 初始化权限中间件
	noAuthPaths := config.Conf.NoAuth
	authMiddleware := middleware.NewAuthMiddleware(jwtService, redisClient, noAuthPaths)

	// 8. 初始化服务上下文
	serviceContext := services.NewServiceContext(db, redisClient, jwtService, userService, authMiddleware)

	// 10. 创建Gin实例（根据配置设置模式）
	if config.Conf.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// 11. 注册全局中间件
	registerCoreMiddleware(router)

	// 12. 注册路由
	routes.UserAPIRoutes(router, serviceContext)

	// 13. 配置HTTP服务器
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.Conf.Server.Port),
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 14. 启动服务器
	go func() {
		logger.Info("服务器启动", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器异常退出", "error", err)
		}
	}()

	// 15. 优雅关机处理
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("正在关闭服务...", "signal", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("强制关闭服务", "error", err)
	}

	// 关闭服务上下文（包含所有服务连接）
	if err := serviceContext.Close(); err != nil {
		logger.Error("关闭服务上下文失败", "error", err)
	}

	logger.Info("服务已安全停止")
}

// 核心中间件注册（与网关功能强相关）
func registerCoreMiddleware(r *gin.Engine) {
	r.Use(
		gin.Recovery(),        // 官方恢复中间件
		logger.GinZapLogger(), // 自定义日志中间件
	)
}

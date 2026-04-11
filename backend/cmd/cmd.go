package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"webGL-720yun/config"
	"webGL-720yun/internal/core/middleware"
	"webGL-720yun/internal/core/router"
	"webGL-720yun/internal/core/setup"
	"webGL-720yun/internal/resource/service"
	upload_service "webGL-720yun/internal/upload/service"
	"webGL-720yun/internal/user"
	"webGL-720yun/pkg/jwt"
	"webGL-720yun/pkg/logger"
	"webGL-720yun/pkg/minio_client"
	"webGL-720yun/pkg/redis"
	"webGL-720yun/pkg/websocket"
)

func Run() {
	if err := config.Init(); err != nil {
		fmt.Printf("配置加载失败: %v\n", err)
		os.Exit(1)
	}

	logger.Setup(&config.Conf.Logger)

	var db *setup.Database
	var err error
	db, err = setup.NewDatabase(&config.Conf.Database)
	if err != nil {
		logger.Warnf("数据库初始化失败（非致命）: %v", err)
	}

	if db != nil {
		logger.Infof("数据库配置: host=%s, port=%d, database=%s", config.Conf.Database.Host, config.Conf.Database.Port, config.Conf.Database.Database)
	}

	redisClient := redis.NewRedisService(&config.Conf.Redis)

	minioClient, err := minio_client.NewMinIOClient(&config.MinIOConfig{
		Endpoint:  config.Conf.MinIO.Endpoint,
		AccessKey: config.Conf.MinIO.AccessKey,
		SecretKey: config.Conf.MinIO.SecretKey,
		Bucket:    config.Conf.MinIO.Bucket,
		UseSSL:    config.Conf.MinIO.UseSSL,
		Region:    config.Conf.MinIO.Region,
	})
	if err != nil {
		logger.Warnf("MinIO初始化失败（非致命）: %v", err)
	}

	jwtService := jwt.NewJWTService(&config.Conf.JWT)

	var userService *user.UserService
	var spaceService *service.SpaceService
	var sceneService *service.SceneService
	var hotspotService *service.HotspotService
	var uploadService *upload_service.UploadService

	if db != nil {
		userService = user.NewUserService(db.GetDB(), jwtService, redisClient)
		spaceService = service.NewSpaceService(db.GetDB(), minioClient)
		sceneService = service.NewSceneService(db.GetDB(), minioClient)
		hotspotService = service.NewHotspotService(db.GetDB())
		
		if err := setup.SeedTestData(db.GetDB()); err != nil {
			logger.Warnf("种子数据初始化失败（非致命）: %v", err)
		}

		if minioClient != nil {
			resourceInitService := service.NewResourceInitService(db.GetDB(), minioClient, ".")
			if initErr := resourceInitService.SeedResourcesIfNeeded(); initErr != nil {
				logger.Warnf("全景资源初始化失败（非致命）: %v", initErr)
			}
		}
	} else {
		// 使用空实现，允许服务启动但功能受限
		userService = &user.UserService{}
		spaceService = &service.SpaceService{}
		sceneService = &service.SceneService{}
		hotspotService = &service.HotspotService{}
	}

	wsHub := websocket.NewHub()
	go wsHub.Run()

	if db != nil {
		uploadService = router.NewUploadService(db.GetDB(), minioClient, redisClient, wsHub)
	} else {
		uploadService = upload_service.NewUploadService(
			nil,
			nil,
			nil,
			minioClient,
			wsHub,
		)
	}

	noAuthPaths := config.Conf.NoAuth
	authMiddleware := middleware.NewAuthMiddleware(jwtService, redisClient, noAuthPaths)

	rbacMiddleware := middleware.NewRBACMiddleware(db.GetDB())

	serviceContext := &router.ServiceContext{
		UserService:    userService,
		SpaceService:   spaceService,
		SceneService:   sceneService,
		HotspotService: hotspotService,
		UploadService:  uploadService,
		AuthMiddleware: authMiddleware,
		RBACMiddleware: rbacMiddleware,
		MinIOClient:    minioClient,
		RedisService:   redisClient,
		JWTService:     jwtService,
		WsHub:          wsHub,
	}

	if config.Conf.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	RegisterCoreMiddleware(r)

	router.RegisterRoutes(r, serviceContext)

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.Conf.Server.Port),
		Handler:        r,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		logger.Info("服务器启动", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器异常退出", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("正在关闭服务...", "signal", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("强制关闭服务", "error", err)
	}

	if err := redisClient.Close(); err != nil {
		logger.Error("关闭Redis连接失败", "error", err)
	}

	if err := db.Close(); err != nil {
		logger.Error("关闭数据库连接失败", "error", err)
	}

	logger.Info("服务已安全停止")
}

func RegisterCoreMiddleware(r *gin.Engine) {
	r.Use(
		gin.Recovery(),
		logger.GinZapLogger(),
		middleware.CORS(),
	)
}

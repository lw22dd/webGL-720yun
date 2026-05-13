package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"webGL-720yun/config"
	"webGL-720yun/internal/core/middleware"
	"webGL-720yun/internal/core/router"
	"webGL-720yun/internal/core/setup"
	"webGL-720yun/internal/model"
	resService "webGL-720yun/internal/resource/service"
	sliceService "webGL-720yun/internal/slice/service"
	userService "webGL-720yun/internal/user/service"
	"webGL-720yun/pkg/database"
	"webGL-720yun/pkg/jwt"
	"webGL-720yun/pkg/logger"
	"webGL-720yun/pkg/minio_client"
	"webGL-720yun/pkg/redis"
	"webGL-720yun/pkg/websocket"
)

const DefaultWorkerCount = 3

func Run() {
	// 设置 Go 运行时内存上限为 4GiB (Go 1.19+)
	// 这会使 GC 在接近此限制时变得更积极，任务完成后也能更快回落
	debug.SetMemoryLimit(4 * 1024 * 1024 * 1024)

	// 1. 初始化配置 (Internal) - 加载应用所有配置项
	if err := config.Init(); err != nil {
		fmt.Printf("配置加载失败: %v\n", err)
		os.Exit(1)
	}

	logger.Setup(&config.Conf.Logger)

	// 2. 初始化基础设施 (Pkg) - 实例化所有外部服务客户端
	db, err := database.NewMySQL(&config.Conf.Database)
	if err != nil {
		logger.Fatal("初始化数据库失败:", err)
	}
	logger.Infof("数据库配置: host=%s, port=%d, database=%s", config.Conf.Database.Host, config.Conf.Database.Port, config.Conf.Database.Database)

	// 自动迁移数据库表结构
	if err := db.AutoMigrate(); err != nil {
		logger.Fatal("数据库迁移失败:", err)
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
		logger.Warnf("MinIO客户端创建失败（非致命）: %v", err)
	}

	jwtService := jwt.NewJWTService(&config.Conf.JWT)

	// 3. 执行业务初始化 (Internal/Core/Setup) - 仅处理业务数据初始化
	var seededScenes []model.ResScene
	if config.Conf.App.InitData {
		if err := setup.SeedTestData(db.DB); err != nil {
			logger.Fatal("种子数据初始化失败:", err)
		}

		if minioClient != nil {
			var initErr error
			seededScenes, initErr = setup.SeedResourcesIfNeeded(db.DB, minioClient, ".")
			if initErr != nil {
				logger.Warnf("全景资源初始化失败（非致命）: %v", initErr)
			}
		}
	}

	// 4. 初始化服务层
	userSvc := userService.NewUserService(db.DB, jwtService, redisClient)
	spaceService := resService.NewSpaceService(db.DB, minioClient)
	sliceQueue := router.NewSliceQueue(redisClient)
	sceneService := resService.NewSceneService(db.DB, minioClient, sliceQueue, redisClient)
	hotspotService := resService.NewHotspotService(db.DB)

	// 5. 为需要切片的场景创建切片任务
	// 包括：新创建的场景 + 已有但未切片的场景
	var scenesToSlice []model.ResScene
	if err := db.DB.Where("slice_status = ? OR slice_status = ?", model.SliceStatusPending, model.SliceStatusFailed).Find(&scenesToSlice).Error; err != nil {
		logger.Warnf("查询待切片场景失败: %v", err)
	}

	// 合并新创建的场景（避免重复）
	sceneMap := make(map[uint]bool)
	for _, scene := range scenesToSlice {
		sceneMap[scene.ID] = true
	}
	for _, scene := range seededScenes {
		if !sceneMap[scene.ID] {
			scenesToSlice = append(scenesToSlice, scene)
			sceneMap[scene.ID] = true
		}
	}

	if len(scenesToSlice) > 0 {
		logger.Infof("🔄 为 %d 个场景创建切片任务...", len(scenesToSlice))
		for _, scene := range scenesToSlice {
			var space model.ResSpace
			if err := db.DB.First(&space, scene.SpaceID).Error; err != nil {
				logger.Warnf("获取景区信息失败: %v", err)
				continue
			}

			task := &sliceService.SliceTask{
				SceneID:   scene.ID,
				SceneCode: scene.SceneCode,
				FileID:    scene.FileID,
				SpaceName: space.Name,
				SpaceSlug: space.Slug,
				UserID:    0, // 系统任务
			}
			if err := sliceQueue.PushTask(task); err != nil {
				logger.Warnf("创建切片任务失败 [%s]: %v", scene.SceneCode, err)
			} else {
				logger.Infof("✅ 切片任务已创建: %s", scene.SceneCode)
			}
		}
	}

	wsHub := websocket.NewHub()
	go wsHub.Run()

	uploadService := router.NewUploadService(db.DB, minioClient, redisClient, wsHub, sliceQueue)

	scheduler := router.NewSliceScheduler(sliceQueue, db.DB, minioClient, wsHub, DefaultWorkerCount)
	go scheduler.Start()

	// 5. 初始化中间件
	noAuthPaths := config.Conf.NoAuth
	authMiddleware := middleware.NewAuthMiddleware(jwtService, redisClient, noAuthPaths)
	rbacMiddleware := middleware.NewRBACMiddleware(db.DB)

	// 6. 启动路由
	serviceContext := &router.ServiceContext{
		UserService:    userSvc,
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
		SliceQueue:     sliceQueue,
		Scheduler:      scheduler,
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

	scheduler.Stop()
	logger.Info("Scheduler 已停止")

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

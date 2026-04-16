package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"webGL-720yun/internal/core/middleware"
	"webGL-720yun/internal/resource/handler"
	"webGL-720yun/internal/resource/service"
	"webGL-720yun/internal/resource/upload"
	"webGL-720yun/internal/slice"
	"webGL-720yun/internal/user"
	"webGL-720yun/pkg/jwt"
	"webGL-720yun/pkg/minio_client"
	"webGL-720yun/pkg/redis"
	"webGL-720yun/pkg/websocket"
)

type ServiceContext struct {
	UserService    *user.UserService
	SpaceService   *service.SpaceService
	SceneService   *service.SceneService
	HotspotService *service.HotspotService
	UploadService  *upload.UploadService
	AuthMiddleware *middleware.AuthMiddleware
	RBACMiddleware *middleware.RBACMiddleware
	MinIOClient    *minio_client.MinIOClient
	RedisService   *redis.RedisService
	JWTService     *jwt.JWTService
	WsHub          *websocket.Hub
	SliceQueue     *slice.SliceQueue
	WorkerPool     *slice.WorkerPool
}

func RegisterRoutes(r *gin.Engine, ctx *ServiceContext) {
	userService := ctx.UserService
	spaceService := ctx.SpaceService
	sceneService := ctx.SceneService
	hotspotService := ctx.HotspotService
	uploadService := ctx.UploadService
	authMiddleware := ctx.AuthMiddleware
	wsHub := ctx.WsHub

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", user.Login(userService))
			auth.POST("/refresh", user.RefreshToken(userService))
		}

		userGroup := api.Group("/user")
		{
			userGroup.POST("/register", user.Register(userService))

			userGroup.Use(authMiddleware.RequireAuth())
			{
				userGroup.GET("/profile", user.GetProfile(userService))
				userGroup.PUT("/profile", user.UpdateProfile(userService))
				userGroup.POST("/change-password", user.ChangePassword(userService))
				userGroup.POST("/logout", user.Logout(userService))

				admin := userGroup.Group("/admin")
				admin.Use(authMiddleware.RequireAdmin())
				{
					admin.GET("/list", user.GetUserList(userService))
					admin.GET("/:id", user.GetUserByID(userService))
					admin.POST("/create", user.Register(userService))
					admin.PUT("/:id", user.UpdateUser(userService))
					admin.DELETE("/:id", user.DeleteUser(userService))
					admin.DELETE("/batch", user.DeleteUserBatch(userService))
					admin.POST("/batch-register", user.BatchRegister(userService))
				}
			}
		}

		resourceGroup := api.Group("/resource")
		resourceGroup.Use(authMiddleware.RequireAuth())
		{
			spaces := resourceGroup.Group("/spaces")
			{
				spaces.GET("", handler.GetSpaceList(spaceService))
				spaces.GET("/:id", handler.GetSpaceDetail(spaceService))
				spaces.GET("/:id/graph", handler.GetSpaceGraph(sceneService))
				spaces.POST("", authMiddleware.RequireAdmin(), handler.CreateSpace(spaceService))
				spaces.PUT("/:id", handler.UpdateSpace(spaceService))
				spaces.DELETE("/:id", handler.DeleteSpace(spaceService))
				spaces.DELETE("/batch", authMiddleware.RequireAdmin(), handler.DeleteSpaceBatch(spaceService))
			}

			scenes := resourceGroup.Group("/scenes")
			{
				scenes.GET("", handler.GetSceneList(sceneService))
				scenes.GET("/:id", handler.GetSceneDetail(sceneService))
				scenes.POST("", authMiddleware.RequireAdmin(), handler.CreateScene(sceneService))
				scenes.PUT("/:id", handler.UpdateScene(sceneService))
				scenes.PUT("/:id/position", handler.UpdateScenePosition(sceneService))
				scenes.PUT("/batch-position", handler.BatchUpdateScenePosition(sceneService))
				scenes.DELETE("/:id", handler.DeleteScene(sceneService))
				scenes.POST("/batch-import", authMiddleware.RequireAdmin(), handler.BatchImportScenes(sceneService))
			}

			hotspots := resourceGroup.Group("/hotspots")
			{
				hotspots.GET("", handler.GetHotspotList(hotspotService))
				hotspots.GET("/:id", handler.GetHotspotDetail(hotspotService))
				hotspots.POST("", authMiddleware.RequireAdmin(), handler.CreateHotspot(hotspotService))
				hotspots.PUT("/:id", handler.UpdateHotspot(hotspotService))
				hotspots.DELETE("/:id", handler.DeleteHotspot(hotspotService))
			}
		}

		uploadGroup := api.Group("/upload")
		uploadGroup.Use(authMiddleware.RequireAuth())
		{
			uploadGroup.POST("/init", upload.InitUpload(uploadService))
			uploadGroup.POST("/chunk", upload.UploadChunk(uploadService))
			uploadGroup.POST("/complete", upload.CompleteUpload(uploadService))
			uploadGroup.GET("/status/:upload_id", upload.GetUploadStatus(uploadService))
			uploadGroup.DELETE("/:upload_id", upload.CancelUpload(uploadService))
			uploadGroup.GET("/file/:file_id", upload.GetFileInfo(uploadService))
		}

		api.GET("/ws", websocket.HandleWebSocket(wsHub))

		// 资源流式读取 —— 公开路由，无需鉴权
		// 瓦片是高频请求（每个场景 500+ 张），不做 JWT 校验
		resGroup := api.Group("/res")
		{
			// GET /api/v1/res/tiles/:sceneCode/:face/:level/:x/:y
			resGroup.GET("/tiles/:sceneCode/:face/:level/:x/:y", handler.GetTile(sceneService))
			resGroup.HEAD("/tiles/:sceneCode/:face/:level/:x/:y", handler.GetTile(sceneService))
			// GET /api/v1/res/cubemap/:sceneCode/:face
			resGroup.GET("/cubemap/:sceneCode/:face", handler.GetCubemapFace(sceneService))
			resGroup.HEAD("/cubemap/:sceneCode/:face", handler.GetCubemapFace(sceneService))
			// GET /api/v1/res/previews/:sceneCode
			resGroup.GET("/previews/:sceneCode", handler.GetPreview(sceneService))
			resGroup.HEAD("/previews/:sceneCode", handler.GetPreview(sceneService))
			// GET /api/v1/res/sources/:sceneCode
			resGroup.GET("/sources/:sceneCode", handler.GetSource(sceneService))
			resGroup.HEAD("/sources/:sceneCode", handler.GetSource(sceneService))
			// GET /api/v1/res/covers/:spaceName
			resGroup.GET("/covers/:spaceName", handler.GetCover(sceneService))
			resGroup.HEAD("/covers/:spaceName", handler.GetCover(sceneService))
		}
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})
}

func NewUploadService(db *gorm.DB, minioClient *minio_client.MinIOClient, redisService *redis.RedisService) *upload.UploadService {
	uploadRepo := upload.NewUploadRepository(redisService, db)
	return upload.NewUploadService(uploadRepo, db, minioClient)
}

func NewSliceQueue(redisService *redis.RedisService) *slice.SliceQueue {
	return slice.NewSliceQueue(redisService)
}

func NewWorkerPool(queue *slice.SliceQueue, db *gorm.DB, minioClient *minio_client.MinIOClient, wsHub *websocket.Hub, workerCount int) *slice.WorkerPool {
	return slice.NewWorkerPool(queue, db, minioClient, wsHub, workerCount)
}

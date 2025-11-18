package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"webGL-720yun/pkg/services"
)


// UserAPIRoutes 用户相关路由注册
func UserAPIRoutes(r *gin.Engine, serviceManager *services.ServiceManager) {
	// 从服务管理器获取所需的服务实例
	userHandler := serviceManager.GetUserHandler()
	authMiddleware := serviceManager.GetAuthMiddleware()
	
	// API路由组
	api := r.Group("/api/v1")
	{
		// 认证相关路由
		auth := api.Group("/auth")
		{
			auth.POST("/login", userHandler.Login)
			auth.POST("/refresh", userHandler.RefreshToken)
		}

		// 用户相关路由
		user := api.Group("/user")
		{
			// 公开路由
			user.POST("/register", userHandler.CreateUser)

			// 需要认证的路由
			user.Use(authMiddleware.RequireAuth())
			{
				user.GET("/profile", userHandler.GetProfile)
				user.PUT("/profile", userHandler.UpdateProfile)
				user.POST("/change-password", userHandler.ChangePassword)
				user.POST("/logout", userHandler.Logout)

				// 管理员专属路由
				admin := user.Group("/admin")
				admin.Use(authMiddleware.RequireAdmin())
				{
					admin.GET("/list", userHandler.GetUserList)
					admin.GET("/:id", userHandler.GetUserByID)
					admin.POST("/create", userHandler.CreateUser)
					admin.PUT("/:id", userHandler.UpdateUser)
					admin.DELETE("/:id", userHandler.DeleteUser)
				}
			}
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})
}
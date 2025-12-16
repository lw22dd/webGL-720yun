package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"webGL-720yun/app/user"
	"webGL-720yun/pkg/services"
)

// UserAPIRoutes 用户相关路由注册
func UserAPIRoutes(r *gin.Engine, serviceContext *services.ServiceContext) {
	// 从服务上下文获取所需的服务实例
	userService := serviceContext.GetUserService()
	authMiddleware := serviceContext.GetAuthMiddleware()

	// API路由组
	api := r.Group("/api/v1")
	{
		// 认证相关路由
		auth := api.Group("/auth")
		{
			auth.POST("/login", user.Login(userService))
			auth.POST("/refresh", user.RefreshToken(userService))
		}

		// 用户相关路由
		userGroup := api.Group("/user")
		{
			// 公开路由
			userGroup.POST("/register", user.Register(userService))

			// 需要认证的路由
			userGroup.Use(authMiddleware.RequireAuth())
			{
				userGroup.GET("/profile", user.GetProfile(userService))
				userGroup.PUT("/profile", user.UpdateProfile(userService))
				userGroup.POST("/change-password", user.ChangePassword(userService))
				userGroup.POST("/logout", user.Logout(userService))

				// 管理员专属路由
				admin := userGroup.Group("/admin")
				admin.Use(authMiddleware.RequireAdmin())
				{
					admin.GET("/list", user.GetUserList(userService))
					admin.GET("/:id", user.GetUserByID(userService))
					admin.POST("/create", user.Register(userService))
					admin.PUT("/:id", user.UpdateUser(userService))
					admin.DELETE("/:id", user.DeleteUser(userService))
					// 批量注册路由
					admin.POST("/batch-register", user.BatchRegister(userService))
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
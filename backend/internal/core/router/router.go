package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"webGL-720yun/internal/core/middleware"
	"webGL-720yun/internal/user"
)

type ServiceContext struct {
	UserService    *user.UserService
	AuthMiddleware *middleware.AuthMiddleware
	RBACMiddleware *middleware.RBACMiddleware
}

func RegisterRoutes(r *gin.Engine, ctx *ServiceContext) {
	userService := ctx.UserService
	authMiddleware := ctx.AuthMiddleware

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
					admin.POST("/batch-register", user.BatchRegister(userService))
				}
			}
		}
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})
}

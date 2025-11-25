package user

import (
	"net/http"
	"strconv"
	"strings"
	"webGL-720yun/app/models"
	"webGL-720yun/pkg/middleware"
	"webGL-720yun/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Login 用户登录处理函数
func Login(service *UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c.Writer, "请求参数错误")
			return
		}

		response, err := service.Login(&req)
		if err != nil {
			utils.Unauthorized(c.Writer, err.Error())
			return
		}

		utils.Success(c.Writer, response)
	}
}

// Register 用户注册处理函数（管理员创建学生账户）
func Register(service *UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req models.RegisterRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(ctx.Writer, "请求参数错误")
			return
		}

		user, err := service.Register(&req)
		if err != nil {
			utils.Error(ctx.Writer, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(ctx.Writer, user)
	}
}

// Logout 用户登出处理函数
func Logout(service *UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)

		// 获取访问令牌
		authHeader := c.GetHeader("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.BadRequest(c.Writer, "认证令牌格式错误")
			return
		}

		if err := service.Logout(userID, parts[1]); err != nil {
			utils.InternalServerError(c.Writer, "登出失败")
			return
		}

		utils.Success(c.Writer, gin.H{"message": "登出成功"})
	}
}

// RefreshToken 刷新访问令牌处理函数
func RefreshToken(service *UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.RefreshTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c.Writer, "请求参数错误")
			return
		}

		response, err := service.RefreshToken(req.RefreshToken)
		if err != nil {
			utils.Unauthorized(c.Writer, err.Error())
			return
		}

		utils.Success(c.Writer, response)
	}
}

// GetProfile 获取个人资料处理函数
func GetProfile(service *UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)

		user, err := service.GetUserByID(userID)
		if err != nil {
			utils.NotFound(c.Writer, "用户不存在")
			return
		}

		utils.Success(c.Writer, user)
	}
}

// UpdateProfile 更新个人资料处理函数
func UpdateProfile(service *UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)

		var req models.UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c.Writer, "请求参数错误")
			return
		}

		user, err := service.UpdateUser(userID, &req)
		if err != nil {
			utils.Error(c.Writer, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c.Writer, user)
	}
}

// ChangePassword 修改密码处理函数
func ChangePassword(service *UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)

		var req models.ChangePasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c.Writer, "请求参数错误")
			return
		}

		if err := service.ChangePassword(userID, &req); err != nil {
			if err.Error() == "旧密码错误" {
				utils.Error(c.Writer, http.StatusBadRequest, err.Error())
			} else {
				utils.InternalServerError(c.Writer, err.Error())
			}
			return
		}

		utils.Success(c.Writer, gin.H{"message": "密码修改成功"})
	}
}

// GetUserByID 根据ID获取用户处理函数
func GetUserByID(service *UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			utils.BadRequest(c.Writer, "用户ID格式错误")
			return
		}

		user, err := service.GetUserByID(uint(id))
		if err != nil {
			utils.NotFound(c.Writer, "用户不存在")
			return
		}

		utils.Success(c.Writer, user)
	}
}

// GetUserList 获取用户列表处理函数
func GetUserList(service *UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.UserListRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			utils.BadRequest(c.Writer, "请求参数错误")
			return
		}

		// 设置默认值
		if req.Page == 0 {
			req.Page = 1
		}
		if req.PageSize == 0 {
			req.PageSize = 10
		}

		response, err := service.GetUserList(&req)
		if err != nil {
			utils.InternalServerError(c.Writer, "获取用户列表失败")
			return
		}

		utils.Success(c.Writer, gin.H{
			"total":     response.Total,
			"users":     response.Users,
			"page":      req.Page,
			"page_size": req.PageSize,
		})
	}
}

// UpdateUser 更新用户信息处理函数
func UpdateUser(service *UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			utils.BadRequest(c.Writer, "用户ID格式错误")
			return
		}

		var req models.UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c.Writer, "请求参数错误")
			return
		}

		user, err := service.UpdateUser(uint(id), &req)
		if err != nil {
			utils.Error(c.Writer, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c.Writer, user)
	}
}

// DeleteUser 删除用户处理函数
func DeleteUser(service *UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			utils.BadRequest(c.Writer, "用户ID格式错误")
			return
		}

		if err := service.DeleteUser(uint(id)); err != nil {
			utils.InternalServerError(c.Writer, "删除用户失败")
			return
		}

		utils.Success(c.Writer, gin.H{"message": "删除成功"})
	}
}

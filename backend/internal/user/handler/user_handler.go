package handler

import (
	"net/http"
	"strings"

	"webGL-720yun/internal/core/middleware"
	"webGL-720yun/internal/user/dto"
	"webGL-720yun/internal/user/service"
	"webGL-720yun/pkg/utils"

	"github.com/gin-gonic/gin"
)

func Login(service *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c, "请求参数错误")
			return
		}

		response, err := service.Login(&req)
		if err != nil {
			utils.Unauthorized(c, err.Error())
			return
		}

		utils.Success(c, response)
	}
}

func Register(service *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c, "请求参数错误")
			return
		}

		user, err := service.Register(&req)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c, user)
	}
}

func Logout(service *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)

		authHeader := c.GetHeader("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.BadRequest(c, "认证令牌格式错误")
			return
		}

		if err := service.Logout(userID, parts[1]); err != nil {
			utils.InternalServerError(c, "登出失败")
			return
		}

		utils.Success(c, gin.H{"message": "登出成功"})
	}
}

func RefreshToken(service *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.RefreshTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c, "请求参数错误")
			return
		}

		response, err := service.RefreshToken(req.RefreshToken)
		if err != nil {
			utils.Unauthorized(c, err.Error())
			return
		}

		utils.Success(c, response)
	}
}

func GetProfile(service *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)

		user, err := service.GetUserByID(userID)
		if err != nil {
			utils.NotFound(c, "用户不存在")
			return
		}

		utils.Success(c, user)
	}
}

func UpdateProfile(service *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)

		var req dto.UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c, "请求参数错误")
			return
		}

		user, err := service.UpdateUser(userID, &req)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c, user)
	}
}

func ChangePassword(service *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)

		var req dto.ChangePasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c, "请求参数错误")
			return
		}

		if err := service.ChangePassword(userID, &req); err != nil {
			if err.Error() == "旧密码错误" {
				utils.Error(c, http.StatusBadRequest, err.Error())
			} else {
				utils.InternalServerError(c, err.Error())
			}
			return
		}

		utils.Success(c, gin.H{"message": "密码修改成功"})
	}
}

type IDParamRequest struct {
	ID uint `uri:"id" binding:"required"`
}

func GetUserByID(service *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req IDParamRequest
		if err := c.ShouldBindUri(&req); err != nil {
			utils.BadRequest(c, "参数错误")
			return
		}

		user, err := service.GetUserByID(req.ID)
		if err != nil {
			utils.NotFound(c, "用户不存在")
			return
		}

		utils.Success(c, user)
	}
}

func GetUserList(service *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.UserListRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			utils.BadRequest(c, "请求参数错误")
			return
		}

		if req.Page == 0 {
			req.Page = 1
		}
		if req.PageSize == 0 {
			req.PageSize = 10
		}

		response, err := service.GetUserList(&req)
		if err != nil {
			utils.InternalServerError(c, "获取用户列表失败")
			return
		}

		utils.Success(c, gin.H{
			"total":     response.Total,
			"users":     response.Users,
			"page":      req.Page,
			"page_size": req.PageSize,
		})
	}
}

func UpdateUser(service *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req IDParamRequest
		if err := c.ShouldBindUri(&req); err != nil {
			utils.BadRequest(c, "参数错误")
			return
		}

		var updateReq dto.UpdateUserRequest
		if err := c.ShouldBindJSON(&updateReq); err != nil {
			utils.BadRequest(c, "请求参数错误")
			return
		}

		user, err := service.UpdateUser(req.ID, &updateReq)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c, user)
	}
}

func DeleteUser(service *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req IDParamRequest
		if err := c.ShouldBindUri(&req); err != nil {
			utils.BadRequest(c, "参数错误")
			return
		}

		if err := service.DeleteUser(req.ID); err != nil {
			utils.InternalServerError(c, "删除用户失败")
			return
		}

		utils.Success(c, gin.H{"message": "删除成功"})
	}
}

// DeleteUserBatchRequest 批量删除用户请求
type DeleteUserBatchRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}

func DeleteUserBatch(service *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req DeleteUserBatchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c, "请求参数错误: 需要提供ids数组")
			return
		}

		var failedIDs []uint
		var successCount int

		for _, id := range req.IDs {
			if err := service.DeleteUser(id); err != nil {
				failedIDs = append(failedIDs, id)
			} else {
				successCount++
			}
		}

		if len(failedIDs) > 0 {
			utils.Success(c, gin.H{
				"message":       "批量删除部分成功",
				"success_count": successCount,
				"failed_count":  len(failedIDs),
				"failed_ids":    failedIDs,
			})
			return
		}

		utils.Success(c, gin.H{
			"message":       "批量删除成功",
			"success_count": successCount,
		})
	}
}

func BatchRegister(service *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			utils.BadRequest(c, "获取上传文件失败: "+err.Error())
			return
		}
		defer file.Close()

		students, err := service.ParseFile(file, header.Filename)
		if err != nil {
			utils.BadRequest(c, "解析文件失败: "+err.Error())
			return
		}

		if len(students) == 0 {
			utils.BadRequest(c, "文件中没有有效数据")
			return
		}

		response, err := service.BatchRegister(students)
		if err != nil {
			utils.InternalServerError(c, "批量注册失败: "+err.Error())
			return
		}

		utils.Success(c, response)
	}
}

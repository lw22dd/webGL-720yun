package handler

import (
	"net/http"
	"strconv"

	"webGL-720yun/internal/core/middleware"
	"webGL-720yun/internal/resource/dto"
	"webGL-720yun/internal/resource/service"
	"webGL-720yun/pkg/utils"

	"github.com/gin-gonic/gin"
)

type SpaceHandler struct {
	service *service.SpaceService
}

func NewSpaceHandler(service *service.SpaceService) *SpaceHandler {
	return &SpaceHandler{service: service}
}

func GetSpaceList(svc *service.SpaceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SpaceListRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			utils.BadRequest(c.Writer, "请求参数错误")
			return
		}

		response, err := svc.GetSpaceList(&req)
		if err != nil {
			utils.InternalServerError(c.Writer, err.Error())
			return
		}

		utils.Success(c.Writer, response)
	}
}

func GetSpaceDetail(svc *service.SpaceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")
		if slug == "" {
			utils.BadRequest(c.Writer, "缺少slug参数")
			return
		}

		response, err := svc.GetSpaceDetail(slug)
		if err != nil {
			utils.NotFound(c.Writer, err.Error())
			return
		}

		utils.Success(c.Writer, response)
	}
}

func CreateSpace(svc *service.SpaceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)
		isSuperAdmin, _ := c.Get("is_super_admin")
		isAdmin := isSuperAdmin.(bool)

		var req dto.CreateSpaceRequest
		if err := c.ShouldBind(&req); err != nil {
			utils.BadRequest(c.Writer, "请求参数错误: "+err.Error())
			return
		}

		var coverFile interface{}
		if file, err := c.FormFile("cover"); err == nil {
			coverFile = file
		}

		space, err := svc.CreateSpace(&req, userID, coverFile)
		if err != nil {
			utils.Error(c.Writer, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c.Writer, space)
	}
}

func UpdateSpace(svc *service.SpaceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)
		isSuperAdmin, _ := c.Get("is_super_admin")
		isAdmin := isSuperAdmin.(bool)

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			utils.BadRequest(c.Writer, "ID格式错误")
			return
		}

		var req dto.UpdateSpaceRequest
		if err := c.ShouldBind(&req); err != nil {
			utils.BadRequest(c.Writer, "请求参数错误: "+err.Error())
			return
		}

		var coverFile interface{}
		if file, err := c.FormFile("cover"); err == nil {
			coverFile = file
		}

		space, err := svc.UpdateSpace(uint(id), &req, userID, isAdmin, coverFile)
		if err != nil {
			utils.Error(c.Writer, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c.Writer, space)
	}
}

func DeleteSpace(svc *service.SpaceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)
		isSuperAdmin, _ := c.Get("is_super_admin")
		isAdmin := isSuperAdmin.(bool)

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			utils.BadRequest(c.Writer, "ID格式错误")
			return
		}

		if err := svc.DeleteSpace(uint(id), userID, isAdmin); err != nil {
			utils.Error(c.Writer, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c.Writer, gin.H{"message": "删除成功"})
	}
}

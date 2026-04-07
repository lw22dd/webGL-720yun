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

type HotspotHandler struct {
	service *service.HotspotService
}

func NewHotspotHandler(service *service.HotspotService) *HotspotHandler {
	return &HotspotHandler{service: service}
}

func GetHotspotList(svc *service.HotspotService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.HotspotListRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			utils.BadRequest(c.Writer, "请求参数错误")
			return
		}

		response, err := svc.GetHotspotList(&req)
		if err != nil {
			utils.InternalServerError(c.Writer, err.Error())
			return
		}

		utils.Success(c.Writer, response)
	}
}

func GetHotspotDetail(svc *service.HotspotService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			utils.BadRequest(c.Writer, "ID格式错误")
			return
		}

		response, err := svc.GetHotspotDetail(uint(id))
		if err != nil {
			utils.NotFound(c.Writer, err.Error())
			return
		}

		utils.Success(c.Writer, response)
	}
}

func CreateHotspot(svc *service.HotspotService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)
		isSuperAdmin, _ := c.Get("is_super_admin")
		isAdmin := isSuperAdmin.(bool)

		var req dto.CreateHotspotRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c.Writer, "请求参数错误: "+err.Error())
			return
		}

		hotspot, err := svc.CreateHotspot(&req, userID, isAdmin)
		if err != nil {
			utils.Error(c.Writer, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c.Writer, hotspot)
	}
}

func UpdateHotspot(svc *service.HotspotService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)
		isSuperAdmin, _ := c.Get("is_super_admin")
		isAdmin := isSuperAdmin.(bool)

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			utils.BadRequest(c.Writer, "ID格式错误")
			return
		}

		var req dto.UpdateHotspotRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c.Writer, "请求参数错误: "+err.Error())
			return
		}

		hotspot, err := svc.UpdateHotspot(uint(id), &req, userID, isAdmin)
		if err != nil {
			utils.Error(c.Writer, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c.Writer, hotspot)
	}
}

func DeleteHotspot(svc *service.HotspotService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)
		isSuperAdmin, _ := c.Get("is_super_admin")
		isAdmin := isSuperAdmin.(bool)

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			utils.BadRequest(c.Writer, "ID格式错误")
			return
		}

		if err := svc.DeleteHotspot(uint(id), userID, isAdmin); err != nil {
			utils.Error(c.Writer, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c.Writer, gin.H{"message": "删除成功"})
	}
}

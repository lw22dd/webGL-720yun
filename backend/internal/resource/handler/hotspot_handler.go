package handler

import (
	"net/http"

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

type HotspotIDParamRequest struct {
	ID uint `uri:"id" binding:"required"`
}

func GetHotspotList(svc *service.HotspotService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.HotspotListRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			utils.BadRequest(c, "请求参数错误")
			return
		}

		response, err := svc.GetHotspotList(&req)
		if err != nil {
			utils.InternalServerError(c, err.Error())
			return
		}

		utils.Success(c, response)
	}
}

func GetHotspotDetail(svc *service.HotspotService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req HotspotIDParamRequest
		if err := c.ShouldBindUri(&req); err != nil {
			utils.BadRequest(c, "参数错误")
			return
		}

		response, err := svc.GetHotspotDetail(req.ID)
		if err != nil {
			utils.NotFound(c, err.Error())
			return
		}

		utils.Success(c, response)
	}
}

func CreateHotspot(svc *service.HotspotService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)
		isAdmin := middleware.GetIsAdmin(c)

		var req dto.CreateHotspotRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c, "请求参数错误: "+err.Error())
			return
		}

		hotspot, err := svc.CreateHotspot(&req, userID, isAdmin)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c, hotspot)
	}
}

func UpdateHotspot(svc *service.HotspotService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)
		isAdmin := middleware.GetIsAdmin(c)

		var req HotspotIDParamRequest
		if err := c.ShouldBindUri(&req); err != nil {
			utils.BadRequest(c, "参数错误")
			return
		}

		var updateReq dto.UpdateHotspotRequest
		if err := c.ShouldBindJSON(&updateReq); err != nil {
			utils.BadRequest(c, "请求参数错误: "+err.Error())
			return
		}

		hotspot, err := svc.UpdateHotspot(req.ID, &updateReq, userID, isAdmin)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c, hotspot)
	}
}

func DeleteHotspot(svc *service.HotspotService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)
		isAdmin := middleware.GetIsAdmin(c)

		var req HotspotIDParamRequest
		if err := c.ShouldBindUri(&req); err != nil {
			utils.BadRequest(c, "参数错误")
			return
		}

		if err := svc.DeleteHotspot(req.ID, userID, isAdmin); err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c, gin.H{"message": "删除成功"})
	}
}

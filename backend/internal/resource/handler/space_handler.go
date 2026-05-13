package handler

import (
	"net/http"

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

type SpaceIDRequest struct {
	ID uint `uri:"id" binding:"required"`
}

func GetSpaceList(svc *service.SpaceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SpaceListRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			utils.BadRequest(c, "请求参数错误")
			return
		}

		response, err := svc.GetSpaceList(&req)
		if err != nil {
			utils.InternalServerError(c, err.Error())
			return
		}

		utils.Success(c, response)
	}
}

func GetSpaceDetail(svc *service.SpaceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SpaceIDRequest
		if err := c.ShouldBindUri(&req); err != nil {
			utils.BadRequest(c, "参数错误")
			return
		}

		response, err := svc.GetSpaceDetail(req.ID)
		if err != nil {
			utils.NotFound(c, err.Error())
			return
		}

		utils.Success(c, response)
	}
}

func CreateSpace(svc *service.SpaceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)

		var req dto.CreateSpaceRequest
		if err := c.ShouldBind(&req); err != nil {
			utils.BadRequest(c, "请求参数错误: "+err.Error())
			return
		}

		coverFile, _ := c.FormFile("cover")

		space, err := svc.CreateSpace(&req, userID, coverFile)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c, space)
	}
}

func UpdateSpace(svc *service.SpaceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)
		isAdmin := middleware.GetIsAdmin(c)

		var req SpaceIDRequest
		if err := c.ShouldBindUri(&req); err != nil {
			utils.BadRequest(c, "参数错误")
			return
		}

		var updateReq dto.UpdateSpaceRequest
		if err := c.ShouldBind(&updateReq); err != nil {
			utils.BadRequest(c, "请求参数错误: "+err.Error())
			return
		}

		coverFile, _ := c.FormFile("cover")

		space, err := svc.UpdateSpace(req.ID, &updateReq, userID, isAdmin, coverFile)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c, space)
	}
}

func DeleteSpace(svc *service.SpaceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)
		isAdmin := middleware.GetIsAdmin(c)

		var req SpaceIDRequest
		if err := c.ShouldBindUri(&req); err != nil {
			utils.BadRequest(c, "参数错误")
			return
		}

		if err := svc.DeleteSpace(req.ID, userID, isAdmin); err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c, gin.H{"message": "删除成功"})
	}
}

type DeleteSpaceBatchRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}

func DeleteSpaceBatch(svc *service.SpaceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)
		isAdmin := middleware.GetIsAdmin(c)

		var req DeleteSpaceBatchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c, "请求参数错误: ids不能为空")
			return
		}

		if err := svc.DeleteSpaceBatch(req.IDs, userID, isAdmin); err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c, gin.H{"message": "批量删除成功"})
	}
}

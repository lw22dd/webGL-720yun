package handler

import (
	"net/http"

	"webGL-720yun/internal/core/middleware"
	"webGL-720yun/internal/resource/dto"
	"webGL-720yun/internal/resource/service"
	"webGL-720yun/pkg/utils"

	"github.com/gin-gonic/gin"
)

type SceneHandler struct {
	service *service.SceneService
}

func NewSceneHandler(service *service.SceneService) *SceneHandler {
	return &SceneHandler{service: service}
}

type SceneIDParamRequest struct {
	ID uint `uri:"id" binding:"required"`
}

type SpaceIDParamRequest struct {
	SpaceID uint `uri:"space_id" binding:"required"`
}

func GetSceneList(svc *service.SceneService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SceneListRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			utils.BadRequest(c, "请求参数错误")
			return
		}

		response, err := svc.GetSceneList(&req)
		if err != nil {
			utils.InternalServerError(c, err.Error())
			return
		}

		utils.Success(c, response)
	}
}

func GetSceneDetail(svc *service.SceneService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SceneIDParamRequest
		if err := c.ShouldBindUri(&req); err != nil {
			utils.BadRequest(c, "参数错误")
			return
		}

		response, err := svc.GetSceneDetail(req.ID)
		if err != nil {
			utils.NotFound(c, err.Error())
			return
		}

		utils.Success(c, response)
	}
}

func CreateScene(svc *service.SceneService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)
		isSuperAdmin, _ := c.Get("is_super_admin")
		isAdmin := isSuperAdmin.(bool)

		var req dto.CreateSceneRequest
		if err := c.ShouldBind(&req); err != nil {
			utils.BadRequest(c, "请求参数错误: "+err.Error())
			return
		}

		scene, err := svc.CreateSceneWithFileID(&req, userID, isAdmin)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c, scene)
	}
}

func UpdateScene(svc *service.SceneService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)
		isSuperAdmin, _ := c.Get("is_super_admin")
		isAdmin := isSuperAdmin.(bool)

		var req SceneIDParamRequest
		if err := c.ShouldBindUri(&req); err != nil {
			utils.BadRequest(c, "参数错误")
			return
		}

		var updateReq dto.UpdateSceneRequest
		if err := c.ShouldBind(&updateReq); err != nil {
			utils.BadRequest(c, "请求参数错误: "+err.Error())
			return
		}

		scene, err := svc.UpdateScene(req.ID, &updateReq, userID, isAdmin)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c, scene)
	}
}

func DeleteScene(svc *service.SceneService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)
		isSuperAdmin, _ := c.Get("is_super_admin")
		isAdmin := isSuperAdmin.(bool)

		var req SceneIDParamRequest
		if err := c.ShouldBindUri(&req); err != nil {
			utils.BadRequest(c, "参数错误")
			return
		}

		if err := svc.DeleteScene(req.ID, userID, isAdmin); err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c, gin.H{"message": "删除成功"})
	}
}

func BatchImportScenes(svc *service.SceneService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)
		isSuperAdmin, _ := c.Get("is_super_admin")
		isAdmin := isSuperAdmin.(bool)

		var req SpaceIDParamRequest
		if err := c.ShouldBindUri(&req); err != nil {
			utils.BadRequest(c, "space_id参数错误")
			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			utils.BadRequest(c, "请上传文件")
			return
		}

		response, err := svc.BatchImport(req.SpaceID, file, userID, isAdmin)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c, response)
	}
}

func GetSpaceGraph(svc *service.SceneService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SpaceIDParamRequest
		if err := c.ShouldBindUri(&req); err != nil {
			utils.BadRequest(c, "参数错误")
			return
		}

		response, err := svc.GetSpaceGraphData(req.SpaceID)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c, response)
	}
}

func UpdateScenePosition(svc *service.SceneService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)
		isSuperAdmin, _ := c.Get("is_super_admin")
		isAdmin := isSuperAdmin.(bool)

		var req SceneIDParamRequest
		if err := c.ShouldBindUri(&req); err != nil {
			utils.BadRequest(c, "参数错误")
			return
		}

		var updateReq dto.UpdatePositionRequest
		if err := c.ShouldBindJSON(&updateReq); err != nil {
			utils.BadRequest(c, "请求参数错误: "+err.Error())
			return
		}

		if err := svc.UpdateScenePosition(req.ID, updateReq.Longitude, updateReq.Latitude, userID, isAdmin); err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c, gin.H{"message": "坐标更新成功"})
	}
}

func BatchUpdateScenePosition(svc *service.SceneService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)
		isSuperAdmin, _ := c.Get("is_super_admin")
		isAdmin := isSuperAdmin.(bool)

		var req dto.BatchUpdatePositionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c, "请求参数错误: "+err.Error())
			return
		}

		if err := svc.BatchUpdateScenePosition(req.Positions, userID, isAdmin); err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c, gin.H{"message": "批量坐标更新成功"})
	}
}

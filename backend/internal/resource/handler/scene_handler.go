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

type SceneHandler struct {
	service *service.SceneService
}

func NewSceneHandler(service *service.SceneService) *SceneHandler {
	return &SceneHandler{service: service}
}

func GetSceneList(svc *service.SceneService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SceneListRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			utils.BadRequest(c.Writer, "请求参数错误")
			return
		}

		response, err := svc.GetSceneList(&req)
		if err != nil {
			utils.InternalServerError(c.Writer, err.Error())
			return
		}

		utils.Success(c.Writer, response)
	}
}

func GetSceneDetail(svc *service.SceneService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			utils.BadRequest(c.Writer, "ID格式错误")
			return
		}

		response, err := svc.GetSceneDetail(uint(id))
		if err != nil {
			utils.NotFound(c.Writer, err.Error())
			return
		}

		utils.Success(c.Writer, response)
	}
}

func CreateScene(svc *service.SceneService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)
		isSuperAdmin, _ := c.Get("is_super_admin")
		isAdmin := isSuperAdmin.(bool)

		var req dto.CreateSceneRequest
		if err := c.ShouldBind(&req); err != nil {
			utils.BadRequest(c.Writer, "请求参数错误: "+err.Error())
			return
		}

		panoramaFile, err := c.FormFile("panorama")
		if err != nil {
			utils.BadRequest(c.Writer, "请上传全景图")
			return
		}

		scene, err := svc.CreateScene(&req, userID, isAdmin, panoramaFile)
		if err != nil {
			utils.Error(c.Writer, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c.Writer, scene)
	}
}

func UpdateScene(svc *service.SceneService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)
		isSuperAdmin, _ := c.Get("is_super_admin")
		isAdmin := isSuperAdmin.(bool)

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			utils.BadRequest(c.Writer, "ID格式错误")
			return
		}

		var req dto.UpdateSceneRequest
		if err := c.ShouldBind(&req); err != nil {
			utils.BadRequest(c.Writer, "请求参数错误: "+err.Error())
			return
		}

		scene, err := svc.UpdateScene(uint(id), &req, userID, isAdmin)
		if err != nil {
			utils.Error(c.Writer, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c.Writer, scene)
	}
}

func DeleteScene(svc *service.SceneService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)
		isSuperAdmin, _ := c.Get("is_super_admin")
		isAdmin := isSuperAdmin.(bool)

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			utils.BadRequest(c.Writer, "ID格式错误")
			return
		}

		if err := svc.DeleteScene(uint(id), userID, isAdmin); err != nil {
			utils.Error(c.Writer, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c.Writer, gin.H{"message": "删除成功"})
	}
}

func BatchImportScenes(svc *service.SceneService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)
		isSuperAdmin, _ := c.Get("is_super_admin")
		isAdmin := isSuperAdmin.(bool)

		spaceID, err := strconv.ParseUint(c.PostForm("space_id"), 10, 32)
		if err != nil {
			utils.BadRequest(c.Writer, "space_id参数错误")
			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			utils.BadRequest(c.Writer, "请上传文件")
			return
		}

		response, err := svc.BatchImport(uint(spaceID), file, userID, isAdmin)
		if err != nil {
			utils.Error(c.Writer, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c.Writer, response)
	}
}

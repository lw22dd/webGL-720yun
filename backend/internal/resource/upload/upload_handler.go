package upload

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"webGL-720yun/internal/core/middleware"
	"webGL-720yun/pkg/utils"
)

func InitUpload(svc *UploadService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)

		var req InitUploadRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c, "请求参数错误: "+err.Error())
			return
		}

		response, err := svc.InitUpload(&req, userID)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c, response)
	}
}

func UploadChunk(svc *UploadService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)

		uploadID := c.PostForm("upload_id")
		if uploadID == "" {
			utils.BadRequest(c, "缺少 upload_id 参数")
			return
		}

		chunkIndexStr := c.PostForm("chunk_index")
		if chunkIndexStr == "" {
			utils.BadRequest(c, "缺少 chunk_index 参数")
			return
		}

		chunkIndex, err := strconv.Atoi(chunkIndexStr)
		if err != nil {
			utils.BadRequest(c, "chunk_index 参数格式错误")
			return
		}

		chunkData, err := c.FormFile("chunk_data")
		if err != nil {
			utils.BadRequest(c, "缺少分片数据")
			return
		}

		response, err := svc.UploadChunk(uploadID, chunkIndex, chunkData, userID)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c, response)
	}
}

func CompleteUpload(svc *UploadService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)

		var req CompleteUploadRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c, "请求参数错误: "+err.Error())
			return
		}

		response, err := svc.CompleteUpload(&req, userID)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c, response)
	}
}

func GetUploadStatus(svc *UploadService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)

		uploadID := c.Param("upload_id")
		if uploadID == "" {
			utils.BadRequest(c, "缺少 upload_id 参数")
			return
		}

		response, err := svc.GetUploadStatus(uploadID, userID)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c, response)
	}
}

func CancelUpload(svc *UploadService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)

		uploadID := c.Param("upload_id")
		if uploadID == "" {
			utils.BadRequest(c, "缺少 upload_id 参数")
			return
		}

		err := svc.CancelUpload(uploadID, userID)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c, gin.H{"message": "上传已取消"})
	}
}

func GetFileInfo(svc *UploadService) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileID := c.Param("file_id")
		if fileID == "" {
			utils.BadRequest(c, "缺少 file_id 参数")
			return
		}

		fileInfo, err := svc.GetFileInfo(fileID)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "获取文件信息失败")
			return
		}

		if fileInfo == nil {
			utils.Error(c, http.StatusNotFound, "文件不存在")
			return
		}

		utils.Success(c, fileInfo)
	}
}

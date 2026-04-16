package handler

import (
	"io"
	"net/http"
	"strconv"

	"webGL-720yun/internal/resource/service"

	"github.com/gin-gonic/gin"
)

// GetCubemapFace 流式返回完整的 cubemap 面图片
// GET /api/v1/res/cubemap/:sceneCode/:face
func GetCubemapFace(svc *service.SceneService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sceneCode := c.Param("sceneCode")
		face := c.Param("face")

		// 校验 face
		if !service.IsValidFace(face) {
			c.Status(http.StatusBadRequest)
			return
		}

		// 流式获取完整的 cubemap 面
		result, err := svc.GetCubemapFaceStream(c.Request.Context(), sceneCode, face)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		defer result.Stream.Close()

		// 设置响应头
		streamResponse(c, result)
	}
}

// GetTile 流式返回瓦片图片
// GET /api/v1/res/tiles/:sceneCode/:face/:level/:x/:y
func GetTile(svc *service.SceneService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sceneCode := c.Param("sceneCode")
		face := c.Param("face")

		if !service.IsValidFace(face) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid face"})
			return
		}

		level, err := strconv.Atoi(c.Param("level"))
		if err != nil || level < 0 || level > 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid level"})
			return
		}

		x, err := strconv.Atoi(c.Param("x"))
		if err != nil || x < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid x"})
			return
		}

		y, err := strconv.Atoi(c.Param("y"))
		if err != nil || y < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid y"})
			return
		}

		result, err := svc.GetTileStream(c.Request.Context(), sceneCode, face, level, x, y)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		defer result.Stream.Close()

		streamResponse(c, result)
	}
}

// GetPreview 流式返回场景预览图
// GET /api/v1/res/previews/:sceneCode
func GetPreview(svc *service.SceneService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sceneCode := c.Param("sceneCode")

		result, err := svc.GetPreviewStream(c.Request.Context(), sceneCode)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		defer result.Stream.Close()

		streamResponse(c, result)
	}
}

// GetSource 流式返回场景原始全景图
// GET /api/v1/res/sources/:sceneCode
func GetSource(svc *service.SceneService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sceneCode := c.Param("sceneCode")

		result, err := svc.GetSourceStream(c.Request.Context(), sceneCode)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		defer result.Stream.Close()

		streamResponse(c, result)
	}
}

// GetCover 流式返回空间封面
// GET /api/v1/res/covers/:spaceName
func GetCover(svc *service.SceneService) gin.HandlerFunc {
	return func(c *gin.Context) {
		spaceName := c.Param("spaceName")

		result, err := svc.GetCoverStream(c.Request.Context(), spaceName)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		defer result.Stream.Close()

		streamResponse(c, result)
	}
}

// streamResponse 通用的流式响应写入
// 设置强缓存头 + ETag，然后用 io.Copy 零拷贝转发
func streamResponse(c *gin.Context, result *service.ResourceStreamResult) {
	c.Header("Content-Type", result.ContentType)
	c.Header("Content-Length", strconv.FormatInt(result.Size, 10))
	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.Header("ETag", result.ETag)

	// 支持条件请求：如果客户端的 If-None-Match 匹配 ETag，返回 304
	if match := c.GetHeader("If-None-Match"); match == result.ETag {
		c.Status(http.StatusNotModified)
		return
	}

	c.Status(http.StatusOK)

	// HEAD 请求不需要响应体
	if c.Request.Method == "HEAD" {
		return
	}

	// 🔥 核心：io.Copy 使用固定 32KB 缓冲，不会将整个文件读入内存
	io.Copy(c.Writer, result.Stream)
}

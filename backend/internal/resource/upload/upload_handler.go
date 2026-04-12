package upload

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	gorillaws "github.com/gorilla/websocket"

	"webGL-720yun/internal/core/middleware"
	"webGL-720yun/pkg/jwt"
	"webGL-720yun/pkg/utils"
	"webGL-720yun/pkg/websocket"
)

var upgrader = gorillaws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type UploadHandler struct {
	uploadService *UploadService
	wsHub         *websocket.Hub
	jwtService    *jwt.JWTService
}

func NewUploadHandler(uploadService *UploadService, wsHub *websocket.Hub, jwtService *jwt.JWTService) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
		wsHub:         wsHub,
		jwtService:    jwtService,
	}
}

func InitUpload(svc *UploadService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)

		var req InitUploadRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c.Writer, "请求参数错误: "+err.Error())
			return
		}

		response, err := svc.InitUpload(&req, userID)
		if err != nil {
			utils.Error(c.Writer, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c.Writer, response)
	}
}

func UploadChunk(svc *UploadService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)

		uploadID := c.PostForm("upload_id")
		if uploadID == "" {
			utils.BadRequest(c.Writer, "缺少 upload_id 参数")
			return
		}

		chunkIndexStr := c.PostForm("chunk_index")
		if chunkIndexStr == "" {
			utils.BadRequest(c.Writer, "缺少 chunk_index 参数")
			return
		}

		chunkIndex, err := strconv.Atoi(chunkIndexStr)
		if err != nil {
			utils.BadRequest(c.Writer, "chunk_index 参数格式错误")
			return
		}

		chunkData, err := c.FormFile("chunk_data")
		if err != nil {
			utils.BadRequest(c.Writer, "缺少分片数据")
			return
		}

		response, err := svc.UploadChunk(uploadID, chunkIndex, chunkData, userID)
		if err != nil {
			utils.Error(c.Writer, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c.Writer, response)
	}
}

func MergeChunks(svc *UploadService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)

		var req MergeUploadRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c.Writer, "请求参数错误: "+err.Error())
			return
		}

		response, err := svc.MergeChunks(req.UploadID, userID)
		if err != nil {
			utils.Error(c.Writer, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c.Writer, response)
	}
}

func GetUploadStatus(svc *UploadService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)

		uploadID := c.Param("upload_id")
		if uploadID == "" {
			utils.BadRequest(c.Writer, "缺少 upload_id 参数")
			return
		}

		response, err := svc.GetUploadStatus(uploadID, userID)
		if err != nil {
			utils.Error(c.Writer, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c.Writer, response)
	}
}

func CancelUpload(svc *UploadService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _, _ := middleware.GetCurrentUser(c)

		uploadID := c.Param("upload_id")
		if uploadID == "" {
			utils.BadRequest(c.Writer, "缺少 upload_id 参数")
			return
		}

		err := svc.CancelUpload(uploadID, userID)
		if err != nil {
			utils.Error(c.Writer, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(c.Writer, gin.H{"message": "上传已取消"})
	}
}

func HandleWebSocket(wsHub *websocket.Hub, jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少认证token"})
			return
		}

		claims, err := jwtService.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		client := websocket.NewClient(wsHub, conn, claims.UserID)
		wsHub.Register(client)

		go client.WritePump()
		go client.ReadPump()
	}
}

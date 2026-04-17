package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			log.Printf("[WebSocket] Unauthorized: user_id not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			return
		}

		log.Printf("[WebSocket] New connection request from user_id=%d", userID.(uint))

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("[WebSocket] Upgrade failed: %v", err)
			return
		}

		client := NewClient(hub, conn, userID.(uint))
		hub.Register(client)
		log.Printf("[WebSocket] Client registered: user_id=%d", userID.(uint))

		go client.WritePump()
		go client.ReadPump()
	}
}

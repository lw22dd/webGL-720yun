package slice

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type QueueStatusResponse struct {
	TaskID               string `json:"task_id"`
	UserQueuePosition    int    `json:"user_queue_position"`
	GlobalQueuePosition  int    `json:"global_queue_position"`
	QueueAheadCount      int    `json:"queue_ahead_count"`
	EstimatedWaitSeconds int    `json:"estimated_wait_seconds"`
	ActiveUsers          int    `json:"active_users"`
	Found                bool   `json:"found"`
}

func GetSliceQueueStatus(queue *SliceQueue) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID := c.Param("task_id")
		if taskID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "task_id is required",
			})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		position, err := queue.GetQueuePosition(taskID, userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, &QueueStatusResponse{
			TaskID:               position.TaskID,
			UserQueuePosition:    position.UserQueuePosition,
			GlobalQueuePosition:  position.GlobalQueuePosition,
			QueueAheadCount:      position.QueueAheadCount,
			EstimatedWaitSeconds: position.EstimatedWaitSec,
			ActiveUsers:          position.ActiveUsers,
			Found:                position.Found,
		})
	}
}

func GetSliceQueueStats(queue *SliceQueue) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		stats, err := queue.GetQueueStats(userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, &QueueStatusResponse{
			UserQueuePosition:    stats.UserQueuePosition,
			GlobalQueuePosition:  stats.GlobalQueuePosition,
			QueueAheadCount:      stats.QueueAheadCount,
			EstimatedWaitSeconds: stats.EstimatedWaitSec,
			ActiveUsers:          stats.ActiveUsers,
		})
	}
}

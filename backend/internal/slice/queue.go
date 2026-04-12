package slice

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"

	redisSvc "webGL-720yun/pkg/redis"
)

const (
	SliceStreamKey   = "slice:tasks"
	SliceGroup       = "slice_workers"
	SliceConsumer    = "worker_1"
	SliceTaskTimeout = 24 * time.Hour
)

type SliceTask struct {
	TaskID    string `json:"task_id"`
	SceneID   uint   `json:"scene_id"`
	SceneCode string `json:"scene_code"`
	FileID    string `json:"file_id"`
	UserID    uint   `json:"user_id"`
	CreatedAt int64  `json:"created_at"`
}

type SliceQueue struct {
	redis *redisSvc.RedisService
}

func NewSliceQueue(redisService *redisSvc.RedisService) *SliceQueue {
	return &SliceQueue{redis: redisService}
}

func (q *SliceQueue) PushTask(task *SliceTask) error {
	if task.TaskID == "" {
		task.TaskID = uuid.New().String()
	}
	task.CreatedAt = time.Now().Unix()

	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("序列化任务失败: %w", err)
	}

	client := q.redis.GetClient()

	err = client.XAdd(&redis.XAddArgs{
		Stream: SliceStreamKey,
		Values: map[string]interface{}{
			"data": string(data),
		},
	}).Err()
	if err != nil {
		return fmt.Errorf("发送任务到Stream失败: %w", err)
	}

	return nil
}

func (q *SliceQueue) PopTask() (*SliceTask, string, error) {
	client := q.redis.GetClient()

	err := client.XGroupCreateMkStream(SliceStreamKey, SliceGroup, "0").Err()
	if err != nil {
	}

	streams, err := client.XReadGroup(&redis.XReadGroupArgs{
		Group:    SliceGroup,
		Consumer: SliceConsumer,
		Streams:  []string{SliceStreamKey, ">"},
		Count:    1,
		Block:    time.Second * 5,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, "", nil
		}
		return nil, "", fmt.Errorf("读取任务失败: %w", err)
	}

	if len(streams) == 0 || len(streams[0].Messages) == 0 {
		return nil, "", nil
	}

	msg := streams[0].Messages[0]
	messageID := msg.ID

	dataStr, ok := msg.Values["data"].(string)
	if !ok {
		client.XDel(SliceStreamKey, messageID)
		return nil, "", fmt.Errorf("任务数据格式错误")
	}

	var task SliceTask
	if err := json.Unmarshal([]byte(dataStr), &task); err != nil {
		client.XDel(SliceStreamKey, messageID)
		return nil, "", fmt.Errorf("解析任务数据失败: %w", err)
	}

	return &task, messageID, nil
}

func (q *SliceQueue) AckTask(messageID string) error {
	client := q.redis.GetClient()

	return client.XAck(SliceStreamKey, SliceGroup, messageID).Err()
}

func (q *SliceQueue) DeleteTask(messageID string) error {
	client := q.redis.GetClient()

	return client.XDel(SliceStreamKey, messageID).Err()
}

func (q *SliceQueue) GetPendingTasks() ([]SliceTask, error) {
	client := q.redis.GetClient()

	pending, err := client.XPendingExt(&redis.XPendingExtArgs{
		Stream: SliceStreamKey,
		Group:  SliceGroup,
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("获取待处理任务失败: %w", err)
	}

	var tasks []SliceTask
	for _, p := range pending {
		msgs, err := client.XRange(SliceStreamKey, p.Id, p.Id).Result()
		if err != nil || len(msgs) == 0 {
			continue
		}

		dataStr, ok := msgs[0].Values["data"].(string)
		if !ok {
			continue
		}

		var task SliceTask
		if err := json.Unmarshal([]byte(dataStr), &task); err != nil {
			continue
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

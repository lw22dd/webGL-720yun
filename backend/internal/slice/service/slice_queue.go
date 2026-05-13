package service

import (
	"encoding/json"
	"fmt"
	"strconv"
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

	UserStreamPrefix  = "slice:user:"
	ActiveUsersKey    = "slice:active_users"
	BusyUsersKey      = "slice:busy_users"
	GlobalCounterKey  = "slice:global_counter"
	AvgProcessTimeSec = 30
)

type SliceTask struct {
	TaskID      string `json:"task_id"`
	SceneID     uint   `json:"scene_id"`
	SceneCode   string `json:"scene_code"`
	FileID      string `json:"file_id"`
	SpaceName   string `json:"space_name"`
	SpaceSlug   string `json:"space_slug"`
	UserID      uint   `json:"user_id"`
	CreatedAt   int64  `json:"created_at"`
	GlobalIndex int64  `json:"global_index"`
}

type QueuePosition struct {
	TaskID              string `json:"task_id"`
	UserQueuePosition   int    `json:"user_queue_position"`
	GlobalQueuePosition int    `json:"global_queue_position"`
	QueueAheadCount     int    `json:"queue_ahead_count"`
	EstimatedWaitSec    int    `json:"estimated_wait_seconds"`
	ActiveUsers         int    `json:"active_users"`
	Found               bool   `json:"found"`
}

type SliceQueue struct {
	redis *redisSvc.RedisService
}

func NewSliceQueue(redisService *redisSvc.RedisService) *SliceQueue {
	return &SliceQueue{redis: redisService}
}

func (q *SliceQueue) getUserStreamKey(userID uint) string {
	return fmt.Sprintf("%s%d:tasks", UserStreamPrefix, userID)
}

func (q *SliceQueue) getUserGroupName(userID uint) string {
	return fmt.Sprintf("user_%d_group", userID)
}

func (q *SliceQueue) SetUserBusy(userID uint) error {
	return q.redis.GetClient().SAdd(BusyUsersKey, userID).Err()
}

func (q *SliceQueue) ClearUserBusy(userID uint) error {
	return q.redis.GetClient().SRem(BusyUsersKey, userID).Err()
}

func (q *SliceQueue) ClearAllBusyUsers() error {
	return q.redis.GetClient().Del(BusyUsersKey).Err()
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

	globalIndex, err := client.Incr(GlobalCounterKey).Result()
	if err != nil {
		return fmt.Errorf("生成全局序号失败: %w", err)
	}
	task.GlobalIndex = globalIndex

	userStreamKey := q.getUserStreamKey(task.UserID)
	groupName := q.getUserGroupName(task.UserID)

	err = client.XAdd(&redis.XAddArgs{
		Stream: userStreamKey,
		Values: map[string]interface{}{
			"data": string(data),
		},
	}).Err()
	if err != nil {
		return fmt.Errorf("发送任务到用户Stream失败: %w", err)
	}

	err = client.XGroupCreateMkStream(userStreamKey, groupName, "0").Err()
	if err != nil {
	}

	err = client.ZAdd(ActiveUsersKey, redis.Z{
		Score:  float64(task.CreatedAt),
		Member: task.UserID,
	}).Err()
	if err != nil {
		return fmt.Errorf("更新活跃用户列表失败: %w", err)
	}

	return nil
}

func (q *SliceQueue) PopTask() (*SliceTask, string, error) {
	client := q.redis.GetClient()

	userIDs, err := client.ZRange(ActiveUsersKey, 0, -1).Result()
	if err != nil {
		return nil, "", fmt.Errorf("获取活跃用户列表失败: %w", err)
	}

	if len(userIDs) == 0 {
		return nil, "", nil
	}

	for _, userIDStr := range userIDs {
		userID, _ := strconv.ParseUint(userIDStr, 10, 64)
		uID := uint(userID)

		// 检查用户是否已经有正在处理的任务 (忙碌检查)
		isBusy, err := client.SIsMember(BusyUsersKey, uID).Result()
		if err == nil && isBusy {
			continue // 该用户已忙，跳过，保证公平性（一户一岗）
		}

		userStreamKey := q.getUserStreamKey(uID)
		groupName := q.getUserGroupName(uID)

		err = client.XGroupCreateMkStream(userStreamKey, groupName, "0").Err()
		if err != nil {
		}

		// 优先尝试读取已经分配但未确认的任务 (PEL 恢复)
		// ID="0" 表示读取分配给当前消费者但还未 ACK 的消息
		task, msgID, err := q.readFromStream(userStreamKey, groupName, "0")
		if err == nil && task != nil {
			q.SetUserBusy(uID) // 标记为忙碌
			return task, msgID, nil
		}

		// 尝试读取新任务 (">" 表示未分配给任何消费者的消息)
		task, msgID, err = q.readFromStream(userStreamKey, groupName, ">")
		if err == nil && task != nil {
			q.SetUserBusy(uID) // 标记为忙碌
			return task, msgID, nil
		}
	}

	return nil, "", nil
}

func (q *SliceQueue) readFromStream(streamKey, groupName, id string) (*SliceTask, string, error) {
	client := q.redis.GetClient()
	streams, err := client.XReadGroup(&redis.XReadGroupArgs{
		Group:    groupName,
		Consumer: SliceConsumer,
		Streams:  []string{streamKey, id},
		Count:    1,
		Block:    10 * time.Millisecond,
	}).Result()

	if err != nil {
		return nil, "", err
	}

	if len(streams) == 0 || len(streams[0].Messages) == 0 {
		return nil, "", nil
	}

	msg := streams[0].Messages[0]
	messageID := msg.ID

	dataStr, ok := msg.Values["data"].(string)
	if !ok {
		client.XDel(streamKey, messageID)
		return nil, "", fmt.Errorf("任务数据格式错误")
	}

	var task SliceTask
	if err := json.Unmarshal([]byte(dataStr), &task); err != nil {
		client.XDel(streamKey, messageID)
		return nil, "", fmt.Errorf("解析任务数据失败: %w", err)
	}

	q.updateUserActiveTime(task.UserID)
	return &task, messageID, nil
}


func (q *SliceQueue) checkAndRemoveEmptyUser(userID uint, userStreamKey string) {
	client := q.redis.GetClient()

	length, err := client.XLen(userStreamKey).Result()
	if err != nil {
		return
	}

	if length == 0 {
		client.ZRem(ActiveUsersKey, userID)
	}
}

func (q *SliceQueue) updateUserActiveTime(userID uint) {
	client := q.redis.GetClient()

	userStreamKey := q.getUserStreamKey(userID)
	length, err := client.XLen(userStreamKey).Result()
	if err != nil || length == 0 {
		client.ZRem(userStreamKey, userID)
		return
	}

	msgs, err := client.XRange(userStreamKey, "-", "+").Result()
	if err != nil || len(msgs) == 0 {
		return
	}

	dataStr, ok := msgs[0].Values["data"].(string)
	if !ok {
		return
	}

	var task SliceTask
	if err := json.Unmarshal([]byte(dataStr), &task); err != nil {
		return
	}

	client.ZAdd(ActiveUsersKey, redis.Z{
		Score:  float64(task.CreatedAt),
		Member: userID,
	})
}

func (q *SliceQueue) AckTask(messageID string, userID uint) error {
	client := q.redis.GetClient()

	userStreamKey := q.getUserStreamKey(userID)
	groupName := q.getUserGroupName(userID)

	return client.XAck(userStreamKey, groupName, messageID).Err()
}

func (q *SliceQueue) DeleteTask(messageID string, userID uint) error {
	client := q.redis.GetClient()

	userStreamKey := q.getUserStreamKey(userID)

	err := client.XDel(userStreamKey, messageID).Err()
	if err != nil {
		return err
	}

	length, err := client.XLen(userStreamKey).Result()
	if err == nil && length == 0 {
		client.ZRem(ActiveUsersKey, userID)
	}

	return nil
}

func (q *SliceQueue) GetQueuePosition(taskID string, userID uint) (*QueuePosition, error) {
	client := q.redis.GetClient()

	result := &QueuePosition{
		TaskID: taskID,
		Found:  false,
	}

	userStreamKey := q.getUserStreamKey(userID)

	msgs, err := client.XRange(userStreamKey, "-", "+").Result()
	if err != nil {
		return nil, fmt.Errorf("获取用户队列失败: %w", err)
	}

	userPosition := 0
	found := false
	for i, msg := range msgs {
		dataStr, ok := msg.Values["data"].(string)
		if !ok {
			continue
		}

		var task SliceTask
		if err := json.Unmarshal([]byte(dataStr), &task); err != nil {
			continue
		}

		if task.TaskID == taskID {
			userPosition = i
			found = true
			break
		}
	}

	if !found {
		return result, nil
	}

	result.Found = true
	result.UserQueuePosition = userPosition + 1

	userIDs, err := client.ZRange(ActiveUsersKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("获取活跃用户列表失败: %w", err)
	}

	result.ActiveUsers = len(userIDs)

	globalPosition := 0
	for _, uidStr := range userIDs {
		uid, err := strconv.ParseUint(uidStr, 10, 64)
		if err != nil {
			continue
		}

		otherUserStreamKey := q.getUserStreamKey(uint(uid))

		if uint(uid) == userID {
			globalPosition += userPosition
			break
		}

		length, err := client.XLen(otherUserStreamKey).Result()
		if err != nil {
			continue
		}

		globalPosition += int(length)
	}

	result.GlobalQueuePosition = globalPosition + 1
	result.QueueAheadCount = globalPosition
	result.EstimatedWaitSec = globalPosition * AvgProcessTimeSec

	return result, nil
}

func (q *SliceQueue) GetQueueStats(userID uint) (*QueuePosition, error) {
	client := q.redis.GetClient()

	result := &QueuePosition{}

	userStreamKey := q.getUserStreamKey(userID)

	length, err := client.XLen(userStreamKey).Result()
	if err != nil {
		return nil, fmt.Errorf("获取用户队列长度失败: %w", err)
	}

	result.UserQueuePosition = int(length)

	userIDs, err := client.ZRange(ActiveUsersKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("获取活跃用户列表失败: %w", err)
	}

	result.ActiveUsers = len(userIDs)

	globalPosition := 0
	for _, uidStr := range userIDs {
		uid, err := strconv.ParseUint(uidStr, 10, 64)
		if err != nil {
			continue
		}

		otherUserStreamKey := q.getUserStreamKey(uint(uid))

		if uint(uid) == userID {
			break
		}

		otherLength, err := client.XLen(otherUserStreamKey).Result()
		if err != nil {
			continue
		}

		globalPosition += int(otherLength)
	}

	result.GlobalQueuePosition = globalPosition + 1
	result.QueueAheadCount = globalPosition
	result.EstimatedWaitSec = globalPosition * AvgProcessTimeSec

	return result, nil
}

func (q *SliceQueue) GetPendingTasks() ([]SliceTask, error) {
	client := q.redis.GetClient()

	userIDs, err := client.ZRange(ActiveUsersKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("获取活跃用户列表失败: %w", err)
	}

	var tasks []SliceTask
	for _, userIDStr := range userIDs {
		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			continue
		}

		userStreamKey := q.getUserStreamKey(uint(userID))

		msgs, err := client.XRange(userStreamKey, "-", "+").Result()
		if err != nil {
			continue
		}

		for _, msg := range msgs {
			dataStr, ok := msg.Values["data"].(string)
			if !ok {
				continue
			}

			var task SliceTask
			if err := json.Unmarshal([]byte(dataStr), &task); err != nil {
				continue
			}

			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

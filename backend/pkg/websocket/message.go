package websocket

import "time"

type MessageType string

const (
	MessageTypeProgress      MessageType = "progress"
	MessageTypeComplete      MessageType = "complete"
	MessageTypeError         MessageType = "error"
	MessageTypeMergeStart    MessageType = "merge_start"
	MessageTypeMergeProgress MessageType = "merge_progress"

	MessageTypeSliceProgress    MessageType = "slice_progress"
	MessageTypeSliceComplete    MessageType = "slice_complete"
	MessageTypeSliceError       MessageType = "slice_error"
	MessageTypeSliceQueueStatus MessageType = "slice_queue_status"
)

type Message struct {
	Type      MessageType `json:"type"`
	UploadID  string      `json:"upload_id,omitempty"`
	UserID    uint        `json:"user_id"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

type ProgressData struct {
	UploadedChunks int     `json:"uploaded_chunks"`
	TotalChunks    int     `json:"total_chunks"`
	Percentage     float64 `json:"percentage"`
	Speed          string  `json:"speed"`
	UploadedBytes  int64   `json:"uploaded_bytes"`
	TotalBytes     int64   `json:"total_bytes"`
}

type MergeProgressData struct {
	Stage      string  `json:"stage"`
	Percentage float64 `json:"percentage"`
	Message    string  `json:"message"`
}

type CompleteData struct {
	SceneID      uint   `json:"scene_id"`
	SourceURL    string `json:"source_url"`
	ThumbnailURL string `json:"thumbnail_url"`
}

type ErrorData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SliceProgressData struct {
	TaskID               string `json:"task_id"`
	SceneID              uint   `json:"scene_id"`
	SceneCode            string `json:"scene_code"`
	Status               string `json:"status"`
	Progress             int    `json:"progress"`
	Stage                string `json:"stage"`
	Message              string `json:"message"`
	Timestamp            int64  `json:"timestamp"`
	QueueAheadCount      int    `json:"queue_ahead_count"`
	EstimatedWaitSeconds int    `json:"estimated_wait_seconds"`
	ActiveUsers          int    `json:"active_users"`
}

type SliceCompleteData struct {
	TaskID     string `json:"task_id"`
	SceneID    uint   `json:"scene_id"`
	SceneCode  string `json:"scene_code"`
	Status     string `json:"status"`
	TileURL    string `json:"tile_url"`
	PreviewURL string `json:"preview_url"`
	Timestamp  int64  `json:"timestamp"`
}

type SliceErrorData struct {
	TaskID    string `json:"task_id"`
	SceneID   uint   `json:"scene_id"`
	SceneCode string `json:"scene_code"`
	Error     string `json:"error"`
	Timestamp int64  `json:"timestamp"`
}

type SliceQueueStatusData struct {
	TaskID               string `json:"task_id"`
	UserQueuePosition    int    `json:"user_queue_position"`
	GlobalQueuePosition  int    `json:"global_queue_position"`
	QueueAheadCount      int    `json:"queue_ahead_count"`
	EstimatedWaitSeconds int    `json:"estimated_wait_seconds"`
	ActiveUsers          int    `json:"active_users"`
	Found                bool   `json:"found"`
}

func NewMessage(msgType MessageType, uploadID string, userID uint, data interface{}) *Message {
	return &Message{
		Type:      msgType,
		UploadID:  uploadID,
		UserID:    userID,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

func NewProgressMessage(uploadID string, userID uint, data *ProgressData) *Message {
	return NewMessage(MessageTypeProgress, uploadID, userID, data)
}

func NewCompleteMessage(uploadID string, userID uint, data *CompleteData) *Message {
	return NewMessage(MessageTypeComplete, uploadID, userID, data)
}

func NewErrorMessage(uploadID string, userID uint, data *ErrorData) *Message {
	return NewMessage(MessageTypeError, uploadID, userID, data)
}

func NewMergeStartMessage(uploadID string, userID uint) *Message {
	return NewMessage(MessageTypeMergeStart, uploadID, userID, nil)
}

func NewMergeProgressMessage(uploadID string, userID uint, data *MergeProgressData) *Message {
	return NewMessage(MessageTypeMergeProgress, uploadID, userID, data)
}

func NewSliceProgressMessage(taskID string, userID uint, data *SliceProgressData) *Message {
	return &Message{
		Type:      MessageTypeSliceProgress,
		UploadID:  taskID,
		UserID:    userID,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

func NewSliceCompleteMessage(taskID string, userID uint, data *SliceCompleteData) *Message {
	return &Message{
		Type:      MessageTypeSliceComplete,
		UploadID:  taskID,
		UserID:    userID,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

func NewSliceErrorMessage(taskID string, userID uint, data *SliceErrorData) *Message {
	return &Message{
		Type:      MessageTypeSliceError,
		UploadID:  taskID,
		UserID:    userID,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

func NewSliceQueueStatusMessage(taskID string, userID uint, data *SliceQueueStatusData) *Message {
	return &Message{
		Type:      MessageTypeSliceQueueStatus,
		UploadID:  taskID,
		UserID:    userID,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

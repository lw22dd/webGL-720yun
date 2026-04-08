package websocket

import "time"

type MessageType string

const (
	MessageTypeProgress      MessageType = "progress"
	MessageTypeComplete      MessageType = "complete"
	MessageTypeError         MessageType = "error"
	MessageTypeMergeStart    MessageType = "merge_start"
	MessageTypeMergeProgress MessageType = "merge_progress"
)

type Message struct {
	Type      MessageType `json:"type"`
	UploadID  string      `json:"upload_id"`
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

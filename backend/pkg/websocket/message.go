package websocket

import "time"

type MessageType string

const (
	MessageTypeProgress      MessageType = "progress"
	MessageTypeComplete      MessageType = "complete"
	MessageTypeError         MessageType = "error"
	MessageTypeMergeStart    MessageType = "merge_start"
	MessageTypeMergeProgress MessageType = "merge_progress"
<<<<<<< HEAD
=======

	MessageTypeSliceProgress MessageType = "slice_progress"
	MessageTypeSliceComplete MessageType = "slice_complete"
	MessageTypeSliceError    MessageType = "slice_error"
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
)

type Message struct {
	Type      MessageType `json:"type"`
<<<<<<< HEAD
	UploadID  string      `json:"upload_id"`
=======
	UploadID  string      `json:"upload_id,omitempty"`
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
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

<<<<<<< HEAD
=======
type SliceProgressData struct {
	TaskID    string `json:"task_id"`
	SceneID   uint   `json:"scene_id"`
	Status    string `json:"status"`
	Progress  int    `json:"progress"`
	Stage     string `json:"stage"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type SliceCompleteData struct {
	TaskID     string `json:"task_id"`
	SceneID    uint   `json:"scene_id"`
	Status     string `json:"status"`
	TileURL    string `json:"tile_url"`
	PreviewURL string `json:"preview_url"`
	Timestamp  int64  `json:"timestamp"`
}

type SliceErrorData struct {
	TaskID    string `json:"task_id"`
	SceneID   uint   `json:"scene_id"`
	Error     string `json:"error"`
	Timestamp int64  `json:"timestamp"`
}

>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
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
<<<<<<< HEAD
=======

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
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503

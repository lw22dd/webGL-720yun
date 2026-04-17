package slice

import "time"

type SliceProgressData struct {
	TaskID    string `json:"task_id"`
	SceneID   uint   `json:"scene_id"`
	SceneCode string `json:"scene_code"`
	Status    string `json:"status"`
	Progress  int    `json:"progress"`
	Stage     string `json:"stage"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
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

type TaskStatusResponse struct {
	TaskID       string    `json:"task_id"`
	SceneID      uint      `json:"scene_id"`
	SceneCode    string    `json:"scene_code"`
	Status       string    `json:"status"`
	Progress     int       `json:"progress"`
	Stage        string    `json:"stage"`
	Message      string    `json:"message"`
	CreatedAt    time.Time `json:"created_at"`
	StartedAt    time.Time `json:"started_at,omitempty"`
	CompletedAt  time.Time `json:"completed_at,omitempty"`
	ErrorMessage string    `json:"error_message,omitempty"`
}

const (
	StageDownloading = "downloading"
	StagePreview     = "preview"
	StageE2C         = "e2c"
	StageTiles       = "tiles"
	StageUploading   = "uploading"
	StageComplete    = "complete"
)

const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)

func NewSliceProgressData(taskID string, sceneID uint, sceneCode string, progress int, stage string, message string) *SliceProgressData {
	return &SliceProgressData{
		TaskID:    taskID,
		SceneID:   sceneID,
		SceneCode: sceneCode,
		Status:    StatusProcessing,
		Progress:  progress,
		Stage:     stage,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}
}

func NewSliceCompleteData(taskID string, sceneID uint, sceneCode string, tileURL string, previewURL string) *SliceCompleteData {
	return &SliceCompleteData{
		TaskID:     taskID,
		SceneID:    sceneID,
		SceneCode:  sceneCode,
		Status:     StatusCompleted,
		TileURL:    tileURL,
		PreviewURL: previewURL,
		Timestamp:  time.Now().Unix(),
	}
}

func NewSliceErrorData(taskID string, sceneID uint, sceneCode string, errorMsg string) *SliceErrorData {
	return &SliceErrorData{
		TaskID:    taskID,
		SceneID:   sceneID,
		SceneCode: sceneCode,
		Error:     errorMsg,
		Timestamp: time.Now().Unix(),
	}
}

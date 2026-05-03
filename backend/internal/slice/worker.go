package slice

import (
	"context"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"

	"webGL-720yun/internal/model"
	"webGL-720yun/pkg/minio_client"
	"webGL-720yun/pkg/websocket"
)

type SliceScheduler struct {
	queue       *SliceQueue
	db          *gorm.DB
	processor   *SliceProcessor
	workerCount int
	semaphore   chan struct{}
	stopChan    chan struct{}
	wg          sync.WaitGroup
}

func NewSliceScheduler(
	redisQueue *SliceQueue,
	db *gorm.DB,
	minioClient *minio_client.MinIOClient,
	wsHub *websocket.Hub,
	workerCount int,
) *SliceScheduler {
	processor := NewSliceProcessor(db, minioClient, wsHub)
	processor.SetQueue(redisQueue)

	return &SliceScheduler{
		queue:       redisQueue,
		db:          db,
		processor:   processor,
		workerCount: workerCount,
		semaphore:   make(chan struct{}, workerCount),
		stopChan:    make(chan struct{}),
	}
}

func (s *SliceScheduler) Start() {
	log.Printf("Slice Scheduler started with %d concurrent workers", s.workerCount)

	for {
		select {
		case <-s.stopChan:
			log.Println("Slice Scheduler stopping...")
			s.wg.Wait()
			log.Println("Slice Scheduler stopped")
			return
		default:
			task, messageID, err := s.queue.PopTask()
			if err != nil {
				log.Printf("Error popping task: %v", err)
				time.Sleep(time.Second)
				continue
			}

			if task == nil {
				time.Sleep(500 * time.Millisecond)
				continue
			}

			s.semaphore <- struct{}{}

			s.wg.Add(1)
			go func(t *SliceTask, msgID string) {
				defer s.wg.Done()
				defer func() { <-s.semaphore }()

				s.processTask(t, msgID)
			}(task, messageID)
		}
	}
}

func (s *SliceScheduler) Stop() {
	close(s.stopChan)
}

func (s *SliceScheduler) processTask(task *SliceTask, messageID string) {
	log.Printf("Processing slice task: task_id=%s, scene_id=%d, scene_code=%s, user_id=%d",
		task.TaskID, task.SceneID, task.SceneCode, task.UserID)

	ctx := context.Background()

	if err := s.updateSceneSlicingStatus(task.SceneID, task.TaskID); err != nil {
		log.Printf("Failed to update scene slicing status: %v", err)
	}

	if err := s.processor.Process(ctx, task); err != nil {
		log.Printf("Slice task failed: task_id=%s, error=%v", task.TaskID, err)
		s.handleTaskError(task, err)
	}

	if err := s.queue.AckTask(messageID, task.UserID); err != nil {
		log.Printf("Failed to ack task: %v", err)
	}

	if err := s.queue.DeleteTask(messageID, task.UserID); err != nil {
		log.Printf("Failed to delete task: %v", err)
	}

	log.Printf("Slice task completed: task_id=%s", task.TaskID)
}

func (s *SliceScheduler) updateSceneSlicingStatus(sceneID uint, taskID string) error {
	return s.db.Model(&model.ResScene{}).Where("id = ?", sceneID).Updates(map[string]interface{}{
		"slice_status": model.SliceStatusSlicing,
		"task_id":      taskID,
		"updated_at":   time.Now(),
	}).Error
}

func (s *SliceScheduler) handleTaskError(task *SliceTask, err error) {
	if dbErr := s.db.Model(&model.ResScene{}).Where("id = ?", task.SceneID).Updates(map[string]interface{}{
		"slice_status": model.SliceStatusFailed,
		"updated_at":   time.Now(),
	}).Error; dbErr != nil {
		log.Printf("Failed to update scene failed status: %v", dbErr)
	}

	s.processor.notifyError(task, err.Error())
}

func (s *SliceScheduler) GetActiveWorkerCount() int {
	return len(s.semaphore)
}

func (s *SliceScheduler) GetQueueStats(userID uint) (*QueuePosition, error) {
	return s.queue.GetQueueStats(userID)
}

func (s *SliceScheduler) GetQueuePosition(taskID string, userID uint) (*QueuePosition, error) {
	return s.queue.GetQueuePosition(taskID, userID)
}

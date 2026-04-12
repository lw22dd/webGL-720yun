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

type SliceWorker struct {
	queue     *SliceQueue
	processor *SliceProcessor
	db        *gorm.DB
	stopChan  chan struct{}
	wg        sync.WaitGroup
}

func NewSliceWorker(
	redisQueue *SliceQueue,
	db *gorm.DB,
	minioClient *minio_client.MinIOClient,
	wsHub *websocket.Hub,
) *SliceWorker {
	return &SliceWorker{
		queue:     redisQueue,
		processor: NewSliceProcessor(db, minioClient, wsHub),
		db:        db,
		stopChan:  make(chan struct{}),
	}
}

func (w *SliceWorker) Start() {
	log.Println("Slice Worker started")

	for {
		select {
		case <-w.stopChan:
			log.Println("Slice Worker stopping...")
			w.wg.Wait()
			log.Println("Slice Worker stopped")
			return
		default:
			task, messageID, err := w.queue.PopTask()
			if err != nil {
				log.Printf("Error popping task: %v", err)
				time.Sleep(time.Second)
				continue
			}

			if task == nil {
				time.Sleep(500 * time.Millisecond)
				continue
			}

			w.wg.Add(1)
			go w.processTask(task, messageID)
		}
	}
}

func (w *SliceWorker) Stop() {
	close(w.stopChan)
}

func (w *SliceWorker) processTask(task *SliceTask, messageID string) {
	defer w.wg.Done()

	log.Printf("Processing slice task: task_id=%s, scene_id=%d, scene_code=%s",
		task.TaskID, task.SceneID, task.SceneCode)

	ctx := context.Background()

	if err := w.updateSceneSlicingStatus(task.SceneID, task.TaskID); err != nil {
		log.Printf("Failed to update scene slicing status: %v", err)
	}

	if err := w.processor.Process(ctx, task); err != nil {
		log.Printf("Slice task failed: task_id=%s, error=%v", task.TaskID, err)
		w.handleTaskError(task, err)
	}

	if err := w.queue.AckTask(messageID); err != nil {
		log.Printf("Failed to ack task: %v", err)
	}

	if err := w.queue.DeleteTask(messageID); err != nil {
		log.Printf("Failed to delete task: %v", err)
	}

	log.Printf("Slice task completed: task_id=%s", task.TaskID)
}

func (w *SliceWorker) updateSceneSlicingStatus(sceneID uint, taskID string) error {
	return w.db.Model(&model.ResScene{}).Where("id = ?", sceneID).Updates(map[string]interface{}{
		"slice_status": model.SliceStatusSlicing,
		"task_id":      taskID,
		"updated_at":   "NOW()",
	}).Error
}

func (w *SliceWorker) handleTaskError(task *SliceTask, err error) {
	if dbErr := w.db.Model(&model.ResScene{}).Where("id = ?", task.SceneID).Updates(map[string]interface{}{
		"slice_status": model.SliceStatusFailed,
		"updated_at":   "NOW()",
	}).Error; dbErr != nil {
		log.Printf("Failed to update scene failed status: %v", dbErr)
	}

	w.processor.notifyError(task, err.Error())
}

type WorkerPool struct {
	workers  []*SliceWorker
	stopChan chan struct{}
}

func NewWorkerPool(
	queue *SliceQueue,
	db *gorm.DB,
	minioClient *minio_client.MinIOClient,
	wsHub *websocket.Hub,
	workerCount int,
) *WorkerPool {
	pool := &WorkerPool{
		workers:  make([]*SliceWorker, workerCount),
		stopChan: make(chan struct{}),
	}

	for i := 0; i < workerCount; i++ {
		pool.workers[i] = NewSliceWorker(queue, db, minioClient, wsHub)
	}

	return pool
}

func (p *WorkerPool) Start() {
	log.Printf("Starting worker pool with %d workers", len(p.workers))

	for _, worker := range p.workers {
		go worker.Start()
	}
}

func (p *WorkerPool) Stop() {
	log.Println("Stopping worker pool...")

	for _, worker := range p.workers {
		worker.Stop()
	}

	close(p.stopChan)
	log.Println("Worker pool stopped")
}

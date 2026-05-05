package slice

import (
	"context"
	"log"
	"sync"
	"time"

	"webGL-720yun/internal/slice/service"
)

type SliceScheduler struct {
	queue       *service.SliceQueue
	processor   *SliceProcessor
	workerCount int
	semaphore   chan struct{}
	stopChan    chan struct{}
	wg          sync.WaitGroup
}

func NewSliceScheduler(queue *service.SliceQueue, processor *SliceProcessor, workerCount int) *SliceScheduler {
	return &SliceScheduler{
		queue:       queue,
		processor:   processor,
		workerCount: workerCount,
		semaphore:   make(chan struct{}, workerCount),
		stopChan:    make(chan struct{}),
	}
}

func (s *SliceScheduler) Start() {
	log.Printf("[Scheduler] 启动切片调度器，并发数: %d", s.workerCount)

	for {
		select {
		case <-s.stopChan:
			log.Println("[Scheduler] 调度器收到停止信号")
			return
		default:
		}

		task, messageID, err := s.queue.PopTask()
		if err != nil {
			log.Printf("[Scheduler] 获取任务失败: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if task == nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		s.semaphore <- struct{}{}
		s.wg.Add(1)

		go func(t *service.SliceTask, msgID string) {
			defer func() {
				<-s.semaphore
				s.wg.Done()
			}()

			log.Printf("[Scheduler] 开始处理任务: %s (场景: %s)", t.TaskID, t.SceneCode)

			ctx := context.Background()
			if err := s.processor.ProcessWithRetry(ctx, t); err != nil {
				log.Printf("[Scheduler] 任务处理失败: %s, error: %v", t.TaskID, err)
				s.processor.notifyError(t, err.Error())
			} else {
				log.Printf("[Scheduler] 任务处理成功: %s", t.TaskID)
			}

			if err := s.queue.AckTask(msgID, t.UserID); err != nil {
				log.Printf("[Scheduler] Ack任务失败: %s, error: %v", t.TaskID, err)
			}

			if err := s.queue.DeleteTask(msgID, t.UserID); err != nil {
				log.Printf("[Scheduler] 删除任务失败: %s, error: %v", t.TaskID, err)
			}
		}(task, messageID)
	}
}

func (s *SliceScheduler) Stop() {
	close(s.stopChan)
	s.wg.Wait()
	log.Println("[Scheduler] 调度器已停止")
}

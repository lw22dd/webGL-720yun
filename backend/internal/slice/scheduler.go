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
	stopChan    chan struct{}
	wg          sync.WaitGroup
}

func NewSliceScheduler(queue *service.SliceQueue, processor *SliceProcessor, workerCount int) *SliceScheduler {
	return &SliceScheduler{
		queue:       queue,
		processor:   processor,
		workerCount: workerCount,
		stopChan:    make(chan struct{}),
	}
}

func (s *SliceScheduler) Start() {
	log.Printf("[Scheduler] 启动切片调度器，并发数: %d", s.workerCount)

	// 启动前先清理所有忙碌状态，防止上次宕机残留
	if err := s.queue.ClearAllBusyUsers(); err != nil {
		log.Printf("[Scheduler] 清理忙碌用户状态失败: %v", err)
	}

	// 启动固定数量的 Worker 协程
	for i := 0; i < s.workerCount; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}

	// 等待停止信号
	<-s.stopChan
	log.Println("[Scheduler] 调度器正在关闭，等待所有 Worker 完成...")
	s.wg.Wait()
	log.Println("[Scheduler] 调度器已彻底停止")
}

func (s *SliceScheduler) worker(id int) {
	defer s.wg.Done()
	log.Printf("[Worker %d] 已启动", id)

	for {
		select {
		case <-s.stopChan:
			log.Printf("[Worker %d] 收到停止信号，退出", id)
			return
		default:
		}

		// 阻塞式获取任务
		task, messageID, err := s.queue.PopTask()
		if err != nil {
			log.Printf("[Worker %d] 获取任务失败: %v", id, err)
			time.Sleep(1 * time.Second)
			continue
		}

		if task == nil {
			// 没有任务，休眠一会儿继续
			time.Sleep(500 * time.Millisecond)
			continue
		}

		s.handleTask(id, task, messageID)
	}
}

func (s *SliceScheduler) handleTask(workerID int, t *service.SliceTask, msgID string) {
	// 无论任务成功与否，最后都要清理该用户的忙碌状态
	defer func() {
		if err := s.queue.ClearUserBusy(t.UserID); err != nil {
			log.Printf("[Worker %d] 清理用户忙碌状态失败: %s, error: %v", workerID, t.TaskID, err)
		}
	}()

	log.Printf("[Worker %d] 开始处理任务: %s (场景: %s)", workerID, t.TaskID, t.SceneCode)

	ctx := context.Background()
	if err := s.processor.ProcessWithRetry(ctx, t); err != nil {
		log.Printf("[Worker %d] 任务处理失败: %s, error: %v", workerID, t.TaskID, err)
		s.processor.notifyError(t, err.Error())
	} else {
		log.Printf("[Worker %d] 任务处理成功: %s", workerID, t.TaskID)
	}

	// 确认并删除任务
	if err := s.queue.AckTask(msgID, t.UserID); err != nil {
		log.Printf("[Worker %d] Ack任务失败: %s, error: %v", workerID, t.TaskID, err)
	}

	if err := s.queue.DeleteTask(msgID, t.UserID); err != nil {
		log.Printf("[Worker %d] 删除任务失败: %s, error: %v", workerID, t.TaskID, err)
	}
}

func (s *SliceScheduler) Stop() {
	close(s.stopChan)
}


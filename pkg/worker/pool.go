package worker

import (
	"context"
	"sync"
)

// Task 任务接口
type Task interface {
	Execute() error
}

// Pool 工作池
type Pool struct {
	maxWorkers int
	taskQueue  chan Task
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewPool 创建工作池
func NewPool(maxWorkers int, queueSize int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())

	return &Pool{
		maxWorkers: maxWorkers,
		taskQueue:  make(chan Task, queueSize),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start 启动工作池
func (p *Pool) Start() {
	for i := 0; i < p.maxWorkers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

// worker 工作协程
func (p *Pool) worker() {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		case task, ok := <-p.taskQueue:
			if !ok {
				return
			}
			if err := task.Execute(); err != nil {
				// 错误已在任务内部处理
				continue
			}
		}
	}
}

// Submit 提交任务
func (p *Pool) Submit(task Task) error {
	select {
	case <-p.ctx.Done():
		return p.ctx.Err()
	case p.taskQueue <- task:
		return nil
	}
}

// Stop 停止工作池
func (p *Pool) Stop() {
	p.cancel()
	close(p.taskQueue)
	p.wg.Wait()
}

// QueueSize 返回队列中待处理任务数量
func (p *Pool) QueueSize() int {
	return len(p.taskQueue)
}

// IsFull 检查队列是否已满
func (p *Pool) IsFull() bool {
	return cap(p.taskQueue) == len(p.taskQueue)
}

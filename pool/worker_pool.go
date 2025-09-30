package pool

import (
	"fmt"
	"sync/atomic"
	"time"
)

type requiredWorkerPool interface {
	Submit(task func()) error
	Stop() error
}

type WorkerPool struct {
	workersCount     int
	freeWorkersCount atomic.Uint32
	isStopped        bool

	taskQueueSize int
	taskQueue     chan func()

	doneCallback func()
}

var _ requiredWorkerPool = (*WorkerPool)(nil)

func NewWorkerPool(workersCount, queueSize int) (*WorkerPool, error) {
	p := &WorkerPool{}

	err := p.setTaskQueueSize(queueSize)
	if err != nil {
		return nil, fmt.Errorf("set task queue size: %w", err)
	}

	err = p.setWorkersCount(workersCount)
	if err != nil {
		return nil, fmt.Errorf("set workers count: %w", err)
	}

	p.setTaskQueue()
	p.startWorkers()

	return p, nil
}

func (p *WorkerPool) WorkersCount() int {
	return p.workersCount
}

func (p *WorkerPool) TaskQueueSize() int {
	return p.taskQueueSize
}

func (p *WorkerPool) IsStopped() bool {
	return p.isStopped
}

func (p *WorkerPool) FreeWorkersCount() int {
	return int(p.freeWorkersCount.Load())
}

func (p *WorkerPool) setWorkersCount(count int) error {
	err := validateWorkersCount(count)
	if err != nil {
		return fmt.Errorf("validate workers count: %w", err)
	}

	p.workersCount = count

	return nil
}

func (p *WorkerPool) setTaskQueueSize(taskQueueSize int) error {
	err := validateTaskQueueSize(taskQueueSize)
	if err != nil {
		return fmt.Errorf("validate task queue size: %w", err)
	}

	p.taskQueueSize = taskQueueSize

	return nil
}

func (p *WorkerPool) SetDoneCallback(doneCallback func()) error {
	err := validateDoneCallback(doneCallback)
	if err != nil {
		return fmt.Errorf("validate done callback: %w", err)
	}

	p.doneCallback = doneCallback

	return nil
}

func (p *WorkerPool) setTaskQueue() {
	p.taskQueue = make(chan func(), p.taskQueueSize)
}

func (p *WorkerPool) Submit(task func()) error {
	if task == nil {
		return ErrInvalidTask
	}

	if p.IsStopped() {
		return ErrPoolStopped
	}

	select {
	case p.taskQueue <- task:
		return nil
	default:
		return ErrQueueFull
	}
}

func (p *WorkerPool) Stop() error {
	if p.isStopped {
		return ErrPoolStopped
	}

	p.isStopped = true

	for len(p.taskQueue) > 0 {
		time.Sleep(5 * time.Millisecond)
	}

	for p.FreeWorkersCount() != p.WorkersCount() {
		time.Sleep(5 * time.Millisecond)
	}

	close(p.taskQueue)

	return nil
}

func (p *WorkerPool) startWorkers() {
	p.freeWorkersCount.Store(uint32(p.workersCount))

	// Возможно использовать wait-групп или err-групп,
	// но тогда решение было бы слишком верхнеуровневым для тестового задания
	for i := 0; i < p.workersCount; i++ {
		go worker(p)
	}
}

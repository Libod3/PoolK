package pool

import (
	"poolk/pool"

	"sync"
	"testing"
	"time"
)

func TestNewWorkerPoolInvalidWorkers(t *testing.T) {
	_, err := pool.NewWorkerPool(0, 1)
	if err == nil {
		t.Errorf("expected error for workers=0, got nil")
	}
}

func TestNewWorkerPoolInvalidQueueSize(t *testing.T) {
	_, err := pool.NewWorkerPool(1, 0)
	if err == nil {
		t.Errorf("expected error for queueSize=0, got nil")
	}
}

func TestSubmitNilTask(t *testing.T) {
	p, err := pool.NewWorkerPool(1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer p.Stop()

	err = p.Submit(nil)
	if err != pool.ErrInvalidTask {
		t.Errorf("expected ErrInvalidTask, got %v", err)
	}
}

func TestSubmitAndExecute(t *testing.T) {
	p, err := pool.NewWorkerPool(2, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer p.Stop()

	var mu sync.Mutex
	executed := 0

	task := func() {
		mu.Lock()
		executed++
		mu.Unlock()
	}

	err = p.Submit(task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = p.Submit(task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if executed != 2 {
		t.Errorf("expected 2 executed tasks, got %d", executed)
	}
}

func TestQueueFull(t *testing.T) {
	p, err := pool.NewWorkerPool(1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer p.Stop()

	started := make(chan struct{})
	release := make(chan struct{})

	blockingTask := func() {
		close(started)
		<-release
	}

	err = p.Submit(blockingTask)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	<-started

	err = p.Submit(func() {})
	if err != nil {
		t.Fatalf("unexpected error on second submit: %v", err)
	}

	err = p.Submit(func() {})
	if err != pool.ErrQueueFull {
		t.Errorf("expected ErrQueueFull, got %v", err)
	}

	close(release)
}

func TestDoneCallback(t *testing.T) {
	p, err := pool.NewWorkerPool(1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer p.Stop()

	var wg sync.WaitGroup
	wg.Add(2)

	cbErr := p.SetDoneCallback(func() {
		wg.Done()
	})
	if cbErr != nil {
		t.Fatalf("unexpected error: %v", cbErr)
	}

	task := func() {}
	err = p.Submit(task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = p.Submit(task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	waitCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitCh)
	}()

	select {
	case <-waitCh:
	case <-time.After(1 * time.Second):
		t.Fatal("done callback was not called expected number of times")
	}
}

func TestStopTwice(t *testing.T) {
	p, err := pool.NewWorkerPool(1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = p.Stop()
	if err != nil {
		t.Errorf("unexpected error on first Stop: %v", err)
	}

	err = p.Stop()
	if err != pool.ErrPoolStopped {
		t.Errorf("expected ErrPoolStopped on second Stop, got %v", err)
	}
}

func TestSubmitAfterStop(t *testing.T) {
	p, err := pool.NewWorkerPool(1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = p.Stop()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = p.Submit(func() {})
	if err != pool.ErrPoolStopped {
		t.Errorf("expected ErrPoolStopped, got %v", err)
	}
}

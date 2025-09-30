package pool

import (
	"errors"
)

var (
	ErrQueueFull   = errors.New("queue is full")
	ErrPoolStopped = errors.New("pool stopped")
	ErrValidation  = errors.New("validation error")
	ErrInvalidTask = errors.New("task error")
)

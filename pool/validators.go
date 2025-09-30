package pool

import "fmt"

func validateWorkersCount(count int) error {
	if count <= 0 {
		return fmt.Errorf(
			"workers count must be greater than 0, got %d: %w",
			count,
			ErrValidation,
		)
	}

	return nil
}

func validateTaskQueueSize(size int) error {
	if size <= 0 {
		return fmt.Errorf(
			"queue size must be greater than 0, got %d: %w",
			size,
			ErrValidation,
		)
	}

	return nil
}

func validateDoneCallback(fn func()) error {
	if fn == nil {
		return fmt.Errorf(
			"done callback should not be nil: %w",
			ErrValidation,
		)
	}

	return nil
}

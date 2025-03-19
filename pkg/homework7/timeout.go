package homework7

import (
	"context"
	"errors"
	"time"
)

// Timeout calls a func with context.Timeout and either returns it's result or an error after the timeout
func Timeout(operation func() error, timeout uint) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancelFunc()

	resultChan := make(chan error, 1)

	go func() {
		resultChan <- operation()
	}()

	select {
	case <-ctx.Done():
		return errors.New("timeout")
	case result := <-resultChan:
		return result
	}
}

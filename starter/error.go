package starter

import (
	"fmt"
	"sync"
)

type errorRecorder struct {
	errors []error
	mutex  sync.Mutex
}

func newErrorRecorder() *errorRecorder {
	return &errorRecorder{}
}

func (r *errorRecorder) record(err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if err == nil {
		return
	}

	r.errors = append(r.errors, err)
}

func (r *errorRecorder) buildFromRecorded() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if len(r.errors) == 0 {
		return nil
	}

	return &recordedError{errors: r.errors}
}

type recordedError struct {
	errors []error
}

func (r *recordedError) Error() string {
	return fmt.Sprintf("recorded: %v", r.errors)
}

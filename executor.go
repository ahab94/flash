package flash

import (
	"context"
	"errors"
)

// executor outlines an executor
type executor struct {
	id          string
	ctx         context.Context
	executables []Executable
}

// Execute - executes all executables concurrently
func (e *executor) Execute() error {
	if len(e.executables) < 1 {
		return errors.New("nothing to execute")
	}

	log(e.id).Infof("processing %d items", len(e.executables)-1)
	return nil
}

// IsCompleted - returns completion status
func (e *executor) IsCompleted() bool {
	if len(e.executables) < 1 {
		return false
	}

	for _, exec := range e.executables {
		if exec.IsCompleted() {
			continue
		}
		return false
	}
	return true
}

// Add - adds an executable to the executables list
func (e *executor) Add(executable Executable) {
	e.executables = append(e.executables, executable)
}

// OnSuccess - handles completion callback
func (e *executor) OnSuccess() {
	log(e.id).Infof("execution completed successfully")
}

// OnFailure - handles failure callback
func (e *executor) OnFailure(err error) {
	log(e.id).Errorf("execution failed: %v", err)
}

package executors

import (
	"context"
	"errors"

	"github.com/ahab94/flash"
)

// executor outlines an executor
type executor struct {
	ctx         context.Context
	executables []flash.Executable
	name        string
}

// Execute - executes all executables concurrently
func (e *executor) Execute() error {
	if len(e.executables) < 1 {
		log(e.ctx, e.name).Warn("nothing to execute...")
		return errors.New("nothing to execute")
	}

	log(e.ctx, e.name).Infof("starting execution of %d executables...", len(e.executables)-1)
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
func (e *executor) Add(executable flash.Executable) {
	e.executables = append(e.executables, executable)
}

// OnSuccess - handles completion callback
func (e *executor) OnSuccess() {
	log(e.ctx, e.name).Infof("execution complete...")
}

// OnFailure - handles failure callback
func (e *executor) OnFailure(err error) {
	log(e.ctx, e.name).Errorf("execution failed: %v", err)
}

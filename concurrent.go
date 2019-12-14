package flash

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// Concurrent is an executor for concurrent executions
type Concurrent struct {
	executor
	dispatcher *Dispatcher
}

// NewConcurrent - initializes concurrent executor
func NewConcurrent(ctx context.Context, dispatcher *Dispatcher) *Concurrent {
	return &Concurrent{
		executor: executor{
			id:  fmt.Sprintf("%s-%s", "concurrent", uuid.NewV4().String()),
			ctx: ctx,
		},
		dispatcher: dispatcher,
	}
}

// Execute - executes all executables concurrently
func (c *Concurrent) Execute() error {
	if err := c.executor.Execute(); err != nil {
		return err
	}

	c.executeDispatch()
	return nil
}

func (c *Concurrent) executeDispatch() {
	doneChans := make([]chan struct{}, 0)
	for _, exec := range c.executables {
		if !exec.IsCompleted() {
			done := make(chan struct{})
			c.dispatcher.input <- Work{
				Executable: exec,
				done:       done,
			}
			doneChans = append(doneChans, done)
		}
	}

	for _, done := range doneChans {
		<-done
	}
}

package executors

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/ahab94/flash"
)

// Concurrent is an executor for concurrent executions
type Concurrent struct {
	executor
	dispatcher *flash.Dispatcher
}

// NewConcurrent - initializes concurrent executor
func NewConcurrent(ctx context.Context, dispatcher *flash.Dispatcher) *Concurrent {
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
	for _, exec := range c.executables {
		if !exec.IsCompleted() {
			c.dispatcher.Input() <- exec
		}
	}
}

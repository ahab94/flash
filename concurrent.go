package flash

import (
	"context"
	"fmt"

	"github.com/ahab94/engine"
	uuid "github.com/satori/go.uuid"
)

// Concurrent is an executor for concurrent executions
type Concurrent struct {
	executor
	engine *engine.Engine
	block  bool
}

// NewConcurrent - initializes concurrent executor; if completionBlock=true, it will block main routine until all tasks completed
func NewConcurrent(ctx context.Context, engine *engine.Engine, completionBlock bool) *Concurrent {
	return &Concurrent{
		executor: executor{
			id:  fmt.Sprintf("%s-%s", "concurrent", uuid.NewV4().String()),
			ctx: ctx,
		},
		engine: engine,
		block:  completionBlock,
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
	doneChans := make([]<-chan struct{}, 0)
	for _, exec := range c.executables {
		if !exec.IsCompleted() {
			done := c.engine.Do(exec)
			doneChans = append(doneChans, done)
		}
	}

	if c.block {
		for _, done := range doneChans {
			<-done
		}
	}
}

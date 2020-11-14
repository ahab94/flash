package flash

import (
	"context"
	"errors"
	"fmt"
	"github.com/ahab94/engine"
	uuid "github.com/satori/go.uuid"
)

// Concurrent is an executor for concurrent executions
type Concurrent struct {
	executor
	engine           *engine.Engine
	successHandler   func()
	failHandler      func(err error)
	aggregateInput   chan chan bool
	aggregateOutput  chan bool
	dispatchComplete bool
	block            bool
}

type ConcurrentOption func(*Concurrent)

// NewConcurrent - initializes concurrent executor; if completionBlock=true, it will block main routine until all tasks completed
func NewConcurrent(ctx context.Context, engine *engine.Engine, completionBlock bool, opts ...ConcurrentOption) *Concurrent {
	con := &Concurrent{
		executor: executor{
			id:  fmt.Sprintf("%s-%s", "concurrent", uuid.NewV4().String()),
			ctx: ctx,
		},
		aggregateInput:  make(chan chan bool, engine.WorkerCount()),
		aggregateOutput: make(chan bool),
		engine:          engine,
		block:           completionBlock,
	}

	for _, opt := range opts {
		opt(con)
	}

	return con
}

// ConcurrentFailHandler - inits fail handler
func ConcurrentFailHandler(fail func(err error)) ConcurrentOption {
	return func(s *Concurrent) {
		s.failHandler = fail
	}
}

// ConcurrentSuccessHandler - inits success handler
func ConcurrentSuccessHandler(success func()) ConcurrentOption {
	return func(c *Concurrent) {
		c.successHandler = success
	}
}

// Execute - executes all executables concurrently
func (c *Concurrent) Execute() error {
	if err := c.executor.Execute(); err != nil {
		return err
	}

	go c.aggregate()

	return c.executeDispatch()
}

func (c *Concurrent) executeDispatch() error {
	for _, exec := range c.executables {
		if !exec.IsCompleted() {
			c.aggregateInput <- c.engine.Do(exec)
		}
	}

	c.dispatchComplete = true

	if c.block {
		if !<-c.aggregateOutput {
			return errors.New("failed to execute all tasks")
		}
	}

	return nil
}

// OnSuccess - handles completion callback
func (c *Concurrent) OnSuccess() {
	c.executor.OnSuccess()
	if c.successHandler != nil {
		c.successHandler()
	}
}

// OnFailure - handles failure callback
func (c *Concurrent) OnFailure(err error) {
	c.executor.OnFailure(err)
	if c.failHandler != nil {
		c.failHandler(err)
	}
}

func (c *Concurrent) aggregate() {
	aggRes := true

	for {
		select {
		case agg := <-c.aggregateInput:
			log(c.id).Debugf("awaiting result %p", agg)
			if !aggRes {
				<-agg
				continue
			}

			aggRes = <-agg
		default:
			if c.dispatchComplete {
				log(c.id).Debugf("stopping aggregration...")
				c.aggregateOutput <- aggRes
				return
			}
		}
	}
}

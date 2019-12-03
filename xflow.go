package xflow

import (
	"context"
	"sync"
)

// Xflow represents a workflow
type Xflow struct {
	ctx         context.Context
	wg          *sync.WaitGroup
	executables []Executable
	concurrent  bool
}

// NewXFlow - initializes and returns Xflow struct
func NewXFlow(ctx context.Context, concurrent bool) *Xflow {
	return &Xflow{ctx: ctx, wg: &sync.WaitGroup{}, concurrent: concurrent}
}

// Execute - executes all executables either sequentially or concurrently
func (x *Xflow) Execute() error {
	if len(x.executables) < 1 {
		log(x.ctx).Warn("nothing to execute...")
		return nil
	}
	log(x.ctx).Infof("starting execution of %d executables...", len(x.executables))

	if x.concurrent {
		x.wg.Add(len(x.executables))
		for i := 0; i < len(x.executables); i++ {
			go func(i int) {
				defer x.wg.Done()
				if err := x.executables[i].Execute(); err != nil {
					log(x.ctx).Errorf("error encountered: %+v", err)
				}
			}(i)
		}
		x.wg.Wait()
		return nil
	}

	for _, exec := range x.executables {
		if err := exec.Execute(); err != nil {
			log(x.ctx).Errorf("error encountered: %+v", err)
			return err
		}
	}
	return nil
}

// IsCompleted - returns completion status
func (x *Xflow) IsCompleted() bool {
	if len(x.executables) < 1 {
		return false
	}

	for _, exec := range x.executables {
		if exec.IsCompleted() {
			continue
		}
		return false
	}
	return true
}

// Add - adds an executable to the executables list
func (x *Xflow) Add(executable Executable) {
	x.executables = append(x.executables, executable)
}

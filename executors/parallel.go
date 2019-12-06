package executors

import (
	"context"
	"sync"

	"github.com/ahab94/flash"
)

// Parallel is an executor for parallel executions
type Parallel struct {
	executor
	wg *sync.WaitGroup
}

// NewParallel - initializes a parallel executor
func NewParallel(ctx context.Context) *Parallel {
	return &Parallel{
		executor: executor{
			name: "parallel",
			ctx:  ctx,
		},
		wg: &sync.WaitGroup{},
	}
}

// Execute - executes all executables in parallel
func (p *Parallel) Execute() error {
	if err := p.executor.Execute(); err != nil {
		return err
	}

	p.executeWg()
	return nil
}

func (p *Parallel) executeWg() {
	p.wg.Add(len(p.executables))
	for i := 0; i < len(p.executables); i++ {
		go func(i int) {
			defer flash.RecoverPanic(p.ctx, p.executables[i])
			defer p.wg.Done()
			if !p.executables[i].IsCompleted() {
				if err := p.executables[i].Execute(); err != nil {
					log(p.ctx, p.name).Errorf("error encountered: %+v", err)
					p.executables[i].OnFailure(err)
				}
			}
			p.executables[i].OnSuccess()
		}(i)
	}
	p.wg.Wait()
}

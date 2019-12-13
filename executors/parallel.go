package executors

import (
	"context"
	"fmt"
	"sync"

	uuid "github.com/satori/go.uuid"

	"github.com/ahab94/flash/utils"
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
			id:  fmt.Sprintf("%s-%s", "parallel", uuid.NewV4().String()),
			ctx: ctx,
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
			defer utils.RecoverPanic(p.ctx)
			defer p.wg.Done()
			if !p.executables[i].IsCompleted() {
				if err := p.executables[i].Execute(); err != nil {
					log(p.id).Errorf("error while executing: %+v", err)
					p.executables[i].OnFailure(err)
					return
				}
				log(p.id).Infof("completed executing: %+v", p.executables[i])
				p.executables[i].OnSuccess()
			}
		}(i)
	}
	p.wg.Wait()
}

package flash

import (
	"context"
	"fmt"
	"sync"

	uuid "github.com/satori/go.uuid"
)

// Parallel is an executor for parallel executions
type Parallel struct {
	executor
	wg             *sync.WaitGroup
	successHandler func()
	failHandler    func(err error)
}

type ParallelOption func(parallel *Parallel)

// NewParallel - initializes a parallel executor
func NewParallel(ctx context.Context, opts ...ParallelOption) *Parallel {
	par := &Parallel{
		executor: executor{
			id:  fmt.Sprintf("%s-%s", "parallel", uuid.NewV4().String()),
			ctx: ctx,
		},
		wg: &sync.WaitGroup{},
	}

	for _, opt := range opts {
		opt(par)
	}

	return par
}

// ParallelFailHandler - inits fail handler
func ParallelFailHandler(fail func(err error)) ParallelOption {
	return func(p *Parallel) {
		p.failHandler = fail
	}
}

// ParallelSuccessHandler - inits success handler
func ParallelSuccessHandler(success func()) ParallelOption {
	return func(p *Parallel) {
		p.successHandler = success
	}
}

// Execute - executes all executables In parallel
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

// OnSuccess - handles completion callback
func (p *Parallel) OnSuccess() {
	p.executor.OnSuccess()
	if p.successHandler != nil {
		p.successHandler()
	}
}

// OnFailure - handles failure callback
func (p *Parallel) OnFailure(err error) {
	p.executor.OnFailure(err)
	if p.failHandler != nil {
		p.failHandler(err)
	}
}

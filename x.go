package xcruncher

import (
	"context"
	"sync"
)

// X represents an executable
type X struct {
	ctx         context.Context
	wg          *sync.WaitGroup
	executables []Executable
	dispatcher  *Dispatcher
	concurrent  bool
}

// NewExecutable - initializes and returns X struct
func NewExecutable(ctx context.Context, concurrent bool) *X {
	return &X{ctx: ctx, wg: &sync.WaitGroup{}, concurrent: concurrent}
}

// WithDispatcher - adds optional dispatcher for concurrent execution
func (x *X) WithDispatcher(dispatcher *Dispatcher) *X {
	if dispatcher != nil {
		x.dispatcher = dispatcher
	}
	return x
}

// Execute - executes all executables either sequentially or concurrently
func (x *X) Execute() error {
	if len(x.executables) < 1 {
		log(x.ctx).Warn("nothing to execute...")
		return nil
	}
	log(x.ctx).Infof("starting execution of %d executables...", len(x.executables)-1)

	if x.concurrent {
		if x.dispatcher != nil {
			x.executeDispatch()
			return nil
		}
		x.executeWg()
		return nil
	}
	for _, exec := range x.executables {
		if !exec.IsCompleted() {
			if err := exec.Execute(); err != nil {
				log(x.ctx).Errorf("error encountered: %+v", err)
				x.OnFailure(err)
				return err
			}
		}
	}
	return nil
}

// IsCompleted - returns completion status
func (x *X) IsCompleted() bool {
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

func (x *X) executeDispatch() {
	for _, exec := range x.executables {
		if !exec.IsCompleted() {
			x.dispatcher.Input() <- exec
		}
	}
}

func (x *X) executeWg() {
	x.wg.Add(len(x.executables))
	for i := 0; i < len(x.executables); i++ {
		go func(i int) {
			defer recoverPanic(x.ctx, x.executables[i])
			defer x.wg.Done()
			if !x.executables[i].IsCompleted() {
				if err := x.executables[i].Execute(); err != nil {
					log(x.ctx).Errorf("error encountered: %+v", err)
					x.executables[i].OnFailure(err)
				}
			}
			x.executables[i].OnSuccess()
		}(i)
	}
	x.wg.Wait()
}

// Add - adds an executable to the executables list
func (x *X) Add(executable Executable) {
	x.executables = append(x.executables, executable)
}

// OnSuccess - handles completion callback
func (x *X) OnSuccess() {
	log(x.ctx).Infof("execution complete...")
}

// OnFailure - handles failure callback
func (x *X) OnFailure(err error) {
	log(x.ctx).Errorf("execution failed: %v", err)
}

package types

import (
	"context"
	"sync"

	"github.com/ahab94/flash"
)

// Flow represents an executable
type Flow struct {
	ctx         context.Context
	wg          *sync.WaitGroup
	executables []flash.Executable
	dispatcher  *flash.Dispatcher
	concurrent  bool
}

// NewFlow - initializes and returns Flow struct
func NewFlow(ctx context.Context, concurrent bool) *Flow {
	return &Flow{ctx: ctx, wg: &sync.WaitGroup{}, concurrent: concurrent}
}

// WithDispatcher - adds optional dispatcher for concurrent execution
func (f *Flow) WithDispatcher(dispatcher *flash.Dispatcher) *Flow {
	if dispatcher != nil {
		f.dispatcher = dispatcher
	}
	return f
}

// Execute - executes all executables either sequentially or concurrently
func (f *Flow) Execute() error {
	if len(f.executables) < 1 {
		log(f.ctx).Warn("nothing to execute...")
		return nil
	}
	log(f.ctx).Infof("starting execution of %d executables...", len(f.executables)-1)

	if f.concurrent {
		if f.dispatcher != nil {
			f.executeDispatch()
			return nil
		}
		f.executeWg()
		return nil
	}
	for _, exec := range f.executables {
		if !exec.IsCompleted() {
			if err := exec.Execute(); err != nil {
				log(f.ctx).Errorf("error encountered: %+v", err)
				f.OnFailure(err)
				return err
			}
		}
	}
	return nil
}

// IsCompleted - returns completion status
func (f *Flow) IsCompleted() bool {
	if len(f.executables) < 1 {
		return false
	}

	for _, exec := range f.executables {
		if exec.IsCompleted() {
			continue
		}
		return false
	}
	return true
}

func (f *Flow) executeDispatch() {
	for _, exec := range f.executables {
		if !exec.IsCompleted() {
			f.dispatcher.Input() <- exec
		}
	}
}

func (f *Flow) executeWg() {
	f.wg.Add(len(f.executables))
	for i := 0; i < len(f.executables); i++ {
		go func(i int) {
			defer flash.RecoverPanic(f.ctx, f.executables[i])
			defer f.wg.Done()
			if !f.executables[i].IsCompleted() {
				if err := f.executables[i].Execute(); err != nil {
					log(f.ctx).Errorf("error encountered: %+v", err)
					f.executables[i].OnFailure(err)
				}
			}
			f.executables[i].OnSuccess()
		}(i)
	}
	f.wg.Wait()
}

// Add - adds an executable to the executables list
func (f *Flow) Add(executable flash.Executable) {
	f.executables = append(f.executables, executable)
}

// OnSuccess - handles completion callback
func (f *Flow) OnSuccess() {
	log(f.ctx).Infof("execution complete...")
}

// OnFailure - handles failure callback
func (f *Flow) OnFailure(err error) {
	log(f.ctx).Errorf("execution failed: %v", err)
}

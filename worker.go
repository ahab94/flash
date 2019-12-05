package flash

import (
	"context"
)

// Worker - a unit task executor
type Worker struct {
	ID      int
	ctx     context.Context
	counter *counter
	stop    chan struct{}
	input   chan Executable
	pool    chan chan Executable
}

// NewWorker - initializes a new Worker
func NewWorker(ctx context.Context, id int, pool chan chan Executable, counter *counter) *Worker {
	return &Worker{
		ID:      id,
		ctx:     ctx,
		pool:    pool,
		counter: counter,
		input:   make(chan Executable),
		stop:    make(chan struct{}),
	}
}

// Start - readies Worker for execution
func (w *Worker) Start() {
	log(w.ctx).Debugf("Worker [%d] is starting...", w.ID)
	go w.work()
}

// Stop - stops the Worker routine
func (w *Worker) Stop() {
	close(w.stop)
}

func (w *Worker) execute(exec Executable) {
	if !exec.IsCompleted() {
		defer RecoverPanic(w.ctx, exec)
		defer w.counter.Done()
		w.counter.Add()
		if err := exec.Execute(); err != nil {
			log(w.ctx).Errorf("Worker [%d]: error while executing: %v", w.ID, exec)
			exec.OnFailure(err)
			return
		}
		log(w.ctx).Infof("Worker [%d]: completed executed: %v", w.ID, exec)
		exec.OnSuccess()
	}
}

func (w *Worker) work() {
	for {
		select {
		case w.pool <- w.input:
			log(w.ctx).Debugf("Worker [%d] back in queue...", w.ID)
		case exec := <-w.input:
			log(w.ctx).Debugf("Worker [%d] executing %v...", w.ID, exec)
			w.execute(exec)
		case <-w.stop:
			log(w.ctx).Debugf("Worker [%d] stopping...", w.ID)
			return
		}
	}
}

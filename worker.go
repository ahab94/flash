package flash

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// Worker - a unit task executor
type Worker struct {
	id      string
	ctx     context.Context
	counter *counter
	stop    chan struct{}
	input   chan Executable
	pool    chan chan Executable
}

// NewWorker - initializes a new Worker
func NewWorker(ctx context.Context, pool chan chan Executable, counter *counter) *Worker {
	return &Worker{
		id:      fmt.Sprintf("%s-%s", "worker", uuid.NewV4().String()),
		ctx:     ctx,
		pool:    pool,
		counter: counter,
		input:   make(chan Executable),
		stop:    make(chan struct{}),
	}
}

// Start - readies Worker for execution
func (w *Worker) Start() {
	log(w.id).Debugf("starting...")
	go w.work()
}

// Stop - stops the Worker routine
func (w *Worker) Stop() {
	defer RecoverPanic(w.ctx)
	close(w.stop)
}

func (w *Worker) execute(exec Executable) {
	defer RecoverPanic(w.ctx)
	defer w.counter.Done()
	if !exec.IsCompleted() {
		if err := exec.Execute(); err != nil {
			log(w.id).Errorf("error while executing: %v", exec)
			exec.OnFailure(err)
			return
		}
		log(w.id).Infof("completed executing: %+v", exec)
		exec.OnSuccess()
	}
}

func (w *Worker) work() {
	for {
		select {
		case w.pool <- w.input:
			log(w.id).Debugf("back in queue")
		case exec := <-w.input:
			log(w.id).Debugf("executing: %+v", exec)
			w.execute(exec)
		case <-w.stop:
			log(w.id).Debugf("stopping...")
			return
		}
	}
}

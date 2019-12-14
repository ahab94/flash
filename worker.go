package flash

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// Worker - a unit task executor
type Worker struct {
	id    string
	ctx   context.Context
	stop  chan struct{}
	input chan Work
	pool  chan chan Work
}

// NewWorker - initializes a new Worker
func NewWorker(ctx context.Context, pool chan chan Work) *Worker {
	return &Worker{
		id:    fmt.Sprintf("%s-%s", "worker", uuid.NewV4().String()),
		ctx:   ctx,
		pool:  pool,
		input: make(chan Work),
		stop:  make(chan struct{}),
	}
}

// Start - readies Worker for execution
func (w *Worker) Start() {
	log(w.id).Debugf("starting...")
	go w.work()
}

// Stop - stops the Worker routine
func (w *Worker) Stop() {
	close(w.stop)
}

func (w *Worker) execute(execute Work) {
	defer close(execute.done)
	if !execute.IsCompleted() {
		if err := execute.Execute(); err != nil {
			log(w.id).Errorf("error while executing: %v", execute)
			execute.OnFailure(err)
			return
		}
		log(w.id).Infof("completed executing: %+v", execute)
		execute.OnSuccess()
	}
}

func (w *Worker) work() {
	for {
		select {
		case w.pool <- w.input:
			log(w.id).Debugf("back In queue")
		case execute := <-w.input:
			log(w.id).Debugf("executing: %+v", execute)
			w.execute(execute)
		case <-w.stop:
			log(w.id).Debugf("stopping...")
			return
		}
	}
}

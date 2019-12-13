package flash

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"sync"
)

// Dispatcher - for creating workers and distributing jobs
type Dispatcher struct {
	id      string
	ctx     context.Context
	stop    chan struct{}
	pool    chan chan Executable
	input   chan Executable
	workers []*Worker
	counter *counter
	start   *sync.Once
}

// NewDispatcher - initializing a new dispatcher
func NewDispatcher(ctx context.Context) *Dispatcher {
	return &Dispatcher{
		id:    fmt.Sprintf("%s-%s", "dispatcher", uuid.NewV4().String()),
		ctx:   ctx,
		start: new(sync.Once),
	}
}

// Start - starting workers and setting up dispatcher for use
func (d *Dispatcher) Start(workerCount uint) {
	d.start.Do(func() {
		d.stop = make(chan struct{})
		d.pool = make(chan chan Executable)
		d.input = make(chan Executable)
		d.workers = make([]*Worker, 0)
		d.counter = new(counter)

		for i := 0; i <= int(workerCount); i++ {
			worker := NewWorker(d.ctx, d.pool, d.counter)
			d.workers = append(d.workers, worker)
			worker.Start()
		}

		go d.dispatch()
	})
}

// Input - returns input channel for receiving tasks
func (d *Dispatcher) Input() chan Executable {
	return d.input
}

// Stop - closes channels/goroutines
func (d *Dispatcher) Stop() {
	defer RecoverPanic(d.ctx)
	defer func() { d.start = new(sync.Once) }()
	for _, worker := range d.workers {
		worker.Stop()
	}
	close(d.stop)
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case work := <-d.input:
			log(d.id).Debugf("dispatching: %v", work)
			d.counter.Add()
			worker := <-d.pool
			worker <- work

		case <-d.stop:
			log(d.id).Debugf("stopping...")
			return
		}
	}
}

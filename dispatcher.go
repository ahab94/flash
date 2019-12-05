package xcruncher

import (
	"context"
	"sync"
	"time"
)

// Dispatcher - for creating workers and distributing jobs
type Dispatcher struct {
	ctx     context.Context
	stop    chan struct{}
	pool    chan chan Executable
	input   chan Executable
	workers []*worker
	counter *counter
	start   *sync.Once
}

// NewDispatcher - initializing a new dispatcher
func NewDispatcher(ctx context.Context) *Dispatcher {
	return &Dispatcher{
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
		d.workers = make([]*worker, 0)
		d.counter = new(counter)

		for i := 0; i <= int(workerCount); i++ {
			worker := NewWorker(d.ctx, i, d.pool, d.counter)
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

// IsWorking - returns true if dispatcher is working
func (d *Dispatcher) IsWorking() bool {
	return d.counter.Count() != 0
}

// Wait - waits until dispatcher completes all outstanding tasks
func (d *Dispatcher) Wait() {
	for d.IsWorking() {
		time.Sleep(100 * time.Millisecond)
		continue
	}
}

// Stop - closes channels/goroutines
func (d *Dispatcher) Stop() {
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
			log(d.ctx).Debugf("dispatching work: %v", work)
			worker := <-d.pool
			worker <- work

		case <-d.stop:
			log(d.ctx).Debugf("dispatcher stopping...")
			return
		}
	}
}
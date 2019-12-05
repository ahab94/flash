package xcruncher

import (
	"context"
	"errors"
)

// Dispatcher - for creating workers and distributing jobs
type Dispatcher struct {
	ctx     context.Context
	workers []*worker
	pool    chan chan Executable
	input   chan Executable
	counter *counter
}

// NewDispatcher - initializing a new dispatcher
func NewDispatcher(ctx context.Context, workerCount int) (*Dispatcher, error) {
	if workerCount < 1 {
		return nil, errors.New("worker count must be > 0")
	}

	dispatcher := &Dispatcher{
		ctx:     ctx,
		pool:    make(chan chan Executable),
		input:   make(chan Executable),
		counter: &counter{},
	}

	for i := 0; i <= workerCount; i++ {
		worker := newWorker(ctx, i, dispatcher.pool, dispatcher.counter)
		dispatcher.workers = append(dispatcher.workers, worker)
	}

	return dispatcher, nil
}

// Start - starting workers and setting up dispatcher for use
func (d *Dispatcher) Start() {
	for _, worker := range d.workers {
		worker.Start()
	}
	go d.dispatch()
}

// Input - returns input channel for receiving tasks
func (d *Dispatcher) Input() chan Executable {
	return d.input
}

// IsWorking - returns true if dispatcher is working
func (d *Dispatcher) IsWorking() bool {
	return d.counter.Count() != 0
}

// Stop - closes channels/goroutines
func (d *Dispatcher) Stop() {
	for _, worker := range d.workers {
		worker.Stop()
	}
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case work := <-d.input:
			worker := <-d.pool
			worker <- work
		}
	}
}

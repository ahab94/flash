package xcruncher

import "sync"

type counter struct {
	count      int
	totalCount int
	lock       sync.Mutex
}

func (t *counter) Add() {
	defer t.lock.Unlock()
	t.lock.Lock()
	t.count++
	t.totalCount++
}

func (t *counter) Done() {
	defer t.lock.Unlock()
	t.lock.Lock()
	t.count--
}

func (t *counter) Count() int {
	defer t.lock.Unlock()
	t.lock.Lock()
	return t.count
}

package task

import (
	"sync"
	"sync/atomic"
)

type Latch struct {
	atomic.Int32
	done chan struct{}
	once sync.Once
}

func (w *Latch) Add(delta int) {
	w.Int32.Add(int32(delta))
}

func (w *Latch) Done() {
	val := w.Int32.Add(-1)
	// state machine: can operate without locks.
	if val <= 0 {
		w.once.Do(func() {
			close(w.done)
		})
	}
}

func (w *Latch) Wait() <-chan struct{} {
	return w.done
}

func NewLatch() *Latch {
	return &Latch{
		done: make(chan struct{}),
	}
}

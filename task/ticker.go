package task

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Ticker struct {
	// using atomic for concurrent read and write
	startOnce atomic.Pointer[sync.Once]
	// using atomic for concurrent read and write
	stopOnce atomic.Pointer[sync.Once]

	cancelFunc context.CancelFunc
	wg         sync.WaitGroup

	taskFunc func()
	interval time.Duration
}

// Start starts ticker task
func (s *Ticker) Start() {
	startOnce := s.startOnce.Load()

	if startOnce == nil {
		return
	}
	startOnce.Do(func() {
		ctx, cc := context.WithCancel(context.Background())
		s.cancelFunc = cc
		s.wg.Add(1)
		go s.run(ctx)
		// assign Stop "Once"
		var so sync.Once
		s.stopOnce.Store(&so) // store after running
	})
}

func (s *Ticker) run(ctx context.Context) {
	defer s.wg.Done()
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if s.taskFunc != nil {
				s.taskFunc()
			}
		case <-ctx.Done():
			return
		}
	}
}

// Stop stops the task
func (s *Ticker) Stop() {
	stopOnce := s.stopOnce.Load()

	if stopOnce == nil {
		return // not running
	}

	// 并发只调用一次
	stopOnce.Do(func() {
		s.cancelFunc() // cancel task
		s.wg.Wait()    // wait task completed
		var so sync.Once
		s.startOnce.Store(&so)
	})
}
func NewTicker(task func(), duration time.Duration) *Ticker {
	s := &Ticker{
		taskFunc: task,
		interval: duration,
	}
	var so sync.Once
	s.startOnce.Store(&so)
	return s
}

package io

import (
	"errors"
	"net"
	"time"
)

var (
	ErrBlockQueueWaitTimeout = errors.New("block queue waiting timeout")
)

type BlockingQueue[T any] struct {
	enqueueCh chan T
	dequeueCh chan chan T
	closeCh   chan struct{}
	closed    bool
}

func NewBlockingQueue[T any]() *BlockingQueue[T] {
	q := &BlockingQueue[T]{
		enqueueCh: make(chan T),
		dequeueCh: make(chan chan T),
		closeCh:   make(chan struct{}),
	}
	go q.run()
	return q
}

func (q *BlockingQueue[T]) run() {
	var buf []T
	for {
		select {
		case v := <-q.enqueueCh:
			buf = append(buf, v)
		case resp := <-q.dequeueCh:
			if len(buf) == 0 {
				// close or blocking if empty
				select {
				case v := <-q.enqueueCh:
					buf = append(buf, v)
				case <-q.closeCh:
					close(resp)
					return
				}
			}
			if len(buf) > 0 {
				resp <- buf[0]
				buf = buf[1:]
			}
		case <-q.closeCh:
			return
		}
	}
}

func (q *BlockingQueue[T]) Enqueue(v T) error {
	select {
	case q.enqueueCh <- v:
		return nil
	case <-q.closeCh:
		return net.ErrClosed
	}
}

func (q *BlockingQueue[T]) Dequeue() (T, error) {
	resp := make(chan T)
	select {
	case q.dequeueCh <- resp:
		v, ok := <-resp
		if !ok {
			var zero T
			return zero, net.ErrClosed
		}
		return v, nil
	case <-q.closeCh:
		var zero T
		return zero, net.ErrClosed
	}
}

func (q *BlockingQueue[T]) DequeueWithTimeout(timeout time.Duration) (T, error) {
	resp := make(chan T)
	select {
	case q.dequeueCh <- resp:
		select {
		case v, ok := <-resp:
			if !ok {
				var zero T
				return zero, net.ErrClosed
			}
			return v, nil
		case <-time.After(timeout):
			var zero T
			return zero, ErrBlockQueueWaitTimeout
		}
	case <-q.closeCh:
		var zero T
		return zero, net.ErrClosed
	case <-time.After(timeout):
		var zero T
		return zero, ErrBlockQueueWaitTimeout
	}
}

func (q *BlockingQueue[T]) Close() {
	if !q.closed {
		q.closed = true
		close(q.closeCh)
	}
}

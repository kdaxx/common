package task

import (
	"context"
	"fmt"
	"sync"
)

type task struct {
	Name string
	Run  func(ctx context.Context) error
}

type Group struct {
	tasks          []task
	cleanup        func(err error)
	waitForCleanup bool
	fastErrReturn  bool
	parallelQueue  chan struct{}
}

func (g *Group) AppendWithName(name string, f func(ctx context.Context) error) *Group {
	g.tasks = append(g.tasks, task{
		Name: name,
		Run:  f,
	})
	return g
}

func (g *Group) Append(f func(ctx context.Context) error) *Group {
	g.tasks = append(g.tasks, task{
		Run: f,
	})
	return g
}

func (g *Group) Cleanup(f func(err error)) *Group {
	g.cleanup = f
	return g
}

func (g *Group) WaitForCleanup() *Group {
	g.waitForCleanup = true
	return g
}

func (g *Group) FastErrReturn() *Group {
	g.fastErrReturn = true
	return g
}

func (g *Group) Parallel(n int) *Group {
	g.parallelQueue = make(chan struct{}, n)
	for i := 0; i < n; i++ {
		g.parallelQueue <- struct{}{}
	}
	return g
}

func (g *Group) Run(ctx context.Context) error {
	if len(g.tasks) == 0 {
		return nil
	}
	taskWait := NewLatch()
	taskWait.Add(len(g.tasks))

	// cancel by upstream or fastErrReturn
	taskCancelContext, taskCancel := context.WithCancel(ctx)
	defer taskCancel()

	var errorAccess sync.Mutex
	var returnError error

	for _, t := range g.tasks {
		go func() {
			defer func() {
				taskWait.Done()
			}()

			if g.parallelQueue != nil {
				select {
				// task canceled
				case <-taskCancelContext.Done():
					return
				case <-g.parallelQueue:
				}
			}
			// run task
			err := t.Run(taskCancelContext)
			if err != nil {
				if t.Name != "" {
					err = fmt.Errorf("%s: %w", t.Name, err)
				}
				errorAccess.Lock()
				if returnError == nil {
					returnError = err
				}
				errorAccess.Unlock()
				if g.fastErrReturn {
					// cancel task
					taskCancel()
				}
			}
			if g.parallelQueue != nil {
				g.parallelQueue <- struct{}{}
			}
		}()
	}
	//  task completed
	select {
	case <-taskWait.Wait():
		// all task done
	case <-taskCancelContext.Done():
		// task group cancelled
		if ctx.Err() != nil {
			returnError = ctx.Err()
		}
	}

	// cleanup
	if g.cleanup != nil {
		g.cleanup(returnError)
	}

	if g.waitForCleanup {
		<-taskWait.Wait()
	}
	return returnError
}

func NewGroup() *Group {
	return &Group{}
}

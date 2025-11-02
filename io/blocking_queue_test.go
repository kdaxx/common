package io

import "testing"

func TestBlockingQueue(t *testing.T) {
	queue := NewBlockingQueue[int]()
	err := queue.Enqueue(1)
	if err != nil {
		t.Fatal(err)
	}
	err = queue.Enqueue(2)
	if err != nil {
		t.Fatal(err)
	}
	dequeue, err := queue.Dequeue()
	if err != nil {
		t.Fatal(err)
	}
	if dequeue != 1 {
		t.Fatalf("dequeue error, expected 1, got %d", dequeue)
	}
	dequeue, err = queue.Dequeue()
	if err != nil {
		t.Fatal(err)
	}
	if dequeue != 2 {
		t.Fatalf("dequeue error, expected 2, got %d", dequeue)
	}
}

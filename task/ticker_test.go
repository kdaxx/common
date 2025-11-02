package task

import (
	"fmt"
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	times := 0
	ticker := NewTicker(func() {
		times++
		fmt.Printf("tick %d times\n", times)
	}, 2*time.Second)
	ticker.Start()
	time.Sleep(5 * time.Second)
	ticker.Stop()
	if times != 2 {
		t.Fatal("times should be 2")
	}
}

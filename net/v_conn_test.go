package net

import (
	"fmt"
	"sync"
	"testing"
)

func TestConnPipe(t *testing.T) {
	conn1, conn2 := NewConnPipe()
	defer conn1.Close()
	defer conn2.Close()
	group := sync.WaitGroup{}
	group.Add(2)
	go func() {
		defer group.Done()
		_, err := conn1.Write([]byte("hello world"))
		if err != nil {
			t.Fatal(err)
		}

		t.Log("conn1 write completed:")
		buf := make([]byte, 1024)
		n, err := conn1.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("conn1 reading:" + string(buf[:n]))
	}()

	go func() {
		defer group.Done()

		buf := make([]byte, 1024)
		n, err := conn2.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("conn1 reading:" + string(buf[:n]))

		_, err = conn2.Write([]byte("hello world"))
		if err != nil {
			t.Fatal(err)
		}

		t.Log("conn1 write completed:")
	}()
	group.Wait()
}

func TestCloseVConn1Conn2Write(t *testing.T) {
	conn1, conn2 := NewConnPipe()
	err := conn1.Close()
	if err != nil {
		t.Fatal("close VConn failed", err)
	}

	_, err = conn2.Write([]byte("hello world"))
	if err == nil {
		t.Fatal("write to closed conn?", err)
	}
	fmt.Println("failed to write conn", err)
}

func TestCloseVConn2Conn1Write(t *testing.T) {
	conn1, conn2 := NewConnPipe()
	err := conn2.Close()
	if err != nil {
		t.Fatal("close VConn failed", err)
	}

	_, err = conn1.Write([]byte("hello world"))
	if err == nil {
		t.Fatal("write to closed conn?", err)
	}
	fmt.Println("failed to write conn", err)
}

func TestCloseVConn1(t *testing.T) {
	conn1, _ := NewConnPipe()
	err := conn1.Close()
	if err != nil {
		t.Fatal("close VConn failed", err)
	}

	_, err = conn1.Write([]byte("hello world"))
	if err == nil {
		t.Fatal("write to closed conn?", err)
	}
	fmt.Println("failed to write conn", err)
}

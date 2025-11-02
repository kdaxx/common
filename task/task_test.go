package task

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	now := time.Now()
	err := NewGroup().
		AppendWithName("task1", func(ctx context.Context) error {
			time.Sleep(1 * time.Second)
			fmt.Println("task1")
			return nil
		}).
		AppendWithName("task2", func(ctx context.Context) error {
			time.Sleep(1 * time.Second)
			fmt.Println("task2")
			return nil
		}).Run(context.Background())
	if err != nil {
		t.Error("TestTask err:", err)
	}
	if time.Now().Sub(now).Seconds() > 2 {
		t.Error("TestTask time out")
	}
}
func TestSyncTask(t *testing.T) {
	now := time.Now()
	err := NewGroup().
		AppendWithName("task1", func(ctx context.Context) error {
			time.Sleep(1 * time.Second)
			fmt.Println("task1")
			return nil
		}).AppendWithName("task2", func(ctx context.Context) error {
		time.Sleep(1 * time.Second)
		fmt.Println("task2")
		return nil
	}).
		Parallel(1).
		Run(context.Background())
	if err != nil {
		t.Fatal("TestSyncTask err:", err)
	}
	if time.Now().Sub(now).Seconds() < 2 {
		t.Error("TestTask time out")
	}
}

func TestParallelTask(t *testing.T) {
	err := NewGroup().
		AppendWithName("task1", func(ctx context.Context) error {
			time.Sleep(1 * time.Second)
			fmt.Println("task1")
			return nil
		}).AppendWithName("task2", func(ctx context.Context) error {
		time.Sleep(1 * time.Second)
		fmt.Println("task2")
		return nil
	}).Append(func(ctx context.Context) error {
		time.Sleep(1 * time.Second)
		fmt.Println("task3")
		return nil
	}).Append(func(ctx context.Context) error {
		time.Sleep(1 * time.Second)
		fmt.Println("task4")
		return nil
	}).
		Parallel(2).
		Run(context.Background())
	if err != nil {
		t.Fatal("TestSyncTask err:", err)
	}
}

func TestFastFailTask(t *testing.T) {
	err := NewGroup().
		AppendWithName("task1", func(ctx context.Context) error {
			time.Sleep(1 * time.Second)
			fmt.Println("task1")
			return errors.New("task1 failed")
		}).AppendWithName("task2", func(ctx context.Context) error {
		time.Sleep(1 * time.Second)
		fmt.Println("task2")
		return errors.New("task2 failed")
	}).Append(func(ctx context.Context) error {
		time.Sleep(1 * time.Second)
		fmt.Println("task3")
		return errors.New("task3 failed")
	}).Append(func(ctx context.Context) error {
		time.Sleep(1 * time.Second)
		fmt.Println("task4")
		return errors.New("task4 failed")
	}).
		Parallel(2).
		FastErrReturn().
		Run(context.Background())
	if err == nil {
		t.Fatal("TestSyncTask failed to fast fallback:", err)
	}
	t.Log("fast fallback success! TestSyncTask err:", err)
}

func TestCleanupTask(t *testing.T) {
	cleanup := false
	err := NewGroup().
		AppendWithName("task1", func(ctx context.Context) error {
			time.Sleep(1 * time.Second)
			fmt.Println("task1")
			return errors.New("task1 failed")
		}).AppendWithName("task2", func(ctx context.Context) error {
		time.Sleep(1 * time.Second)
		fmt.Println("task2")
		return errors.New("task2 failed")
	}).Cleanup(func(err error) {
		cleanup = true
	}).Run(context.Background())
	if err != nil && !cleanup {
		t.Fatal("TestCleanupTask err:", err)
	}
}

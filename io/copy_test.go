package io

import (
	"bytes"
	"context"
	"github.com/kdaxx/common/net"
	"github.com/kdaxx/common/task"
	"testing"
)

func TestCopy(t *testing.T) {
	lreader, lwriter := net.NewPacketConnPipe()
	rreader, rwriter := net.NewPacketConnPipe()
	defer lreader.Close()
	defer lwriter.Close()
	defer rreader.Close()
	defer rwriter.Close()

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	go func() {
		defer cancelFunc()
		err := task.NewGroup().FastErrReturn().Append(func(ctx context.Context) error {
			_, err := CopyPacket(rwriter, lreader)
			return err
		}).Run(ctx)
		if err != nil {
			t.Error(err)
			return
		}
	}()
	msg := "hello"
	proxyAddr := net.ProxyAddr{
		Host: net.ProxyFQDN("127.0.0.1"),
		Port: 1233,
	}
	_, err := lwriter.WriteTo([]byte(msg), proxyAddr)
	if err != nil {
		t.Error(err)
		return
	}
	data := make([]byte, len(msg))
	_, addr, err := rreader.ReadFrom(data)
	if err != nil {
		t.Error(err)
		return
	}
	if addr.String() != proxyAddr.String() {
		t.Errorf("expect %s, got %s", proxyAddr.String(), addr.String())
		return
	}

	if !bytes.Equal(data, []byte(msg)) {
		t.Errorf("expect %s, got %s", msg, string(data))
		return
	}

}

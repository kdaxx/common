package net

import (
	"common/errs"
	"io"
	"net"
	"sync"
)

var (
	_ net.PacketConn = (*virtualPacketConn)(nil)
	_ net.Conn       = (*virtualConn)(nil)
)

type virtualConn struct {
	net.Conn
	reader     io.ReadCloser
	peerWriter io.WriteCloser

	once sync.Once
}

func (v *virtualConn) Read(p []byte) (n int, err error) {
	return v.reader.Read(p)
}

func (v *virtualConn) Write(p []byte) (n int, err error) {
	return v.peerWriter.Write(p)
}

func (v *virtualConn) Close() error {
	return errs.Combine(
		v.reader.Close(),
		v.peerWriter.Close())
}

func NewConnPipe() (net.Conn, net.Conn) {
	lReader, lWriter := io.Pipe()
	rReader, rWriter := io.Pipe()
	return &virtualConn{
			reader:     lReader,
			peerWriter: rWriter,
		}, &virtualConn{
			reader:     rReader,
			peerWriter: lWriter,
		}
}

type packet struct {
	data []byte
	addr net.Addr
}

type virtualPacketConn struct {
	net.PacketConn
	peer *virtualPacketConn
	ch   chan packet

	done chan struct{}
	once sync.Once
}

func (c *virtualPacketConn) ReadFrom(p []byte) (int, net.Addr, error) {
	select {
	case <-c.done:
		return 0, nil, net.ErrClosed
	case pkt := <-c.ch:
		n := copy(p, pkt.data)
		return n, pkt.addr, nil
	}
}

func (c *virtualPacketConn) WriteTo(p []byte, addr net.Addr) (int, error) {
	select {
	case <-c.done:
		return 0, net.ErrClosed
	default:
	}

	// write to peer
	c.peer.ch <- packet{
		data: p,
		addr: addr,
	}
	return len(p), nil
}

func (c *virtualPacketConn) Close() error {
	c.once.Do(func() {
		close(c.done) // close pipe
	})
	return nil
}

// NewPacketConnPipe returns a virtual packet conn
func NewPacketConnPipe() (net.PacketConn, net.PacketConn) {
	c1 := &virtualPacketConn{
		ch:   make(chan packet),
		done: make(chan struct{}),
	}
	c2 := &virtualPacketConn{
		ch:   make(chan packet),
		done: make(chan struct{}),
	}
	c1.peer = c2
	c2.peer = c1
	return c1, c2
}

func NewPacketConnPipeWithCache(cache int) (net.PacketConn, net.PacketConn) {
	c1 := &virtualPacketConn{
		ch:   make(chan packet, cache),
		done: make(chan struct{}),
	}
	c2 := &virtualPacketConn{
		ch:   make(chan packet, cache),
		done: make(chan struct{}),
	}
	c1.peer = c2
	c2.peer = c1
	return c1, c2
}

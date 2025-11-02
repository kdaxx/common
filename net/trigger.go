package net

import (
	"io"
)

var (
	_ io.Writer    = (*WriterLambda)(nil)
	_ io.Reader    = (*ReaderLambda)(nil)
	_ PacketWriter = (*PacketWriterLambda)(nil)
	_ PacketReader = (*PacketReaderLambda)(nil)
)

type WriterLambda struct {
	writerFun func(p []byte) (n int, err error)
}

func (h *WriterLambda) Write(p []byte) (n int, err error) {
	return h.writerFun(p)
}

func NewWriterLambda(writerFun func(p []byte) (n int, err error)) *WriterLambda {
	return &WriterLambda{
		writerFun: writerFun,
	}
}

type PacketWriterLambda struct {
	packetWriterFun func(b []byte, addr ProxyAddr) (int, error)
}

func (p *PacketWriterLambda) WriteTo(b []byte, addr ProxyAddr) (int, error) {
	return p.packetWriterFun(b, addr)
}

func NewPacketWriterLambda(packetWriterFun func(b []byte, addr ProxyAddr) (int, error)) *PacketWriterLambda {
	return &PacketWriterLambda{
		packetWriterFun: packetWriterFun,
	}
}

type ReaderLambda struct {
	readerFun func(p []byte) (n int, err error)
}

func (e *ReaderLambda) Read(p []byte) (n int, err error) {
	return e.readerFun(p)
}

func NewReaderLambda(handleEchoFun func(p []byte) (n int, err error)) *ReaderLambda {
	return &ReaderLambda{
		readerFun: handleEchoFun,
	}
}

type PacketReaderLambda struct {
	packetReaderFun func(b []byte) (int, ProxyAddr, error)
}

func (p *PacketReaderLambda) ReadFrom(b []byte) (int, ProxyAddr, error) {
	return p.packetReaderFun(b)
}

func NewPacketReaderLambda(handleEchoFun func(b []byte) (int, ProxyAddr, error)) *PacketReaderLambda {
	return &PacketReaderLambda{
		packetReaderFun: handleEchoFun,
	}
}

package net

import "net"

type PacketReader interface {
	ReadFrom(b []byte) (int, net.Addr, error)
}

type PacketWriter interface {
	WriteTo(b []byte, addr net.Addr) (int, error)
}

type PacketConn interface {
	PacketReader
	PacketWriter
}

package net

type PacketReader interface {
	ReadFrom(b []byte) (int, ProxyAddr, error)
}

type PacketWriter interface {
	WriteTo(b []byte, addr ProxyAddr) (int, error)
}

type PacketConn interface {
	PacketReader
	PacketWriter
}

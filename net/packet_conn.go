package net

import (
	"net"
)

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

type udpPacketWriter struct {
	udpConn *net.UDPConn
}

func (u *udpPacketWriter) WriteTo(b []byte, addr net.Addr) (int, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr.String())
	if err != nil {
		return 0, err
	}
	return u.udpConn.WriteToUDP(b, udpAddr)
}

func NewUDPPacketWriter(udpConn *net.UDPConn) PacketWriter {
	return &udpPacketWriter{udpConn: udpConn}
}

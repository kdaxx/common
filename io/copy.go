package io

import (
	"github.com/kdaxx/common/errs"
	NET "github.com/kdaxx/common/net"
	"net"
)

const BufferSize = 8 * 1024

func CopyPacketWithBufferSize(dest NET.PacketWriter, src NET.PacketReader, buf []byte) (total int64, err error) {
	if len(buf) <= 0 {
		return 0, errs.New("buffer size must be greater than zero")
	}
	buffer := buf
	if udpConn, isUDPConn := dest.(*net.UDPConn); isUDPConn {
		dest = NET.NewUDPPacketWriter(udpConn)
	}
	for {
		nr, addr, err := src.ReadFrom(buffer)
		if err != nil {
			return total, err
		}
		nw, err := dest.WriteTo(buffer[:nr], addr)
		if err != nil {
			return total, err
		}
		total += int64(nw)
	}
}

func CopyPacket(dest NET.PacketWriter, src NET.PacketReader) (total int64, err error) {
	return CopyPacketWithBufferSize(dest, src, make([]byte, BufferSize))
}

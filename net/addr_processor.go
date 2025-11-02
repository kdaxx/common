package net

import (
	"encoding/binary"
	"errors"
	"io"
	"net/netip"
)

var (
	ErrInvalidHost = errors.New("invalid host")
	ErrTooLongFQDN = errors.New("fqdn too long")
	ErrHostFlag    = errors.New("unknown host flag")
)

const invalidAddrFlagByte = byte(0x11)

type portProcessor struct {
}

func (p *portProcessor) writePortTo(w io.Writer, port uint16) error {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, port)
	_, err := w.Write(b)
	return err
}

func (p *portProcessor) readPortFrom(r io.Reader) (port uint16, err error) {
	b := make([]byte, 2)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return
	}
	port = binary.BigEndian.Uint16(b)
	return
}

type HostFlagByte struct {
	V4   byte
	V6   byte
	FQDN byte
}

var DefaultFlagByte = HostFlagByte{
	V4:   0x00,
	V6:   0x01,
	FQDN: 0x02,
}

type hostProcessor struct {
	hostFlag HostFlagByte
}

// WriteTo writes host to writer
func (a *hostProcessor) writeHostTo(w io.Writer, host ProxyHost) error {
	af := invalidAddrFlagByte
	if host.IsFQDN() {
		if len(host.FQDN()) > 255 {
			return ErrTooLongFQDN
		}
		af = a.hostFlag.FQDN
	}
	if host.IP().Is4() {
		af = a.hostFlag.V4
	}
	if host.IP().Is6() {
		af = a.hostFlag.V6
	}
	switch af {
	case a.hostFlag.V4, a.hostFlag.V6:
		var b []byte
		if _, err := w.Write([]byte{af}); err != nil {
			return err
		}
		b = host.IP().AsSlice()
		if _, err := w.Write(b); err != nil {
			return err
		}
		return nil
	case a.hostFlag.FQDN:
		if _, err := w.Write([]byte{af, byte(len(host.FQDN()))}); err != nil {
			return err
		}
		if _, err := w.Write([]byte(host.FQDN())); err != nil {
			return err
		}
		return nil
	default:
		return ErrInvalidHost
	}
}

func (a *hostProcessor) readHostFrom(r io.Reader) (ProxyHost, error) {
	addrByte := make([]byte, 1)
	_, err := io.ReadFull(r, addrByte)
	if err != nil {
		return nil, err
	}
	switch addrByte[0] {
	case a.hostFlag.V4:
		var v4Bytes [4]byte
		if _, err = io.ReadFull(r, v4Bytes[:]); err != nil {
			return nil, err
		}
		addrFrom4 := netip.AddrFrom4(v4Bytes)
		if !addrFrom4.IsValid() {
			return nil, ErrInvalidHost
		}
		return &ProxyIP{
			NetIP: addrFrom4,
		}, nil
	case a.hostFlag.V6:
		var v6Bytes [16]byte
		if _, err = io.ReadFull(r, v6Bytes[:]); err != nil {
			return nil, err
		}
		addrFrom16 := netip.AddrFrom16(v6Bytes)
		if !addrFrom16.IsValid() {
			return nil, ErrInvalidHost
		}
		return &ProxyIP{
			NetIP: addrFrom16,
		}, nil
	case a.hostFlag.FQDN:
		hostLenByte := make([]byte, 1)
		if _, err = io.ReadFull(r, hostLenByte); err != nil {
			return nil, err
		}
		length := int(hostLenByte[0])
		hostBytes := make([]byte, length)
		if _, err := io.ReadFull(r, hostBytes); err != nil {
			return nil, err
		}
		var host ProxyHost
		if host, err = ParseProxyHost(string(hostBytes)); err != nil {
			return nil, err
		}
		return host, nil
	default:
		return nil, ErrHostFlag
	}
}

type AddressProcessor interface {
	ReadAddrFrom(reader io.Reader) (ProxyAddr, error)
	WriteAddrTo(writer io.Writer, addr ProxyAddr) error
}

type portFirstAddrProcessor struct {
	portProcessor portProcessor
	hostProcessor hostProcessor
}

func (a *portFirstAddrProcessor) WriteAddrTo(w io.Writer, addr ProxyAddr) error {
	if err := a.portProcessor.writePortTo(w, addr.Port); err != nil {
		return err
	}
	if err := a.hostProcessor.writeHostTo(w, addr.Host); err != nil {
		return err
	}
	return nil
}

func (a *portFirstAddrProcessor) ReadAddrFrom(r io.Reader) (ProxyAddr, error) {
	port, err := a.portProcessor.readPortFrom(r)
	if err != nil {
		return ProxyAddr{}, err
	}
	host, err := a.hostProcessor.readHostFrom(r)
	if err != nil {
		return ProxyAddr{}, err
	}
	return ProxyAddr{
		Host: host,
		Port: port,
	}, nil
}

func NewPortFirstAddrProcessorWithFlagByte(flagByte HostFlagByte) AddressProcessor {
	return &portFirstAddrProcessor{
		hostProcessor: hostProcessor{
			hostFlag: flagByte,
		},
	}
}

func NewPortFirstAddrProcessor() AddressProcessor {
	return &portFirstAddrProcessor{
		hostProcessor: hostProcessor{
			hostFlag: DefaultFlagByte,
		},
	}
}

type hostFirstAddrProcessor struct {
	portProcessor portProcessor
	hostProcessor hostProcessor
}

func (a *hostFirstAddrProcessor) WriteAddrTo(w io.Writer, addr ProxyAddr) error {

	if err := a.hostProcessor.writeHostTo(w, addr.Host); err != nil {
		return err
	}
	if err := a.portProcessor.writePortTo(w, addr.Port); err != nil {
		return err
	}
	return nil
}

func (a *hostFirstAddrProcessor) ReadAddrFrom(r io.Reader) (ProxyAddr, error) {

	host, err := a.hostProcessor.readHostFrom(r)
	if err != nil {
		return ProxyAddr{}, err
	}
	port, err := a.portProcessor.readPortFrom(r)
	if err != nil {
		return ProxyAddr{}, err
	}
	return ProxyAddr{
		Host: host,
		Port: port,
	}, nil
}

func NewHostFirstAddrProcessorWithFlagByte(flagByte HostFlagByte) AddressProcessor {
	return &hostFirstAddrProcessor{
		hostProcessor: hostProcessor{
			hostFlag: flagByte,
		},
	}
}

func NewHostFirstAddrProcessor() AddressProcessor {
	return &hostFirstAddrProcessor{
		hostProcessor: hostProcessor{
			hostFlag: DefaultFlagByte,
		},
	}
}

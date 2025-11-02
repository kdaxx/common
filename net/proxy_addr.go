package net

import (
	"errors"
	"net"
	"net/netip"
	"strconv"
)

type ProxyHost interface {
	IsIP() bool
	IP() netip.Addr
	IsFQDN() bool
	FQDN() string
	String() string
}

type ProxyIP struct {
	NetIP netip.Addr
}

func (a *ProxyIP) IP() netip.Addr {
	return a.NetIP
}

func (a *ProxyIP) FQDN() string {
	return ""
}

func (a *ProxyIP) IsIP() bool {
	return true
}

func (a *ProxyIP) IsFQDN() bool {
	return false
}

func (a *ProxyIP) IsValid() bool {
	return a.IsIP()
}

func (a *ProxyIP) String() string {
	return a.NetIP.String()
}

type ProxyFQDN string

func (p ProxyFQDN) IsIP() bool {
	return false
}

func (p ProxyFQDN) IP() netip.Addr {
	return netip.Addr{}
}

func (p ProxyFQDN) IsFQDN() bool {
	return true
}

func (p ProxyFQDN) FQDN() string {
	return string(p)
}

func (p ProxyFQDN) String() string {
	return p.FQDN()
}

type ProxyAddr struct {
	Host ProxyHost
	Port uint16
}

func (a ProxyAddr) PortString() string {
	return strconv.Itoa(int(a.Port))
}

func (a ProxyAddr) Network() string {
	return "proxy"
}

func (a ProxyAddr) String() string {
	// ensure ipv6 being surrounded by square brackets
	return net.JoinHostPort(a.Host.String(), a.PortString())
}

func ParseProxyAddr(addr string) (*ProxyAddr, error) {
	host, p, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	port, err := strconv.ParseUint(p, 10, 16)
	if err != nil {
		return nil, err
	}
	addrPort, err := ParseProxyHostPort(host, uint16(port))
	if err != nil {
		return nil, err
	}
	return addrPort, nil
}

func ParseProxyHost(host string) (ProxyHost, error) {
	if IsDomainName(host) {
		fqdn := ProxyFQDN(host)
		return &fqdn, nil
	}
	ip, err := netip.ParseAddr(host)
	if err != nil {
		return nil, errors.New("invalid host")
	}

	return &ProxyIP{
		NetIP: ip,
	}, nil
}

func ParseProxyHostPort(host string, port uint16) (*ProxyAddr, error) {
	proxyHost, err := ParseProxyHost(host)
	if err != nil {
		return nil, err
	}
	return &ProxyAddr{
		Host: proxyHost,
		Port: port,
	}, nil
}

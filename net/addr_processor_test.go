package net

import (
	"bytes"
	"testing"
)

func TestIP4PortFirstAddrSerializer(t *testing.T) {
	buffer := bytes.NewBuffer(make([]byte, 0, 5))
	p := NewPortFirstAddrProcessor()
	addrString := "127.0.0.1:8080"
	addr, err := ParseProxyAddr(addrString)
	if err != nil {
		t.Fatal(err)
	}
	err = p.WriteAddrTo(buffer, addr)
	if err != nil {
		t.Fatal(err)
	}
	if len(buffer.Bytes()) != 7 {
		t.Fatalf("addr serialize failed, expected 7 bytes, got %d", len(buffer.Bytes()))
	}
	proxyAddr, err := p.ReadAddrFrom(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if proxyAddr.String() != addrString {
		t.Errorf("proxy addr read expect %s but %s", addrString, proxyAddr.String())
	}
}

func TestIP6PortFirstAddrSerializer(t *testing.T) {
	buffer := bytes.NewBuffer(make([]byte, 0, 19))
	p := NewPortFirstAddrProcessor()
	addrString := "[2001:db8::1]:8080"
	addr, err := ParseProxyAddr(addrString)
	if err != nil {
		t.Fatal(err)
	}
	err = p.WriteAddrTo(buffer, addr)
	if err != nil {
		t.Fatal(err)
	}
	if len(buffer.Bytes()) != 19 {
		t.Fatalf("addr serialize failed, expected 19 bytes, got %d", len(buffer.Bytes()))
	}
	proxyAddr, err := p.ReadAddrFrom(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if proxyAddr.String() != addrString {
		t.Errorf("proxy addr read expect %s but %s", addrString, proxyAddr.String())
	}
}

func TestFQDNPortFirstAddrSerializer(t *testing.T) {
	buffer := bytes.NewBuffer(make([]byte, 0, 19))
	p := NewPortFirstAddrProcessor()
	addrString := "www.google.com:8080"
	addr, err := ParseProxyAddr(addrString)
	if err != nil {
		t.Fatal(err)
	}
	err = p.WriteAddrTo(buffer, addr)
	if err != nil {
		t.Fatal(err)
	}
	if len(buffer.Bytes()) != len("www.google.com")+1+3 { // addr + addr length + (port + addr type)
		t.Fatalf("addr serialize failed, expected %d bytes, got %d", len("www.google.com")+1+3, len(buffer.Bytes()))
	}
	proxyAddr, err := p.ReadAddrFrom(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if proxyAddr.String() != addrString {
		t.Errorf("proxy addr read expect %s but %s", addrString, proxyAddr.String())
	}
}

func TestIP4HostFirstAddrSerializer(t *testing.T) {
	buffer := bytes.NewBuffer(make([]byte, 0, 5))
	p := NewHostFirstAddrProcessor()
	addrString := "127.0.0.1:8080"
	addr, err := ParseProxyAddr(addrString)
	if err != nil {
		t.Fatal(err)
	}
	err = p.WriteAddrTo(buffer, addr)
	if err != nil {
		t.Fatal(err)
	}
	if len(buffer.Bytes()) != 7 {
		t.Fatalf("addr serialize failed, expected 7 bytes, got %d", len(buffer.Bytes()))
	}
	proxyAddr, err := p.ReadAddrFrom(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if proxyAddr.String() != addrString {
		t.Errorf("proxy addr read expect %s but %s", addrString, proxyAddr.String())
	}
}

func TestIP6HostFirstAddrSerializer(t *testing.T) {
	buffer := bytes.NewBuffer(make([]byte, 0, 19))
	p := NewHostFirstAddrProcessor()
	addrString := "[2001:db8::1]:8080"
	addr, err := ParseProxyAddr(addrString)
	if err != nil {
		t.Fatal(err)
	}
	err = p.WriteAddrTo(buffer, addr)
	if err != nil {
		t.Fatal(err)
	}
	if len(buffer.Bytes()) != 19 {
		t.Fatalf("addr serialize failed, expected 19 bytes, got %d", len(buffer.Bytes()))
	}
	proxyAddr, err := p.ReadAddrFrom(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if proxyAddr.String() != addrString {
		t.Errorf("proxy addr read expect %s but %s", addrString, proxyAddr.String())
	}
}

func TestFQDNHostFirstAddrSerializer(t *testing.T) {
	buffer := bytes.NewBuffer(make([]byte, 0, 19))
	p := NewHostFirstAddrProcessor()
	addrString := "www.google.com:8080"
	addr, err := ParseProxyAddr(addrString)
	if err != nil {
		t.Fatal(err)
	}
	err = p.WriteAddrTo(buffer, addr)
	if err != nil {
		t.Fatal(err)
	}
	if len(buffer.Bytes()) != len("www.google.com")+1+3 {
		t.Fatalf("addr serialize failed, expected %d bytes, got %d", len("www.google.com")+1+3, len(buffer.Bytes()))
	}
	proxyAddr, err := p.ReadAddrFrom(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if proxyAddr.String() != addrString {
		t.Errorf("proxy addr read expect %s but %s", addrString, proxyAddr.String())
	}
}

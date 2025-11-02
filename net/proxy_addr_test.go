package net

import (
	"testing"
)

func TestIP4HostPort(t *testing.T) {
	port, err := ParseProxyHostPort("127.0.0.1", 4243)
	if err != nil {
		t.Fatal(err)
	}
	if port.Port != 4243 {
		t.Errorf("port.Port is %d, want %d", port.Port, 4243)
		return
	}
	if port.Host.String() != "127.0.0.1" {
		t.Fatalf("port.Host.String is %s, want 127.0.0.1", port.Host.String())
	}
}

func TestIP6HostPort(t *testing.T) {
	port, err := ParseProxyHostPort("2001:db8::1", 4243)
	if err != nil {
		t.Fatal(err)
	}
	if port.Port != 4243 {
		t.Errorf("port.Port is %d, want %d", port.Port, 4243)
		return
	}
	if port.Host.String() != "2001:db8::1" {
		t.Fatalf("port.Host.String is %s, want 2001:db8::1", port.Host.String())
	}
}

func TestIP6WithSquareHostPort(t *testing.T) {
	_, err := ParseProxyHostPort("[2001:db8::1]", 4243)
	if err == nil {
		t.Fatal(err)
	}
}

func TestFQDNHostPort(t *testing.T) {
	port, err := ParseProxyHostPort("www.google.com", 4243)
	if err != nil {
		t.Fatal(err)
	}
	if port.Port != 4243 {
		t.Errorf("port.Port is %d, want %d", port.Port, 4243)
		return
	}
	if port.Host.String() != "www.google.com" {
		t.Fatalf("port.Host.String is %s, want www.google.com", port.Host.String())
	}
}

func TestIP4Addr(t *testing.T) {
	port, err := ParseProxyAddr("127.0.0.1:4243")
	if err != nil {
		t.Fatal(err)
	}
	if port.Port != 4243 {
		t.Errorf("port.Port is %d, want %d", port.Port, 4243)
		return
	}
	if port.Host.String() != "127.0.0.1" {
		t.Fatalf("port.Host.String is %s, want 127.0.0.1", port.Host.String())
	}
}

func TestIP6Addr(t *testing.T) {
	port, err := ParseProxyAddr("[2001:db8::1]:4243")
	if err != nil {
		t.Fatal(err)
	}
	if port.Port != 4243 {
		t.Errorf("port.Port is %d, want %d", port.Port, 4243)
		return
	}
	if port.Host.String() != "2001:db8::1" {
		t.Fatalf("port.Host.String is %s, want 2001:db8::1", port.Host.String())
	}
}

func TestIP6NoSquareAddr(t *testing.T) {
	_, err := ParseProxyAddr("2001:db8::1:4243")
	if err == nil {
		t.Fatal(err)
	}
}

func TestFQDNAddr(t *testing.T) {
	port, err := ParseProxyAddr("www.google.com:4243")
	if err != nil {
		t.Fatal(err)
	}
	if port.Port != 4243 {
		t.Errorf("port.Port is %d, want %d", port.Port, 4243)
		return
	}
	if port.Host.String() != "www.google.com" {
		t.Fatalf("port.Host.String is %s, want www.google.com", port.Host.String())
	}
}

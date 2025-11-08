package algorithm

import (
	"crypto/cipher"
	"golang.org/x/crypto/chacha20poly1305"
)

func NewChaCha20Poly1305(key []byte) cipher.AEAD {
	cc, err := chacha20poly1305.New(key)
	if err != nil {
		panic("failed to new chacha20poly1305: " + err.Error())
	}
	return cc
}

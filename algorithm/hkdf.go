package algorithm

import (
	"golang.org/x/crypto/hkdf"
	"hash"
	"io"
)

func HKDF(hash func() hash.Hash, byteSize uint32, secret, salt, info []byte) []byte {
	h := hkdf.New(hash, secret, salt, info)
	key := make([]byte, byteSize) // exp:  32byte = 256bit/8
	_, _ = io.ReadFull(h, key)
	return key
}

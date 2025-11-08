package algorithm

import (
	"crypto/sha256"
	"testing"
)

func TestHkdf(t *testing.T) {

	hkdfBytes := HKDF(sha256.New, 16, []byte("secret"), []byte("salt"), []byte("info"))
	if len(hkdfBytes) != 16 {
		t.Fatal("hkdfBytes is not of expected size")
	}
}

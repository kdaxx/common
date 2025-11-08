package algorithm

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"testing"
)

func TestChacha20Poly1305(t *testing.T) {
	key := sha256.New().Sum(nil)
	chaCha20Poly1305 := NewChaCha20Poly1305(key)
	plainText := []byte("hello world")
	cipherBytes := make([]byte, len(plainText)+chaCha20Poly1305.Overhead())

	nonce := make([]byte, chaCha20Poly1305.NonceSize())
	_, err := rand.Read(nonce)
	if err != nil {
		t.Fatal(err)
	}
	chaCha20Poly1305.Seal(cipherBytes[:0], nonce, plainText, nil)
	t.Logf("encrypted by chacha20Poly1305: %x", cipherBytes)

	plainBytes, err := chaCha20Poly1305.Open(cipherBytes[:0], nonce, cipherBytes, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("decrypted by chacha20Poly1305: %s", plainBytes)

	if !bytes.Equal(plainText, plainBytes) {
		t.Fatal("decrypted is not equal to plainText")
	}
}

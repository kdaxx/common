package algorithm

import (
	"bytes"
	"crypto/aes"
	"crypto/md5"
	"crypto/rand"
	"testing"
)

func TestAesCbc(t *testing.T) {
	key := md5.New().Sum(nil)
	iv := make([]byte, aes.BlockSize)
	_, err := rand.Read(iv)
	if err != nil {
		t.Fatal(err)
	}
	cbc := NewAesCbcEncrypter(key, iv)
	plainText := []byte("hello world")
	paddedText := Pkcs7Padding(plainText, 16)
	// encrypted
	cbc.CryptBlocks(paddedText, paddedText)
	t.Logf("encrypted by aescbc128: %x", paddedText)

	decrypter := NewAesCbcDecrypter(key, iv)
	// decrypted
	decrypter.CryptBlocks(paddedText, paddedText)
	plainBytes, err := Pkcs7UnPadding(paddedText)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("decrypted by aescbc128: %s", plainBytes)
	if !bytes.Equal(plainText, plainBytes) {
		t.Fatal("decrypted is not equal to paddedText")
	}
}

func TestAesCfb(t *testing.T) {
	key := md5.New().Sum(nil)
	iv := make([]byte, aes.BlockSize)
	_, err := rand.Read(iv)
	if err != nil {
		t.Fatal(err)
	}
	cfbEncrypter := NewAesCfbEncrypter(key, iv)
	plainText := []byte("hello world")
	cipherBytes := make([]byte, len(plainText))
	// encrypted
	cfbEncrypter.XORKeyStream(cipherBytes, plainText)
	t.Logf("encrypted by aescfb128: %x", cipherBytes)

	decrypter := NewAesCfbDecrypter(key, iv)
	// decrypted
	decrypter.XORKeyStream(cipherBytes, cipherBytes)
	t.Logf("decrypted by aescfb128: %s", cipherBytes)
	if !bytes.Equal(cipherBytes, plainText) {
		t.Fatal("decrypted is not equal to plainText")
	}
}

func TestAesCtr(t *testing.T) {
	key := md5.New().Sum(nil)
	iv := make([]byte, aes.BlockSize)
	_, err := rand.Read(iv)
	if err != nil {
		t.Fatal(err)
	}
	ctr := NewAesCtr(key, iv)
	plainText := []byte("hello world")
	cipherBytes := make([]byte, len(plainText))
	ctr.XORKeyStream(cipherBytes, plainText)
	t.Logf("encrypted by aesctr128: %x", cipherBytes)

	// The CTR stream is one-time use and will change the internal counter state.
	// The CTR must be reconstructed (using the same key + IV) to decrypt it.
	ctrDecrypter := NewAesCtr(key, iv)
	ctrDecrypter.XORKeyStream(cipherBytes, cipherBytes)
	t.Logf("decrypted by aesctr128: %s", cipherBytes)

	if !bytes.Equal(cipherBytes, plainText) {
		t.Fatal("decrypted is not equal to plainText")
	}
}

func TestAesGcm(t *testing.T) {
	key := md5.New().Sum(nil)
	aesGcm := NewAesGcm(key)
	plainText := []byte("hello world")
	cipherBytes := make([]byte, len(plainText)+aesGcm.Overhead())

	nonce := make([]byte, aesGcm.NonceSize())
	_, err := rand.Read(nonce)
	if err != nil {
		t.Fatal(err)
	}
	aesGcm.Seal(cipherBytes[:0], nonce, plainText, nil)
	t.Logf("encrypted by aesgcm128: %x", cipherBytes)

	plainBytes, err := aesGcm.Open(cipherBytes[:0], nonce, cipherBytes, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("decrypted by aesctr128: %s", plainBytes)

	if !bytes.Equal(plainText, plainBytes) {
		t.Fatal("decrypted is not equal to plainText")
	}
}

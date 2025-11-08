package algorithm

import (
	"crypto/aes"
	"crypto/cipher"
)

func NewAesCipher(key []byte) cipher.Block {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic("failed to initialize AES cipher: " + err.Error())
	}
	return block
}

// NewAesCbcEncrypter block mode
// returns an AESCBC encryption instance.
// When using block encryption, the plaintext must be padded first.
func NewAesCbcEncrypter(key []byte, iv []byte) cipher.BlockMode {
	block := NewAesCipher(key)
	return cipher.NewCBCEncrypter(block, iv)
}

// NewAesCbcDecrypter block mode
// returns an AESCBC decryption instance.
// When using block decryption, the decoded plaintext must be unPadded first.
func NewAesCbcDecrypter(key []byte, iv []byte) cipher.BlockMode {
	block := NewAesCipher(key)
	return cipher.NewCBCDecrypter(block, iv)
}

// NewAesCfbEncrypter stream mode
// returns an AESCFB encryption stream instance.
func NewAesCfbEncrypter(key, iv []byte) cipher.Stream {
	block := NewAesCipher(key)
	return cipher.NewCFBEncrypter(block, iv)
}

// NewAesCfbDecrypter stream mode
// returns an AESCFB decryption stream instance.
func NewAesCfbDecrypter(key, iv []byte) cipher.Stream {
	block := NewAesCipher(key)
	return cipher.NewCFBDecrypter(block, iv)
}

// NewAesCtr stream mode
// returns an AESCTR encryption instance.
func NewAesCtr(key, iv []byte) cipher.Stream {
	block := NewAesCipher(key)
	return cipher.NewCTR(block, iv)
}

// NewAesGcm aead mode
// returns an AESGCM encryption instance.
func NewAesGcm(key []byte) cipher.AEAD {
	block := NewAesCipher(key)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic("failed to initialize AES GCM: " + err.Error())
	}
	return gcm
}

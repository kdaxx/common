package algorithm

import (
	"bytes"
	"errors"
)

// Pkcs7Padding padding src to dst according to PKCS7 standard
func Pkcs7Padding(src []byte, blockSize int) []byte {
	paddingSize := blockSize - len(src)%blockSize
	paddingBytes := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)
	return append(src, paddingBytes...)
}

// Pkcs7UnPadding unPadding src to dst according to PKCS7 standard
func Pkcs7UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return src, nil
	}
	padding := int(src[length-1])
	if padding > length {
		return nil, errors.New("invalid padding size: padding size greater than data size")
	}
	return src[:length-padding], nil
}

package core

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"os"
)

func EncryptString(key string, value string) (string, error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	out := ""
	buf := make([]byte, c.BlockSize())
	for l := 0; l < len(value); l += len(buf) {
		end := l + len(buf)
		if end > len(value) {
			end = len(value)
			buf[end-l] = 0
		}
		copy(buf, value[l:end])
		c.Encrypt(buf, buf)
		out += hex.EncodeToString(buf)
	}

	return out, nil
}

func DecryptString(key string, value string) (string, error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext, _ := hex.DecodeString(value)
	blockSize := c.BlockSize()
	for l := 0; l < len(ciphertext); l += blockSize {
		c.Decrypt(ciphertext[l:l+blockSize], ciphertext[l:l+blockSize])
	}

	idx := bytes.IndexByte(ciphertext, 0)
	if idx < 0 || idx >= len(ciphertext) {
		return "", os.ErrInvalid
	}
	return string(ciphertext[0:idx]), nil
}


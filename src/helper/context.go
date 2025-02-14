package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func DecryptPayload(encryptedText, iv, key string) (string, error) {
	// Decode Base64
	ciphertext, _ := base64.StdEncoding.DecodeString(encryptedText)
	ivBytes, _ := base64.StdEncoding.DecodeString(iv)
	keyBytes := []byte(key)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(ciphertext, ciphertext)

	// Remove padding
	ciphertext = unpad(ciphertext)

	return string(ciphertext), nil
}

func unpad(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}

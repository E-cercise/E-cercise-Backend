package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/E-cercise/E-cercise/src/config"
	"github.com/gofiber/fiber/v2/log"
)

// DecryptValue takes an encrypted string in the format "IVBase64|CiphertextBase64"
// and decrypts it using AES-CBC with PKCS#7 padding.
func DecryptPayload(encryptedData string) (string, error) {

	key := []byte(config.SecretKey)
	// Split into IV and ciphertext
	parts := splitEncryptedData(encryptedData)
	if len(parts) != 2 {
		return "", errors.New("invalid encrypted data format")
	}
	ivBase64, ciphertextBase64 := parts[0], parts[1]

	// Decode IV
	iv, err := base64.StdEncoding.DecodeString(ivBase64)
	if err != nil {
		return "", fmt.Errorf("error decoding IV from Base64: %w", err)
	}

	// Decode ciphertext
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", fmt.Errorf("error decoding ciphertext from Base64: %w", err)
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("error creating AES cipher: %w", err)
	}

	// CBC mode needs IV length == block size
	if len(iv) != aes.BlockSize {
		return "", errors.New("invalid IV size")
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	// Decrypt in CBC mode
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	log.Info(plaintext)
	// Remove PKCS#7 padding
	plaintext, err = pkcs7Unpad(plaintext)
	if err != nil {
		return "", fmt.Errorf("error unpadding plaintext: %w", err)
	}

	return string(plaintext), nil
}

// pkcs7Unpad removes the PKCS#7 padding from decrypted data.
func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("invalid padding size")
	}

	// The last byte indicates the number of padding bytes
	padLen := int(data[length-1])
	if padLen == 0 || padLen > length {
		return nil, errors.New("invalid padding size")
	}

	// Check that all padding bytes match padLen
	for i := 0; i < padLen; i++ {
		if data[length-1-i] != byte(padLen) {
			return nil, errors.New("invalid padding")
		}
	}

	// Return the plaintext without the padding
	return data[:length-padLen], nil

}

// splitEncryptedData is a helper function to split "IVBase64|CiphertextBase64"
func splitEncryptedData(encryptedData string) []string {
	// You can just do strings.Split if the format is always "IV|Ciphertext".

	for i := range encryptedData {
		if encryptedData[i] == '|' {
			return []string{encryptedData[:i], encryptedData[i+1:]}
		}
	}
	return []string{}
}

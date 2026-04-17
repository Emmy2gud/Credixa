package utils

import (
	"bytes"
	"crypto/des"      // Go's built-in 3DES library
	"encoding/base64" // to convert encrypted bytes to readable text
	"fmt"
	"os"
)

func pkcs5Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	// repeat the byte based on the value [03 03 03],[02 02],[01]
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func Encryption3des(payload string) (string, error) {

	encryptionKey := os.Getenv("FLW_ENCRYPTION_KEY")

	if encryptionKey == "" {
		return "", fmt.Errorf("FLW_ENCRYPTION_KEY is not set")
	}

	// 3DES requires exactly 24-byte key; pad with null bytes if shorter
	keyBytes := make([]byte, 24)
	copy(keyBytes, []byte(encryptionKey))

	plainText := []byte(payload)

	// Pad the plaintext to a multiple of 8 bytes (3DES block size)
	padText := pkcs5Pad(plainText, des.BlockSize)

	block, err := des.NewTripleDESCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("cipher error: %v", err)
	}

	// Flutterwave requires ECB mode (no IV — each block encrypted independently)
	cipherText := make([]byte, len(padText))
	for i := 0; i < len(padText); i += 8 {
		block.Encrypt(cipherText[i:i+8], padText[i:i+8])
	}

	encoded := base64.StdEncoding.EncodeToString(cipherText)
	return encoded, nil
}

package utils

import (
	"crypto/rand"
	"math/big"
	"fmt"
)

func RandomString(letterBytes string, n int) (string, error) {
	b := make([]byte, n)
	
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %v", err)
		}
		b[i] = letterBytes[num.Int64()]
	}
	return string(b), nil
}
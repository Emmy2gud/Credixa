package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	
	"time"
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

func GenerateRequestID()(string,error){
		// creating requestid using date formatting for strings
	const AlphaNumericBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomString, err := RandomString(AlphaNumericBytes, 10)
	if err != nil {
		return "", err
	}

	loc, _ := time.LoadLocation("Africa/Lagos")
	now := time.Now().In(loc)
	requestid := now.Format("200601021504") + randomString
	fmt.Println("Request ID:", requestid)
	return requestid, nil
}
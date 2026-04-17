package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte("supersecretkey") // secret key to sign tokens

func GenerateToken(userID uint, role string) (string, error) {
	// 1. Create claims (data inside the token)
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,                                  // can be extended for roles like 'admin'
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // expires in 24 hours
	}

	// 2. Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 3. Sign the token using the secret key
	//the reason we use SignedString is to generate a signed token string from the token object using the provided secret key.
	signedToken, err := token.SignedString(JwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
func GeneratePasswordToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(JwtSecret)
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	// Parse and validate token
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure token method is HMAC algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return JwtSecret, nil
	})
}

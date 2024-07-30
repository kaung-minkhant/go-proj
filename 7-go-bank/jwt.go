package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	AccountNumber string `json:"account_number"`
}

func getAndValidateJWTToken(tokenString string) (*jwt.Token, error) {
	if tokenString == "" {
		ErrorLog("authorization token required")
		return nil, fmt.Errorf("permission denied")
	}
	authChunks := strings.Split(strings.TrimSpace(tokenString), " ")
	if len(authChunks) != 2 {
		ErrorLog("malformed authorization token")
		return nil, fmt.Errorf("permission denied")
	}
	if authChunks[0] != "Bearer" {
		ErrorLog("malformed authorization token")
		return nil, fmt.Errorf("permission denied")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	token, err := jwt.ParseWithClaims(authChunks[1], &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			ErrorLog("invalid signing method: %v", t.Header["alg"])
			return nil, fmt.Errorf("permission denied")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		ErrorLog("cannot verify jwt with error %v", err)
		return nil, fmt.Errorf("permission denied")
	}
	// if !token.Valid {
	// 	return nil, fmt.Errorf("invalid token")
	// }
	// claims, ok := token.Claims.(*JWTClaims)
	// if !ok {
	// 	ErrorLog("malformed authorization token")
	// 	return nil, fmt.Errorf("permission denied")
	// }
	// fmt.Printf("Claims: %#v \n", claims)
	return token, nil
}

func makeJWTToken(account *Account) (string, error) {
	claims := &JWTClaims{
		jwt.RegisteredClaims{
			Issuer:    "gobank",
			ID:        account.ID.String(),
			Audience:  []string{account.FirstName + account.LastName},
			Subject:   "account jwt token",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1000000 * time.Hour)),
		},
		account.Number.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	ss, err := token.SignedString([]byte(secret))
	if err != nil {
		ErrorLog("cannot issue jwt with error %v", err)
		return "", fmt.Errorf("permission denied")
	}
	return ss, nil
}

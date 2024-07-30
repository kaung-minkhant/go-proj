package main

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(password string) (string, error) {
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(encryptedPass), err
}

func ComparePassword(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

package main

import "golang.org/x/crypto/bcrypt"

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 10)
}

func checkPassword(password string, hashedPass string) bool {
	return nil == bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
}

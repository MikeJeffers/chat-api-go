package chat

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 10)
}

func CheckPassword(password string, hashedPass string) bool {
	return nil == bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
}

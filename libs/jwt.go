package chat

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func SignToken(user UserRow) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		Id:       user.Id,
		Username: user.Username,
	})
	return token.SignedString([]byte(SECRET_JWT))
}

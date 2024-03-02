package chat

import "github.com/golang-jwt/jwt"

type CustomClaims struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func SignToken(user UserRow) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		Id:       user.Id,
		Username: user.Username,
	})
	return token.SignedString([]byte(SECRET_JWT))
}

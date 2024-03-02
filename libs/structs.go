package chat

type UserRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserTokenData struct {
	Id       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
}

type ReturnedUserData struct {
	Token string        `json:"token"`
	User  UserTokenData `json:"user"`
}

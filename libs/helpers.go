package chat

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func StoreAndRespond(user UserRow, redisClient *redis.Client, c *gin.Context) (ReturnedUserData, error) {
	signed, err := SignToken(user)
	if err != nil {
		return ReturnedUserData{}, err
	}
	status := StoreToken(signed, user.Id, redisClient, c)
	if status.Err() != nil {
		return ReturnedUserData{}, status.Err()
	}
	data := ReturnedUserData{Token: signed, User: UserTokenData{Id: user.Id, Username: user.Username}}
	return data, nil
}

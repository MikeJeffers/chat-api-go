package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/redis/go-redis/v9"
)

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

func storeAndRespond(user UserRow, redisClient *redis.Client, c *gin.Context) (ReturnedUserData, error) {
	signed, err := signToken(user)
	if err != nil {
		return ReturnedUserData{}, err
	}
	status := storeToken(signed, user.Id, redisClient, c)
	if status.Err() != nil {
		return ReturnedUserData{}, status.Err()
	}
	data := ReturnedUserData{Token: signed, User: UserTokenData{user.Id, user.Username}}
	return data, nil
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	db := dbConnect()
	red := redisClient()

	r.POST("/login", func(c *gin.Context) {
		var requestBody UserRequestBody

		if err := c.BindJSON(&requestBody); err != nil {
			return
		} else {
			user, err := getUserRow(requestBody.Username, db, c)
			if err != nil {
				c.JSON(400, err.Error())
				return
			} else if !checkPassword(requestBody.Password, user.Password) {
				c.JSON(400, gin.H{"message": "No such user"})
				return
			}
			data, err := storeAndRespond(user, red, c)
			if err != nil {
				c.JSON(500, gin.H{"message": "Server Error"})
				return
			}
			c.JSON(200, data)
		}
	})

	r.POST("/register", func(c *gin.Context) {
		var requestBody UserRequestBody

		if err := c.BindJSON(&requestBody); err != nil {
			return
		} else {
			hashed, err := hashPassword(requestBody.Password)
			if err != nil {
				c.JSON(500, gin.H{"message": "Server Error"})
				log.Println(err.Error())
				return
			}
			user, err := addUser(requestBody.Username, string(hashed), db, c)
			if err != nil {
				c.JSON(400, err.Error())
				return
			}
			data, err := storeAndRespond(user, red, c)
			if err != nil {
				c.JSON(500, gin.H{"message": "Server Error"})
				return
			}
			c.JSON(200, data)
		}
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":3000")
}

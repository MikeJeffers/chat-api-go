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

type UserRow struct {
	Id       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
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
			users := []UserRow{}
			db.Select(&users, "SELECT id, username, password FROM users WHERE username = $1 LIMIT 1", requestBody.Username)
			if len(users) < 1 {
				c.JSON(400, gin.H{"message": "No such user"})
				return
			}
			if !checkPassword(requestBody.Password, users[0].Password) {
				c.JSON(400, gin.H{"message": "No such user"})
				return
			}
			data, err := storeAndRespond(users[0], red, c)
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
			users := []UserRow{}
			db.Select(&users, "INSERT INTO Users (username, password) VALUES ($1, $2) RETURNING id, username, password", requestBody.Username, string(hashed))
			if len(users) < 1 {
				c.JSON(400, gin.H{"message": "unable to add user"})
				return
			}
			data, err := storeAndRespond(users[0], red, c)
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

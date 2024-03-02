package chat

import (
	libs "chat/libs"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func Register(db *sqlx.DB, red *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody libs.UserRequestBody

		if err := c.BindJSON(&requestBody); err != nil {
			return
		} else {
			hashed, err := libs.HashPassword(requestBody.Password)
			if err != nil {
				c.JSON(500, gin.H{"message": "Server Error"})
				log.Println(err.Error())
				return
			}
			user, err := libs.AddUser(requestBody.Username, string(hashed), db, c)
			if err != nil {
				c.JSON(400, err.Error())
				return
			}
			data, err := libs.StoreAndRespond(user, red, c)
			if err != nil {
				c.JSON(500, gin.H{"message": "Server Error"})
				return
			}
			c.JSON(201, data)
		}
	}
}

package chat

import (
	libs "github.com/mikejeffers/chat-api-go/libs"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func Login(db *sqlx.DB, red *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody libs.UserRequestBody

		if err := c.BindJSON(&requestBody); err != nil {
			return
		} else {
			user, err := libs.GetUserRow(requestBody.Username, db, c)
			if err != nil {
				c.JSON(400, err.Error())
				return
			} else if !libs.CheckPassword(requestBody.Password, user.Password) {
				c.JSON(400, gin.H{"message": "No such user"})
				return
			}
			data, err := libs.StoreAndRespond(user, red, c)
			if err != nil {
				c.JSON(500, gin.H{"message": "Server Error"})
				return
			}
			c.JSON(200, data)
		}
	}
}

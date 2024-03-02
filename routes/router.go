package chat

import (
	libs "github.com/mikejeffers/chat-api-go/libs"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()
	db := libs.DbConnect()
	red := libs.RedisClient()

	r.POST("/login", Login(db, red))
	r.POST("/register", Register(db, red))

	return r
}

package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
)

func dbConnect() *sql.DB {
	connStr := "postgresql://user:password@localhost:5432/test?sslmode=disable"
	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

type UserAuthBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	// db := dbConnect()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.POST("/login", func(c *gin.Context) {
		var requestBody UserAuthBody

		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(400, gin.H{"message": "bad data"})
		} else {
			c.JSON(200, requestBody)
		}
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

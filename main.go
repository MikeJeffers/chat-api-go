package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	DB_USER     = os.Getenv("POSTGRES_USER")
	DB_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	DB_NAME     = os.Getenv("POSTGRES_DB")
	DB_HOST     = os.Getenv("POSTGRES_HOST")
	DB_PORT     = os.Getenv("POSTGRES_PORT")
)

func dbConnect() *sqlx.DB {
	connStr := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	fmt.Println(connStr)
	// Connect to database
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

type UserAuthBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRow struct {
	Id       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	db := dbConnect()

	r.POST("/login", func(c *gin.Context) {
		var requestBody UserAuthBody

		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(400, gin.H{"message": "bad data"})
		} else {
			users := []UserRow{}
			db.Select(&users, "SELECT id, username, password FROM users WHERE username = $1 LIMIT 1", requestBody.Username)
			if len(users) < 1 {
				c.JSON(400, gin.H{"message": "no such user"})
				return
			}
			c.JSON(200, users[0])
		}
	})

	r.POST("/register", func(c *gin.Context) {
		var requestBody UserAuthBody

		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(400, gin.H{"message": "bad data"})
		} else {
			users := []UserRow{}
			db.Select(&users, "INSERT INTO Users (username, password) VALUES ($1, $2) RETURNING id, username, password", requestBody.Username, requestBody.Password)
			if len(users) < 1 {
				c.JSON(400, gin.H{"message": "unable to add user"})
				return
			}
			c.JSON(200, users[0])
		}
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

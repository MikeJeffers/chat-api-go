package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type UserRow struct {
	Id       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

func dbConnect() *sqlx.DB {
	connStr := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func getUserRow(username string, db *sqlx.DB, c *gin.Context) (UserRow, error) {
	users := []UserRow{}
	db.Select(&users, "SELECT id, username, password FROM users WHERE username = $1 LIMIT 1", username)
	if len(users) != 1 {
		return UserRow{}, fmt.Errorf("could not find user with username=%v", username)
	}
	return users[0], nil
}

func addUser(username, hashedPassword string, db *sqlx.DB, c *gin.Context) (UserRow, error) {
	users := []UserRow{}
	db.Select(&users, "INSERT INTO Users (username, password) VALUES ($1, $2) RETURNING id, username, password", username, hashedPassword)
	if len(users) < 1 {
		return UserRow{}, fmt.Errorf("failed to insert new user")
	}
	return users[0], nil
}

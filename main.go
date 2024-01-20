package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) < 1 {
		return defaultValue
	}
	return value
}

var (
	DB_USER     = getEnv("POSTGRES_USER", "user")
	DB_PASSWORD = getEnv("POSTGRES_PASSWORD", "password")
	DB_NAME     = getEnv("POSTGRES_DB", "test")
	DB_HOST     = getEnv("POSTGRES_HOST", "localhost")
	DB_PORT     = getEnv("POSTGRES_PORT", "5432")

	REDIS_HOST     = getEnv("REDIS_HOST", "localhost")
	REDIS_PORT     = getEnv("REDIS_PORT", "6379")
	REDIS_PASSWORD = getEnv("REDIS_PASSWORD", "")

	SECRET_JWT = getEnv("SECRET_JWT", "idk")
)

func dbConnect() *sqlx.DB {
	connStr := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func redisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", REDIS_HOST, REDIS_PORT),
		Password: REDIS_PASSWORD,
		DB:       0,
		Protocol: 3,
	})
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

type UserTokenData struct {
	Id       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
}

type ReturnedUserData struct {
	Token string        `json:"token"`
	User  UserTokenData `json:"user"`
}

type MyCustomClaims struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func signToken(user UserRow) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		Id:       user.Id,
		Username: user.Username,
	})
	return token.SignedString([]byte(SECRET_JWT))
}

func storeToken(token string, userId int64, redisClient *redis.Client, c *gin.Context) *redis.StatusCmd {
	return redisClient.SetEx(c, fmt.Sprintf("jwt:%v", userId), token, time.Hour*1)
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

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 10)
}

func checkPassword(password string, hashedPass string) bool {
	return nil == bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	db := dbConnect()
	red := redisClient()

	r.POST("/login", func(c *gin.Context) {
		var requestBody UserAuthBody

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
		var requestBody UserAuthBody

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

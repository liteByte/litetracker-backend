package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func login(c *gin.Context) {
	data := parseRequestBody(c)

	username := data["username"].(string)
	password := data["password"].(string)

	hashedPassword := hashPassword(username, password)

	var databaseUsername string
	var databasePassword string
	var databaseId string
	var databaseAdmin bool

	err := db.QueryRow("SELECT id, username, password, admin FROM users WHERE username=?", username).Scan(&databaseId, &databaseUsername, &databasePassword, &databaseAdmin)

	if err != nil {
		content := gin.H{"success": "false", "description": "Wrong username"}
		c.JSON(200, content)
		return
	}

	if hashedPassword != databasePassword {
		content := gin.H{"success": "false", "description": "Wrong password"}
		c.JSON(200, content)
		return
	}

	content := gin.H{"success": "true", "description": createToken(databaseId, username, databaseAdmin)}
	c.JSON(200, content)
}

type MyCustomClaims struct {
	Id    string
	Admin bool
	jwt.StandardClaims
}

func createToken(id string, username string, admin bool) string {
	claims := MyCustomClaims{
		id,
		admin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	checkErr(err)

	return tokenString
}

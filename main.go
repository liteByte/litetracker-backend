package main

import (
	"os"
    "github.com/gin-gonic/gin"
    "database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error
var router *gin.Engine;
var jwtSecretKey = []byte(jwtSecret)

func main(){
	connectToDatabase()
  	createRouter()
  	router.Run(":" + getPort())
}

func getPort() string {
    port := os.Getenv("PORT")

    if port == "" {
        port = defaultPort
    }

    return port
}
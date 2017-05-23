package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
)

func connectToDatabase() {
	db, err = sql.Open(dbType, dbUsername+":"+dbPassword+"@tcp("+dbHost+":"+dbPort+")/"+dbName)
	checkErr(err)
	err = db.Ping()
	checkErr(err)
}

func parseRequestBody(c *gin.Context) map[string]interface{} {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var data map[string]interface{}

	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}

	return data
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func hashPassword(username string, password string) string {
	hasher := sha1.New()
	hasher.Write([]byte(username + password))
	hashedPassword := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return hashedPassword
}

func getUsernameWithId(id string) string {
	var databaseUsername string
	err := db.QueryRow("SELECT username FROM users WHERE id=?", id).Scan(&databaseUsername)
	if err != nil {
		return "error"
	}

	return databaseUsername
}

func getProjectWithId(id string) string {
	var databaseProject string
	err := db.QueryRow("SELECT name FROM projects WHERE id=?", id).Scan(&databaseProject)
	if err != nil {
		return "error"
	}

	return databaseProject
}

package main

import(
    "encoding/json"
    "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
    Id string
    Username string
}

func createUser (c *gin.Context){
    data := parseRequestBody(c)

    username := data["username"].(string)
    password := data["password"].(string)

    hashedPassword := hashPassword(username, password)

    _, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
    if err != nil {
        content := gin.H{"success":"false", "description":"Username already exists"}
        c.JSON(200, content)
        return
    }

    content := gin.H{"success": "true", "description":""}
    c.JSON(200, content)
}

func modifyUser (c *gin.Context){
    data := parseRequestBody(c)
    id_user := c.Params.ByName("id_user")

    username := data["username"].(string)

    _, err = db.Exec("UPDATE users SET username = ? WHERE id = ?", username, id_user)
    if err != nil {
        content := gin.H{"success":"false", "description":"Username already exists"}
        c.JSON(200, content)
        return
    }

    content := gin.H{"success": "true", "description":""}
    c.JSON(200, content)
}

func deleteUser (c *gin.Context){
    id_user := c.Params.ByName("id_user")

    _, err = db.Exec("DELETE FROM users WHERE id = ?", id_user)
    checkErr(err)

    content := gin.H{"success": "true", "description":""}
    c.JSON(200, content)
}

func getUserList(c *gin.Context) {
    var name string
    var id string
    userList := make([]User, 0)

    rows, err := db.Query("select id, username from users")
    checkErr(err)
    defer rows.Close()

    for rows.Next() {
        err := rows.Scan(&id, &name)
        checkErr(err)

        user := User{Id: id, Username: name}

        userList = append(userList, user)
    }

    response, err := json.Marshal(userList)

    c.JSON(200, string(response))
}


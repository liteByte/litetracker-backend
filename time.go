package main

import(
    "strings"
    "encoding/json"
    "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Time struct {
    Id string
    IdUsersProjects string
    Hours string
    Description string
    Date string
}

func createTime (c *gin.Context){
    data := parseRequestBody(c)

    id_users_projects := data["id_users_projects"].(string)
    hours := data["hours"].(string)
    description := data["description"].(string)
    date := data["date"].(string)

    _, err = db.Exec("INSERT INTO time(id_users_projects, hours, description, date) VALUES(?, ?, ?, ?)", id_users_projects, hours, description, date)
    checkErr(err)

    content := gin.H{"success": "true", "description":""}
    c.JSON(200, content)
}

func modifyTime (c *gin.Context){
    data := parseRequestBody(c)
    id_time := c.Params.ByName("id_time")

    id_users_projects := data["id_users_projects"].(string)
    hours := data["hours"].(string)
    description := data["description"].(string)

    _, err = db.Exec("UPDATE time SET id_users_projects = ?, hours = ?, description = ? WHERE id = ?", id_users_projects, hours, description, id_time)
    checkErr(err)

    content := gin.H{"success": "true", "description":""}
    c.JSON(200, content)
}

func deleteTime (c *gin.Context){
    id_time := c.Params.ByName("id_time")

    _, err = db.Exec("DELETE FROM time WHERE id = ?", id_time)
    checkErr(err)

    content := gin.H{"success": "true", "description":""}
    c.JSON(200, content)
}

func getTimeList(c *gin.Context) {
    start_date := c.Params.ByName("start_date")
    end_date := c.Params.ByName("end_date")

    var id string
    var id_users_projects string
    var hours string
    var description string
    var date string

    timeList := make([]Time, 0)

    rows, err := db.Query("select id, id_users_projects, hours, description, date from time WHERE (date BETWEEN ? AND ?)", start_date, end_date)
    checkErr(err)
    defer rows.Close()

    for rows.Next() {
        err := rows.Scan(&id, &id_users_projects, &hours, &description, &date)
        checkErr(err)

        time := Time{Id: id, IdUsersProjects: id_users_projects, Hours: hours, Description: description, Date: date}

        timeList = append(timeList, time)
    }

    response, err := json.Marshal(timeList)
    c.JSON(200, string(response))
}

func getTimeListOfUser(c *gin.Context) {
    id_user := c.Params.ByName("id_user")
    start_date := c.Params.ByName("start_date")
    end_date := c.Params.ByName("end_date")

    id_users_projects_db := getUsersProjectsOfUser(id_user)

    var id string
    var id_users_projects string
    var hours string
    var description string
    var date string

    timeList := make([]Time, 0)

    rows, err := db.Query(`select id, id_users_projects, hours, description, date from time WHERE (date BETWEEN ? AND ?) AND id_users_projects IN ("` + strings.Join(id_users_projects_db,`","`) + `")`, start_date, end_date)
    checkErr(err)
    defer rows.Close()

    for rows.Next() {
        err := rows.Scan(&id, &id_users_projects, &hours, &description, &date)
        checkErr(err)

        time := Time{Id: id, IdUsersProjects: id_users_projects, Hours: hours, Description: description, Date: date}

        timeList = append(timeList, time)
    }

    response, err := json.Marshal(timeList)
    c.JSON(200, string(response))
}

func getTimeListOfProject(c *gin.Context) {
    id_project := c.Params.ByName("id_project")
    start_date := c.Params.ByName("start_date")
    end_date := c.Params.ByName("end_date")

    id_users_projects_db := getUsersProjectsOfProject(id_project)

    var id string
    var id_users_projects string
    var hours string
    var description string
    var date string

    timeList := make([]Time, 0)

    rows, err := db.Query("select id, id_users_projects, hours, description, date from time WHERE (date BETWEEN ? AND ?) AND id_users_projects IN (?)", start_date, end_date, id_users_projects_db)
    checkErr(err)
    defer rows.Close()

    for rows.Next() {
        err := rows.Scan(&id, &id_users_projects, &hours, &description, &date)
        checkErr(err)

        time := Time{Id: id, IdUsersProjects: id_users_projects, Hours: hours, Description: description, Date: date}

        timeList = append(timeList, time)
    }

    response, err := json.Marshal(timeList)
    c.JSON(200, string(response))
}

func getUsersProjectsOfUser (id_user string) []string {
    var id string

    slice := []string{}

    rows, err := db.Query("select id from users_projects WHERE id_user = ?", id_user)
    checkErr(err)
    defer rows.Close()

    for rows.Next() {
        err := rows.Scan(&id)
        checkErr(err)

        slice = append(slice, id)
    }

    return slice
}

func getUsersProjectsOfProject (id_project string) string {
    var id string
    var list string

    rows, err := db.Query("select id from users_projects WHERE id_project = ?", id_project)
    checkErr(err)
    defer rows.Close()

    for rows.Next() {
        err := rows.Scan(&id)
        checkErr(err)

        list += id
        list += ", "
    }

    if list == "" {
        return ""
    }

    list = list[:len(list)-2]

    return list
}
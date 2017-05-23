package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Assignment struct {
	Id        string
	IdUser    string
	IdProject string
}

func createAssignment(c *gin.Context) {
	data := parseRequestBody(c)

	username := data["username"].(string)
	id_project := data["id_project"].(string)

	id_user := getUserIdWithUsername(username)

	var databaseUserId string
	var databaseProjectId string

	err := db.QueryRow("SELECT id_user, id_project FROM users_projects WHERE id_user = ? AND id_project = ?", id_user, id_project).Scan(&databaseUserId, &databaseProjectId)

	if err == nil {
		content := gin.H{"success": "false", "description": "Assignment already exists"}
		c.JSON(200, content)
		return
	}

	_, err = db.Exec("INSERT INTO users_projects(id_user, id_project) VALUES(?, ?)", id_user, id_project)
	checkErr(err)

	content := gin.H{"success": "true", "description": ""}
	c.JSON(200, content)
}

func modifyAssignment(c *gin.Context) {
	data := parseRequestBody(c)
	id_assignment := c.Params.ByName("id_assignment")

	name := data["name"].(string)
	color := data["color"].(string)

	_, err = db.Exec("UPDATE users_projects SET name = ?, color = ? WHERE id = ?", name, color, id_assignment)
	if err != nil {
		content := gin.H{"success": "false", "description": "Assignment name already exists"}
		c.JSON(200, content)
		return
	}

	content := gin.H{"success": "true", "description": ""}
	c.JSON(200, content)
}

func deleteAssignment(c *gin.Context) {
	id_assignment := c.Params.ByName("id_assignment")

	_, err = db.Exec("DELETE FROM users_projects WHERE id = ?", id_assignment)
	checkErr(err)

	content := gin.H{"success": "true", "description": ""}
	c.JSON(200, content)
}

func getAssignmentList(c *gin.Context) {
	var id string
	var id_user string
	var id_project string

	assignmentList := make([]Assignment, 0)

	rows, err := db.Query("select id, id_user, id_project from users_projects order by 'date' desc")
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &id_user, &id_project)
		checkErr(err)

		assignment := Assignment{Id: id, IdUser: id_user, IdProject: id_project}

		assignmentList = append(assignmentList, assignment)
	}

	response, err := json.Marshal(assignmentList)
	c.JSON(200, string(response))
}

func getAssignmentListOfUser(c *gin.Context) {
	id_user_db := c.Params.ByName("id_user")

	var id string
	var id_user string
	var id_project string

	assignmentList := make([]Assignment, 0)

	rows, err := db.Query("select id, id_user, id_project from users_projects WHERE id_user = ? order by 'date' desc", id_user_db)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &id_user, &id_project)
		checkErr(err)

		assignment := Assignment{Id: id, IdUser: id_user, IdProject: id_project}

		assignmentList = append(assignmentList, assignment)
	}

	response, err := json.Marshal(assignmentList)
	c.JSON(200, string(response))
}

func getAssignmentListOfProject(c *gin.Context) {
	id_project_db := c.Params.ByName("id_project")

	var id string
	var id_user string
	var id_project string

	assignmentList := make([]Assignment, 0)

	rows, err := db.Query("select id, id_user, id_project from users_projects WHERE id_project = ? order by 'date' desc", id_project_db)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &id_user, &id_project)
		checkErr(err)

		assignment := Assignment{Id: id, IdUser: id_user, IdProject: id_project}

		assignmentList = append(assignmentList, assignment)
	}

	response, err := json.Marshal(assignmentList)
	c.JSON(200, string(response))
}

func getUserIdWithUsername(username string) string {
	var id_user string

	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&id_user)

	if err == nil {
		checkErr(err)
	}

	return id_user
}

package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Project struct {
	Id    string
	Name  string
	Color string
}

func createProject(c *gin.Context) {
	data := parseRequestBody(c)

	name := data["name"].(string)
	color := "#ACACAC" //data["color"].(string)

	_, err = db.Exec("INSERT INTO projects(name, color) VALUES(?, ?)", name, color)
	if err != nil {
		content := gin.H{"success": "false", "description": "Project name already exists"}
		c.JSON(200, content)
		return
	}

	content := gin.H{"success": "true", "description": ""}
	c.JSON(200, content)
}

func modifyProject(c *gin.Context) {
	data := parseRequestBody(c)
	id_project := c.Params.ByName("id_project")

	name := data["name"].(string)
	color := data["color"].(string)

	_, err = db.Exec("UPDATE projects SET name = ?, color = ? WHERE id = ?", name, color, id_project)
	if err != nil {
		content := gin.H{"success": "false", "description": "Project name already exists"}
		c.JSON(200, content)
		return
	}

	content := gin.H{"success": "true", "description": ""}
	c.JSON(200, content)
}

func deleteProject(c *gin.Context) {
	id_project := c.Params.ByName("id_project")

	_, err = db.Exec("DELETE FROM projects WHERE id = ?", id_project)
	checkErr(err)

	content := gin.H{"success": "true", "description": ""}
	c.JSON(200, content)
}

func getProjectList(c *gin.Context) {
	var name string
	var id string
	var color string

	projectList := make([]Project, 0)

	rows, err := db.Query("select id, name, color from projects")
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name, &color)
		checkErr(err)

		project := Project{Id: id, Name: name, Color: color}

		projectList = append(projectList, project)
	}

	response, err := json.Marshal(projectList)
	c.JSON(200, string(response))
}

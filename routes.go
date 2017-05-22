package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TokenData struct {
	Id       string
	Admin    bool
	Username string
}

func verifyToken(tokenString string) TokenData {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		return TokenData{Id: "-1", Admin: false, Username: ""}
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return TokenData{Id: claims.Id, Admin: claims.Admin, Username: claims.StandardClaims.Issuer}
	} else {
		return TokenData{Id: "-1", Admin: false, Username: ""}
	}
}

func createRouter() {
	router = gin.Default()

	router.LoadHTMLGlob("templates/*.tmpl")

	router.POST("/login", login)

	userGroup := router.Group("/users")
	{
		userGroup.GET("/list", getUserList)
		userGroup.POST("/create", createUser)
		userGroup.PUT("/modify/:id_user", modifyUser)
		userGroup.DELETE("/delete/:id_user", deleteUser)
	}

	projectGroup := router.Group("/projects")
	{
		projectGroup.GET("/list", getProjectList)
		projectGroup.POST("/create", createProject)
		projectGroup.PUT("/modify/:id_project", modifyProject)
		projectGroup.DELETE("/delete/:id_project", deleteProject)
	}

	assignmentGroup := router.Group("/assignments")
	{
		assignmentGroup.GET("/list", getAssignmentList)
		assignmentGroup.GET("/list/user/:id_user", getAssignmentListOfUser)
		assignmentGroup.GET("/list/project/:id_project", getAssignmentListOfProject)
		assignmentGroup.POST("/create", createAssignment)
		assignmentGroup.DELETE("/delete/:id_assignment", deleteAssignment)
	}

	timeGroup := router.Group("/time")
	{
		timeGroup.GET("/list/:start_date/:end_date", getTimeList)
		timeGroup.GET("/list/:start_date/:end_date/user/:id_user", getTimeListOfUser)
		timeGroup.GET("/list/:start_date/:end_date/project/:id_project", getTimeListOfProject)
		timeGroup.POST("/create", createTime)
		timeGroup.PUT("/modify/:id_time", modifyTime)
		timeGroup.DELETE("/delete/:id_time", deleteTime)
	}

	appGroup := router.Group("/app")
	{
		appGroup.Static("/css", "templates/css")
		appGroup.Static("/fonts", "templates/fonts")
		appGroup.Static("/js", "templates/js")

		appGroup.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.tmpl", nil)
		})

		appGroup.GET("/calendar/:token", func(c *gin.Context) {
			jwtAuth := verifyToken(c.Params.ByName("token"))

			data := gin.H{
				"id_user": jwtAuth.Id,
				"admin":   jwtAuth.Admin,
			}

			c.HTML(http.StatusOK, "calendar.tmpl", data)
		})

		appGroup.GET("/users", func(c *gin.Context) {
			c.HTML(http.StatusOK, "user_list.tmpl", nil)
		})

		appGroup.GET("/projects", func(c *gin.Context) {
			c.HTML(http.StatusOK, "project_list.tmpl", nil)
		})

		appGroup.GET("/reports", func(c *gin.Context) {
			c.HTML(http.StatusOK, "reports.tmpl", nil)
		})

		appGroup.GET("/project/:id_project", func(c *gin.Context) {
			id_project := c.Params.ByName("id_project")
			project_name := getProjectWithId(id_project)

			data := gin.H{
				"id_project":   id_project,
				"project_name": project_name,
			}

			c.HTML(http.StatusOK, "project.tmpl", data)
		})

		appGroup.GET("/user/:id_user", func(c *gin.Context) {
			id_user := c.Params.ByName("id_user")
			username := getUsernameWithId(id_user)

			data := gin.H{
				"id_user":  id_user,
				"username": username,
			}

			c.HTML(http.StatusOK, "user.tmpl", data)
		})

		appGroup.GET("/profile/:token", func(c *gin.Context) {

			jwtAuth := verifyToken(c.Params.ByName("token"))

			data := gin.H{
				"username": jwtAuth.Username,
				"id_user":  jwtAuth.Id,
				"admin":    jwtAuth.Admin,
			}

			c.HTML(http.StatusOK, "profile.tmpl", data)
		})
	}

}

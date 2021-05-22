package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/private-square/bkst-users-api/controllers"
	"github.com/private-square/bkst-users-api/services"
	"github.com/private-square/bkst-users-api/utils"
	"log"
	"os"
)

func StartApp() {
	gin.SetMode(gin.ReleaseMode)
	r := NewRouter()
	SetupRoutes(r)

	udb := &services.UsersDbConn{
		Hostname: os.Getenv("USERSDB_HOST"),
		Port:     os.Getenv("USERSDB_PORT"),
		Schema:   os.Getenv("USERSDB_SCHEMA"),
		Username: os.Getenv("USERSDB_USERNAME"),
		Password: os.Getenv("USERSDB_PASSWORD"),
	}

	if err := udb.Open(); err != nil {
		log.Fatalln(err)
	} else {
		fmt.Printf("Using external database: %s:%s/%s\n", udb.Hostname, udb.Port, udb.Schema)
	}

	fmt.Println("Starting REST API server on port 8080")

	if err := r.Run(":8080"); err != nil {
		log.Fatalln("Unable to run the web server.")
	}
}

func NewRouter() *gin.Engine {

	r := gin.Default()
	r.NoRoute(utils.NoRoute)
	r.HandleMethodNotAllowed = true
	r.NoMethod(utils.MethodNotAllowed)

	return r
}

func SetupRoutes(r *gin.Engine) *gin.Engine {

	usersApiPath := "/users"

	r.GET("/health", utils.Health)
	r.GET(usersApiPath+"/:userId", controllers.GetUser)
	r.GET(usersApiPath+"/search", controllers.SearchUser)
	r.POST(usersApiPath, controllers.CreateUser)
	r.PUT(usersApiPath+"/:userId", controllers.UpdateUser)
	r.DELETE(usersApiPath+"/:userId", controllers.DeleteUser)

	return r
}

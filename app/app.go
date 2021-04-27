package app

import (
	"github.com/gin-gonic/gin"
	"github.com/private-square/bkst-users-api/controllers"
	"github.com/private-square/bkst-users-api/utils"
	"log"
)

func StartApp() {
	gin.SetMode(gin.ReleaseMode)
	r := NewRouter()
	SetupRoutes(r)

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
	r.PUT(usersApiPath, controllers.UpdateUser)
	r.DELETE(usersApiPath+"/:userId", controllers.DeleteUser)

	return r
}

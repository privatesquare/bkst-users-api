package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/private-square/bkst-users-api/config"
	"github.com/private-square/bkst-users-api/controllers"
	"github.com/private-square/bkst-users-api/domain/users"
	"github.com/private-square/bkst-users-api/utils/httputils"
	"github.com/private-square/bkst-users-api/utils/logger"
)

const (
	defaultWebServerPort   = "8080"
	externalDBMsg          = "Using external database: %s:%s/%s"
	apiServerStartingMsg   = "Starting the API server..."
	apiServerStartedMsg    = "The API server has started and is listening on %s"
	apiServerStartupErrMsg = "Unable to run the web server: %v"

	apiHealthPath     = "/health"
	apiUsersPath      = "/users"
	apiUserIdParamExt = "/:userId"
	apiSearchPathExt  = "/search"
)

func StartApp() {
	gin.SetMode(gin.ReleaseMode)
	r := NewRouter()
	SetupRoutes(r)

	udb := &users.UserDbConn{
		Driver:   config.GlobalCnf.DBDriver,
		Hostname: config.GlobalCnf.DBHost,
		Port:     config.GlobalCnf.DBPort,
		Schema:   config.GlobalCnf.DBSchema,
		Username: config.GlobalCnf.DBUsername,
		Password: config.GlobalCnf.DBPassword,
	}

	logger.Info(fmt.Sprintf(externalDBMsg, udb.Hostname, udb.Port, udb.Schema))
	if err := udb.Open(); err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info(apiServerStartingMsg)
	logger.Info(fmt.Sprintf(apiServerStartedMsg, defaultWebServerPort))
	if err := r.Run(":8080"); err != nil {
		logger.Fatal(fmt.Sprintf(apiServerStartupErrMsg, err))
	}
}

func NewRouter() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(logger.GinZap())
	r.Use(gin.Recovery())
	r.NoRoute(httputils.NoRoute)
	r.HandleMethodNotAllowed = true
	r.NoMethod(httputils.MethodNotAllowed)

	return r
}

func SetupRoutes(r *gin.Engine) *gin.Engine {
	r.GET(apiHealthPath, httputils.Health)
	r.GET(apiUsersPath+apiUserIdParamExt, controllers.GetUser)
	r.GET(apiUsersPath+apiSearchPathExt, controllers.SearchUser)
	r.POST(apiUsersPath, controllers.CreateUser)
	r.PUT(apiUsersPath+apiUserIdParamExt, controllers.UpdateUser)
	r.DELETE(apiUsersPath+apiUserIdParamExt, controllers.DeleteUser)

	return r
}

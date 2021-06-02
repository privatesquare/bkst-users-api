package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/privatesquare/bkst-go-utils/utils/httputils"
	"github.com/privatesquare/bkst-go-utils/utils/logger"
	"github.com/privatesquare/bkst-users-api/config"
	"github.com/privatesquare/bkst-users-api/interfaces/db/mysql"
	"github.com/privatesquare/bkst-users-api/interfaces/rest"
	"github.com/privatesquare/bkst-users-api/services"
	"os"
)

const (
	defaultWebServerPort   = "8080"
	externalDBMsg          = "Using external database: %s:%s/%s"
	apiServerStartingMsg   = "Starting the API server..."
	apiServerStartedMsg    = "The API server has started and is listening on %s"
	apiServerStartupErrMsg = "Unable to run the web server"

	apiHealthPath     = "/health"
	apiUsersPath      = "/users"
	apiUserIdParamExt = "/:userId"
	apiSearchPathExt  = "/search"
)

func StartApp() {
	r := NewRouter()
	SetupRoutes(r)

	udb := &mysql.Cfg{
		Driver:   config.GlobalCnf.DBDriver,
		Hostname: config.GlobalCnf.DBHost,
		Port:     config.GlobalCnf.DBPort,
		Schema:   config.GlobalCnf.DBSchema,
		Username: config.GlobalCnf.DBUsername,
		Password: config.GlobalCnf.DBPassword,
	}

	logger.Info(fmt.Sprintf(externalDBMsg, udb.Hostname, udb.Port, udb.Schema))
	if err := udb.Open(); err != nil {
		logger.Error("", err)
		os.Exit(1)
	}

	logger.Info(apiServerStartingMsg)
	logger.Info(fmt.Sprintf(apiServerStartedMsg, defaultWebServerPort))
	if err := r.Run(":8080"); err != nil {
		logger.Error(apiServerStartupErrMsg, err)
		os.Exit(1)
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
	usersHandler := rest.NewUsersHandler(services.NewUsersService(mysql.NewUsersStore(mysql.UserDbClient)))
	r.GET(apiHealthPath, httputils.Health)
	r.GET(apiUsersPath+apiUserIdParamExt, usersHandler.Get)
	r.GET(apiUsersPath+apiSearchPathExt, usersHandler.Search)
	r.POST(apiUsersPath, usersHandler.Create)
	r.PUT(apiUsersPath+apiUserIdParamExt, usersHandler.Update)
	r.DELETE(apiUsersPath+apiUserIdParamExt, usersHandler.Delete)

	return r
}

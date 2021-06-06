package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/privatesquare/bkst-go-utils/utils/config"
	"github.com/privatesquare/bkst-go-utils/utils/httputils"
	"github.com/privatesquare/bkst-go-utils/utils/logger"
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
	apiLoginPathExt   = "/login"
)

func StartApp() {
	r := httputils.NewRouter()
	setupRoutes(r)
	dbConnect()

	logger.Info(apiServerStartingMsg)
	logger.Info(fmt.Sprintf(apiServerStartedMsg, defaultWebServerPort))
	if err := r.Run(":8080"); err != nil {
		logger.Error(apiServerStartupErrMsg, err)
		os.Exit(1)
	}
}

func dbConnect() {
	cfg := &mysql.Cfg{}
	if err := config.Load(cfg); err != nil {
		logger.Error(err.Error(), err)
		os.Exit(1)
	}

	logger.Info(fmt.Sprintf(externalDBMsg, cfg.Hostname, cfg.Port, cfg.Schema))
	if err := cfg.Open(); err != nil {
		logger.Error("", err)
		os.Exit(1)
	}
}

func setupRoutes(r *gin.Engine) *gin.Engine {
	usersHandler := rest.NewUsersHandler(services.NewUsersService(mysql.NewUsersStore(mysql.UserDbClient)))
	r.GET(apiHealthPath, httputils.Health)
	r.GET(apiUsersPath+apiUserIdParamExt, usersHandler.Get)
	r.GET(apiUsersPath+apiSearchPathExt, usersHandler.Search)
	r.POST(apiUsersPath, usersHandler.Create)
	r.PUT(apiUsersPath+apiUserIdParamExt, usersHandler.Update)
	r.DELETE(apiUsersPath+apiUserIdParamExt, usersHandler.Delete)
	r.POST(apiUsersPath+apiLoginPathExt, usersHandler.Login)

	return r
}

package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/privatesquare/bkst-go-utils/utils/errors"
	"github.com/privatesquare/bkst-go-utils/utils/httputils"
	"github.com/privatesquare/bkst-go-utils/utils/logger"
	"github.com/privatesquare/bkst-users-api/domain"
	"github.com/privatesquare/bkst-users-api/interfaces/db/mysql"
	"github.com/privatesquare/bkst-users-api/services"
	"net/http"
	"strconv"
)

const (
	userCreatedMsg    = "User with id %d was created"
	userUpdatedMsg    = "User with id %d was updated"
	userDeletedMsg    = "User with id %d was deleted"
	invalidUserIdMsg  = "invalid user id"
	invalidPayloadMsg = "invalid payload"
)

func NewUsersHandler(s services.UsersService) UsersHandler {
	return &usersHandler{Service: s}
}

type UsersHandler interface {
	Get(ctx *gin.Context)
	Search(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type usersHandler struct {
	Service services.UsersService
}

func (uh *usersHandler) Get(ctx *gin.Context) {
	userId, err := uh.parseUserId(ctx)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	usersService := services.NewUsersService(mysql.NewUsersStore(mysql.UserDbClient))
	user, restErr := usersService.Get(*userId)
	if restErr != nil {
		ctx.JSON(restErr.Status, restErr)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (uh *usersHandler) Search(ctx *gin.Context) {
	usersService := services.NewUsersService(mysql.NewUsersStore(mysql.UserDbClient))
	usersList, restErr := usersService.FindByStatus(uh.parseStatus(ctx))
	if restErr != nil {
		ctx.JSON(restErr.Status, restErr)
		return
	}
	ctx.JSON(http.StatusOK, usersList)
}

func (uh *usersHandler) Create(ctx *gin.Context) {
	u, err := uh.parseUser(ctx)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	usersService := services.NewUsersService(mysql.NewUsersStore(mysql.UserDbClient))
	user, restErr := usersService.Create(*u)
	if restErr != nil {
		ctx.JSON(restErr.Status, restErr)
		return
	}
	ctx.JSON(http.StatusCreated, httputils.RestMsg{Message: fmt.Sprintf(userCreatedMsg, user.Id)})
}

func (uh *usersHandler) Update(ctx *gin.Context) {
	userId, err := uh.parseUserId(ctx)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	user, err := uh.parseUser(ctx)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	user.Id = *userId
	usersService := services.NewUsersService(mysql.NewUsersStore(mysql.UserDbClient))
	if _, restErr := usersService.Update(*user); restErr != nil {
		ctx.JSON(restErr.Status, restErr)
		return
	}
	ctx.JSON(http.StatusOK, httputils.RestMsg{Message: fmt.Sprintf(userUpdatedMsg, user.Id)})
}

func (uh *usersHandler) Login(ctx *gin.Context) {
	usersService := services.NewUsersService(mysql.NewUsersStore(mysql.UserDbClient))

	login, err := uh.parseLogin(ctx)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	user, restErr := usersService.Login(*login)
	if restErr != nil {
		ctx.JSON(restErr.Status, restErr)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (uh *usersHandler) Delete(ctx *gin.Context) {
	userId, err := uh.parseUserId(ctx)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	usersService := services.NewUsersService(mysql.NewUsersStore(mysql.UserDbClient))
	if restErr := usersService.Delete(*userId); restErr != nil {
		ctx.JSON(restErr.Status, restErr)
		return
	}
	ctx.JSON(http.StatusOK, httputils.RestMsg{Message: fmt.Sprintf(userDeletedMsg, *userId)})
}

func (uh *usersHandler) parseUserId(ctx *gin.Context) (*int64, *errors.RestErr) {
	userId, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		logger.Info(invalidUserIdMsg)
		return nil, errors.BadRequestError(invalidUserIdMsg)
	}
	return &userId, nil
}

func (uh *usersHandler) parseStatus(ctx *gin.Context) string {
	return ctx.Query("status")
}

func (uh *usersHandler) parseUser(ctx *gin.Context) (*domain.User, *errors.RestErr) {
	user := new(domain.User)
	if err := ctx.ShouldBindJSON(user); err != nil {
		logger.Info(invalidPayloadMsg)
		return nil, errors.BadRequestError(invalidPayloadMsg)
	}
	return user, nil
}

func (uh *usersHandler) parseLogin(ctx *gin.Context) (*domain.Login, *errors.RestErr) {
	login := new(domain.Login)
	if err := ctx.ShouldBindJSON(login); err != nil {
		logger.Info(invalidPayloadMsg)
		return nil, errors.BadRequestError(invalidPayloadMsg)
	}
	return login, nil
}

package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/private-square/bkst-users-api/domain/users"
	"github.com/private-square/bkst-users-api/services"
	"github.com/private-square/bkst-users-api/utils/errors"
	"github.com/private-square/bkst-users-api/utils/httputils"
	"github.com/private-square/bkst-users-api/utils/logger"
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

func GetUser(ctx *gin.Context) {
	userId, err := parseUserId(ctx)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	user := new(users.User)
	user.Id = *userId
	user, restErr := services.UsersService.Get(user)
	if restErr != nil {
		ctx.JSON(restErr.Status, restErr)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func SearchUser(ctx *gin.Context) {
	user := users.User{Status: parseStatus(ctx)}
	usersList, restErr := services.UsersService.Find(&user)
	if restErr != nil {
		ctx.JSON(restErr.Status, restErr)
		return
	}
	ctx.JSON(http.StatusOK, usersList)
}

func CreateUser(ctx *gin.Context) {
	user, err := parseUserInfo(ctx)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	user, restErr := services.UsersService.Create(user)
	if restErr != nil {
		ctx.JSON(restErr.Status, restErr)
		return
	}
	ctx.JSON(http.StatusCreated, httputils.RestMsg{Message: fmt.Sprintf(userCreatedMsg, user.Id)})
}

func UpdateUser(ctx *gin.Context) {
	userId, err := parseUserId(ctx)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	user, err := parseUserInfo(ctx)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	user.Id = *userId
	if restErr := services.UsersService.Update(user); restErr != nil {
		ctx.JSON(restErr.Status, restErr)
		return
	}
	ctx.JSON(http.StatusOK, httputils.RestMsg{Message: fmt.Sprintf(userUpdatedMsg, user.Id)})
}

func DeleteUser(ctx *gin.Context) {
	userId, err := parseUserId(ctx)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	user := new(users.User)
	user.Id = *userId
	if restErr := services.UsersService.Delete(user); restErr != nil {
		ctx.JSON(restErr.Status, restErr)
		return
	}
	ctx.JSON(http.StatusOK, httputils.RestMsg{Message: fmt.Sprintf(userDeletedMsg, user.Id)})
}

func parseUserId(ctx *gin.Context) (*int64, *errors.RestErr) {
	userId, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		logger.Info(invalidUserIdMsg)
		return nil, errors.BadRequestError(invalidUserIdMsg)
	}
	return &userId, nil
}

func parseStatus(ctx *gin.Context) string {
	return ctx.Query("status")
}

func parseUserInfo(ctx *gin.Context) (*users.User, *errors.RestErr) {
	user := new(users.User)
	if err := ctx.ShouldBindJSON(user); err != nil {
		logger.Info(invalidPayloadMsg)
		return nil, errors.BadRequestError(invalidPayloadMsg)
	}
	return user, nil
}

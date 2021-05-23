package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/private-square/bkst-users-api/domain/users"
	"github.com/private-square/bkst-users-api/services"
	"github.com/private-square/bkst-users-api/utils"
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
	user, restErr := services.GetUser(user)
	if restErr != nil {
		ctx.JSON(restErr.Status, restErr)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func SearchUser(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, utils.RestMsg{Message: "not implemented yet"})
}

func CreateUser(ctx *gin.Context) {
	user, err := parseUserInfo(ctx)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	user, restErr := services.CreateUser(user)
	if restErr != nil {
		ctx.JSON(restErr.Status, restErr)
		return
	}
	ctx.JSON(http.StatusCreated, utils.RestMsg{Message: fmt.Sprintf(userCreatedMsg, user.Id)})
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
	if restErr := services.UpdateUser(user); restErr != nil {
		ctx.JSON(restErr.Status, restErr)
		return
	}
	ctx.JSON(http.StatusCreated, utils.RestMsg{Message: fmt.Sprintf(userUpdatedMsg, user.Id)})
}

func DeleteUser(ctx *gin.Context) {
	userId, err := parseUserId(ctx)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	user := new(users.User)
	user.Id = *userId
	if restErr := services.DeleteUser(user); restErr != nil {
		ctx.JSON(restErr.Status, restErr)
		return
	}
	ctx.JSON(http.StatusCreated, utils.RestMsg{Message: fmt.Sprintf(userDeletedMsg, user.Id)})
}

func parseUserId(ctx *gin.Context) (*int64, *utils.RestErr) {
	userId, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		return nil, utils.BadRequestError(invalidUserIdMsg)
	}
	return &userId, nil
}

func parseUserInfo(ctx *gin.Context) (*users.User, *utils.RestErr) {
	user := new(users.User)
	if err := ctx.ShouldBindJSON(user); err != nil {
		return nil, utils.BadRequestError(invalidPayloadMsg)
	}
	return user, nil
}

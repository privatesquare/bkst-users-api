package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/private-square/bkst-users-api/domain/users"
	"github.com/private-square/bkst-users-api/utils"
	"net/http"
	"strconv"
)

func GetUser(ctx *gin.Context) {
	userId, err := parseUserId(ctx)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	user := users.User{Id: *userId}
	if restErr := user.Get(); restErr != nil {
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
	if restErr := user.Create(); restErr != nil {
		ctx.JSON(restErr.Status, restErr)
		return
	}
	ctx.JSON(http.StatusCreated, utils.RestMsg{Message: fmt.Sprintf("User with id %d was created", user.Id)})
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
	if restErr := user.Update(); restErr != nil {
		ctx.JSON(restErr.Status, restErr)
		return
	}
	ctx.JSON(http.StatusCreated, utils.RestMsg{Message: fmt.Sprintf("User with id %d was updated", user.Id)})
}

func DeleteUser(ctx *gin.Context) {
	userId, err := parseUserId(ctx)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	user := users.User{Id: *userId}
	if restErr := user.Delete(); restErr != nil {
		ctx.JSON(restErr.Status, restErr)
		return
	}
	ctx.JSON(http.StatusCreated, utils.RestMsg{Message: fmt.Sprintf("User with id %d was deleted", user.Id)})
}

func parseUserId(ctx *gin.Context) (*int64, *utils.RestErr) {
	userId, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		return nil, utils.BadRequestError("invalid user id")
	}
	return &userId, nil
}

func parseUserInfo(ctx *gin.Context) (*users.User, *utils.RestErr) {
	user := new(users.User)
	if err := ctx.ShouldBindJSON(user); err != nil {
		return nil, utils.BadRequestError("invalid payload")
	}
	return user, nil
}

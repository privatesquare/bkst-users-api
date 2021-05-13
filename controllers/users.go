package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/private-square/bkst-users-api/domain/users"
	"github.com/private-square/bkst-users-api/utils"
	"net/http"
	"strconv"
)

func GetUser(ctx *gin.Context) {
	userId, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		restErr := utils.BadRequestError("invalid user id")
		ctx.JSON(restErr.Status, restErr)
		return
	}
	user := users.User{
		Id: userId,
	}
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
	user := new(users.User)
	if err := ctx.ShouldBindJSON(user); err != nil {
		restErr := utils.BadRequestError("invalid payload")
		ctx.JSON(restErr.Status, restErr)
		return
	}
	if restErr := user.Create(); restErr != nil {
		ctx.JSON(restErr.Status, restErr)
		return
	}
	//ctx.JSON(http.StatusCreated, utils.RestMsg{Message: fmt.Sprintf("User with id %d was created", user.Id)})
	ctx.JSON(http.StatusCreated, user)
}

func UpdateUser(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, utils.RestMsg{Message: "not implemented yet"})
}

func DeleteUser(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, utils.RestMsg{Message: "not implemented yet"})
}

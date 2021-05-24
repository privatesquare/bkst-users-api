package services

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/private-square/bkst-users-api/domain/users"
	"github.com/private-square/bkst-users-api/utils"
	"strings"
)

func GetUser(u *users.User) (*users.User, *utils.RestErr) {
	restErr := u.Get()
	return u, restErr
}

func FindUser(u *users.User) (*[]users.User, *utils.RestErr) {
	if err := u.ValidateStatus(); err != nil {
		return nil, utils.BadRequestError(err.Error())
	}
	usersList, restErr := u.FindByStatus()
	return &usersList, restErr
}

func CreateUser(u *users.User) (*users.User, *utils.RestErr) {
	var err error
	if err := u.Validate(); err != nil {
		return nil, utils.BadRequestError(err.Error())
	}
	u.DateCreated = utils.GetDateTimeNowFormat()
	u.DateUpdated = utils.GetDateTimeNowFormat()

	if u.Password, err = utils.EncryptPassword(u.Password, ""); err != nil {
		return nil, utils.InternalServerError(err.Error())
	}

	restErr := u.Create()
	return u, restErr
}

func UpdateUser(u *users.User) *utils.RestErr {
	updateInfo := *u

	if err := u.Get(); err != nil {
		return err
	}

	if strings.TrimSpace(updateInfo.FirstName) != "" {
		u.FirstName = updateInfo.FirstName
	}
	if strings.TrimSpace(updateInfo.Lastname) != "" {
		u.Lastname = updateInfo.Lastname
	}
	if strings.TrimSpace(updateInfo.Email) != "" {
		u.Email = updateInfo.Email
		if err := u.ValidateEmail(); err != nil {
			return utils.BadRequestError(err.Error())
		}
	}

	u.DateUpdated = utils.GetDateTimeNowFormat()
	restErr := u.Update()
	return restErr
}

func DeleteUser(u *users.User) *utils.RestErr {
	restErr := u.Delete()
	return restErr
}

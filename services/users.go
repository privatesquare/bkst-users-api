package services

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/private-square/bkst-users-api/domain/users"
	"github.com/private-square/bkst-users-api/utils"
	"strings"
)

var (
	UsersService usersServiceInterface = &usersService{}
)

type usersService struct{}

type usersServiceInterface interface {
	Get(*users.User) (*users.User, *utils.RestErr)
	Find(*users.User) (*[]users.User, *utils.RestErr)
	Create(*users.User) (*users.User, *utils.RestErr)
	Update(*users.User) *utils.RestErr
	Delete(*users.User) *utils.RestErr
}

func (s *usersService) Get(u *users.User) (*users.User, *utils.RestErr) {
	restErr := u.Get()
	return u, restErr
}

func (s *usersService) Find(u *users.User) (*[]users.User, *utils.RestErr) {
	if err := u.ValidateStatus(); err != nil {
		return nil, utils.BadRequestError(err.Error())
	}
	usersList, restErr := u.FindByStatus()
	return &usersList, restErr
}

func (s *usersService) Create(u *users.User) (*users.User, *utils.RestErr) {
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

func (s *usersService) Update(u *users.User) *utils.RestErr {
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

func (s *usersService) Delete(u *users.User) *utils.RestErr {
	restErr := u.Delete()
	return restErr
}

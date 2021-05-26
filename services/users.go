package services

import (
	"github.com/private-square/bkst-users-api/domain/users"
	"github.com/private-square/bkst-users-api/utils/dateutils"
	"github.com/private-square/bkst-users-api/utils/errors"
	"github.com/private-square/bkst-users-api/utils/logger"
	"github.com/private-square/bkst-users-api/utils/secrets"
	"strings"
)

var (
	UsersService usersServiceInterface = &usersService{}
)

type usersService struct{}

type usersServiceInterface interface {
	Get(*users.User) (*users.User, *errors.RestErr)
	Find(*users.User) (*[]users.User, *errors.RestErr)
	Create(*users.User) (*users.User, *errors.RestErr)
	Update(*users.User) *errors.RestErr
	Delete(*users.User) *errors.RestErr
}

func (s *usersService) Get(u *users.User) (*users.User, *errors.RestErr) {
	restErr := u.Get()
	return u, restErr
}

func (s *usersService) Find(u *users.User) (*[]users.User, *errors.RestErr) {
	if err := u.ValidateStatus(); err != nil {
		logger.Info(err.Error())
		return nil, errors.BadRequestError(err.Error())
	}
	usersList, restErr := u.FindByStatus()
	return &usersList, restErr
}

func (s *usersService) Create(u *users.User) (*users.User, *errors.RestErr) {
	var err error
	if err := u.Validate(); err != nil {
		logger.Info(err.Error())
		return nil, errors.BadRequestError(err.Error())
	}
	u.DateCreated = dateutils.GetDateTimeNowFormat()
	u.DateUpdated = dateutils.GetDateTimeNowFormat()

	if u.Password, err = secrets.EncryptPassword(u.Password, ""); err != nil {
		logger.Error(err.Error(), nil)
		return nil, errors.InternalServerError(users.InternalServerErrMsg)
	}

	restErr := u.Create()
	return u, restErr
}

func (s *usersService) Update(u *users.User) *errors.RestErr {
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
			logger.Info(err.Error())
			return errors.BadRequestError(err.Error())
		}
	}

	u.DateUpdated = dateutils.GetDateTimeNowFormat()
	restErr := u.Update()
	return restErr
}

func (s *usersService) Delete(u *users.User) *errors.RestErr {
	restErr := u.Delete()
	return restErr
}

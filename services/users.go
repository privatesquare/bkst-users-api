package services

import (
	"github.com/privatesquare/bkst-go-utils/utils/dateutils"
	"github.com/privatesquare/bkst-go-utils/utils/errors"
	"github.com/privatesquare/bkst-go-utils/utils/logger"
	"github.com/privatesquare/bkst-go-utils/utils/secrets"
	"github.com/privatesquare/bkst-users-api/domain"
	"github.com/privatesquare/bkst-users-api/interfaces/db/mysql"
	"strings"
)

func NewUsersService(UserStore mysql.UsersStore) UsersService {
	return &usersService{
		UserStore: UserStore,
	}
}

type UsersService interface {
	Get(id int64) (*domain.User, *errors.RestErr)
	FindByStatus(status string) ([]domain.User, *errors.RestErr)
	Create(u domain.User) (*domain.User, *errors.RestErr)
	Update(u domain.User) (*domain.User, *errors.RestErr)
	Delete(id int64) *errors.RestErr
}

type usersService struct {
	UserStore mysql.UsersStore
}

func (s *usersService) Get(id int64) (*domain.User, *errors.RestErr) {
	user, restErr := s.UserStore.Get(id)
	return user, restErr
}

func (s *usersService) FindByStatus(status string) ([]domain.User, *errors.RestErr) {
	u := domain.User{Status: status}
	if err := u.ValidateStatus(); err != nil {
		logger.Info(err.Error())
		return nil, errors.BadRequestError(err.Error())
	}
	usersList, restErr := s.UserStore.FindByStatus(u.Status)
	return usersList, restErr
}

func (s *usersService) Create(u domain.User) (*domain.User, *errors.RestErr) {
	var err error
	u.Status = domain.ActiveStatus
	if err := u.Validate(); err != nil {
		logger.Info(err.Error())
		return nil, errors.BadRequestError(err.Error())
	}
	u.DateCreated = dateutils.GetDateTimeNowFormat()
	u.DateUpdated = dateutils.GetDateTimeNowFormat()

	if u.Password, err = secrets.EncryptPassword(u.Password, ""); err != nil {
		logger.Error(err.Error(), nil)
		return nil, errors.InternalServerError()
	}

	user, restErr := s.UserStore.Create(u)
	return user, restErr
}

func (s *usersService) Update(u domain.User) (*domain.User, *errors.RestErr) {
	updateInfo := u

	user, restErr := s.UserStore.Get(u.Id)
	if restErr != nil {
		return nil, restErr
	}

	if strings.TrimSpace(updateInfo.FirstName) != "" {
		user.FirstName = updateInfo.FirstName
	}
	if strings.TrimSpace(updateInfo.Lastname) != "" {
		user.Lastname = updateInfo.Lastname
	}
	if strings.TrimSpace(updateInfo.Email) != "" {
		user.Email = updateInfo.Email
		if err := user.ValidateEmail(); err != nil {
			logger.Info(err.Error())
			return nil, errors.BadRequestError(err.Error())
		}
	}

	user.DateUpdated = dateutils.GetDateTimeNowFormat()
	user, restErr = s.UserStore.Update(*user)
	return user, restErr
}

func (s *usersService) Delete(id int64) *errors.RestErr {
	restErr := s.UserStore.Delete(id)
	return restErr
}

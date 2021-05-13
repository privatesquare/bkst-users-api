package users

import (
	"fmt"
	"github.com/private-square/bkst-users-api/utils"
	"regexp"
	"strings"
)

const (
	emailRe            = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	invalidEmailErrMsg = "The email id is not valid."
)

var (
	usersDb = make(map[int64]*User)
)

type InvalidEmailError struct{}

func (e InvalidEmailError) Error() string {
	return invalidEmailErrMsg
}

type User struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"firstName"`
	Lastname    string `json:"lastName"`
	Email       string `json:"email"`
	DateCreated string `json:"dateCreated"`
}

func (u *User) Get() *utils.RestErr {
	result := usersDb[u.Id]
	if result == nil {
		return utils.NotFoundError(fmt.Sprintf("User %d was not found", u.Id))
	}

	u.Id = result.Id
	u.FirstName = result.FirstName
	u.Lastname = result.Lastname
	u.Email = result.Email
	u.DateCreated = result.DateCreated

	return nil
}

func (u *User) Create() *utils.RestErr {
	if err := u.validate(); err != nil {
		return utils.BadRequestError(err.Error())
	}
	current := usersDb[u.Id]
	if current != nil {
		return utils.BadRequestError(fmt.Sprintf("User %d already exists", u.Id))
	}
	for _, user := range usersDb {
		if user.Email == u.Email {
			return utils.BadRequestError(fmt.Sprintf("Account already exists for the email id %s", u.Email))
		}
	}
	u.DateCreated = utils.GetDateTimeNowFormat()
	usersDb[u.Id] = u
	return nil
}

func (u *User) validate() error {
	if err := u.validateNotEmpty(); err != nil {
		return err
	}
	if err := u.validateEmail(); err != nil {
		return err
	}
	return nil
}

func (u *User) validateNotEmpty() error {
	var missingParams []string

	if strings.TrimSpace(u.FirstName) == "" {
		missingParams = append(missingParams, utils.GetFieldTagValue(u, &u.FirstName))
	}
	if strings.TrimSpace(u.Lastname) == "" {
		missingParams = append(missingParams, utils.GetFieldTagValue(u, &u.Lastname))
	}
	if strings.TrimSpace(u.Email) == "" {
		missingParams = append(missingParams, utils.GetFieldTagValue(u, &u.Email))
	}

	if len(missingParams) > 0 {
		return utils.MissingMandatoryParamError(missingParams)
	}

	return nil
}

func (u *User) validateEmail() error {
	pattern := regexp.MustCompile(emailRe)
	if !pattern.MatchString(u.Email) {
		return InvalidEmailError{}
	}
	u.Email = strings.ToLower(u.Email)
	return nil
}

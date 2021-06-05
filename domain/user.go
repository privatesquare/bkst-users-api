package domain

import (
	"fmt"
	"github.com/privatesquare/bkst-go-utils/utils/errors"
	"github.com/privatesquare/bkst-go-utils/utils/secrets"
	"github.com/privatesquare/bkst-go-utils/utils/slice"
	"github.com/privatesquare/bkst-go-utils/utils/structutils"
	"regexp"
	"strings"
)

type User struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"firstName"`
	Lastname    string `json:"lastName"`
	Email       string `json:"email"`
	Status      string `json:"status"`
	Password    string `json:"password,omitempty"`
	DateCreated string `json:"dateCreated"`
	DateUpdated string `json:"dateUpdated"`
}

const (
	ActiveStatus = "active"
	InactiveStatus = "inactive"
)

var (
	validStatusList = []string{ActiveStatus, InactiveStatus}
)

const (
	emailRegex          = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	invalidEmailErrMsg  = "The email id is not valid."
	invalidStatusErrMsg = "Invalid status '%s'. Valid Status's: %v"
)

type InvalidEmailError struct{}

func (e InvalidEmailError) Error() string {
	return invalidEmailErrMsg
}

type InvalidStatusError struct {
	invalidStatus   string
	validStatusList []string
}

func (e InvalidStatusError) Error() string {
	return fmt.Sprintf(invalidStatusErrMsg, e.invalidStatus, e.validStatusList)
}

func (u *User) Validate() error {
	if err := u.validateNotEmpty(); err != nil {
		return err
	}
	if err := u.ValidateEmail(); err != nil {
		return err
	}
	if err := u.ValidatePassword(); err != nil {
		return err
	}
	return nil
}

func (u *User) validateNotEmpty() error {
	var missingParams []string

	if strings.TrimSpace(u.FirstName) == "" {
		missingParams = append(missingParams, structutils.GetFieldTagValue(u, &u.FirstName))
	}
	if strings.TrimSpace(u.Lastname) == "" {
		missingParams = append(missingParams, structutils.GetFieldTagValue(u, &u.Lastname))
	}
	if strings.TrimSpace(u.Email) == "" {
		missingParams = append(missingParams, structutils.GetFieldTagValue(u, &u.Email))
	}
	if strings.TrimSpace(u.Status) == "" {
		missingParams = append(missingParams, structutils.GetFieldTagValue(u, &u.Status))
	}
	if strings.TrimSpace(u.Password) == "" {
		missingParams = append(missingParams, structutils.GetFieldTagValue(u, &u.Password))
	}

	if len(missingParams) > 0 {
		return errors.MissingMandatoryParamError(missingParams)
	}

	return nil
}

func (u *User) ValidateEmail() error {
	pattern := regexp.MustCompile(emailRegex)
	if !pattern.MatchString(u.Email) {
		return InvalidEmailError{}
	}
	u.Email = strings.ToLower(u.Email)
	return nil
}

func (u *User) ValidateStatus() error {
	if !slice.EntryExists(validStatusList, u.Status) {
		return InvalidStatusError{
			invalidStatus:   u.Status,
			validStatusList: validStatusList,
		}
	}
	return nil
}

func (u *User) ValidatePassword() error {
	return secrets.VerifyPassword(u.Password)
}

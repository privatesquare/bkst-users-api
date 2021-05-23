package users

import (
	"database/sql"
	"fmt"
	"github.com/private-square/bkst-users-api/services"
	"github.com/private-square/bkst-users-api/utils"
	"regexp"
	"strings"
)

const (
	emailRe            = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	invalidEmailErrMsg = "The email id is not valid."

	querySelectUserById = "SELECT id, first_name, last_name, email, date_created, date_updated FROM users WHERE id=?;"
	queryInsertUser     = "INSERT INTO users(first_name, last_name, email) VALUES(?, ?, ?);"
	queryUpdateUser     = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser		= "DELETE FROM users WHERE id=?;"
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
	DateUpdated string `json:"dateUpdated"`
}

func (u *User) Get() *utils.RestErr {
	stmt, err := services.UsersDbClient.Prepare(querySelectUserById)
	if err != nil {
		return utils.InternalServerError(err.Error())
	}
	defer stmt.Close()

	err = stmt.QueryRow(u.Id).Scan(&u.Id, &u.FirstName, &u.Lastname, &u.Email, &u.DateCreated, &u.DateUpdated)
	switch {
	case err == nil:
		return nil
	case err == sql.ErrNoRows:
		return utils.NotFoundError(fmt.Sprintf("User with id %d was not found", u.Id))
	default:
		return utils.InternalServerError(err.Error())
	}
}

func (u *User) Create() *utils.RestErr {
	if err := u.validate(); err != nil {
		return utils.BadRequestError(err.Error())
	}

	stmt, err := services.UsersDbClient.Prepare(queryInsertUser)
	if err != nil {
		return utils.InternalServerError(err.Error())
	}
	defer stmt.Close()

	result, err := stmt.Exec(u.FirstName, u.Lastname, u.Email)
	if err := u.handleQueryExecError(err); err != nil {
		return err
	}

	if u.Id, err = result.LastInsertId(); err != nil {
		return utils.InternalServerError(err.Error())
	}
	return nil
}

func (u *User) Update() *utils.RestErr {

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
		if err := u.validateEmail(); err != nil {
			return utils.BadRequestError(err.Error())
		}
	}

	stmt, err := services.UsersDbClient.Prepare(queryUpdateUser)
	if err != nil {
		return utils.InternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.FirstName, u.Lastname, u.Email, u.Id)
	if err := u.handleQueryExecError(err); err != nil {
		return err
	}

	return nil
}

func (u *User) Delete() *utils.RestErr {
	stmt, err := services.UsersDbClient.Prepare(queryDeleteUser)
	if err != nil {
		return utils.InternalServerError(err.Error())
	}
	defer stmt.Close()

	result, err := stmt.Exec(u.Id)
	if err := u.handleQueryExecError(err); err != nil {
		return err
	}

	if rowsAffected, err := result.RowsAffected(); err != nil {
		return utils.InternalServerError(err.Error())
	} else if rowsAffected == 0 {
		return utils.NotFoundError(fmt.Sprintf("User with id %d was not found", u.Id))
	} else {
		return nil
	}
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

func (u *User) handleQueryExecError(err error) *utils.RestErr {
	if err == nil {
		return nil
	} else if strings.Contains(err.Error(), "email_UNIQUE") {
		return utils.BadRequestError(fmt.Sprintf("Email id %s is already in use.", u.Email))
	} else {
		return utils.InternalServerError(err.Error())
	}
}

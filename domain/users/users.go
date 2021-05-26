package users

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/private-square/bkst-users-api/utils/errors"
	"github.com/private-square/bkst-users-api/utils/logger"
	"github.com/private-square/bkst-users-api/utils/secrets"
	"github.com/private-square/bkst-users-api/utils/slice"
	"github.com/private-square/bkst-users-api/utils/structutils"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	emailRegex           = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	invalidEmailErrMsg   = "The email id is not valid."
	invalidStatusErrMsg  = "Invalid status '%s'. Valid Status's: %v"
	InternalServerErrMsg = "Unable to process the request due to an internal error. Please contact the systems administrator"

	usersDbDataSourceNameFormat = "%s:%s@tcp(%s:%s)/%s?charset=utf8"
	usersDbConnErrMsg           = "Users Db connection error : %v"
	usersDbConnSuccessMsg       = "Successfully connected to the Users database"
	usersDbEmailUniqueStr       = "email_UNIQUE"
	usersDbPrepareStmtErrMsg    = "Error when trying to prepare a statement"
	usersDbQueryErrMsg          = "Error when trying to run a query on the database"
	usersDbQueryRowsErrMsg      = "Error when trying to get query rows"
	usersDbScanRowsErrMsg       = "Error when trying to scan rows"
	usersDbExecErrMsg           = "Error when trying to execute a statement on the database"
	usersDbLastInsertIdErrMsg   = "Error when trying to get the last insert id"
	usersDbRowsAffectedErrMsg   = "Error when trying to get the number of affected rows"

	userNotFoundMsg     = "User with id %d was not found"
	emailAlreadyUsedMsg = "Email id %s is already in use."

	querySelectUserById     = "SELECT id, first_name, last_name, email, status, date_created, date_updated FROM users WHERE id=?;"
	querySelectUserByStatus = "SELECT id, first_name, last_name, email, status, date_created, date_updated FROM users WHERE status=?;"
	queryInsertUser         = "INSERT INTO users(first_name, last_name, email, status, password, date_created, date_updated) VALUES(?, ?, ?, ?, ?, ?, ?);"
	queryUpdateUser         = "UPDATE users SET first_name=?, last_name=?, email=?, date_updated=? WHERE id=?;"
	queryDeleteUser         = "DELETE FROM users WHERE id=?;"
)

var (
	UserDbClient    *sql.DB
	validStatusList = []string{"active", "inactive"}
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

type UserDbConn struct {
	Driver   string
	Hostname string
	Port     string
	Schema   string
	Username string
	Password string
}

func init() {
	mysqlLogger := log.New(ioutil.Discard, "", 0)
	err := mysql.SetLogger(mysqlLogger)
	if err != nil {
		logger.Error(err.Error(), nil)
		os.Exit(1)
	}
}

func (db *UserDbConn) Open() error {
	var err error
	dataSourceName := fmt.Sprintf(usersDbDataSourceNameFormat,
		db.Username,
		db.Password,
		db.Hostname,
		db.Port,
		db.Schema)
	if UserDbClient, err = sql.Open(db.Driver, dataSourceName); err != nil {
		err = errors.NewError(fmt.Sprintf(usersDbConnErrMsg, err))
		logger.Error(err.Error(), nil)
		return err
	}
	if err := UserDbClient.Ping(); err != nil {
		err = errors.NewError(fmt.Sprintf(usersDbConnErrMsg, err))
		logger.Error(err.Error(), nil)
		return err
	}
	logger.Info(usersDbConnSuccessMsg)
	return nil
}

func (db *UserDbConn) Close() {
	_ = UserDbClient.Close()
}

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

func (u *User) Get() *errors.RestErr {
	stmt, err := UserDbClient.Prepare(querySelectUserById)
	if err != nil {
		logger.Error(usersDbPrepareStmtErrMsg, err)
		return errors.InternalServerError(InternalServerErrMsg)
	}
	defer stmt.Close()

	err = stmt.QueryRow(u.Id).Scan(&u.Id, &u.FirstName, &u.Lastname, &u.Email, &u.Status, &u.DateCreated, &u.DateUpdated)
	switch {
	case err == nil:
		return nil
	case err == sql.ErrNoRows:
		msg := fmt.Sprintf(userNotFoundMsg, u.Id)
		logger.Info(msg)
		return errors.NotFoundError(msg)
	default:
		logger.Error(usersDbQueryRowsErrMsg, err)
		return errors.InternalServerError(err.Error())
	}
}

func (u *User) FindByStatus() ([]User, *errors.RestErr) {
	stmt, err := UserDbClient.Prepare(querySelectUserByStatus)
	if err != nil {
		logger.Error(usersDbPrepareStmtErrMsg, err)
		return nil, errors.InternalServerError(InternalServerErrMsg)
	}
	defer stmt.Close()

	rows, err := stmt.Query(u.Status)
	if err != nil {
		logger.Error(usersDbQueryErrMsg, err)
		return nil, errors.InternalServerError(InternalServerErrMsg)
	}
	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.Lastname, &user.Email, &user.Status, &user.DateCreated, &user.DateUpdated); err != nil {
			logger.Error(usersDbScanRowsErrMsg, err)
			return nil, errors.InternalServerError(InternalServerErrMsg)
		}
		users = append(users, user)
	}
	defer rows.Close()

	if len(users) > 0 {
		return users, nil
	} else {
		return []User{}, nil
	}
}

func (u *User) Create() *errors.RestErr {
	stmt, err := UserDbClient.Prepare(queryInsertUser)
	if err != nil {
		logger.Error(usersDbPrepareStmtErrMsg, err)
		return errors.InternalServerError(InternalServerErrMsg)
	}
	defer stmt.Close()

	result, err := stmt.Exec(u.FirstName, u.Lastname, u.Email, u.Status, u.Password, u.DateCreated, u.DateUpdated)
	if err := u.handleQueryExecError(err); err != nil {
		return err
	}

	if u.Id, err = result.LastInsertId(); err != nil {
		logger.Error(usersDbLastInsertIdErrMsg, err)
		return errors.InternalServerError(InternalServerErrMsg)
	}
	return nil
}

func (u *User) Update() *errors.RestErr {
	stmt, err := UserDbClient.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error(usersDbPrepareStmtErrMsg, err)
		return errors.InternalServerError(InternalServerErrMsg)
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.FirstName, u.Lastname, u.Email, u.DateUpdated, u.Id)
	if err := u.handleQueryExecError(err); err != nil {
		return err
	}

	return nil
}

func (u *User) Delete() *errors.RestErr {
	stmt, err := UserDbClient.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error(usersDbPrepareStmtErrMsg, err)
		return errors.InternalServerError(InternalServerErrMsg)
	}
	defer stmt.Close()

	result, err := stmt.Exec(u.Id)
	if err := u.handleQueryExecError(err); err != nil {
		return err
	}

	if rowsAffected, err := result.RowsAffected(); err != nil {
		logger.Error(usersDbRowsAffectedErrMsg, err)
		return errors.InternalServerError(InternalServerErrMsg)
	} else if rowsAffected == 0 {
		msg := fmt.Sprintf(userNotFoundMsg, u.Id)
		logger.Info(msg)
		return errors.NotFoundError(msg)
	} else {
		return nil
	}
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

func (u *User) handleQueryExecError(err error) *errors.RestErr {
	if err == nil {
		return nil
	} else if strings.Contains(err.Error(), usersDbEmailUniqueStr) {
		msg := fmt.Sprintf(emailAlreadyUsedMsg, u.Email)
		logger.Info(msg)
		return errors.BadRequestError(msg)
	} else {
		logger.Error(usersDbExecErrMsg, err)
		return errors.InternalServerError(InternalServerErrMsg)
	}
}

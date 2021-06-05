package mysql

import (
	"database/sql"
	"fmt"
	"github.com/privatesquare/bkst-go-utils/utils/errors"
	"github.com/privatesquare/bkst-go-utils/utils/logger"
	"github.com/privatesquare/bkst-users-api/domain"
	"strings"
)

const (
	usersDbEmailUniqueStr     = "email_UNIQUE"
	usersDbPrepareStmtErrMsg  = "Error when trying to prepare a statement"
	usersDbQueryErrMsg        = "Error when trying to run a query on the database"
	usersDbQueryRowsErrMsg    = "Error when trying to get query rows"
	usersDbScanRowsErrMsg     = "Error when trying to scan rows"
	usersDbExecErrMsg         = "Error when trying to execute a statement on the database"
	usersDbLastInsertIdErrMsg = "Error when trying to get the last insert id"
	usersDbRowsAffectedErrMsg = "Error when trying to get the number of affected rows"

	userNotFoundMsg     = "User with id %d was not found"
	emailAlreadyUsedMsg = "Email id %s is already in use."

	querySelectUserById     = "SELECT id, first_name, last_name, email, status, date_created, date_updated FROM users WHERE id=?;"
	querySelectUserByStatus = "SELECT id, first_name, last_name, email, status, date_created, date_updated FROM users WHERE status=?;"
	queryInsertUser         = "INSERT INTO users(first_name, last_name, email, status, password, date_created, date_updated) VALUES(?, ?, ?, ?, ?, ?, ?);"
	queryUpdateUser         = "UPDATE users SET first_name=?, last_name=?, email=?, date_updated=? WHERE id=?;"
	queryDeleteUser         = "DELETE FROM users WHERE id=?;"
)

func NewUsersStore(db *sql.DB) UsersStore {
	return &userStore{
		db: db,
	}
}

type UsersStore interface {
	Get(id int64) (*domain.User, *errors.RestErr)
	FindByStatus(status string) ([]domain.User, *errors.RestErr)
	Create(u domain.User) (*domain.User, *errors.RestErr)
	Update(u domain.User) (*domain.User, *errors.RestErr)
	Delete(id int64) *errors.RestErr
}

type userStore struct {
	db *sql.DB
}

func (us *userStore) Get(id int64) (*domain.User, *errors.RestErr) {
	u := new(domain.User)
	stmt, err := UserDbClient.Prepare(querySelectUserById)
	if err != nil {
		logger.Error(usersDbPrepareStmtErrMsg, err)
		return nil, errors.InternalServerError()
	}
	defer stmt.Close()

	err = stmt.QueryRow(&id).Scan(&u.Id, &u.FirstName, &u.Lastname, &u.Email, &u.Status, &u.DateCreated, &u.DateUpdated)
	switch {
	case err == nil:
		return u, nil
	case err == sql.ErrNoRows:
		msg := fmt.Sprintf(userNotFoundMsg, id)
		logger.Info(msg)
		return nil, errors.NotFoundError(msg)
	default:
		logger.Error(usersDbQueryRowsErrMsg, err)
		return nil, errors.InternalServerError()
	}
}

func (us *userStore) FindByStatus(status string) ([]domain.User, *errors.RestErr) {
	stmt, err := UserDbClient.Prepare(querySelectUserByStatus)
	if err != nil {
		logger.Error(usersDbPrepareStmtErrMsg, err)
		return nil, errors.InternalServerError()
	}
	defer stmt.Close()

	rows, err := stmt.Query(&status)
	if err != nil {
		logger.Error(usersDbQueryErrMsg, err)
		return nil, errors.InternalServerError()
	}
	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.Lastname, &user.Email, &user.Status, &user.DateCreated, &user.DateUpdated); err != nil {
			logger.Error(usersDbScanRowsErrMsg, err)
			return nil, errors.InternalServerError()
		}
		users = append(users, user)
	}
	defer rows.Close()

	if len(users) > 0 {
		return users, nil
	} else {
		return []domain.User{}, nil
	}
}

func (us *userStore) Create(u domain.User) (*domain.User, *errors.RestErr) {
	stmt, err := UserDbClient.Prepare(queryInsertUser)
	if err != nil {
		logger.Error(usersDbPrepareStmtErrMsg, err)
		return nil, errors.InternalServerError()
	}
	defer stmt.Close()

	result, err := stmt.Exec(&u.FirstName, &u.Lastname, &u.Email, &u.Status, &u.Password, &u.DateCreated, &u.DateUpdated)
	if err := us.handleQueryExecError(u, err); err != nil {
		return nil, err
	}

	if u.Id, err = result.LastInsertId(); err != nil {
		logger.Error(usersDbLastInsertIdErrMsg, err)
		return nil, errors.InternalServerError()
	}
	return &u, nil
}

func (us *userStore) Update(u domain.User) (*domain.User, *errors.RestErr) {
	stmt, err := UserDbClient.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error(usersDbPrepareStmtErrMsg, err)
		return nil, errors.InternalServerError()
	}
	defer stmt.Close()

	_, err = stmt.Exec(&u.FirstName, &u.Lastname, &u.Email, &u.DateUpdated, &u.Id)
	if err := us.handleQueryExecError(u, err); err != nil {
		return nil, err
	}
	return &u, nil
}

func (us *userStore) Delete(id int64) *errors.RestErr {
	u := domain.User{Id: id}
	stmt, err := UserDbClient.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error(usersDbPrepareStmtErrMsg, err)
		return errors.InternalServerError()
	}
	defer stmt.Close()

	result, err := stmt.Exec(&id)
	if err := us.handleQueryExecError(u, err); err != nil {
		return err
	}

	if rowsAffected, err := result.RowsAffected(); err != nil {
		logger.Error(usersDbRowsAffectedErrMsg, err)
		return errors.InternalServerError()
	} else if rowsAffected == 0 {
		msg := fmt.Sprintf(userNotFoundMsg, id)
		logger.Info(msg)
		return errors.NotFoundError(msg)
	} else {
		return nil
	}
}

func (us *userStore) handleQueryExecError(u domain.User, err error) *errors.RestErr {
	if err == nil {
		return nil
	} else if strings.Contains(err.Error(), usersDbEmailUniqueStr) {
		msg := fmt.Sprintf(emailAlreadyUsedMsg, u.Email)
		logger.Info(msg)
		return errors.BadRequestError(msg)
	} else {
		logger.Error(usersDbExecErrMsg, err)
		return errors.InternalServerError()
	}
}

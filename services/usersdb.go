package services

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

const (
	usersDbDriveName           = "mysql"
	userDbDataSourceNameFormat = "%s:%s@tcp(%s:%s)/%s?charset=utf8"
)

var (
	UsersDbClient *sql.DB
)

type UsersDbConn struct {
	Hostname string
	Port     string
	Schema   string
	Username string
	Password string
}

func (db *UsersDbConn) Open() error {
	var err error
	dataSourceName := fmt.Sprintf(userDbDataSourceNameFormat,
		db.Username,
		db.Password,
		db.Hostname,
		db.Port,
		db.Schema)
	if UsersDbClient, err = sql.Open(usersDbDriveName, dataSourceName); err != nil {
		return errors.New(fmt.Sprintf("Users Db connection error : %v", err))
	}
	if err := UsersDbClient.Ping(); err != nil {
		return errors.New(fmt.Sprintf("Users Db connection error : %v", err))
	}
	fmt.Println("Successfully connected to the Users database.")
	return nil
}

func (db *UsersDbConn) Close() {
	_ = UsersDbClient.Close()
}

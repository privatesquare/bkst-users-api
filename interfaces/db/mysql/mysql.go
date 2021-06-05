package mysql

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/privatesquare/bkst-go-utils/utils/errors"
	"github.com/privatesquare/bkst-go-utils/utils/logger"
	"io/ioutil"
	"log"
	"os"
)

var (
	UserDbClient *sql.DB
)

const (
	usersDbDataSourceNameFormat = "%s:%s@tcp(%s:%s)/%s?charset=utf8"
	usersDbConnErrMsg           = "Users Db connection error : %v"
	usersDbConnSuccessMsg       = "Successfully connected to the Users database"
	internalServerErrMsg        = "Unable to process the request due to an internal error. Please contact the systems administrator"
)

type Cfg struct {
	Driver   string `mapstructure:"DB_DRIVER"`
	Hostname string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	Schema   string `mapstructure:"DB_SCHEMA"`
	Username string `mapstructure:"DB_USERNAME"`
	Password string `mapstructure:"DB_PASSWORD"`
}

func init() {
	mysqlLogger := log.New(ioutil.Discard, "", 0)
	err := mysql.SetLogger(mysqlLogger)
	if err != nil {
		logger.Error(err.Error(), nil)
		os.Exit(1)
	}
}

func (db *Cfg) Open() error {
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

func (db *Cfg) Close() {
	_ = UserDbClient.Close()
}

package config

import (
	"fmt"
	"github.com/private-square/bkst-users-api/utils/errors"
	"github.com/private-square/bkst-users-api/utils/logger"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strings"
)

const (
	defaultConfigPath    = "."
	defaultConfigName    = "config"
	defaultConfigType    = "env"
	configLoadSuccessMsg = "Configuration loaded successfully"
	configLoadErrMsg     = "Unable to load the required configuration"
)

var (
	GlobalCnf = &Config{}
)

type Config struct {
	DBDriver   string `mapstructure:"DB_DRIVER"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBSchema   string `mapstructure:"DB_SCHEMA"`
	DBUsername string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
}

func init() {
	GlobalCnf.Load()
	GlobalCnf.Validate()
}

func (c *Config) Load() {
	viper.AddConfigPath(defaultConfigPath)
	viper.SetConfigName(defaultConfigName)
	viper.SetConfigType(defaultConfigType)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		logger.Error(configLoadErrMsg, err)
		os.Exit(1)
	}

	err = viper.Unmarshal(GlobalCnf)
	if err != nil {
		logger.Error(configLoadErrMsg, err)
		os.Exit(1)
	}

	logger.Info(configLoadSuccessMsg)
}

func (c *Config) Validate() {
	var missingParams []string
	sr := reflect.ValueOf(c).Elem()

	for i := 0; i < sr.NumField(); i++ {
		if strings.TrimSpace(sr.Field(i).String()) == "" {
			missingParams = append(missingParams, sr.Type().Field(i).Tag.Get("mapstructure"))
			fmt.Println(os.Getenv(sr.Type().Field(i).Tag.Get("mapstructure")))
		}
	}
	if len(missingParams) > 0 {
		err := errors.MissingMandatoryParamError(missingParams)
		logger.Error(err.Error(), nil)
		os.Exit(1)
	}
}

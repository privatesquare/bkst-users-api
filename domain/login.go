package domain

import (
	"github.com/privatesquare/bkst-go-utils/utils/errors"
	"github.com/privatesquare/bkst-go-utils/utils/structutils"
	"strings"
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (l *Login) Validate() error {
	var missingParams []string

	l.Username = strings.TrimSpace(l.Username)
	if l.Username == "" {
		missingParams = append(missingParams, structutils.GetFieldTagValue(l, &l.Username))
	}
	l.Password = strings.TrimSpace(l.Password)
	if l.Password == "" {
		missingParams = append(missingParams, structutils.GetFieldTagValue(l, &l.Password))
	}

	if len(missingParams) > 0 {
		return errors.MissingMandatoryParamError(missingParams)
	}
	return nil
}

package errors

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	missingMandatoryParamErrMsg = "Missing mandatory parameter(s) : %v"
	invalidPasswordErrMsg       = "password should be at least 8 characters long with at least one number, one uppercase letter, one lowercase letter and one special character"
	passwordEncryptionErrMsg    = "password encryption error: %v"
	passwordDecryptionErrMsg    = "password decryption error: %v"
)

var (
	InvalidPasswordError = errors.New(invalidPasswordErrMsg)
)

type RestErr struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func NewError(msg string) error {
	return errors.New(msg)
}

func BadRequestError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusBadRequest,
		Error:   http.StatusText(http.StatusBadRequest),
	}
}

func UnauthorizedError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusUnauthorized,
		Error:   http.StatusText(http.StatusUnauthorized),
	}
}

func ForbiddenError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusForbidden,
		Error:   http.StatusText(http.StatusForbidden),
	}
}

func NotFoundError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusNotFound,
		Error:   http.StatusText(http.StatusNotFound),
	}
}

func InternalServerError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusInternalServerError,
		Error:   http.StatusText(http.StatusInternalServerError),
	}
}

// MissingMandatoryParamError represents an error when a mandatory parameter is missing
type MissingMandatoryParamError []string

// Error returns the formatted MissingMandatoryParamError
func (e MissingMandatoryParamError) Error() string {
	return fmt.Sprintf(missingMandatoryParamErrMsg, []string(e))
}

type PasswordEncryptionError struct {
	Err error
}

func (e PasswordEncryptionError) Error() string {
	return fmt.Sprintf(passwordEncryptionErrMsg, e.Err)
}

type PasswordDecryptionError struct {
	Err error
}

func (e PasswordDecryptionError) Error() string {
	return fmt.Sprintf(passwordDecryptionErrMsg, e.Err)
}

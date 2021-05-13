package utils

import (
	"fmt"
	"net/http"
)

const (
	missingMandatoryParamErrMsg = "Missing mandatory parameter(s) : %v"
)

type RestErr struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func NewError(message string, status int, error string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  status,
		Error:   error,
	}
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

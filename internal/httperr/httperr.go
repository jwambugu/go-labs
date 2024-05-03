package httperr

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

type Error struct {
	Message    string
	Errors     []map[string]string
	StatusCode int
}

func (e Error) Code() int {
	return e.StatusCode
}

func (e Error) Error() string {
	return e.Message
}

func ValidationError(err error) Error {
	if err == nil {
		return Error{}
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return Error{}
	}

	errs := make([]map[string]string, len(validationErrors))

	for i, err := range validationErrors {
		field := strings.ToLower(err.Field())

		errs[i] = map[string]string{
			field: err.Error(),
		}
	}

	return Error{
		Errors:     errs,
		StatusCode: http.StatusUnprocessableEntity,
	}
}

func New(code int, msg string) Error {
	return Error{
		StatusCode: code,
		Message:    msg,
	}
}

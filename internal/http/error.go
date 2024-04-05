package httpapi

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

type Error struct {
	err        error
	statusCode int
}

type errorResponse struct {
	Message string              `json:"message,omitempty"`
	Errors  []map[string]string `json:"errors,omitempty"`
}

func (e *Error) code() int {
	return e.statusCode
}

func (e *Error) toJson() *errorResponse {
	if e.err == nil {
		return nil
	}

	resp := &errorResponse{
		Message: e.err.Error(),
	}

	if e.statusCode == 0 {
		e.statusCode = http.StatusInternalServerError
		resp.Message = http.StatusText(e.statusCode)
	}

	var validationErrors validator.ValidationErrors
	if ok := errors.As(e.err, &validationErrors); ok {
		errs := make([]map[string]string, len(validationErrors))
		for i, err := range validationErrors {
			field := strings.ToLower(err.Field())

			errs[i] = map[string]string{
				field: err.Error(),
			}
		}

		e.statusCode = http.StatusBadRequest
		resp.Errors = errs
		resp.Message = "The given data was invalid"
	}

	return resp
}

func (e *Error) Send(w http.ResponseWriter) {
	encode(w, e.statusCode, e.toJson())
}

func NewError(err error, statusCode int) *Error {
	return &Error{
		err:        err,
		statusCode: statusCode,
	}
}

func encode(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		panic(err)
	}
}

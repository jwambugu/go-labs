package api

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

type errorResponse struct {
	Message string              `json:"message,omitempty"`
	Errors  []map[string]string `json:"errors,omitempty"`
}

type httpError struct {
	err    error
	status int
}

func (h *httpError) statusCode() int {
	return h.status
}

func (h *httpError) toJson() *errorResponse {
	if h.err == nil {
		return nil
	}

	resp := &errorResponse{
		Message: h.err.Error(),
	}

	var validationErrors validator.ValidationErrors
	if errors.As(h.err, &validationErrors) {
		errs := make([]map[string]string, len(validationErrors))

		for i, err := range validationErrors {
			field := strings.ToLower(err.Field())

			errs[i] = map[string]string{
				field: err.Error(),
			}
		}

		h.status = http.StatusUnprocessableEntity
		resp.Errors = errs
		resp.Message = "The given data was invalid"
	}

	return resp
}

func (a *Api) errorResponse(w http.ResponseWriter, status int, err error) {
	var (
		httpErr = httpError{
			err:    err,
			status: status,
		}

		payload = httpErr.toJson()
	)

	encode(w, httpErr.statusCode(), payload)
}

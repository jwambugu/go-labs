package api

import (
	"errors"
	"go-labs/internal/httperr"
	"net/http"
)

type errorResponse struct {
	Message string              `json:"message,omitempty"`
	Errors  []map[string]string `json:"errors,omitempty"`
}

func (a *Api) JSONError(w http.ResponseWriter, err error) {
	var httpErr httperr.Error
	if !errors.As(err, &httpErr) {
		httpErr = httperr.ErrServerError
	}

	var payload any

	payload = httpErr.Message
	if httpErr.StatusCode == http.StatusUnprocessableEntity {
		payload = httpErr.Errors
	}

	encode(w, httpErr.StatusCode, payload)
}

package httperr

import "net/http"

var (
	ErrServerError        = New(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	ErrInvalidCredentials = New(http.StatusUnauthorized, "Invalid credentials provided.")
	ErrDuplicateEmail     = New(http.StatusUnprocessableEntity, "Email is already in use.")
)

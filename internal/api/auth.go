package api

import (
	"encoding/json"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email,omitempty" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required"`
}

func (r *LoginRequest) Validate() error {
	return validate.Struct(r)
}

func (a *Api) loginHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest *LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		a.errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	if err := loginRequest.Validate(); err != nil {
		a.errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	a.successResponse(w, http.StatusOK, nil)
}

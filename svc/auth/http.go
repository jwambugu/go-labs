package auth

import (
	"github.com/go-playground/validator/v10"
	"go-labs/internal/model"
)

// use a single instance of Validate, it caches struct info

var validate = validator.New(validator.WithRequiredStructEnabled())

type LoginRequest struct {
	Email    string `json:"email,omitempty" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required"`
}

type validationError struct {
	err any
}

func (e validationError) Error() string {
	return ""
}

func (r *LoginRequest) Validate() error {
	return validate.Struct(r)
}

type LoginResponse struct {
	User        *model.User `json:"user,omitempty"`
	AccessToken string      `json:"access_token,omitempty"`
}

type RegisterRequest struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (r *RegisterRequest) Validate() error {
	return validate.Struct(r)
}

type RegisterResponse struct {
	User        *model.User `json:"user,omitempty"`
	AccessToken string      `json:"access_token,omitempty"`
}

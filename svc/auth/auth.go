package auth

import (
	"context"
	"go-labs/internal/repository"
)

type Authenticator interface {
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
}

type authSvc struct {
	rs repository.Store
}

func (a *authSvc) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	return nil, nil
}

func (a *authSvc) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	return nil, nil
}

func NewAuthSvc(rs repository.Store) Authenticator {
	return &authSvc{
		rs: rs,
	}
}

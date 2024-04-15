package auth

import (
	"context"
	"errors"
	"go-labs/internal/model"
	"go-labs/internal/repository"
	"go-labs/internal/util"
	"go.uber.org/zap"
	"strings"
	"time"
)

type Authenticator interface {
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
}

var ErrInvalidHash = errors.New("auth: invalid hash")

type authSvc struct {
	logger     *zap.Logger
	rs         *repository.Store
	jwtManager JWTManager
}

func (a *authSvc) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	logger := a.logger.Named("Login(_)").With(zap.String("email", req.Email))

	user, err := a.rs.User.FindByEmail(ctx, req.Email)
	if err != nil {
		logger.Error("find by email", zap.Error(err))
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, repository.ErrRecordNotFound
		}

		return nil, err
	}

	if err = util.CompareHash(string(user.Password), req.Password); err != nil {
		return nil, ErrInvalidHash
	}

	user.Password = nil

	accessToken, err := a.jwtManager.Generate(user, JWTTokenExpiresAt)
	if err != nil {
		logger.Error("failed to generate access token", zap.Error(err), zap.Uint64("user", user.ID))
		return nil, err
	}

	resp := &LoginResponse{
		User:        user,
		AccessToken: accessToken,
	}

	logger.Info("user authenticated successfully", zap.Uint64("id", user.ID))
	return resp, nil
}

func (a *authSvc) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	logger := a.logger.Named("Register(_)").With(
		zap.String("name", req.Name),
		zap.String("email", req.Email),
	)

	if err := a.rs.User.EmailExists(ctx, req.Email); err != nil {
		logger.Error("failed to check id email exists", zap.Error(err))

		if errors.Is(err, repository.ErrRecordExists) {
			return nil, err
		}
		return nil, err
	}

	password, err := util.HashString(req.Password)
	if err != nil {
		logger.Error("failed to hash password", zap.Error(err))
		return nil, err
	}

	var (
		now  = util.TimePtr(time.Now())
		user = &model.User{
			CreatedAt: now,
			UpdatedAt: now,
			Name:      util.StrTitle(req.Name),
			Email:     strings.ToLower(req.Email),
			Password:  []byte(password),
		}
	)

	if err = a.rs.User.Create(ctx, user); err != nil {
		logger.Error("failed create new user", zap.Error(err))
		return nil, err
	}

	accessToken, err := a.jwtManager.Generate(user, JWTTokenExpiresAt)
	if err != nil {
		logger.Error("failed to generate access token", zap.Error(err), zap.Uint64("user", user.ID))
		return nil, err
	}

	resp := &RegisterResponse{
		User:        user,
		AccessToken: accessToken,
	}

	return resp, nil
}

func NewAuthSvc(logger *zap.Logger, rs *repository.Store, jwt JWTManager) Authenticator {
	logger = logger.Named("auth-svc")

	return &authSvc{
		logger:     logger,
		rs:         rs,
		jwtManager: jwt,
	}
}

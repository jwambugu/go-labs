package repository

import (
	"context"
	"go-labs/internal/model"
)

type User interface {
	Create(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id uint64) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
}

type userRepo struct {
}

func (u *userRepo) Create(ctx context.Context, user *model.User) error {
	return nil
}

func (u *userRepo) Delete(ctx context.Context, user *model.User) error {
	return nil
}

func (u *userRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepo) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepo) Update(ctx context.Context, user *model.User) error {
	return nil
}

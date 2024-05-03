package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go-labs/internal/model"
	"go-labs/internal/util"
	"time"
)

type User interface {
	Create(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint64) error
	EmailExists(ctx context.Context, email string) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id uint64) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
}

type userRepo struct {
	db *sqlx.DB
}

const (
	_userCreateQuery = `INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`
	_userSelectQuery      = `SELECT id, name, email, password, created_at, updated_at, deleted_at FROM users`
	_userFindByEmailQuery = _userSelectQuery + ` WHERE email = ? AND deleted_at IS NULL;`
	_userFindByIDQuery    = _userSelectQuery + ` WHERE id = ? AND deleted_at IS NULL LIMIT 1;`
	_userUpdateQuery      = `UPDATE users
		SET name       = ?,
			email      = ?,
			password   = ?,
			updated_at = ?,
			deleted_at = ?
		WHERE id = ?;`
	_userEmailExists = `SELECT EXISTS(SELECT 1 FROM users WHERE email = ? AND deleted_at IS NULL)`
)

func (u *userRepo) Create(ctx context.Context, user *model.User) error {
	defer func() {
		if user != nil {
			user.Password = nil
		}
	}()

	user.CreatedAt = util.TimePtr(time.Now())
	user.UpdatedAt = util.TimePtr(time.Now())

	result, err := u.db.ExecContext(
		ctx,
		_userCreateQuery,
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("insert: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get id: %v", err)
	}

	user.ID = uint64(id)
	return nil
}

func (u *userRepo) Delete(ctx context.Context, id uint64) error {
	user, err := u.FindByID(ctx, id)
	if err != nil {
		return err
	}

	user.DeletedAt = util.TimePtr(time.Now())

	return u.Update(ctx, user)
}

func (u *userRepo) EmailExists(ctx context.Context, email string) error {
	var exists bool

	if err := u.db.GetContext(ctx, &exists, _userEmailExists, email); err != nil {
		return err
	}

	if exists {
		return ErrRecordExists
	}

	return nil
}

func (u *userRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	if err := u.db.GetContext(ctx, &user, _userFindByEmailQuery, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("get: %v", err)
	}

	return &user, nil
}

func (u *userRepo) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User

	if err := u.db.GetContext(ctx, &user, _userFindByIDQuery, id); err != nil {
		return nil, fmt.Errorf("get: %v", err)
	}

	return &user, nil
}

func (u *userRepo) Update(ctx context.Context, user *model.User) error {
	user.UpdatedAt = util.TimePtr(time.Now())

	_, err := u.db.ExecContext(
		ctx, _userUpdateQuery, user.Name, user.Email, user.Password, user.UpdatedAt, user.DeletedAt, user.ID,
	)
	if err != nil {
		return fmt.Errorf("update: %v", err)
	}
	return nil
}

func NewUserRepo(db *sqlx.DB) User {
	return &userRepo{db: db}
}

func NewUserTestRepo(ctx context.Context, db *sqlx.DB, users ...*model.User) (User, error) {
	repo := NewUserRepo(db)

	for _, user := range users {
		if err := repo.Create(ctx, user); err != nil {
			return nil, err
		}
	}

	return repo, nil
}

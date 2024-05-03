package repository

import "errors"

type Store struct {
	User User
}

var (
	ErrRecordExists   = errors.New("repository: record exists")
	ErrRecordNotFound = errors.New("repository: record not found")
)

func NewStore() *Store {
	return &Store{}
}

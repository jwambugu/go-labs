package model

import "time"

type User struct {
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	Name      string     `json:"name,omitempty" db:"name"`
	Email     string     `json:"email,omitempty" db:"email"`
	Password  []byte     `json:"password,omitempty" db:"password"`
	ID        uint64     `json:"id,omitempty" db:"id"`
}

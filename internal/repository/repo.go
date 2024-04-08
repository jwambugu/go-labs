package repository

type Store struct {
	User User
}

func NewStore() *Store {
	return &Store{}
}

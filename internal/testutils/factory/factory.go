package factory

import (
	"github.com/jaswdr/faker"
	"go-labs/internal/model"
	"go-labs/internal/util"
	"time"
)

func NewUser() *model.User {
	f := faker.New()

	return &model.User{
		CreatedAt: util.TimePtr(time.Now()),
		UpdatedAt: util.TimePtr(time.Now()),
		Name:      f.Person().Name(),
		Email:     f.Numerify("########") + "." + f.Internet().Email(),
		Password:  []byte("secret"),
	}
}

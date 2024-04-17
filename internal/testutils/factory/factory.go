package factory

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"github.com/jaswdr/faker"
	"go-labs/internal/model"
	"go-labs/internal/util"
	"math/rand"
	"time"
)

const UserPassword = "password"

func Seed() int64 {
	var b [8]byte
	_, err := cryptorand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	return int64(binary.LittleEndian.Uint64(b[:]))
}

func NewUser() *model.User {
	f := faker.NewWithSeed(rand.NewSource(Seed()))

	return &model.User{
		CreatedAt: util.TimePtr(time.Now()),
		UpdatedAt: util.TimePtr(time.Now()),
		Name:      util.StrTitle(f.Person().Name()),
		Email:     f.Numerify("########") + "." + f.Internet().Email(),
		Password:  []byte("$2a$10$rsO27uABCBzkMlv5UasLvux80qeuUbl/0QQxPvGnVMJqCFkoiK5eG"),
	}
}

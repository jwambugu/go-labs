package util

import (
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"time"
)

func TimePtr(t time.Time) *time.Time {
	return &t
}

func PtrToTime(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Time{}
}

func CompareHash(hash string, s string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))
}

func HashString(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func StrTitle(s string) string {
	return cases.Title(language.English, cases.Compact).String(s)
}

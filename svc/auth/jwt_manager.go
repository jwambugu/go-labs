package auth

import (
	"errors"
	"fmt"
	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
	"go-labs/internal/model"
	"strconv"
	"time"
)

var (
	ErrExpiredToken = errors.New("auth:token has expired")
	ErrInvalidToken = errors.New("auth:token is invalid")
)

const JWTTokenExpiresAt = 12 * time.Hour

type Payload struct {
	Expiration time.Time
	IssuedAt   time.Time
	NotBefore  time.Time
	User       *model.User
	Audience   string
	Issuer     string
	Subject    string
}

type pasetoToken struct {
	paseto *paseto.V2
	key    []byte
}

type JWTManager interface {
	Generate(user *model.User, expiresAfter time.Duration) (string, error)
	Verify(token string) (*Payload, error)
}

func (p *pasetoToken) Generate(user *model.User, expiresAfter time.Duration) (string, error) {
	var (
		now     = time.Now()
		payload = &Payload{
			Expiration: now.Add(expiresAfter),
			IssuedAt:   now,
			NotBefore:  now,
			User:       user,
			Audience:   "Todo",
			Issuer:     "Todo",
			Subject:    strconv.FormatUint(user.ID, 10),
		}
	)

	return p.paseto.Encrypt(p.key, payload, nil)
}

func (p *pasetoToken) Verify(token string) (*Payload, error) {
	var payload Payload

	if err := p.paseto.Decrypt(token, p.key, &payload, ""); err != nil {
		return nil, ErrInvalidToken
	}

	if time.Now().After(payload.Expiration) {
		return nil, ErrExpiredToken
	}

	return &payload, nil
}

func NewPasetoToken(key string) (JWTManager, error) {
	if size := chacha20poly1305.KeySize; len(key) != size {
		return nil, fmt.Errorf("auth: invalid key size must be exactly %d characters", size)
	}

	return &pasetoToken{
		paseto: paseto.NewV2(),
		key:    []byte(key),
	}, nil
}

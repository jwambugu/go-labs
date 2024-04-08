package api

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"go-labs/internal/repository"
	"net/http"
)

// validate using a single instance of Validate, it caches struct info
var validate = validator.New(validator.WithRequiredStructEnabled())

type Api struct {
	rs *repository.Store
}

func (a *Api) Serve(port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a.Routes(),
	}
}

func NewApi(rs *repository.Store) *Api {
	return &Api{rs: rs}
}

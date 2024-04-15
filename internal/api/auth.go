package api

import (
	"encoding/json"
	"go-labs/svc/auth"
	"net/http"
)

func (a *Api) loginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq *auth.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		a.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	if err := loginReq.Validate(); err != nil {
		a.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	loginResp, err := a.authSVC.Login(r.Context(), loginReq)
	if err != nil {
		a.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	a.JSON(w, http.StatusOK, loginResp)
}

func (a *Api) registerHandler(w http.ResponseWriter, r *http.Request) {
	var registerReq auth.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
		a.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	if err := registerReq.Validate(); err != nil {
		a.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	registerResp, err := a.authSVC.Register(r.Context(), &registerReq)
	if err != nil {
		a.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	a.JSON(w, http.StatusOK, registerResp)
}

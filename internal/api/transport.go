package api

import (
	"encoding/json"
	"go-labs/internal/model"
	"log"
	"net/http"
)

type successResponse struct {
	User        *model.User `json:"user,omitempty"`
	AccessToken string      `json:"access_token,omitempty"`
}

func encode(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Panicf("encode: %v", err)
	}
}

func (a *Api) JSON(w http.ResponseWriter, status int, payload any) {
	encode(w, status, payload)
}

package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type successResponse struct {
	Data any `json:"data,omitempty"`
}

func encode(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Panicf("encode: %v", err)
	}
}

func (a *Api) JSON(w http.ResponseWriter, status int, data any) {
	resp := &successResponse{
		Data: data,
	}

	encode(w, status, resp)
}

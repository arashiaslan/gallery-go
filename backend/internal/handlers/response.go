package handlers

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

func WriteJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

// LimitBody caps a JSON request body to keep memory use bounded.
func LimitBody(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
}

func WriteError(w http.ResponseWriter, code int, message string) {
	WriteJSON(w, code, ErrorResponse{Error: message})
}

func WriteSuccess(w http.ResponseWriter, code int, data interface{}) {
	WriteJSON(w, code, SuccessResponse{Data: data})
}

package httpcommon

import (
	"encoding/json"
	"net/http"
)

type userIDKey struct{}

type ErrorResponse struct {
	Error ErrorContent `json:"error"`
}

type ErrorContent struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteError(w http.ResponseWriter, status int, code string, msg string) {
	WriteJSON(w, status, ErrorResponse{Error: ErrorContent{Code: code, Message: msg}})
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

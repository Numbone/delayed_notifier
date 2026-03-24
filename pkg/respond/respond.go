package respond

import (
	"encoding/json"
	"net/http"
)

// ErrorBody — формат JSON ошибки: {"error": "текст"}.
type ErrorBody struct {
	Error string `json:"error"`
}

// JSON отправляет JSON ответ с нужным HTTP статусом.
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// Error отправляет JSON ошибку: {"error": "msg"}.
func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, ErrorBody{Error: msg})
}

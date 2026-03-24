package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/delayed_notifier/internal/service"
	"github.com/delayed_notifier/pkg/respond"
)

// UserHandler — обработчик HTTP запросов для пользователей.
type UserHandler struct {
	svc service.UserService
}

// NewUserHandler создаёт обработчик.
func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// RegisterRoutes регистрирует маршруты пользователя на переданном mux.
//
// POST /users        — создать пользователя
// GET  /users        — список всех пользователей
// GET  /users/{id}   — получить пользователя по ID
func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /users", h.CreateUser)
	mux.HandleFunc("GET /users", h.ListUsers)
	mux.HandleFunc("GET /users/{id}", h.GetUser)
}

// createUserRequest — тело запроса для создания пользователя.
type createUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// CreateUser обрабатывает POST /users.
// Принимает JSON: {"email": "...", "name": "..."}
// Возвращает созданного пользователя.
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "неверный формат JSON")
		return
	}

	user, err := h.svc.CreateUser(r.Context(), req.Email, req.Name)
	if err != nil {
		if errors.Is(err, service.ErrEmptyEmail) || errors.Is(err, service.ErrEmptyName) {
			respond.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		log.Printf("ошибка при создании пользователя: %v", err)
		respond.Error(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	respond.JSON(w, http.StatusCreated, user)
}

// ListUsers обрабатывает GET /users.
// Возвращает массив всех пользователей.
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.ListUsers(r.Context())
	if err != nil {
		log.Printf("ошибка при получении списка пользователей: %v", err)
		respond.Error(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	respond.JSON(w, http.StatusOK, users)
}

// GetUser обрабатывает GET /users/{id}.
// Возвращает пользователя по ID из URL.
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "id должен быть числом")
		return
	}

	user, err := h.svc.GetUser(r.Context(), int32(id))
	if err != nil {
		log.Printf("ошибка при получении пользователя id=%d: %v", id, err)
		respond.Error(w, http.StatusNotFound, "пользователь не найден")
		return
	}

	respond.JSON(w, http.StatusOK, user)
}

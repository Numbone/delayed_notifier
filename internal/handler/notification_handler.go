package handler

import (
	"encoding/json"
	"net/http"

	"github.com/delayed_notifier/internal/db"
	"github.com/delayed_notifier/internal/service"
	"github.com/delayed_notifier/pkg/respond"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type NotificationHandler struct {
	svc service.NotificationService
}

func NewNotificationHandler(svc service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		svc: svc,
	}
}

func (h *NotificationHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /notify", h.CreateNotification)
	mux.HandleFunc("GET /notify/{id}", h.GetNotificationById)
	mux.HandleFunc("DELETE /notify/{id}", h.DeleteNotificationById)
}

type createNotificationRequest struct {
	Channel   string             `json:"channel"`
	Recipient string             `json:"recipient"`
	Subject   string             `json:"subject"`
	Body      string             `json:"body"`
	Status    string             `json:"status"`
	SendAt    pgtype.Timestamptz `json:"send_at"`
}

func (h *NotificationHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	var req createNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	params := db.CreateNotificationParams{
		Channel:   req.Channel,
		Recipient: req.Recipient,
		Subject:   req.Subject,
		Body:      req.Body,
		Status:    req.Status,
		SendAt:    req.SendAt,
	}
	notification, err := h.svc.CreateNotification(r.Context(), params)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusCreated, notification)
}

func (h *NotificationHandler) GetNotificationById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "Invalid notification ID")
		return
	}
	notification, err := h.svc.GetNotification(r.Context(), id)
	if err != nil {
		respond.Error(w, http.StatusNotFound, "Notification not found")
		return
	}
	respond.JSON(w, http.StatusOK, notification)
}

func (h *NotificationHandler) DeleteNotificationById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "Invalid notification ID")
		return
	}
	if err := h.svc.DeleteNotificationById(r.Context(), id); err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

package service

import (
	"context"

	"github.com/delayed_notifier/internal/db"
	"github.com/delayed_notifier/internal/repository"
	"github.com/google/uuid"
)

type NotificationService interface {
	CreateNotification(ctx context.Context, params db.CreateNotificationParams) (db.Notification, error)
	GetNotification(ctx context.Context, id uuid.UUID) (db.Notification, error)
	DeleteNotificationById(ctx context.Context, id uuid.UUID) error
}

type notificationService struct {
	repo repository.NotificationRepository
}

func (n *notificationService) CreateNotification(ctx context.Context, params db.CreateNotificationParams) (db.Notification, error) {
	return n.repo.Create(ctx, params)
}

func (n *notificationService) GetNotification(ctx context.Context, id uuid.UUID) (db.Notification, error) {
	return n.repo.GetByID(ctx, id)
}

func (n *notificationService) DeleteNotificationById(ctx context.Context, id uuid.UUID) error {
	return n.repo.DeleteByID(ctx, id)
}

func NewNotificationService(repo repository.NotificationRepository) NotificationService {
	return &notificationService{repo: repo}
}

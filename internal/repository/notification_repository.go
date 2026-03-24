package repository

import (
	"context"
	"github.com/delayed_notifier/internal/db"
	"github.com/google/uuid"
)

type NotificationRepository interface {
	Create(ctx context.Context, params db.CreateNotificationParams) (db.Notification, error)
	GetByID(ctx context.Context, id uuid.UUID) (db.Notification, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type notificationRepo struct {
	q *db.Queries
}

func (n *notificationRepo) Create(ctx context.Context, params db.CreateNotificationParams) (db.Notification, error) {
	return n.q.CreateNotification(ctx, params)
}

func (n *notificationRepo) GetByID(ctx context.Context, id uuid.UUID) (db.Notification, error) {
	return n.q.GetNotificationById(ctx, id)
}

func (n *notificationRepo) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return n.q.DeleteNotification(ctx, id)
}

func NewNotificationRepository(q *db.Queries) NotificationRepository {
	return &notificationRepo{q: q}
}

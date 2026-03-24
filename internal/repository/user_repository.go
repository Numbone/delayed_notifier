package repository

import (
	"context"

	"github.com/delayed_notifier/internal/db"
)

// UserRepository — интерфейс для работы с пользователями в базе данных.
// Нужен чтобы service слой не зависел напрямую от sqlc.
type UserRepository interface {
	Create(ctx context.Context, email, name string) (db.User, error)
	GetByID(ctx context.Context, id int32) (db.User, error)
	List(ctx context.Context) ([]db.User, error)
}

// userRepo — реализация UserRepository, внутри использует sqlc Queries.
type userRepo struct {
	q *db.Queries
}

// NewUserRepository создаёт новый репозиторий пользователей.
func NewUserRepository(q *db.Queries) UserRepository {
	return &userRepo{q: q}
}

func (r *userRepo) Create(ctx context.Context, email, name string) (db.User, error) {
	return r.q.CreateUser(ctx, db.CreateUserParams{
		Email: email,
		Name:  name,
	})
}

func (r *userRepo) GetByID(ctx context.Context, id int32) (db.User, error) {
	return r.q.GetUserByID(ctx, id)
}

func (r *userRepo) List(ctx context.Context) ([]db.User, error) {
	return r.q.ListUsers(ctx)
}

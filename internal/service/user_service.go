package service

import (
	"context"
	"errors"
	"strings"

	"github.com/delayed_notifier/internal/db"
	"github.com/delayed_notifier/internal/repository"
)

// Ошибки валидации — возвращаются клиенту с кодом 400.
var (
	ErrEmptyEmail = errors.New("email не может быть пустым")
	ErrEmptyName  = errors.New("имя не может быть пустым")
)

// UserService — интерфейс бизнес-логики пользователей.
type UserService interface {
	CreateUser(ctx context.Context, email, name string) (db.User, error)
	GetUser(ctx context.Context, id int32) (db.User, error)
	ListUsers(ctx context.Context) ([]db.User, error)
}

// userService — реализация, использует репозиторий.
type userService struct {
	repo repository.UserRepository
}

// NewUserService создаёт сервис пользователей.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, email, name string) (db.User, error) {
	// Убираем пробелы по краям
	email = strings.TrimSpace(email)
	name = strings.TrimSpace(name)

	// Валидация
	if email == "" {
		return db.User{}, ErrEmptyEmail
	}
	if name == "" {
		return db.User{}, ErrEmptyName
	}

	return s.repo.Create(ctx, email, name)
}

func (s *userService) GetUser(ctx context.Context, id int32) (db.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context) ([]db.User, error) {
	return s.repo.List(ctx)
}

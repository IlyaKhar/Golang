package user

import (
	"context"
	"day2/internal/structure/domain"
)

type Service interface {
	CreateUser(ctx context.Context, u *domain.User) error
	GetUser(ctx context.Context, id int) (domain.User, error)
	ListUsers(ctx context.Context, limit, offset int) ([]domain.User, error)
	UpdateUser(ctx context.Context, u *domain.User) error
	DeleteUser(ctx context.Context, id int) error
}

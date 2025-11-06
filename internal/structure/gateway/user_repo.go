package gateway

import (
    "context"
    "day2/internal/structure/domain"
)

type UserRepo interface {
    Create(ctx context.Context, u *domain.User) error
    GetByID(ctx context.Context, id int) (domain.User, error)
    GetByEmail(ctx context.Context, email string) (domain.User, error)
    List(ctx context.Context, limit, offset int) ([]domain.User, error)
    Update(ctx context.Context, u *domain.User) error
    Delete(ctx context.Context, id int) error
}



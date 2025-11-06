package user

import (
	"context"
	"day2/internal/structure/domain"
	"day2/internal/structure/gateway"
	"strings"
)

type serviceImpl struct {
	repo gateway.UserRepo
}

func NewUserService(r gateway.UserRepo) Service { return &serviceImpl{repo: r} }

func (s *serviceImpl) CreateUser(ctx context.Context, u *domain.User) error {
	if u == nil || u.Name == "" || !strings.Contains(u.Email, "@") || u.Age < 0 || u.Age > 150 {
		return domain.ErrInvalidInput
	}
	return s.repo.Create(ctx, u)
}

func (s *serviceImpl) GetUser(ctx context.Context, id int) (domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *serviceImpl) ListUsers(ctx context.Context, limit, offset int) ([]domain.User, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *serviceImpl) UpdateUser(ctx context.Context, u *domain.User) error {
	if u.ID <= 0 || u.Name == "" || !strings.Contains(u.Email, "@") || u.Age < 0 || u.Age > 150 {
		return domain.ErrInvalidInput
	}
	return s.repo.Update(ctx, u)
}

func (s *serviceImpl) DeleteUser(ctx context.Context, id int) error {
	if id <= 0 {
		return domain.ErrInvalidInput
	}
	return s.repo.Delete(ctx, id)
}

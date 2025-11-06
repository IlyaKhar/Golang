package main

import (
	"context"
	"log"
	"net"

	"day2/grpc/gen"
	grpcsrv "day2/internal/transport/grpc/server"

	"google.golang.org/grpc"

	"day2/internal/structure/domain"
	useuser "day2/internal/structure/usecase/user"
)

// Внимание: httpdemo.newMemUserService не экспортируется в текущем виде.
// Для быстрого старта можешь заменить на свою реализацию Service
// или сделать отдельный конструктор ин-мемори сервиса в отдельном пакете.

type memSvc struct{ users map[int]domain.User }

func newMemSvc() useuser.Service { return &memSvc{users: make(map[int]domain.User)} }

func (m *memSvc) CreateUser(ctx context.Context, u *domain.User) error {
	if u == nil || u.Name == "" || u.Email == "" || u.Age < 0 || u.Age > 150 {
		return domain.ErrInvalidInput
	}
	if _, ok := m.users[u.ID]; ok {
		return domain.ErrEmailExists
	}
	m.users[u.ID] = *u
	return nil
}

func (m *memSvc) GetUser(ctx context.Context, id int) (domain.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return domain.User{}, domain.ErrUserNotFound
}

func (m *memSvc) ListUsers(ctx context.Context, limit, offset int) ([]domain.User, error) {
	out := make([]domain.User, 0, len(m.users))
	for _, u := range m.users {
		out = append(out, u)
	}
	return out, nil
}

func (m *memSvc) UpdateUser(ctx context.Context, u *domain.User) error {
	if u == nil || u.ID <= 0 {
		return domain.ErrInvalidInput
	}
	if _, ok := m.users[u.ID]; !ok {
		return domain.ErrUserNotFound
	}
	m.users[u.ID] = *u
	return nil
}

func (m *memSvc) DeleteUser(ctx context.Context, id int) error {
	if _, ok := m.users[id]; !ok {
		return domain.ErrUserNotFound
	}
	delete(m.users, id)
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	svc := newMemSvc()
	gen.RegisterUserServiceServer(s, grpcsrv.NewUserServer(svc))

	log.Println("gRPC listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

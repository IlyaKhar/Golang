package server

import (
	"context"

	"day2/grpc/gen"
	dom "day2/internal/structure/domain"
	useuser "day2/internal/structure/usecase/user"
)

type UserServer struct {
	gen.UnimplementedUserServiceServer
	svc useuser.Service
}

func NewUserServer(svc useuser.Service) *UserServer { return &UserServer{svc: svc} }

func (s *UserServer) GetUser(ctx context.Context, req *gen.GetUserRequest) (*gen.GetUserResponse, error) {
	u, err := s.svc.GetUser(ctx, int(req.GetId()))
	if err != nil {
		return nil, err
	}
	return &gen.GetUserResponse{User: toProto(u)}, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *gen.CreateUserRequest) (*gen.CreateUserResponse, error) {
	du := fromProto(req.GetUser())
	if err := s.svc.CreateUser(ctx, &du); err != nil {
		return nil, err
	}
	return &gen.CreateUserResponse{User: toProto(du)}, nil
}

func toProto(u dom.User) *gen.User {
	return &gen.User{Id: int32(u.ID), Name: u.Name, Email: u.Email, Age: int32(u.Age)}
}

func fromProto(u *gen.User) dom.User {
	if u == nil {
		return dom.User{}
	}
	return dom.User{ID: int(u.GetId()), Name: u.GetName(), Email: u.GetEmail(), Age: int(u.GetAge())}
}

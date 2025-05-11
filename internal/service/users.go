package service

import (
	"context"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userService struct {
	repository repository.Db
}

type UserService interface {
	CreateUser(ctx context.Context, dto *model.CreateUserDTO) error
	Login(ctx context.Context, user *model.User) (string, error)
}

func NewUserService(repository repository.Db) (UserService, error) {
	return &userService{
		repository: repository,
	}, nil
}
func (s *userService) CreateUser(ctx context.Context, dto *model.CreateUserDTO) error {
	u := model.NewUser(dto)
	err := s.repository.CreateUser(ctx, u)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	return nil
}

func (s *userService) Login(ctx context.Context, user *model.User) (string, error) {
	token, err := s.repository.Login(ctx, user)
	if err != nil {
		return "", status.Errorf(codes.Internal, err.Error())
	}
	return token, nil
}

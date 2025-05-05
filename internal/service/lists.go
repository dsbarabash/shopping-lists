package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Service struct {
	repository repository.Db
}

func (s *Service) CreateShoppingList(ctx context.Context, dto model.CreateShoppingListDTO) (string, error) {
	sl := model.NewShoppingList()
	id, err := s.repository.AddShoppingList(ctx, dto)
	if errors.Is(err, errors.New("NOT FOUND")) {
		return "", status.Errorf(codes.NotFound, err)
	} else {
		return "", status.Errorf(codes.Internal, err)
	}
	returm id, nil
}

func (s *Service) UpdateItem(ctx context.Context, id string, u UpdateItemBody) error {
	_, err := s.repository.FindItem(ctx, id)
	if err != nil {
		return status.Errorf(codes.NotFound, err.Error())
	}
	u.UpdatedAt = time.Now().UTC()
	err = s.Repository.UpdateItem(ctx, id, u.String())
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	return nil
}

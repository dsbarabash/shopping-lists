package service

import (
	"context"
	"errors"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	repository repository.Db
}

type Service interface {
	CreateShoppingList(ctx context.Context, dto *model.CreateShoppingListDTO) error
	GetShoppingListById(ctx context.Context, id string) (*model.ShoppingList, error)
	GetShoppingLists(ctx context.Context) ([]*model.ShoppingList, error)
	UpdateShoppingList(ctx context.Context, id string, dto *model.UpdateShoppingListDTO) error
	DeleteShoppingListById(ctx context.Context, id string) error
	CreateItem(ctx context.Context, dto *model.CreateItemDTO) error
	GetItemById(ctx context.Context, id string) (*model.Item, error)
	GetItems(ctx context.Context) ([]*model.Item, error)
	UpdateItem(ctx context.Context, id string, dto *model.UpdateItemDTO) error
	DeleteItemById(ctx context.Context, id string) error
}

func NewService(repository repository.Db) (Service, error) {
	return &service{
		repository: repository,
	}, nil
}

func (s *service) CreateShoppingList(ctx context.Context, dto *model.CreateShoppingListDTO) error {
	sl := model.NewShoppingList(dto)
	err := s.repository.AddShoppingList(ctx, sl)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return status.Errorf(codes.NotFound, err.Error())
		} else {
			return status.Errorf(codes.Internal, err.Error())
		}
	}
	return nil
}

func (s *service) GetShoppingListById(ctx context.Context, id string) (*model.ShoppingList, error) {
	sl, err := s.repository.GetSlById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return sl, status.Errorf(codes.NotFound, err.Error())
		} else {
			return sl, status.Errorf(codes.Internal, err.Error())
		}
	}
	return sl, nil
}

func (s *service) GetShoppingLists(ctx context.Context) ([]*model.ShoppingList, error) {
	sl, err := s.repository.GetSls(ctx)
	if err != nil {
		if errors.Is(err, errors.New("NOT FOUND")) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		} else {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}
	return sl, nil
}

func (s *service) UpdateShoppingList(ctx context.Context, id string, dto *model.UpdateShoppingListDTO) error {
	if dto.Title == "" && len(dto.Items) == 0 && dto.UserId == "" {
		return status.Errorf(codes.InvalidArgument, "nothing to update")
	}
	_, err := s.repository.GetSlById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return status.Errorf(codes.NotFound, err.Error())
		} else {
			return status.Errorf(codes.Internal, err.Error())
		}
	}
	sl := model.UpdateShoppingList(dto)
	err = s.repository.UpdateSl(ctx, id, sl)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return status.Errorf(codes.NotFound, err.Error())
		} else {
			return status.Errorf(codes.Internal, err.Error())
		}
	}
	return nil
}

func (s *service) DeleteShoppingListById(ctx context.Context, id string) error {
	_, err := s.repository.GetSlById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return status.Errorf(codes.NotFound, err.Error())
		} else {
			return status.Errorf(codes.Internal, err.Error())
		}
	}
	err = s.repository.DeleteSlById(ctx, id)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	return nil
}

func (s *service) CreateItem(ctx context.Context, dto *model.CreateItemDTO) error {
	i := model.NewItem(dto)
	err := s.repository.AddItem(ctx, i)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return status.Errorf(codes.NotFound, err.Error())
		} else {
			return status.Errorf(codes.Internal, err.Error())
		}
	}
	return nil
}

func (s *service) GetItemById(ctx context.Context, id string) (*model.Item, error) {
	i, err := s.repository.GetItemById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return i, status.Errorf(codes.NotFound, err.Error())
		} else {
			return i, status.Errorf(codes.Internal, err.Error())
		}
	}
	return i, nil
}

func (s *service) GetItems(ctx context.Context) ([]*model.Item, error) {
	i, err := s.repository.GetItems(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		} else {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}
	return i, nil
}

func (s *service) UpdateItem(ctx context.Context, id string, dto *model.UpdateItemDTO) error {
	if dto.Title == "" && dto.Comment == "" && dto.UserId == "" {
		return status.Errorf(codes.InvalidArgument, "nothing to update")
	}
	_, err := s.repository.GetItemById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return status.Errorf(codes.NotFound, err.Error())
		} else {
			return status.Errorf(codes.Internal, err.Error())
		}
	}
	i := model.UpdateItem(dto)
	err = s.repository.UpdateItem(ctx, id, i)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return status.Errorf(codes.NotFound, err.Error())
		} else {
			return status.Errorf(codes.Internal, err.Error())
		}
	}
	return nil
}

func (s *service) DeleteItemById(ctx context.Context, id string) error {
	_, err := s.repository.GetItemById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return status.Errorf(codes.NotFound, err.Error())
		} else {
			return status.Errorf(codes.Internal, err.Error())
		}
	}
	err = s.repository.DeleteItemById(ctx, id)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	return nil
}

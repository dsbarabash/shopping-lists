package service

import (
	"context"
	"errors"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	repository repository.Db
}

func (s *Service) CreateShoppingList(ctx context.Context, dto model.CreateShoppingListDTO) error {
	sl := model.NewShoppingList1(dto)
	err := s.repository.AddShoppingList(ctx, sl)
	if err != nil {
		if errors.Is(err, errors.New("NOT FOUND")) {
			return status.Errorf(codes.NotFound, err.Error())
		} else {
			return status.Errorf(codes.Internal, err.Error())
		}
	}
	return nil
}

func (s *Service) GetShoppingListById(ctx context.Context, id string) (*model.ShoppingList, error) {
	sl, err := s.repository.GetSlById(ctx, id)
	if err != nil {
		if errors.Is(err, errors.New("NOT FOUND")) {
			return sl, status.Errorf(codes.NotFound, err.Error())
		} else {
			return sl, status.Errorf(codes.Internal, err.Error())
		}
	}
	return sl, nil
}

func (s *Service) GetShoppingLists(ctx context.Context) ([]*model.ShoppingList, error) {
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

func (s *Service) UpdateShoppingList(ctx context.Context, id string, dto model.UpdateShoppingListDTO) error {
	if dto.Title == "" && len(dto.Items) == 0 && dto.UserId == "" {
		return status.Errorf(codes.InvalidArgument, "nothing to update")
	}
	sl := model.UpdateShoppingList1(dto)
	err := s.repository.UpdateSl(ctx, id, sl)
	if err != nil {
		if errors.Is(err, errors.New("NOT FOUND")) {
			return status.Errorf(codes.NotFound, err.Error())
		} else {
			return status.Errorf(codes.Internal, err.Error())
		}
	}
	return nil
}

func (s *Service) CreateItem(ctx context.Context, dto model.CreateItemDTO) error {
	i := model.NewItem1(dto)
	err := s.repository.AddItem(ctx, i)
	if err != nil {
		if errors.Is(err, errors.New("NOT FOUND")) {
			return status.Errorf(codes.NotFound, err.Error())
		} else {
			return status.Errorf(codes.Internal, err.Error())
		}
	}
	return nil
}

func (s *Service) GetItemById(ctx context.Context, id string) (*model.Item, error) {
	i, err := s.repository.GetItemById(ctx, id)
	if err != nil {
		if errors.Is(err, errors.New("NOT FOUND")) {
			return i, status.Errorf(codes.NotFound, err.Error())
		} else {
			return i, status.Errorf(codes.Internal, err.Error())
		}
	}
	return i, nil
}

func (s *Service) GetItems(ctx context.Context) ([]*model.Item, error) {
	i, err := s.repository.GetItems(ctx)
	if err != nil {
		if errors.Is(err, errors.New("NOT FOUND")) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		} else {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}
	return i, nil
}

func (s *Service) UpdateItem(ctx context.Context, id string, dto model.UpdateItemDTO) error {
	if dto.Title == "" && dto.Comment == "" && dto.UserId == "" {
		return status.Errorf(codes.InvalidArgument, "nothing to update")
	}
	i := model.UpdateItem1(dto)
	err := s.repository.UpdateItem(ctx, id, i)
	if err != nil {
		if errors.Is(err, errors.New("NOT FOUND")) {
			return status.Errorf(codes.NotFound, err.Error())
		} else {
			return status.Errorf(codes.Internal, err.Error())
		}
	}
	return nil
}

//func (s *Service) UpdateItem(ctx context.Context, id string, u UpdateItemBody) error {
//	_, err := s.repository.FindItem(ctx, id)
//	if err != nil {
//		return status.Errorf(codes.NotFound, err.Error())
//	}
//	u.UpdatedAt = time.Now().UTC()
//	err = s.Repository.UpdateItem(ctx, id, u.String())
//	if err != nil {
//		return status.Errorf(codes.Internal, err.Error())
//	}
//	return nil
//}

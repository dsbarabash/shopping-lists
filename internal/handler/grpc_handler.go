package handler

import (
	"context"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"github.com/dsbarabash/shopping-lists/internal/proto_api/pkg/grpc/v1/shopping_list_api"
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"time"
)

type Server struct {
	shopping_list_api.ShoppingListServiceServer
}

func (s *Server) CreateShoppingList(
	ctx context.Context,
	req *shopping_list_api.CreateShoppingListRequest,
) (*shopping_list_api.CreateShoppingListResponse, error) {
	if req.Title == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Title не должен быть пустым")
	}
	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "UserId не должен быть пустым")
	}
	iID, err := uuid.NewUUID()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}
	sl := &model.ShoppingList{
		Id:     iID.String(),
		Title:  req.Title,
		UserId: req.UserId,
		Items:  req.Items,
	}
	repository.ShoppingList.Add(sl)
	return &shopping_list_api.CreateShoppingListResponse{
		ShoppingList: &shopping_list_api.ShoppingList{
			Id:     iID.String(),
			Title:  req.Title,
			UserId: req.UserId,
			//CreatedAt:       time.Now(),
			//UpdatedAt:      time.Now(),
			Items: req.Items,
		},
	}, nil

}

func (s *Server) UpdateShoppingList(
	ctx context.Context,
	req *shopping_list_api.UpdateShoppingListRequest,
) (*shopping_list_api.UpdateShoppingListResponse, error) {
	if req.GetId() <= "" {
		return nil, status.Errorf(codes.InvalidArgument, "id не должен быть пустым")
	}
	for _, i := range repository.ShoppingList.Store {
		if i.Id == req.GetId() {
			i.Title = req.Title
			i.UserId = req.UserId
			i.UpdatedAt = time.Now()
			i.Items = req.Items

			repository.ItemList.SaveSliceToFile(repository.ItemList.Store)

			return &shopping_list_api.UpdateShoppingListResponse{
				ShoppingList: &shopping_list_api.ShoppingList{
					Id:     i.Id,
					Title:  i.Title,
					UserId: i.UserId,
					Items:  i.Items,
				},
			}, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "Запись не найдена")
}

func (s *Server) DeleteShoppingList(
	ctx context.Context,
	req *shopping_list_api.DeleteShoppingListRequest,
) (*shopping_list_api.DeleteShoppingListResponse, error) {
	if req.GetId() <= "" {
		return nil, status.Errorf(codes.InvalidArgument, "id не должен быть пустым")
	}

	for idx, sl := range repository.ShoppingList.Store {
		if sl.Id == req.GetId() {
			copy(repository.ShoppingList.Store[idx:], repository.ShoppingList.Store[idx+1:])
			repository.ShoppingList.Store = repository.ShoppingList.Store[:len(repository.ShoppingList.Store)-1]
			repository.ShoppingList.SaveSliceToFile(repository.ShoppingList.Store)
			return nil, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "Запись не найдена")
}
func (s *Server) GetShoppingList(
	ctx context.Context,
	req *shopping_list_api.GetShoppingListRequest,
) (*shopping_list_api.GetShoppingListResponse, error) {
	if req.GetId() <= "" {
		return nil, status.Errorf(codes.InvalidArgument, "id не должен быть пустым")
	}
	for _, i := range repository.ShoppingList.Store {
		if i.Id == req.GetId() {
			return &shopping_list_api.GetShoppingListResponse{
				ShoppingList: &shopping_list_api.ShoppingList{
					Id:     i.Id,
					Title:  i.Title,
					UserId: i.UserId,
					Items:  i.Items,
				},
			}, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "Запись не найдена")
}

func (s *Server) GetShoppingLists(ctx context.Context, _ *emptypb.Empty) (*shopping_list_api.GetShoppingListsResponse, error) {
	var sl []*shopping_list_api.ShoppingList
	for _, i := range repository.ShoppingList.Store {
		sl = append(sl, &shopping_list_api.ShoppingList{
			Id:     i.Id,
			Title:  i.Title,
			UserId: i.UserId,
			Items:  i.Items,
		})
	}
	return &shopping_list_api.GetShoppingListsResponse{
		ShoppingList: sl,
	}, nil
}

func (s *Server) CreateItem(
	ctx context.Context,
	req *shopping_list_api.CreateItemRequest,
) (*shopping_list_api.CreateItemResponse, error) {
	if req.Title == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Title не должен быть пустым")
	}
	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "UserId не должен быть пустым")
	}
	if req.ShoppingListId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "ShoppingListId не должен быть пустым")
	}
	iID, err := uuid.NewUUID()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}
	item := &model.Item{
		Id:             iID.String(),
		Title:          req.Title,
		Comment:        req.Comment,
		IsDone:         false,
		UserId:         req.UserId,
		ShoppingListId: req.ShoppingListId,
	}
	repository.ItemList.Add(item)
	return &shopping_list_api.CreateItemResponse{
		Item: &shopping_list_api.Item{
			Id:      iID.String(),
			Title:   req.Title,
			Comment: req.Comment,
			IsDone:  false,
			UserId:  req.UserId,
			//CreatedAt:       time.Now(),
			//UpdatedAt:      time.Now(),
			ShoppingListId: req.ShoppingListId,
		},
	}, nil

}

func (s *Server) UpdateItem(
	ctx context.Context,
	req *shopping_list_api.UpdateItemRequest,
) (*shopping_list_api.UpdateItemResponse, error) {
	if req.GetId() <= "" {
		return nil, status.Errorf(codes.InvalidArgument, "id не должен быть пустым")
	}
	for _, i := range repository.ItemList.Store {
		if i.Id == req.GetId() {
			i.Title = req.Title
			i.Comment = req.Comment
			i.IsDone = false
			i.UserId = req.UserId
			i.UpdatedAt = time.Now()
			i.ShoppingListId = req.ShoppingListId
			repository.ItemList.SaveSliceToFile(repository.ItemList.Store)

			return &shopping_list_api.UpdateItemResponse{
				Item: &shopping_list_api.Item{
					Id:             i.Id,
					Title:          i.Title,
					Comment:        i.Comment,
					IsDone:         i.IsDone,
					UserId:         i.UserId,
					ShoppingListId: i.ShoppingListId,
				},
			}, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "Запись не найдена")
}

func (s *Server) DeleteItem(
	ctx context.Context,
	req *shopping_list_api.DeleteItemRequest,
) (*shopping_list_api.DeleteItemResponse, error) {
	if req.GetId() <= "" {
		return nil, status.Errorf(codes.InvalidArgument, "id не должен быть пустым")
	}

	for idx, i := range repository.ItemList.Store {
		if i.Id == req.GetId() {
			copy(repository.ItemList.Store[idx:], repository.ItemList.Store[idx+1:])
			repository.ItemList.Store = repository.ItemList.Store[:len(repository.ItemList.Store)-1]
			repository.ItemList.SaveSliceToFile(repository.ItemList.Store)
			return nil, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "Запись не найдена")
}
func (s *Server) GetItem(
	ctx context.Context,
	req *shopping_list_api.GetItemRequest,
) (*shopping_list_api.GetItemResponse, error) {
	if req.GetId() <= "" {
		return nil, status.Errorf(codes.InvalidArgument, "id не должен быть пустым")
	}

	for _, i := range repository.ItemList.Store {
		if i.Id == req.GetId() {
			return &shopping_list_api.GetItemResponse{
				Item: &shopping_list_api.Item{
					Id:             i.Id,
					Title:          i.Title,
					Comment:        i.Comment,
					IsDone:         i.IsDone,
					UserId:         i.UserId,
					ShoppingListId: i.ShoppingListId,
				},
			}, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "Запись не найдена")
}

func (s *Server) GetItems(ctx context.Context, _ *emptypb.Empty) (*shopping_list_api.GetItemsResponse, error) {
	var items []*shopping_list_api.Item
	for _, i := range repository.ItemList.Store {
		items = append(items, &shopping_list_api.Item{
			Id:             i.Id,
			Title:          i.Title,
			Comment:        i.Comment,
			IsDone:         i.IsDone,
			UserId:         i.UserId,
			ShoppingListId: i.ShoppingListId,
		})
	}
	return &shopping_list_api.GetItemsResponse{
		Items: items,
	}, nil
}

func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	start := time.Now()
	resp, err = handler(ctx, req)
	st, _ := status.FromError(err)

	var reqJSON, respJSON string

	if m, ok := req.(proto.Message); ok {
		b, _ := protojson.Marshal(m)
		reqJSON = string(b)
	} else {
		reqJSON = "<non-proto request>"
	}

	if m, ok := resp.(proto.Message); ok && resp != nil {
		b, _ := protojson.Marshal(m)
		respJSON = string(b)
	} else {
		respJSON = "<non-proto response or nil>"
	}

	log.Printf(
		"[gRPC] method=%s status=%s error=%v duration=%s request=%s response=%s",
		info.FullMethod,
		st.Code(),
		err,
		time.Since(start),
		reqJSON,
		respJSON,
	)

	return resp, err
}

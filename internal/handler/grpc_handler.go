package handler

import (
	"context"
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"github.com/dsbarabash/shopping-lists/internal/proto_api/pkg/grpc/v1/shopping_list_api"
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

func (c *Controller) CreateShoppingList(
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
	c.MongoDb.AddShoppingList(ctx, sl)
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

func (c *Controller) UpdateShoppingList(
	ctx context.Context,
	req *shopping_list_api.UpdateShoppingListRequest,
) (*shopping_list_api.UpdateShoppingListResponse, error) {
	if req.GetId() <= "" {
		return nil, status.Errorf(codes.InvalidArgument, "id не должен быть пустым")
	}
	_, err := c.MongoDb.UpdateSl(ctx, req.GetId(), []byte(fmt.Sprintf(`{"Title": %s, "UserId": %s, "Items": %s}`, req.Title, req.UserId, req.Items)))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return &shopping_list_api.UpdateShoppingListResponse{
		ShoppingList: &shopping_list_api.ShoppingList{
			Id:     req.Id,
			Title:  req.Title,
			UserId: req.UserId,
			Items:  req.Items,
		},
	}, nil
}

func (c *Controller) DeleteShoppingList(
	ctx context.Context,
	req *shopping_list_api.DeleteShoppingListRequest,
) (*shopping_list_api.DeleteShoppingListResponse, error) {
	if req.GetId() <= "" {
		return nil, status.Errorf(codes.InvalidArgument, "id не должен быть пустым")
	}
	_, err := c.MongoDb.DeleteSlById(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return nil, nil
}
func (c *Controller) GetShoppingList(
	ctx context.Context,
	req *shopping_list_api.GetShoppingListRequest,
) (*shopping_list_api.GetShoppingListResponse, error) {
	if req.GetId() <= "" {
		return nil, status.Errorf(codes.InvalidArgument, "id не должен быть пустым")
	}
	sl, err := c.MongoDb.GetSlById(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return &shopping_list_api.GetShoppingListResponse{
		ShoppingList: &shopping_list_api.ShoppingList{
			Id:     sl[0].Id,
			Title:  sl[0].Title,
			UserId: sl[0].UserId,
			Items:  sl[0].Items,
		},
	}, nil

}

func (c *Controller) GetShoppingListsGrpc(ctx context.Context, _ *emptypb.Empty) (*shopping_list_api.GetShoppingListsResponse, error) {
	var sl []*shopping_list_api.ShoppingList
	list := c.MongoDb.GetSls(ctx)
	for _, i := range list {
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

func (c *Controller) CreateItem(
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
	c.MongoDb.AddItem(ctx, item)
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

func (c *Controller) UpdateItem(
	ctx context.Context,
	req *shopping_list_api.UpdateItemRequest,
) (*shopping_list_api.UpdateItemResponse, error) {
	if req.GetId() <= "" {
		return nil, status.Errorf(codes.InvalidArgument, "id не должен быть пустым")
	}
	_, err := c.MongoDb.UpdateItem(ctx, req.GetId(), []byte(fmt.Sprintf(`{"Title": %s, "Comment": %s, "IsDone": %s, "UserId": %s, "ShoppingListId": %s}`, req.Title, req.Comment, false, req.UserId, req.ShoppingListId)))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return &shopping_list_api.UpdateItemResponse{
		Item: &shopping_list_api.Item{
			Id:             req.Id,
			Title:          req.Title,
			Comment:        req.Comment,
			IsDone:         false,
			UserId:         req.UserId,
			ShoppingListId: req.ShoppingListId,
		},
	}, nil

}

func (c *Controller) DeleteItem(
	ctx context.Context,
	req *shopping_list_api.DeleteItemRequest,
) (*shopping_list_api.DeleteItemResponse, error) {
	if req.GetId() <= "" {
		return nil, status.Errorf(codes.InvalidArgument, "id не должен быть пустым")
	}
	_, err := c.MongoDb.DeleteItemById(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return nil, nil
}

func (c *Controller) GetItem(
	ctx context.Context,
	req *shopping_list_api.GetItemRequest,
) (*shopping_list_api.GetItemResponse, error) {
	if req.GetId() <= "" {
		return nil, status.Errorf(codes.InvalidArgument, "id не должен быть пустым")
	}
	item, err := c.MongoDb.GetItemById(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return &shopping_list_api.GetItemResponse{
		Item: &shopping_list_api.Item{
			Id:             item[0].Id,
			Title:          item[0].Title,
			Comment:        item[0].Comment,
			IsDone:         item[0].IsDone,
			UserId:         item[0].UserId,
			ShoppingListId: item[0].ShoppingListId,
		},
	}, nil
}

func (c *Controller) GetItemsGrpc(ctx context.Context, _ *emptypb.Empty) (*shopping_list_api.GetItemsResponse, error) {
	var items []*shopping_list_api.Item
	list := c.MongoDb.GetItems(ctx)
	for _, i := range list {
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

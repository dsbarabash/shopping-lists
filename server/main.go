package main

import (
	"context"
	"encoding/json"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"github.com/dsbarabash/shopping-lists/internal/proto_api/pkg/grpc/v1/shopping_list_api"
	"github.com/dsbarabash/shopping-lists/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type ItemStore struct {
	mu                  sync.Mutex
	store               []*model.Item
	printedElementCount int
	filePath            string
}

type ShoppingListStore struct {
	mu                  sync.Mutex
	store               []*model.ShoppingList
	printedElementCount int
	filePath            string
}

var ShoppingList = ShoppingListStore{
	sync.Mutex{},
	make([]*model.ShoppingList, 0),
	0,
	"./shoppingLists.json",
}

var ItemList = ItemStore{
	sync.Mutex{},
	make([]*model.Item, 0),
	0,
	"./items.json",
}

func ReadJson(fileName string) ([]byte, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (s *ShoppingListStore) LoadFromFile() {
	s.mu.Lock()
	defer s.mu.Unlock()
	items, err := ReadJson(s.filePath)
	if err == io.EOF {
		return
	} else if err != nil {
		log.Fatal(err)
	}
	if len(items) != 0 {
		if err := json.Unmarshal(items, &s.store); err != nil {
			log.Fatal(err)
		}
	}
}

func (s *ShoppingListStore) SaveToFile(sl *model.ShoppingList) {
	f, err := os.OpenFile(s.filePath, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	// Перемещаемся в конец файла
	stat, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	// Если это начала файла, начинаем массив json
	if stat.Size() == 0 {
		_, err := f.WriteString("[")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Если файл не пуст, то заменчем последнюю закрывающую скобку массива на запятую
		if _, err := f.Seek(-1, io.SeekEnd); err != nil {
			log.Fatal(err)
		}
		if _, err := f.WriteString(","); err != nil {
			log.Fatal(err)
		}
	}
	// Добавляем объект в файл и закрываем массив
	e := json.NewEncoder(f)
	if err := e.Encode(sl); err != nil {
		log.Fatal(err)
	}
	// Добавляем закрывающую скобку
	_, err = f.WriteString("]")
	if err != nil {
		log.Fatal(err)
	}
}

func (s *ShoppingListStore) SaveSliceToFile(sls []*model.ShoppingList) {
	f, err := os.OpenFile(s.filePath, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	e := json.NewEncoder(f)
	if err := e.Encode(sls); err != nil {
		log.Fatal(err)
	}
}

func (s *ShoppingListStore) Add(sl *model.ShoppingList) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store = append(s.store, sl)
	s.SaveToFile(sl)
}

func (i *ItemStore) LoadFromFile() {
	i.mu.Lock()
	defer i.mu.Unlock()
	items, err := ReadJson(i.filePath)
	if err == io.EOF {
		return
	} else if err != nil {
		log.Fatal(err)
	}
	if len(items) != 0 {
		if err := json.Unmarshal(items, &i.store); err != nil {
			log.Fatal(err)
		}
	}
}

func (i *ItemStore) SaveToFile(item *model.Item) {
	hu := i.filePath
	f, err := os.OpenFile(hu, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	// Перемещаемся в конец файла
	stat, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	// Если это начала файла, начинаем массив json
	if stat.Size() == 0 {
		_, err := f.WriteString("[")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Если файл не пуст, то заменчем последнюю закрывающую скобку массива на запятую
		if _, err := f.Seek(-1, io.SeekEnd); err != nil {
			log.Fatal(err)
		}
		if _, err := f.WriteString(","); err != nil {
			log.Fatal(err)
		}
	}
	// Добавляем объект в файл и закрываем массив
	e := json.NewEncoder(f)
	if err := e.Encode(item); err != nil {
		log.Fatal(err)
	}
	// Добавляем закрывающую скобку
	_, err = f.WriteString("]")
	if err != nil {
		log.Fatal(err)
	}
}

func (i *ItemStore) SaveSliceToFile(sls []*model.Item) {
	f, err := os.OpenFile(i.filePath, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	e := json.NewEncoder(f)
	if err := e.Encode(sls); err != nil {
		log.Fatal(err)
	}
}

func (i *ItemStore) Add(item *model.Item) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.store = append(i.store, item)
	i.SaveToFile(item)
}

func FillSlices() {
	ItemList.LoadFromFile()
	ShoppingList.LoadFromFile()
	service.UserList.LoadFromFile()
}

type server struct {
	shopping_list_api.ShoppingListServiceServer
}

func (s *server) CreateShoppingList(
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
	ShoppingList.Add(sl)
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

func (s *server) UpdateShoppingList(
	ctx context.Context,
	req *shopping_list_api.UpdateShoppingListRequest,
) (*shopping_list_api.UpdateShoppingListResponse, error) {
	if req.GetId() <= "" {
		return nil, status.Errorf(codes.InvalidArgument, "id не должен быть пустым")
	}
	for _, i := range ShoppingList.store {
		if i.Id == req.GetId() {
			i.Title = req.Title
			i.UserId = req.UserId
			i.UpdatedAt = time.Now()
			i.Items = req.Items

			ItemList.SaveSliceToFile(ItemList.store)

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

func (s *server) DeleteShoppingList(
	ctx context.Context,
	req *shopping_list_api.DeleteShoppingListRequest,
) (*shopping_list_api.DeleteShoppingListResponse, error) {
	if req.GetId() <= "" {
		return nil, status.Errorf(codes.InvalidArgument, "id не должен быть пустым")
	}

	for idx, sl := range ShoppingList.store {
		if sl.Id == req.GetId() {
			copy(ShoppingList.store[idx:], ShoppingList.store[idx+1:])
			ShoppingList.store = ShoppingList.store[:len(ShoppingList.store)-1]
			ShoppingList.SaveSliceToFile(ShoppingList.store)
			return nil, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "Запись не найдена")
}
func (s *server) GetShoppingList(
	ctx context.Context,
	req *shopping_list_api.GetShoppingListRequest,
) (*shopping_list_api.GetShoppingListResponse, error) {
	if req.GetId() <= "" {
		return nil, status.Errorf(codes.InvalidArgument, "id не должен быть пустым")
	}
	for _, i := range ShoppingList.store {
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

func (s *server) GetShoppingLists(ctx context.Context, _ *emptypb.Empty) (*shopping_list_api.GetShoppingListsResponse, error) {
	var sl []*shopping_list_api.ShoppingList
	for _, i := range ShoppingList.store {
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

func (s *server) CreateItem(
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
	ItemList.Add(item)
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

func (s *server) UpdateItem(
	ctx context.Context,
	req *shopping_list_api.UpdateItemRequest,
) (*shopping_list_api.UpdateItemResponse, error) {
	if req.GetId() <= "" {
		return nil, status.Errorf(codes.InvalidArgument, "id не должен быть пустым")
	}
	for _, i := range ItemList.store {
		if i.Id == req.GetId() {
			i.Title = req.Title
			i.Comment = req.Comment
			i.IsDone = false
			i.UserId = req.UserId
			i.UpdatedAt = time.Now()
			i.ShoppingListId = req.ShoppingListId
			ItemList.SaveSliceToFile(ItemList.store)

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

func (s *server) DeleteItem(
	ctx context.Context,
	req *shopping_list_api.DeleteItemRequest,
) (*shopping_list_api.DeleteItemResponse, error) {
	if req.GetId() <= "" {
		return nil, status.Errorf(codes.InvalidArgument, "id не должен быть пустым")
	}

	for idx, i := range ItemList.store {
		if i.Id == req.GetId() {
			copy(ItemList.store[idx:], ItemList.store[idx+1:])
			ItemList.store = ItemList.store[:len(ItemList.store)-1]
			ItemList.SaveSliceToFile(ItemList.store)
			return nil, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "Запись не найдена")
}
func (s *server) GetItem(
	ctx context.Context,
	req *shopping_list_api.GetItemRequest,
) (*shopping_list_api.GetItemResponse, error) {
	if req.GetId() <= "" {
		return nil, status.Errorf(codes.InvalidArgument, "id не должен быть пустым")
	}

	for _, i := range ItemList.store {
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

func (s *server) GetItems(ctx context.Context, _ *emptypb.Empty) (*shopping_list_api.GetItemsResponse, error) {
	var items []*shopping_list_api.Item
	for _, i := range ItemList.store {
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

func main() {
	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	FillSlices()
	s := grpc.NewServer()
	shopping_list_api.RegisterShoppingListServiceServer(s, &server{})

	reflection.Register(s)

	log.Println("Server is running at :5001")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

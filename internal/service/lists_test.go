package service

import (
	"context"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"reflect"
	"testing"
)

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := repository.NewMockDb(ctrl)
	type args struct {
		repository repository.Db
	}
	tests := []struct {
		name    string
		args    args
		want    Service
		wantErr bool
	}{
		{
			"Valid",
			args{mockService},
			&service{
				repository: mockService,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewService(tt.args.repository)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_CreateItem(t *testing.T) {
	type args struct {
		ctx context.Context
		dto *model.CreateItemDTO
	}
	ID, _ := uuid.NewUUID()
	timeNow := timestamppb.Now()
	var createItemDTO = &model.CreateItemDTO{
		Id:             ID.String(),
		Title:          "Test New Item",
		Comment:        "test comment",
		IsDone:         false,
		UserId:         "12345",
		CreatedAt:      timeNow,
		UpdatedAt:      timeNow,
		ShoppingListId: "4321",
	}
	var item = &model.Item{
		Id:             ID.String(),
		Title:          "Test New Item",
		Comment:        "test comment",
		IsDone:         false,
		UserId:         "12345",
		CreatedAt:      timeNow,
		UpdatedAt:      timeNow,
		ShoppingListId: "4321",
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		prepare func(mockService *repository.MockDb)
	}{
		{
			"Valid",
			args{context.Background(), createItemDTO},
			false,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().AddItem(context.Background(), item).Return(nil).Times(1)
			},
		},
		{
			"InValidErrInternal",
			args{context.Background(), createItemDTO},
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().AddItem(context.Background(), item).Return(io.EOF).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := repository.NewMockDb(ctrl)
			s := &service{
				repository: mockService,
			}
			tt.prepare(mockService)
			if err := s.CreateItem(tt.args.ctx, tt.args.dto); (err != nil) != tt.wantErr {
				t.Errorf("CreateItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_GetItemById(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	ID, _ := uuid.NewUUID()
	timeNow := timestamppb.Now()
	var item = &model.Item{
		Id:             ID.String(),
		Title:          "Test New Item",
		Comment:        "test comment",
		IsDone:         false,
		UserId:         "12345",
		CreatedAt:      timeNow,
		UpdatedAt:      timeNow,
		ShoppingListId: "4321",
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Item
		wantErr bool
		prepare func(mockService *repository.MockDb)
	}{
		{
			"Valid",
			args{context.Background(), "123"},
			item,
			false,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetItemById(context.Background(), "123").Return(item, nil).Times(1)
			},
		},
		{
			"InvalidErrNotFound",
			args{context.Background(), "123"},
			item,
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetItemById(context.Background(), "123").Return(item, repository.ErrNotFound).Times(1)
			},
		},
		{
			"InvalidErrInternal",
			args{context.Background(), "123"},
			item,
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetItemById(context.Background(), "123").Return(item, io.EOF).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := repository.NewMockDb(ctrl)
			s := &service{
				repository: mockService,
			}
			tt.prepare(mockService)
			got, err := s.GetItemById(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetItemById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetItemById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_GetItems(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	ID, _ := uuid.NewUUID()
	timeNow := timestamppb.Now()
	var items = []*model.Item{
		{Id: ID.String(),
			Title:          "Test New Item",
			Comment:        "test comment",
			IsDone:         false,
			UserId:         "12345",
			CreatedAt:      timeNow,
			UpdatedAt:      timeNow,
			ShoppingListId: "4321"},
		{
			Id:             "23122",
			Title:          "Test New Item 2",
			Comment:        "test comment 2",
			IsDone:         false,
			UserId:         "555",
			CreatedAt:      timeNow,
			UpdatedAt:      timeNow,
			ShoppingListId: "4444",
		},
	}
	tests := []struct {
		name    string
		args    args
		want    []*model.Item
		wantErr bool
		prepare func(mockService *repository.MockDb)
	}{
		{
			"Valid",
			args{context.Background()},
			items,
			false,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetItems(context.Background()).Return(items, nil).Times(1)
			},
		},
		{
			"InvalidErrNotFound",
			args{context.Background()},
			nil,
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetItems(context.Background()).Return(nil, repository.ErrNotFound).Times(1)
			},
		},
		{
			"InvalidErrInternal",
			args{context.Background()},
			nil,
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetItems(context.Background()).Return(nil, io.EOF).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := repository.NewMockDb(ctrl)
			s := &service{
				repository: mockService,
			}
			tt.prepare(mockService)
			got, err := s.GetItems(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetItems() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_UpdateItem(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
		dto *model.UpdateItemDTO
	}
	ID, _ := uuid.NewUUID()
	timeNow := timestamppb.Now()
	var item = &model.Item{
		Id:             ID.String(),
		Title:          "Test New Item",
		Comment:        "test comment",
		IsDone:         false,
		UserId:         "12345",
		CreatedAt:      timeNow,
		UpdatedAt:      timeNow,
		ShoppingListId: "4321",
	}
	var updateItemDTO = &model.UpdateItemDTO{
		Id:             ID.String(),
		Title:          "Updated Test New Item",
		Comment:        "Updated test comment",
		IsDone:         false,
		UserId:         "123456",
		CreatedAt:      timeNow,
		UpdatedAt:      timestamppb.Now(),
		ShoppingListId: "4321",
	}
	updateItem := model.UpdateItem(updateItemDTO)
	tests := []struct {
		name    string
		args    args
		wantErr bool
		prepare func(mockService *repository.MockDb)
	}{
		{
			"Valid",
			args{context.Background(), ID.String(), updateItemDTO},
			false,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetItemById(context.Background(), ID.String()).Return(item, nil).Times(1)
				mockService.EXPECT().UpdateItem(context.Background(), ID.String(), updateItem).Return(nil).Times(1)
			},
		},
		{
			"InvalidEmptyParameters",
			args{context.Background(), ID.String(), &model.UpdateItemDTO{}},
			true,
			func(mockService *repository.MockDb) {
			},
		},
		{
			"InvalidErrNotFoundGet",
			args{context.Background(), ID.String(), updateItemDTO},
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetItemById(context.Background(), ID.String()).Return(nil, repository.ErrNotFound).Times(1)
			},
		},
		{
			"InvalidErrInternalGet",
			args{context.Background(), ID.String(), updateItemDTO},
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetItemById(context.Background(), ID.String()).Return(nil, io.EOF).Times(1)
			},
		},
		{
			"InvalidErrNotFoundUpdate",
			args{context.Background(), ID.String(), updateItemDTO},
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetItemById(context.Background(), ID.String()).Return(nil, nil).Times(1)
				mockService.EXPECT().UpdateItem(context.Background(), ID.String(), updateItem).Return(repository.ErrNotFound).Times(1)
			},
		},
		{
			"InvalidErrInternalUpdate",
			args{context.Background(), ID.String(), updateItemDTO},
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetItemById(context.Background(), ID.String()).Return(nil, nil).Times(1)
				mockService.EXPECT().UpdateItem(context.Background(), ID.String(), updateItem).Return(io.EOF).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := repository.NewMockDb(ctrl)
			s := &service{
				repository: mockService,
			}
			tt.prepare(mockService)
			if err := s.UpdateItem(tt.args.ctx, tt.args.id, tt.args.dto); (err != nil) != tt.wantErr {
				t.Errorf("UpdateItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_DeleteItemById(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	ID, _ := uuid.NewUUID()
	timeNow := timestamppb.Now()
	var item = &model.Item{
		Id:             ID.String(),
		Title:          "Test New Item",
		Comment:        "test comment",
		IsDone:         false,
		UserId:         "12345",
		CreatedAt:      timeNow,
		UpdatedAt:      timeNow,
		ShoppingListId: "4321",
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		prepare func(mockService *repository.MockDb)
	}{
		{
			"Valid",
			args{context.Background(), "123"},
			false,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetItemById(context.Background(), "123").Return(item, nil).Times(1)
				mockService.EXPECT().DeleteItemById(context.Background(), "123").Return(nil).Times(1)
			},
		},
		{
			"InvalidErrNotFound",
			args{context.Background(), "123"},
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetItemById(context.Background(), "123").Return(item, repository.ErrNotFound).Times(1)
			},
		},
		{
			"InvalidErrInternalGet",
			args{context.Background(), "123"},
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetItemById(context.Background(), "123").Return(item, io.EOF).Times(1)
			},
		},
		{
			"InvalidErrInternalDelete",
			args{context.Background(), "123"},
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetItemById(context.Background(), "123").Return(item, nil).Times(1)
				mockService.EXPECT().DeleteItemById(context.Background(), "123").Return(io.EOF).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := repository.NewMockDb(ctrl)
			s := &service{
				repository: mockService,
			}
			tt.prepare(mockService)
			if err := s.DeleteItemById(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteItemById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_CreateShoppingList(t *testing.T) {
	type args struct {
		ctx context.Context
		dto *model.CreateShoppingListDTO
	}
	ID, _ := uuid.NewUUID()
	timeNow := timestamppb.Now()
	var createSLDTO = &model.CreateShoppingListDTO{
		Id:        ID.String(),
		Title:     "Test New SL",
		UserId:    "12345",
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		Items:     []string{"234", "567"},
	}
	var sl = &model.ShoppingList{
		Id:        ID.String(),
		Title:     "Test New SL",
		UserId:    "12345",
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		Items:     []string{"234", "567"},
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		prepare func(mockService *repository.MockDb)
	}{
		{
			"Valid",
			args{context.Background(), createSLDTO},
			false,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().AddShoppingList(context.Background(), sl).Return(nil).Times(1)
			},
		},
		{
			"InvalidErrInternal",
			args{context.Background(), createSLDTO},
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().AddShoppingList(context.Background(), sl).Return(io.EOF).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := repository.NewMockDb(ctrl)
			s := &service{
				repository: mockService,
			}
			tt.prepare(mockService)
			if err := s.CreateShoppingList(tt.args.ctx, tt.args.dto); (err != nil) != tt.wantErr {
				t.Errorf("CreateShoppingList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_DeleteShoppingListById(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	ID, _ := uuid.NewUUID()
	timeNow := timestamppb.Now()
	var sl = &model.ShoppingList{
		Id:        ID.String(),
		Title:     "Test New SL",
		UserId:    "12345",
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		Items:     []string{"234", "567"},
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		prepare func(mockService *repository.MockDb)
	}{
		{
			"Valid",
			args{context.Background(), "123"},
			false,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetSlById(context.Background(), "123").Return(sl, nil).Times(1)
				mockService.EXPECT().DeleteSlById(context.Background(), "123").Return(nil).Times(1)
			},
		},
		{
			"InvalidErrNotFound",
			args{context.Background(), "123"},
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetSlById(context.Background(), "123").Return(sl, repository.ErrNotFound).Times(1)
			},
		},
		{
			"InvalidErrInternalGet",
			args{context.Background(), "123"},
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetSlById(context.Background(), "123").Return(sl, io.EOF).Times(1)
			},
		},
		{
			"InvalidErrInternalGDelete",
			args{context.Background(), "123"},
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetSlById(context.Background(), "123").Return(sl, nil).Times(1)
				mockService.EXPECT().DeleteSlById(context.Background(), "123").Return(io.EOF).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := repository.NewMockDb(ctrl)
			s := &service{
				repository: mockService,
			}
			tt.prepare(mockService)
			if err := s.DeleteShoppingListById(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteShoppingListById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_GetShoppingListById(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	ID, _ := uuid.NewUUID()
	timeNow := timestamppb.Now()
	var sl = &model.ShoppingList{
		Id:        ID.String(),
		Title:     "Test New SL",
		UserId:    "12345",
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		Items:     []string{"234", "567"},
	}
	tests := []struct {
		name    string
		args    args
		want    *model.ShoppingList
		wantErr bool
		prepare func(mockService *repository.MockDb)
	}{
		{
			"Valid",
			args{context.Background(), "123"},
			sl,
			false,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetSlById(context.Background(), "123").Return(sl, nil).Times(1)
			},
		},
		{
			"InvalidErrNotFound",
			args{context.Background(), "123"},
			nil,
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetSlById(context.Background(), "123").Return(nil, repository.ErrNotFound).Times(1)
			},
		},
		{
			"InvalidErrInternal",
			args{context.Background(), "123"},
			nil,
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetSlById(context.Background(), "123").Return(nil, io.EOF).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := repository.NewMockDb(ctrl)
			s := &service{
				repository: mockService,
			}
			tt.prepare(mockService)
			got, err := s.GetShoppingListById(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetShoppingListById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetShoppingListById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_GetShoppingLists(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	ID, _ := uuid.NewUUID()
	timeNow := timestamppb.Now()
	var shoppingLists = []*model.ShoppingList{
		{
			Id:        ID.String(),
			Title:     "Test New SL 1",
			UserId:    "12345",
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			Items:     []string{"2341", "5671"},
		},
		{
			Id:        "1231",
			Title:     "Test New SL 2",
			UserId:    "12345",
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			Items:     []string{"2342", "5672"},
		},
	}
	tests := []struct {
		name    string
		args    args
		want    []*model.ShoppingList
		wantErr bool
		prepare func(mockService *repository.MockDb)
	}{
		{
			"Valid",
			args{context.Background()},
			shoppingLists,
			false,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetSls(context.Background()).Return(shoppingLists, nil).Times(1)
			},
		},
		{
			"InvalidErrNotFound",
			args{context.Background()},
			nil,
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetSls(context.Background()).Return(nil, repository.ErrNotFound).Times(1)
			},
		},
		{
			"InvalidErrInternal",
			args{context.Background()},
			nil,
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetSls(context.Background()).Return(nil, io.EOF).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := repository.NewMockDb(ctrl)
			s := &service{
				repository: mockService,
			}
			tt.prepare(mockService)
			got, err := s.GetShoppingLists(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetShoppingLists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetShoppingLists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_UpdateShoppingList(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
		dto *model.UpdateShoppingListDTO
	}
	ID, _ := uuid.NewUUID()
	timeNow := timestamppb.Now()
	var sl = &model.ShoppingList{
		Id:        ID.String(),
		Title:     "Test New SL",
		UserId:    "12345",
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		Items:     []string{"234", "567"},
	}
	var updateSlDTO = &model.UpdateShoppingListDTO{
		Id:        ID.String(),
		Title:     "Updated Test New SL",
		UserId:    "123456",
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		Items:     []string{"234", "567"},
	}
	updateSl := model.UpdateShoppingList(updateSlDTO)
	tests := []struct {
		name    string
		args    args
		wantErr bool
		prepare func(mockService *repository.MockDb)
	}{
		{
			"Valid",
			args{context.Background(), ID.String(), updateSlDTO},
			false,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetSlById(context.Background(), ID.String()).Return(sl, nil).Times(1)
				mockService.EXPECT().UpdateSl(context.Background(), ID.String(), updateSl).Return(nil).Times(1)
			},
		},
		{
			"InvalidEmptyParameters",
			args{context.Background(), ID.String(), &model.UpdateShoppingListDTO{}},
			true,
			func(mockService *repository.MockDb) {
			},
		},
		{
			"InvalidErrNotFoundGet",
			args{context.Background(), ID.String(), updateSlDTO},
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetSlById(context.Background(), ID.String()).Return(nil, repository.ErrNotFound).Times(1)
			},
		},
		{
			"InvalidErrInternalGet",
			args{context.Background(), ID.String(), updateSlDTO},
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetSlById(context.Background(), ID.String()).Return(nil, io.EOF).Times(1)
			},
		},
		{
			"InvalidErrNotFoundUpdate",
			args{context.Background(), ID.String(), updateSlDTO},
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetSlById(context.Background(), ID.String()).Return(nil, nil).Times(1)
				mockService.EXPECT().UpdateSl(context.Background(), ID.String(), updateSl).Return(repository.ErrNotFound).Times(1)
			},
		},
		{
			"InvalidErrInternalUpdate",
			args{context.Background(), ID.String(), updateSlDTO},
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().GetSlById(context.Background(), ID.String()).Return(nil, nil).Times(1)
				mockService.EXPECT().UpdateSl(context.Background(), ID.String(), updateSl).Return(io.EOF).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := repository.NewMockDb(ctrl)
			s := &service{
				repository: mockService,
			}
			tt.prepare(mockService)
			if err := s.UpdateShoppingList(tt.args.ctx, tt.args.id, tt.args.dto); (err != nil) != tt.wantErr {
				t.Errorf("UpdateShoppingList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

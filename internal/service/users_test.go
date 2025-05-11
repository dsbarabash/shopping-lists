package service

import (
	"context"
	"errors"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"io"
	"reflect"
	"testing"
)

func TestNewUserService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := repository.NewMockDb(ctrl)
	type args struct {
		repository repository.Db
	}
	tests := []struct {
		name    string
		args    args
		want    UserService
		wantErr bool
	}{
		{
			"Valid",
			args{mockService},
			&userService{
				repository: mockService,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUserService(tt.args.repository)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUserService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_CreateUser(t *testing.T) {
	type args struct {
		ctx context.Context
		dto *model.CreateUserDTO
	}
	ID, _ := uuid.NewUUID()
	var userDTO = &model.CreateUserDTO{
		Id:       ID.String(),
		Name:     "Test user",
		Password: "123456",
		State:    1,
	}
	var user = model.NewUser(userDTO)
	tests := []struct {
		name    string
		args    args
		wantErr bool
		prepare func(mockService *repository.MockDb)
	}{
		{
			"Valid",
			args{context.Background(), userDTO},
			false,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().CreateUser(context.Background(), user).Return(nil).Times(1)
			},
		},
		{
			"InvalidInternal",
			args{context.Background(), userDTO},
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().CreateUser(context.Background(), user).Return(io.EOF).Times(1)
			},
		},
		{
			"InvalidUserAlreadyExist",
			args{context.Background(), userDTO},
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().CreateUser(context.Background(), user).Return(errors.New("USER ALREADY EXISTS")).Times(1)
			},
		},
		{
			"InvalidErrUserNotActive",
			args{context.Background(), userDTO},
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().CreateUser(context.Background(), user).Return(errors.New("USER NOT ACTIVE")).Times(1)
			},
		},
		{
			"InvalidErrNotFound",
			args{context.Background(), userDTO},
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().CreateUser(context.Background(), user).Return(repository.ErrNotFound).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := repository.NewMockDb(ctrl)
			s := &userService{
				repository: mockService,
			}
			tt.prepare(mockService)
			if err := s.CreateUser(tt.args.ctx, tt.args.dto); (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userService_Login(t *testing.T) {
	type args struct {
		ctx  context.Context
		user *model.User
	}
	ID, _ := uuid.NewUUID()
	var userDTO = &model.CreateUserDTO{
		Id:       ID.String(),
		Name:     "Test user",
		Password: "123456",
		State:    1,
	}
	var user = model.NewUser(userDTO)
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		prepare func(mockService *repository.MockDb)
	}{
		{
			"Valid",
			args{context.Background(), user},
			"123",
			false,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().Login(context.Background(), user).Return("123", nil).Times(1)
			},
		},
		{
			"InvalidInternal",
			args{context.Background(), user},
			"",
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().Login(context.Background(), user).Return("", io.EOF).Times(1)
			},
		},
		{
			"InvalidErrUserNotActive",
			args{context.Background(), user},
			"",
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().Login(context.Background(), user).Return("", errors.New("USER NOT ACTIVE")).Times(1)
			},
		},
		{
			"InvalidErrNotFound",
			args{context.Background(), user},
			"",
			true,
			func(mockService *repository.MockDb) {
				mockService.EXPECT().Login(context.Background(), user).Return("", repository.ErrNotFound).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := repository.NewMockDb(ctrl)
			s := &userService{
				repository: mockService,
			}
			tt.prepare(mockService)
			got, err := s.Login(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Login() got = %v, want %v", got, tt.want)
			}
		})
	}
}

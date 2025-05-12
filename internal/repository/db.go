package repository

import (
	"context"
	"errors"
	"github.com/dsbarabash/shopping-lists/internal/model"
)

//go:generate mockgen -source=db.go -destination=db_mock.go -package=repository

type Db interface {
	AddItem(context.Context, *model.Item) error
	AddShoppingList(context.Context, *model.ShoppingList) error
	GetItems(context.Context) ([]*model.Item, error)
	GetItemById(context.Context, string) (*model.Item, error)
	DeleteItemById(context.Context, string) error
	GetSls(context.Context) ([]*model.ShoppingList, error)
	GetSlById(context.Context, string) (*model.ShoppingList, error)
	DeleteSlById(context.Context, string) error
	UpdateSl(context.Context, string, *model.ShoppingList) error
	UpdateItem(context.Context, string, *model.Item) error
	CreateUser(context.Context, *model.User) error
	Login(context.Context, *model.User) (string, error)
}

var ErrNotFound = errors.New("NOT FOUND")

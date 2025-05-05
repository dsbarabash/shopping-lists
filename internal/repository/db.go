package repository

import (
	"context"
	"github.com/dsbarabash/shopping-lists/internal/model"
)

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
	Registration(context.Context, string, string) (*model.User, error)
	Login(context.Context, *model.User) (string, error)
}

package repository

import (
	"context"
	"github.com/dsbarabash/shopping-lists/internal/model"
)

type Db interface {
	AddItem(context.Context, *model.Item)
	AddShoppingList(context.Context, *model.ShoppingList)
	GetItems(context.Context) []model.Item
	GetItemById(context.Context, string) ([]model.Item, error)
	DeleteItemById(context.Context, string) error
	GetSls(context.Context) []model.ShoppingList
	GetSlById(context.Context, string) ([]model.ShoppingList, error)
	DeleteSlById(context.Context, string) error
	UpdateSl(context.Context, string, string) error
	FindItem(context.Context, string) (*model.Item, error)
	UpdateItem(context.Context, string, string) error
	Registration(context.Context, string, string) *model.User
	Login(context.Context, *model.User) (string, error)
}

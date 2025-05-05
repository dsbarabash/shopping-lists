package repository

import (
	"context"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"time"
)

type Db interface {
	AddItem(context.Context, *Item)
	AddShoppingList(context.Context, *ShoppingList)
	GetItems(context.Context) []Item
	GetItemById(context.Context, string) ([]Item, error)
	DeleteItemById(context.Context, string) error
	GetSls(context.Context) []ShoppingList
	GetSlById(context.Context, string) ([]ShoppingList, error)
	DeleteSlById(context.Context, string) error
	UpdateSl(context.Context, string, string) error
	FindItem(context.Context, string) (*Item, error)
	UpdateItem(context.Context, string, string) error
	Registration(context.Context, string, string) *model.User
	Login(context.Context, *model.User) (string, error)
}

type Item struct {
	Id             string    `json:"id"`
	Title          string    `json:"title"`
	Comment        string    `json:"comment"`
	IsDone         bool      `json:"is_done"`
	UserId         string    `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	ShoppingListId string    `json:"shopping_list_id"`
}

type ShoppingList struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	UserId    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Items     []string  `json:"items"`
	State     int       `json:"state"`
}

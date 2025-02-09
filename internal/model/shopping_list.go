package model

import (
	"time"
)

type ShoppingList struct {
	id        string
	title     string
	userId    string
	createdAt time.Time
	updatedAt time.Time
	items     []string
	state     State
}

func (s ShoppingList) SetId(id string) {
	s.id = id
}

func (s ShoppingList) GetId() string {
	return s.id
}

func (s ShoppingList) SetTitle(title string) {
	s.title = title
}

func (s ShoppingList) GetTitle() string {
	return s.title
}

func (s ShoppingList) AddItemToList(item string) {
	s.items = append(s.items, item)
}

type Item struct {
	id             string
	title          string
	comment        string
	isDone         bool
	userId         string
	createdAt      time.Time
	updatedAt      time.Time
	ShoppingListId string
}

func (i Item) SetId(id string) {
	i.id = id
}

func (i Item) GetId() string {
	return i.id
}

func (i Item) SetTitle(title string) {
	i.title = title
}

func (i Item) GetTitle() string {
	return i.title
}

package model

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type ShoppingLists interface {
	UpdateShoppingList(string, []string)
	String() string
}

type ShoppingList struct {
	id        string
	title     string
	userId    string
	createdAt time.Time
	updatedAt time.Time
	items     []string
	state     State
}

func NewShoppingList(title string, userId string) (*ShoppingList, error) {
	if title == "" {
		return nil, errors.New("title must not be empty")
	} else if userId == "" {
		return nil, errors.New("userId must not be empty")
	}
	id := uuid.New()
	return &ShoppingList{
		id:        id.String(),
		title:     title,
		userId:    userId,
		createdAt: time.Now(),
		updatedAt: time.Now(),
		items:     make([]string, 0),
		state:     1,
	}, nil
}

func (s *ShoppingList) UpdateShoppingList(title string, items []string) {
	s.title = title
	s.updatedAt = time.Now()
	for _, i := range items {
		s.items = append(s.items, i)
	}
}

func (s ShoppingList) String() string {
	return fmt.Sprintf("id: \"%s\", title: \"%s\", userId: \"%s\", createdAt: \"%s\", updatedAt: \"%s\"", s.id, s.title, s.userId, s.createdAt.Format(time.DateTime), s.updatedAt.Format(time.DateTime))
}

type Items interface {
	UpdateItem(string, string, bool)
	String() string
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

func NewItem(title string, comment string, userId string, shoppingListId string) (*Item, error) {
	if title == "" {
		return nil, errors.New("title must not be empty")
	} else if userId == "" {
		return nil, errors.New("userId must not be empty")
	} else if shoppingListId == "" {
		return nil, errors.New("shoppingListId must not be empty")
	}
	id := uuid.New()
	return &Item{
		id:             id.String(),
		title:          title,
		comment:        comment,
		isDone:         false,
		userId:         userId,
		createdAt:      time.Now(),
		updatedAt:      time.Now(),
		ShoppingListId: shoppingListId,
	}, nil
}

func (i *Item) UpdateItem(title string, comment string, isDone bool) {
	i.title = title
	i.comment = comment
	i.isDone = isDone
	i.updatedAt = time.Now()
}

func (i Item) String() string {
	return fmt.Sprintf("id: \"%s\", title: \"%s\", comment: \"%s\", isDone: \"%v\", userId: \"%s\", createdAt: \"%s\", updatedAt: \"%s\", ShoppingListId: \"%s\"", i.id, i.title, i.comment, i.isDone, i.userId, i.createdAt.Format(time.DateTime), i.updatedAt.Format(time.DateTime), i.ShoppingListId)
}

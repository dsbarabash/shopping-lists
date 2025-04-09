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
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	UserId    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Items     []string  `json:"items"`
	State     State     `json:"state"`
}

func NewShoppingList(title string, userId string) (*ShoppingList, error) {
	if title == "" {
		return nil, errors.New("title must not be empty")
	} else if userId == "" {
		return nil, errors.New("userId must not be empty")
	}
	id := uuid.New()
	return &ShoppingList{
		Id:        id.String(),
		Title:     title,
		UserId:    userId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Items:     make([]string, 0),
		State:     1,
	}, nil
}

func (s *ShoppingList) UpdateShoppingList(title string, items []string) {
	s.Title = title
	s.UpdatedAt = time.Now()
	for _, i := range items {
		s.Items = append(s.Items, i)
	}
}

func (s ShoppingList) String() string {
	return fmt.Sprintf("id: \"%s\", title: \"%s\", userId: \"%s\", createdAt: \"%s\", updatedAt: \"%s\"", s.Id, s.Title, s.UserId, s.CreatedAt.Format(time.DateTime), s.UpdatedAt.Format(time.DateTime))
}

type Items interface {
	UpdateItem(string, string, bool)
	String() string
}

type Item struct {
	Id             string    `json:"id"`
	Title          string    `json:"title"`
	Comment        string    `json:"comment"`
	IsDone         bool      `json:"is_done"`
	UserId         string    `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updatedA_at"`
	ShoppingListId string    `json:"shopping_list_id"`
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
		Id:             id.String(),
		Title:          title,
		Comment:        comment,
		IsDone:         false,
		UserId:         userId,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ShoppingListId: shoppingListId,
	}, nil
}

func (i *Item) UpdateItem(title string, comment string, isDone bool) {
	i.Title = title
	i.Comment = comment
	i.IsDone = isDone
	i.UpdatedAt = time.Now()
}

func (i Item) String() string {
	return fmt.Sprintf("id: \"%s\", title: \"%s\", comment: \"%s\", isDone: \"%v\", userId: \"%s\", createdAt: \"%s\", updatedAt: \"%s\", ShoppingListId: \"%s\"", i.Id, i.Title, i.Comment, i.IsDone, i.UserId, i.CreatedAt.Format(time.DateTime), i.UpdatedAt.Format(time.DateTime), i.ShoppingListId)
}

type CreateItemRequest struct {
	Title          string `json:"title"`
	Comment        string `json:"comment"`
	UserId         string `json:"user_id"`
	ShoppingListId string `json:"shopping_list_id"`
}

type CreateShoppingListRequest struct {
	Title  string   `json:"title"`
	UserId string   `json:"user_id"`
	Items  []string `json:"items"`
}

type UpdateShoppingListRequest struct {
	Title  string   `json:"title"`
	UserId string   `json:"user_id"`
	Items  []string `json:"items"`
	State  State    `json:"state"`
}

type UpdateItemRequest struct {
	Title          string `json:"title"`
	Comment        string `json:"comment"`
	IsDone         bool   `json:"is_done"`
	UserId         string `json:"user_id"`
	ShoppingListId string `json:"shopping_list_id"`
}

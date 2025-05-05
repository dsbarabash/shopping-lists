package model

import (
	"context"
	"errors"
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type ShoppingLists interface {
	UpdateShoppingList(string, []string)
	String() string
}

type DobryniaModel struct {
	Repository repository.Db
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

func (s *ShoppingList) NewShoppingList(title string, userId string) (*ShoppingList, error) {
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
	UpdatedAt      time.Time `json:"updated_at"`
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

type UpdateItemBody struct {
	Title          string    `json:"title"`
	Comment        string    `json:"comment"`
	IsDone         bool      `json:"is_done"`
	UserId         string    `json:"user_id"`
	ShoppingListId string    `json:"shopping_list_id"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (u UpdateItemBody) String() string {
	return fmt.Sprintf("title: \"%s\", comment: \"%s\", isDone: \"%v\", userId: \"%s\", updatedAt: \"%s\", ShoppingListId: \"%s\"", u.Title, u.Comment, u.IsDone, u.UserId, u.UpdatedAt.Format(time.DateTime), u.ShoppingListId)
}
func (m *DobryniaModel) UpdateItem(ctx context.Context, id string, u UpdateItemBody) error {
	_, err := m.Repository.FindItem(ctx, id)
	if err != nil {
		return status.Errorf(codes.NotFound, err.Error())
	}
	u.UpdatedAt = time.Now().UTC()
	err = m.Repository.UpdateItem(ctx, id, u.String())
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	return nil
}

func (i Item) String() string {
	return fmt.Sprintf("id: \"%s\", title: \"%s\", comment: \"%s\", isDone: \"%v\", userId: \"%s\", createdAt: \"%s\", updatedAt: \"%s\", ShoppingListId: \"%s\"", i.Id, i.Title, i.Comment, i.IsDone, i.UserId, i.CreatedAt.Format(time.DateTime), i.UpdatedAt.Format(time.DateTime), i.ShoppingListId)
}

type UpdateShoppingListBody struct {
	Title     string    `json:"title"`
	UserId    string    `json:"user_id"`
	Items     []string  `json:"items"`
	State     State     `json:"state"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u UpdateShoppingListBody) String() string {
	return fmt.Sprintf("Title: \"%s\", UserId: \"%s\", Items: \"%v\", State: \"%s\", UpdatedAt: \"%s\"", u.Title, u.UserId, u.Items, u.State, u.UpdatedAt.Format(time.DateTime))
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
	Title     string    `json:"title"`
	UserId    string    `json:"user_id"`
	Items     []string  `json:"items"`
	State     State     `json:"state"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateItemRequest struct {
	Title          string    `json:"title"`
	Comment        string    `json:"comment"`
	IsDone         bool      `json:"is_done"`
	UserId         string    `json:"user_id"`
	ShoppingListId string    `json:"shopping_list_id"`
	UpdatedAt      time.Time `json:"updated_at"`
}

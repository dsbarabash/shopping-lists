package model

import (
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ShoppingList struct {
	Id        string                 `json:"id,omitempty"`
	Title     string                 `json:"title,omitempty"`
	UserId    string                 `json:"user_id,omitempty"`
	CreatedAt *timestamppb.Timestamp `json:"created_at,omitempty"`
	UpdatedAt *timestamppb.Timestamp `json:"updated_at"`
	Items     []string               `json:"items,omitempty"`
	State     State                  `json:"state,omitempty"`
}

func NewShoppingList(dto *CreateShoppingListDTO) *ShoppingList {
	return &ShoppingList{
		Id:        dto.Id,
		Title:     dto.Title,
		UserId:    dto.UserId,
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
		Items:     make([]string, 0),
		State:     1,
	}
}

func UpdateShoppingList(dto *UpdateShoppingListDTO) *ShoppingList {
	return &ShoppingList{
		Id:        dto.Id,
		Title:     dto.Title,
		UserId:    dto.UserId,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: timestamppb.Now(),
		Items:     dto.Items,
		State:     dto.State,
	}
}

type CreateShoppingListDTO struct {
	Id        string                 `json:"id,omitempty" bson:"id,omitempty"`
	Title     string                 `json:"title,omitempty" bson:"title,omitempty"`
	UserId    string                 `json:"user_id,omitempty" bson:"user_id,omitempty"`
	CreatedAt *timestamppb.Timestamp `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *timestamppb.Timestamp `json:"updated_at" bson:"updated_at"`
	Items     []string               `json:"items,omitempty" bson:"items,omitempty"`
	State     State                  `json:"state,omitempty" bson:"state,omitempty"`
}

type UpdateShoppingListDTO struct {
	Id        string                 `json:"id,omitempty" bson:"id,omitempty"`
	Title     string                 `json:"title,omitempty" bson:"title,omitempty"`
	UserId    string                 `json:"user_id,omitempty" bson:"user_id,omitempty"`
	CreatedAt *timestamppb.Timestamp `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *timestamppb.Timestamp `json:"updated_at" bson:"updated_at"`
	Items     []string               `json:"items,omitempty" bson:"items,omitempty"`
	State     State                  `json:"state,omitempty" bson:"state,omitempty"`
}

func (s ShoppingList) String() string {
	return fmt.Sprintf("id: \"%s\", title: \"%s\", userId: \"%s\", createdAt: \"%s\", updatedAt: \"%s\"", s.Id, s.Title, s.UserId, s.CreatedAt.String(), s.UpdatedAt.String())
}

type Item struct {
	Id             string                 `json:"id,omitempty"`
	Title          string                 `json:"title,omitempty"`
	Comment        string                 `json:"comment,omitempty"`
	IsDone         bool                   `json:"is_done,omitempty"`
	UserId         string                 `json:"user_id,omitempty"`
	CreatedAt      *timestamppb.Timestamp `json:"created_at,omitempty"`
	UpdatedAt      *timestamppb.Timestamp `json:"updated_at"`
	ShoppingListId string                 `json:"shopping_list_id,omitempty"`
}

type CreateItemDTO struct {
	Id             string                 `json:"id,omitempty" bson:"id,omitempty"`
	Title          string                 `json:"title,omitempty" bson:"title,omitempty"`
	Comment        string                 `json:"comment,omitempty" bson:"Comment,omitempty"`
	IsDone         bool                   `json:"is_done,omitempty" bson:"IsDone,omitempty"`
	UserId         string                 `json:"user_id,omitempty" bson:"user_id,omitempty"`
	CreatedAt      *timestamppb.Timestamp `json:"created_at" bson:"created_at"`
	UpdatedAt      *timestamppb.Timestamp `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	ShoppingListId string                 `json:"shopping_list_id,omitempty" bson:"shopping_list_id,omitempty"`
}

type UpdateItemDTO struct {
	Id             string                 `json:"id,omitempty" bson:"id,omitempty"`
	Title          string                 `json:"title,omitempty" bson:"title,omitempty"`
	Comment        string                 `json:"comment,omitempty"`
	IsDone         bool                   `json:"is_done,omitempty" bson:"IsDone,omitempty"`
	UserId         string                 `json:"user_id,omitempty" bson:"user_id,omitempty"`
	CreatedAt      *timestamppb.Timestamp `json:"created_at" bson:"created_at"`
	UpdatedAt      *timestamppb.Timestamp `json:"updated_at" bson:"updated_at"`
	ShoppingListId string                 `json:"shopping_list_id,omitempty" bson:"shopping_list_id,omitempty"`
}

func NewItem(dto *CreateItemDTO) *Item {
	return &Item{
		Id:             dto.Id,
		Title:          dto.Title,
		Comment:        dto.Comment,
		IsDone:         dto.IsDone,
		UserId:         dto.UserId,
		CreatedAt:      timestamppb.Now(),
		UpdatedAt:      timestamppb.Now(),
		ShoppingListId: dto.ShoppingListId,
	}
}

func UpdateItem(dto *UpdateItemDTO) *Item {
	return &Item{
		Id:             dto.Id,
		Title:          dto.Title,
		Comment:        dto.Comment,
		IsDone:         dto.IsDone,
		UserId:         dto.UserId,
		CreatedAt:      dto.CreatedAt,
		UpdatedAt:      timestamppb.Now(),
		ShoppingListId: dto.ShoppingListId,
	}
}

func (i Item) String() string {
	return fmt.Sprintf("id: \"%s\", title: \"%s\", comment: \"%s\", isDone: \"%v\", userId: \"%s\", createdAt: \"%s\", updatedAt: \"%s\", ShoppingListId: \"%s\"", i.Id, i.Title, i.Comment, i.IsDone, i.UserId, i.CreatedAt.String(), i.UpdatedAt.String(), i.ShoppingListId)
}

type UpdateShoppingListBody struct {
	Title     string                 `json:"title"`
	UserId    string                 `json:"user_id"`
	Items     []string               `json:"items"`
	State     State                  `json:"state"`
	UpdatedAt *timestamppb.Timestamp `json:"updated_at"`
}

func (u UpdateShoppingListBody) String() string {
	return fmt.Sprintf("Title: \"%s\", UserId: \"%s\", Items: \"%v\", State: \"%s\", UpdatedAt: \"%s\"", u.Title, u.UserId, u.Items, u.State, u.UpdatedAt.String())
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
	Title     string                 `json:"title"`
	UserId    string                 `json:"user_id"`
	Items     []string               `json:"items"`
	State     State                  `json:"state"`
	UpdatedAt *timestamppb.Timestamp `json:"updated_at"`
}

type UpdateItemRequest struct {
	Title          string                 `json:"title"`
	Comment        string                 `json:"comment"`
	IsDone         bool                   `json:"is_done"`
	UserId         string                 `json:"user_id"`
	ShoppingListId string                 `json:"shopping_list_id"`
	UpdatedAt      *timestamppb.Timestamp `json:"updated_at"`
}

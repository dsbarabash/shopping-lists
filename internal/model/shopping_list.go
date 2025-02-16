package model

import (
	"crypto/rand"
	"fmt"
	"log"
	"time"
)

type ShoppingLists interface {
	UpdateShoppingList(string, []string) ShoppingList
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

func NewShoppingList(title string, userId string) ShoppingList {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return ShoppingList{
		id:        uuid,
		title:     title,
		userId:    userId,
		createdAt: time.Now(),
		updatedAt: time.Now(),
		items:     make([]string, 0),
		state:     1,
	}
}

func (s ShoppingList) UpdateShoppingList(title string, items []string) ShoppingList {
	s.title = title
	s.updatedAt = time.Now()
	for _, i := range items {
		s.items = append(s.items, i)
	}
	return s
}

type Items interface {
	UpdateItem(string, string, bool) Item
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

func NewItem(title string, comment string, userId string, shoppingListId string) Item {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return Item{
		id:             uuid,
		title:          title,
		comment:        comment,
		isDone:         false,
		userId:         userId,
		createdAt:      time.Now(),
		updatedAt:      time.Now(),
		ShoppingListId: shoppingListId,
	}
}

func (i Item) UpdateItem(title string, comment string, isDone bool) Item {
	i.title = title
	i.comment = comment
	i.isDone = isDone
	i.updatedAt = time.Now()
	return i
}

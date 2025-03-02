package service

import "github.com/dsbarabash/shopping-lists/internal/model"

func CreateShoppingList(name string, userId string) *model.ShoppingList {
	ShoppingList, err := model.NewShoppingList(name, userId)
	if err != nil {
		panic(err)
	}
	return ShoppingList
}

func CreateItem(name string, comment string, userId string, shoppingListId string) *model.Item {
	Item, err := model.NewItem(name, comment, userId, shoppingListId)
	if err != nil {
		panic(err)
	}
	return Item
}

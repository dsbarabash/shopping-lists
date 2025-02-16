package service

import "github.com/dsbarabash/shopping-lists/internal/model"

func CreateShoppingList(name string, userId string) model.ShoppingList {
	ShoppingList := model.NewShoppingList(name, userId)
	return ShoppingList
}

func CreateItem(name string, comment string, userId string, shoppingListId string) model.Item {
	Item := model.NewItem(name, comment, shoppingListId, userId)
	return Item
}

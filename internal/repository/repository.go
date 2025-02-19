package repository

import (
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/model"
)

var ShoppingListSlice = make([]model.ShoppingList, 0)
var ItemSlice = make([]model.Item, 0)

func CheckInterface(arg interface{}) {
	switch arg.(type) {
	case model.ShoppingLists:
		ShoppingListSlice = append(ShoppingListSlice, *arg.(*model.ShoppingList))
	case model.Items:
		ItemSlice = append(ItemSlice, *arg.(*model.Item))
	default:
		fmt.Println("Неизвестный тип ")
	}
	fmt.Println("ShoppingList: ", ShoppingListSlice)
	fmt.Println("Item: ", ItemSlice)
}

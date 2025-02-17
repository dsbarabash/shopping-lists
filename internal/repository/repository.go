package repository

import (
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/model"
)

var ShoppingList = make([]model.ShoppingList, 0)
var Item = make([]model.Item, 0)

func CheckInterface(arg interface{}) {
	switch arg.(type) {
	case model.ShoppingLists:
		ShoppingList = append(ShoppingList, arg.(model.ShoppingList))
	case model.Items:
		Item = append(Item, arg.(model.Item))
	default:
		fmt.Println("Неизвестный тип ")
	}
	fmt.Println("ShoppingList: ", ShoppingList)
	fmt.Println("Item: ", Item)
}

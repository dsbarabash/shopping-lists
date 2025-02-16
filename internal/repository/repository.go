package repository

import (
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/model"
)

func CheckInterface(arg interface{}) {
	ShoppingList := make([]interface{}, 0)
	Item := make([]interface{}, 0)
	switch arg.(type) {
	case model.ShoppingLists:

		ShoppingList = append(ShoppingList, arg)
	case model.Items:
		Item = append(Item, arg)
	default:
		fmt.Println("Неизвестный тип ")
	}
	fmt.Println("ShoppingList: ", ShoppingList)
	fmt.Println("Item: ", Item)
}

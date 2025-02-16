package repository

import (
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/model"
)

var ShoppingList = make([]interface{}, 0)
var Item = make([]interface{}, 0)

func CheckInterface(arg interface{}) {
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

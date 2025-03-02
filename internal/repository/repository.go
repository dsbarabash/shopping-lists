package repository

import (
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"log"
)

var ShoppingListSlice = make([]model.ShoppingList, 0)
var ItemSlice = make([]model.Item, 0)
var lenSLSlice = len(ShoppingListSlice)
var lenISlice = len(ItemSlice)

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

func CheckInterface2(ch chan interface{}) {
	arg := <-ch
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

func LoggingSlice() {
	if len(ShoppingListSlice) != lenSLSlice {
		log.Println(ShoppingListSlice[len(ShoppingListSlice)-1:])
		lenSLSlice = len(ShoppingListSlice)
	}
	if len(ItemSlice) != lenISlice {
		log.Println(ItemSlice[len(ItemSlice)-1:])
		lenISlice = len(ItemSlice)
	}
}

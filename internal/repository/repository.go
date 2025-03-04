package repository

import (
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"log"
	"sync"
)

var ShoppingListSlice = make([]*model.ShoppingList, 0)
var ItemSlice = make([]*model.Item, 0)
var lenSLSlice = len(ShoppingListSlice)
var lenISlice = len(ItemSlice)

func CheckInterface(arg interface{}) {
	mu := sync.Mutex{}
	mu.Lock()
	switch arg.(type) {
	case model.ShoppingLists:
		ShoppingListSlice = append(ShoppingListSlice, arg.(*model.ShoppingList))
	case model.Items:
		ItemSlice = append(ItemSlice, arg.(*model.Item))
	default:
		fmt.Println("Неизвестный тип ")
	}
	mu.Unlock()
	fmt.Println("ShoppingList: ", ShoppingListSlice)
	fmt.Println("Item: ", ItemSlice)
}

func LoggingSlice() {
	mu := sync.Mutex{}
	if len(ShoppingListSlice) != lenSLSlice {
		mu.Lock()
		for i := lenSLSlice; i < len(ShoppingListSlice); i++ {
			log.Println(ShoppingListSlice[i])
		}
		lenSLSlice = len(ShoppingListSlice)
		mu.Unlock()
	}
	if len(ItemSlice) != lenISlice {
		mu.Lock()
		for i := lenISlice; i < len(ItemSlice); i++ {
			log.Println(ItemSlice[i])
		}
		lenISlice = len(ItemSlice)
		mu.Unlock()
	}
}

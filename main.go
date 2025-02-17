package main

import (
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"github.com/dsbarabash/shopping-lists/internal/service"
)

func main() {
	NewList := service.CreateShoppingList("test_list", "test_user")
	repository.CheckInterface(NewList)
	NewItem := service.CreateItem("test_item", "test_comment", "test_user", "")
	repository.CheckInterface(NewItem)
}

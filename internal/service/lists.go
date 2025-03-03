package service

import (
	"github.com/dsbarabash/shopping-lists/internal/model"
	"math/rand/v2"
	"strings"
)

func CreateRandomStructs() any {
	rnd := rand.IntN(100)
	var result any
	var err error
	if rnd%2 == 0 {
		result, err = model.NewShoppingList(randomSLData())

	} else {
		result, err = model.NewItem(randomIData())
	}
	if err != nil {
		panic(err)
	}
	return result
}

func randomSLData() (string, string) {
	return randomString(10), randomString(10)
}

func randomIData() (string, string, string, string) {
	return randomString(10), randomString(10), randomString(10), randomString(10)
}

func randomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.IntN(len(charset))])
	}
	return sb.String()
}

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

package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"github.com/dsbarabash/shopping-lists/internal/repository"
)

func CheckInterface(arg interface{}) {
	switch arg.(type) {
	case model.ShoppingLists:
		repository.ShoppingList.Add(arg.(*model.ShoppingList))
		repository.ShoppingList.PrintNewElement()
	case model.Items:
		repository.ItemList.Add(arg.(*model.Item))
		repository.ItemList.PrintNewElement()
	default:
		fmt.Println("Неизвестный тип ")
	}
}

func GetItems() string {
	repository.ItemList.Mu.Lock()
	defer repository.ItemList.Mu.Unlock()
	iString := ""
	for _, i := range repository.ItemList.Store {
		iString = iString + i.String()
	}
	return iString
}

func GetItemById(id string) (string, error) {
	repository.ItemList.Mu.Lock()
	defer repository.ItemList.Mu.Unlock()
	for _, i := range repository.ItemList.Store {
		if i.Id == id {
			return i.String(), nil
		}
	}
	return "", errors.New("NOT FOUND")
}

func DeleteItemById(id string) error {
	repository.ItemList.Mu.Lock()
	defer repository.ItemList.Mu.Unlock()
	for idx, i := range repository.ItemList.Store {
		if i.Id == id {
			copy(repository.ItemList.Store[idx:], repository.ItemList.Store[idx+1:])
			repository.ItemList.Store = repository.ItemList.Store[:len(repository.ItemList.Store)-1]
			repository.ItemList.SaveSliceToFile(repository.ItemList.Store)
			return nil
		}
	}
	return errors.New("NOT FOUND")
}

func GetSls() string {
	repository.ShoppingList.Mu.Lock()
	defer repository.ShoppingList.Mu.Unlock()
	slString := ""
	for _, l := range repository.ShoppingList.Store {
		slString = slString + l.String()
	}
	return slString
}

func GetSlById(id string) (string, error) {
	repository.ShoppingList.Mu.Lock()
	defer repository.ShoppingList.Mu.Unlock()
	for _, sl := range repository.ShoppingList.Store {
		if sl.Id == id {
			return sl.String(), nil
		}
	}
	return "", errors.New("NOT FOUND")
}

func DeleteSlById(id string) error {
	repository.ShoppingList.Mu.Lock()
	defer repository.ShoppingList.Mu.Unlock()
	for idx, sl := range repository.ShoppingList.Store {
		if sl.Id == id {
			copy(repository.ShoppingList.Store[idx:], repository.ShoppingList.Store[idx+1:])
			repository.ShoppingList.Store = repository.ShoppingList.Store[:len(repository.ShoppingList.Store)-1]
			repository.ShoppingList.SaveSliceToFile(repository.ShoppingList.Store)
			return nil
		}
	}
	return errors.New("NOT FOUND")
}

func UpdateSl(id string, body []byte) error {
	repository.ShoppingList.Mu.Lock()
	defer repository.ShoppingList.Mu.Unlock()
	for _, sl := range repository.ShoppingList.Store {
		if sl.Id == id {
			err := json.Unmarshal(body, &sl)
			if err != nil {
				return err
			}
			repository.ShoppingList.SaveSliceToFile(repository.ShoppingList.Store)
			return nil
		}
	}
	return errors.New("NOT FOUND")
}

func UpdateItem(id string, body []byte) error {
	repository.ItemList.Mu.Lock()
	defer repository.ItemList.Mu.Unlock()
	for _, item := range repository.ItemList.Store {
		if item.Id == id {
			err := json.Unmarshal(body, &item)
			if err != nil {
				return err
			}
			repository.ItemList.SaveSliceToFile(repository.ItemList.Store)
			return nil
		}
	}
	return errors.New("NOT FOUND")
}

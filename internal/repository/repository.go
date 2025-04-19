package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"github.com/dsbarabash/shopping-lists/internal/service"
	"io"
	"log"
	"os"
	"sync"
)

type ItemStore struct {
	Mu                  sync.Mutex
	Store               []*model.Item
	PrintedElementCount int
	FilePath            string
}

type ShoppingListStore struct {
	Mu                  sync.Mutex
	Store               []*model.ShoppingList
	PrintedElementCount int
	FilePath            string
}

var ShoppingList = ShoppingListStore{
	sync.Mutex{},
	make([]*model.ShoppingList, 0),
	0,
	"./shoppingLists.json",
}

var ItemList = ItemStore{
	sync.Mutex{},
	make([]*model.Item, 0),
	0,
	"./items.json",
}

func (i *ItemStore) LoadFromFile() {
	i.Mu.Lock()
	defer i.Mu.Unlock()
	items, err := ReadJson(i.FilePath)
	if err == io.EOF {
		return
	} else if err != nil {
		log.Fatal(err)
	}
	if len(items) != 0 {
		if err := json.Unmarshal(items, &i.Store); err != nil {
			log.Fatal(err)
		}
	}
}

func (i *ItemStore) SaveToFile(item *model.Item) {
	hu := i.FilePath
	f, err := os.OpenFile(hu, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	// Перемещаемся в конец файла
	stat, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	// Если это начала файла, начинаем массив json
	if stat.Size() == 0 {
		_, err := f.WriteString("[")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Если файл не пуст, то заменчем последнюю закрывающую скобку массива на запятую
		if _, err := f.Seek(-1, io.SeekEnd); err != nil {
			log.Fatal(err)
		}
		if _, err := f.WriteString(","); err != nil {
			log.Fatal(err)
		}
	}
	// Добавляем объект в файл и закрываем массив
	e := json.NewEncoder(f)
	if err := e.Encode(item); err != nil {
		log.Fatal(err)
	}
	// Добавляем закрывающую скобку
	_, err = f.WriteString("]")
	if err != nil {
		log.Fatal(err)
	}
}

func (i *ItemStore) SaveSliceToFile(sls []*model.Item) {
	f, err := os.OpenFile(i.FilePath, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	e := json.NewEncoder(f)
	if err := e.Encode(sls); err != nil {
		log.Fatal(err)
	}
}

func (i *ItemStore) Add(item *model.Item) {
	i.Mu.Lock()
	defer i.Mu.Unlock()
	i.Store = append(i.Store, item)
	i.SaveToFile(item)
}

func (i *ItemStore) PrintNewElement() {
	i.Mu.Lock()
	defer i.Mu.Unlock()
	if len(i.Store) > i.PrintedElementCount {
		for j := i.PrintedElementCount; j < len(i.Store); j++ {
			log.Println(i.Store[j])
		}
		i.PrintedElementCount = len(i.Store)
	}
}

func (s *ShoppingListStore) LoadFromFile() {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	items, err := ReadJson(s.FilePath)
	if err == io.EOF {
		return
	} else if err != nil {
		log.Fatal(err)
	}
	if len(items) != 0 {
		if err := json.Unmarshal(items, &s.Store); err != nil {
			log.Fatal(err)
		}
	}
}

func (s *ShoppingListStore) SaveToFile(sl *model.ShoppingList) {
	f, err := os.OpenFile(s.FilePath, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	// Перемещаемся в конец файла
	stat, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	// Если это начала файла, начинаем массив json
	if stat.Size() == 0 {
		_, err := f.WriteString("[")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Если файл не пуст, то заменчем последнюю закрывающую скобку массива на запятую
		if _, err := f.Seek(-1, io.SeekEnd); err != nil {
			log.Fatal(err)
		}
		if _, err := f.WriteString(","); err != nil {
			log.Fatal(err)
		}
	}
	// Добавляем объект в файл и закрываем массив
	e := json.NewEncoder(f)
	if err := e.Encode(sl); err != nil {
		log.Fatal(err)
	}
	// Добавляем закрывающую скобку
	_, err = f.WriteString("]")
	if err != nil {
		log.Fatal(err)
	}
}

func (s *ShoppingListStore) SaveSliceToFile(sls []*model.ShoppingList) {
	f, err := os.OpenFile(s.FilePath, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	e := json.NewEncoder(f)
	if err := e.Encode(sls); err != nil {
		log.Fatal(err)
	}
}

func (s *ShoppingListStore) Add(sl *model.ShoppingList) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Store = append(s.Store, sl)
	s.SaveToFile(sl)
}

func (s *ShoppingListStore) PrintNewElement() {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if len(s.Store) > s.PrintedElementCount {
		for i := s.PrintedElementCount; i < len(s.Store); i++ {
			log.Println(s.Store[i])
		}
		s.PrintedElementCount = len(s.Store)
	}
}

func CheckInterface(arg interface{}) {
	switch arg.(type) {
	case model.ShoppingLists:
		ShoppingList.Add(arg.(*model.ShoppingList))
		ShoppingList.PrintNewElement()
	case model.Items:
		ItemList.Add(arg.(*model.Item))
		ItemList.PrintNewElement()
	default:
		fmt.Println("Неизвестный тип ")
	}
}

func ReadJson(fileName string) ([]byte, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func FillSlices() {
	ItemList.LoadFromFile()
	ShoppingList.LoadFromFile()
	service.UserList.LoadFromFile()
}

func GetItems() string {
	ItemList.Mu.Lock()
	defer ItemList.Mu.Unlock()
	iString := ""
	for _, i := range ItemList.Store {
		iString = iString + i.String()
	}
	return iString
}

func GetItemById(id string) (string, error) {
	ItemList.Mu.Lock()
	defer ItemList.Mu.Unlock()
	for _, i := range ItemList.Store {
		if i.Id == id {
			return i.String(), nil
		}
	}
	return "", errors.New("NOT FOUND")
}

func DeleteItemById(id string) error {
	ItemList.Mu.Lock()
	defer ItemList.Mu.Unlock()
	for idx, i := range ItemList.Store {
		if i.Id == id {
			copy(ItemList.Store[idx:], ItemList.Store[idx+1:])
			ItemList.Store = ItemList.Store[:len(ItemList.Store)-1]
			ItemList.SaveSliceToFile(ItemList.Store)
			return nil
		}
	}
	return errors.New("NOT FOUND")
}

func GetSls() string {
	ShoppingList.Mu.Lock()
	defer ShoppingList.Mu.Unlock()
	slString := ""
	for _, l := range ShoppingList.Store {
		slString = slString + l.String()
	}
	return slString
}

func GetSlById(id string) (string, error) {
	ShoppingList.Mu.Lock()
	defer ShoppingList.Mu.Unlock()
	for _, sl := range ShoppingList.Store {
		if sl.Id == id {
			return sl.String(), nil
		}
	}
	return "", errors.New("NOT FOUND")
}

func DeleteSlById(id string) error {
	ShoppingList.Mu.Lock()
	defer ShoppingList.Mu.Unlock()
	for idx, sl := range ShoppingList.Store {
		if sl.Id == id {
			copy(ShoppingList.Store[idx:], ShoppingList.Store[idx+1:])
			ShoppingList.Store = ShoppingList.Store[:len(ShoppingList.Store)-1]
			ShoppingList.SaveSliceToFile(ShoppingList.Store)
			return nil
		}
	}
	return errors.New("NOT FOUND")
}

func UpdateSl(id string, body []byte) error {
	ShoppingList.Mu.Lock()
	defer ShoppingList.Mu.Unlock()
	for _, sl := range ShoppingList.Store {
		if sl.Id == id {
			err := json.Unmarshal(body, &sl)
			if err != nil {
				return err
			}
			ShoppingList.SaveSliceToFile(ShoppingList.Store)
			return nil
		}
	}
	return errors.New("NOT FOUND")
}

func UpdateItem(id string, body []byte) error {
	ItemList.Mu.Lock()
	defer ItemList.Mu.Unlock()
	for _, item := range ItemList.Store {
		if item.Id == id {
			err := json.Unmarshal(body, &item)
			if err != nil {
				return err
			}
			ItemList.SaveSliceToFile(ItemList.Store)
			return nil
		}
	}
	return errors.New("NOT FOUND")
}

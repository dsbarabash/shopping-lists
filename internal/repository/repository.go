package repository

import (
	"encoding/json"
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"io"
	"log"
	"os"
	"sync"
)

type ItemStore struct {
	mu                  sync.Mutex
	store               []*model.Item
	printedElementCount int
	filePath            string
}

type ShoppingListStore struct {
	mu                  sync.Mutex
	store               []*model.ShoppingList
	printedElementCount int
	filePath            string
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
	i.mu.Lock()
	defer i.mu.Unlock()
	items, err := readJson(i.filePath)
	if err == io.EOF {
		return
	} else if err != nil {
		log.Fatal(err)
	}
	if len(items) != 0 {
		if err := json.Unmarshal(items, &i.store); err != nil {
			log.Fatal(err)
		}
	}
}

func (i *ItemStore) SaveToFile(item *model.Item) {
	hu := i.filePath
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

func (i *ItemStore) Add(item *model.Item) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.store = append(i.store, item)
	i.SaveToFile(item)

}

func (i *ItemStore) PrintNewElement() {
	i.mu.Lock()
	defer i.mu.Unlock()
	if len(i.store) > i.printedElementCount {
		for j := i.printedElementCount; j < len(i.store); j++ {
			log.Println(i.store[j])
		}
		i.printedElementCount = len(i.store)
	}
}

func (s *ShoppingListStore) LoadFromFile() {
	s.mu.Lock()
	defer s.mu.Unlock()
	items, err := readJson(s.filePath)
	if err == io.EOF {
		return
	} else if err != nil {
		log.Fatal(err)
	}
	if len(items) != 0 {
		if err := json.Unmarshal(items, &s.store); err != nil {
			log.Fatal(err)
		}
	}
}

func (s *ShoppingListStore) SaveToFile(sl *model.ShoppingList) {
	f, err := os.OpenFile(s.filePath, os.O_RDWR, 0644)
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

func (s *ShoppingListStore) Add(sl *model.ShoppingList) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store = append(s.store, sl)
	s.SaveToFile(sl)

}

func (s *ShoppingListStore) PrintNewElement() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.store) > s.printedElementCount {
		for i := s.printedElementCount; i < len(s.store); i++ {
			log.Println(s.store[i])
		}
		s.printedElementCount = len(s.store)
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

func readJson(fileName string) ([]byte, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func FillSlices() {
	ItemList.LoadFromFile()
	ShoppingList.LoadFromFile()
}

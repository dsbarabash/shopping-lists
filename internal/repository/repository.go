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

var ShoppingListSlice = make([]*model.ShoppingList, 0)
var ItemSlice = make([]*model.Item, 0)
var lenSLSlice = len(ShoppingListSlice)
var lenISlice = len(ItemSlice)
var mu = sync.Mutex{}

func CheckInterface(arg interface{}) {
	mu.Lock()
	switch arg.(type) {
	case model.ShoppingLists:
		ShoppingListSlice = append(ShoppingListSlice, arg.(*model.ShoppingList))
		addItemToFile(arg)
	case model.Items:
		ItemSlice = append(ItemSlice, arg.(*model.Item))
		addItemToFile(arg)
	default:
		fmt.Println("Неизвестный тип ")
	}
	fmt.Println("ShoppingList: ", ShoppingListSlice)
	fmt.Println("Item: ", ItemSlice)
	mu.Unlock()
}

func LoggingSlice() {
	mu.Lock()
	if len(ShoppingListSlice) != lenSLSlice {
		for i := lenSLSlice; i < len(ShoppingListSlice); i++ {
			log.Println(ShoppingListSlice[i])
		}
		lenSLSlice = len(ShoppingListSlice)
	}
	if len(ItemSlice) != lenISlice {
		for i := lenISlice; i < len(ItemSlice); i++ {
			log.Println(ItemSlice[i])
		}
		lenISlice = len(ItemSlice)
	}
	mu.Unlock()
}

func FillSlices() {
	mu.Lock()
	lists, err := readJson("./shoppingLists.json")
	if err == io.EOF {
		return
	} else if err != nil {
		log.Fatal(err)
	}

	if len(lists) != 0 {
		if err := json.Unmarshal(lists, &ShoppingListSlice); err != nil {
			log.Fatal(err)
		}
	}
	items, err := readJson("./items.json")
	if err == io.EOF {
		return
	} else if err != nil {
		log.Fatal(err)
	}

	if len(items) != 0 {
		if err := json.Unmarshal(items, &ItemSlice); err != nil {
			log.Fatal(err)
		}
	}
	mu.Unlock()
}

func readJson(fileName string) ([]byte, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func addItemToFile(arg interface{}) {
	switch arg.(type) {
	case model.ShoppingLists:
		log.Println("addSLToFile start")
		f, err := os.OpenFile("./shoppingLists.json", os.O_RDWR, 0644)
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
			if _, err := f.Seek(-2, io.SeekEnd); err != nil {
				log.Fatal(err)
			}
			if _, err := f.WriteString(","); err != nil {
				log.Fatal(err)
			}
		}
		// Добавляем объект в файл и закрываем массив
		e := json.NewEncoder(f)
		if err := e.Encode(arg.(*model.ShoppingList)); err != nil {
			log.Fatal(err)
		}
		// Добавляем закрывающую скобку
		_, err = f.WriteString("]")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("addSLToFile end")

	case model.Items:
		log.Println("addItemToFile start")
		f, err := os.OpenFile("./items.json", os.O_RDWR, 0644)
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
			if _, err := f.Seek(-2, io.SeekEnd); err != nil {
				log.Fatal(err)
			}
			if _, err := f.WriteString(","); err != nil {
				log.Fatal(err)
			}
		}
		// Добавляем объект в файл и закрываем массив
		e := json.NewEncoder(f)
		if err := e.Encode(arg.(*model.Item)); err != nil {
			log.Fatal(err)
		}
		// Добавляем закрывающую скобку
		_, err = f.WriteString("]")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("addItemToFile end")
	default:
		log.Fatal("Неизвестный тип.")
	}
}

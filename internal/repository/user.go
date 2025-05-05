package repository

import (
	"encoding/json"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"github.com/dsbarabash/shopping-lists/internal/repository/mongo"
	"io"
	"log"
	"os"
	"sync"
)

type UserStore struct {
	Mu       sync.Mutex
	Store    []*model.User
	FilePath string
}

var UserList = UserStore{
	sync.Mutex{},
	make([]*model.User, 0),
	"./users.json",
}

func (u *UserStore) LoadFromFile() {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	items, err := mongo.ReadJson(u.FilePath)
	if err == io.EOF {
		return
	} else if err != nil {
		log.Fatal(err)
	}
	if len(items) != 0 {
		if err := json.Unmarshal(items, &u.Store); err != nil {
			log.Fatal(err)
		}
	}
}

func (u *UserStore) SaveToFile(sl *model.User) {
	f, err := os.OpenFile(u.FilePath, os.O_RDWR, 0644)
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
		// Если файл не пуст, то заменяем последнюю закрывающую скобку массива на запятую
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

package service

import (
	"encoding/json"
	"errors"
	"github.com/dsbarabash/shopping-lists/internal/config"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"log"
	"os"
	"sync"
)

type UserStore struct {
	mu       sync.Mutex
	store    []*model.User
	filePath string
}

var UserList = UserStore{
	sync.Mutex{},
	make([]*model.User, 0),
	"./users.json",
}

func Registration(name, password string) *model.User {
	var user *model.User
	user = model.NewUser(name, password)
	UserList.SaveToFile(user)
	return user
}

func Login(user *model.User) (string, error) {
	UserList.mu.Lock()
	defer UserList.mu.Unlock()
	if user.State != 1 {
		return "", errors.New("USER NOT ACTIVE")
	}
	for _, u := range UserList.store {
		if u.Id == user.Id {
			if u.Password == user.Password && u.Name == user.Name {
				secretKey := []byte(config.Cfg.Secret)
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"id":       user.Id,
					"name":     user.Name,
					"password": user.Password,
					"state":    1,
				})
				tokenString, err := token.SignedString(secretKey)
				if err != nil {
					return "", err
				}
				return tokenString, nil
			} else {
				return "", errors.New("LOGIN OR PASSWORD IS INCORRECT")
			}

		}
	}
	return "", errors.New("NOT FOUND")

}

func (u *UserStore) LoadFromFile() {
	u.mu.Lock()
	defer u.mu.Unlock()
	items, err := ReadJson(u.filePath)
	if err == io.EOF {
		return
	} else if err != nil {
		log.Fatal(err)
	}
	if len(items) != 0 {
		if err := json.Unmarshal(items, &u.store); err != nil {
			log.Fatal(err)
		}
	}
}

func (u *UserStore) SaveToFile(sl *model.User) {
	f, err := os.OpenFile(u.filePath, os.O_RDWR, 0644)
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

func ReadJson(fileName string) ([]byte, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return data, nil
}

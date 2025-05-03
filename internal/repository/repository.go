package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dsbarabash/shopping-lists/internal/config"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type MongoDb struct {
	ShoppingListCollection *mongo.Collection
	ItemCollection         *mongo.Collection
	UserCollection         *mongo.Collection
}

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

func (m *MongoDb) AddItem(ctx context.Context, item *model.Item) {
	_, err := m.ItemCollection.InsertOne(ctx, item)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *MongoDb) AddShoppingList(ctx context.Context, sl *model.ShoppingList) {
	_, err := m.ItemCollection.InsertOne(ctx, sl)
	if err != nil {
		log.Fatal(err)
	}
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

func ReadJson(fileName string) ([]byte, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func FillSlices() {
	//ItemList.LoadFromFile()
	//ShoppingList.LoadFromFile()
	UserList.LoadFromFile()
}

func ConnectRedisDb() (*redis.Client, error) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Адрес и порт Redis-сервера
		Password: "",               // Пароль (если есть)
		DB:       0,                // Номер базы данных
	})

	// Проверка соединения
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func ConnectMongoDb() (*MongoDb, error) {
	ctx := context.Background()
	// Подключение к MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Пинг сервера для проверки соединения
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Создание или переключение на базу данных
	dbName := "shopping_lists_db"
	db := client.Database(dbName)

	// Создание коллекции

	shoppingListCollection := db.Collection("lists")

	itemCollection := db.Collection("items")

	userCollection := db.Collection("users")

	return &MongoDb{
		ShoppingListCollection: shoppingListCollection,
		ItemCollection:         itemCollection,
		UserCollection:         userCollection,
	}, nil
}

func (m *MongoDb) GetItems(ctx context.Context) []model.Item {
	items, err := m.ItemCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var ls []model.Item
	err = items.All(ctx, &ls)
	if err != nil {
		log.Fatal(err)
	}
	return ls
}

func (m *MongoDb) GetItemById(ctx context.Context, id string) ([]model.Item, error) {
	items, err := m.ItemCollection.Find(ctx, bson.D{{"id", id}})
	if err != nil {
		log.Fatal(err)
	}
	var ls []model.Item
	err = items.All(ctx, &ls)
	if err != nil {
		log.Fatal(err)
	}
	if len(ls) == 0 {
		return nil, errors.New("NOT FOUND")
	}
	return ls, nil
}

func (m *MongoDb) DeleteItemById(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	items, err := m.ItemCollection.Find(ctx, bson.D{{"id", id}})
	if err != nil {
		log.Fatal(err)
	}
	var ls []model.Item
	err = items.All(ctx, &ls)
	if err != nil {
		log.Fatal(err)
	}
	if len(ls) == 0 {
		return nil, errors.New("NOT FOUND")
	}
	res, err := m.ItemCollection.DeleteOne(ctx, bson.D{{"id", id}})
	if err != nil {
		log.Fatal(err)
	}
	return res, nil
}

func (m *MongoDb) GetSls(ctx context.Context) []model.ShoppingList {
	lists, err := m.ShoppingListCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var ls []model.ShoppingList
	err = lists.All(ctx, &ls)
	if err != nil {
		log.Fatal(err)
	}
	return ls
}

func (m *MongoDb) GetSlById(ctx context.Context, id string) ([]model.ShoppingList, error) {
	items, err := m.ShoppingListCollection.Find(ctx, bson.D{{"id", id}})
	if err != nil {
		log.Fatal(err)
	}
	var ls []model.ShoppingList
	err = items.All(ctx, &ls)
	if err != nil {
		log.Fatal(err)
	}
	if len(ls) == 0 {
		return nil, errors.New("NOT FOUND")
	}
	return ls, nil
}

func (m *MongoDb) DeleteSlById(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	lists, err := m.ShoppingListCollection.Find(ctx, bson.D{{"id", id}})
	if err != nil {
		log.Fatal(err)
	}
	var ls []model.ShoppingList
	err = lists.All(ctx, &ls)
	if err != nil {
		log.Fatal(err)
	}
	if len(ls) == 0 {
		return nil, errors.New("NOT FOUND")
	}
	res, err := m.ShoppingListCollection.DeleteOne(ctx, bson.D{{"id", id}})
	if err != nil {
		log.Fatal(err)
	}
	return res, nil
}

func (m *MongoDb) UpdateSl(ctx context.Context, id string, body []byte) (*mongo.UpdateResult, error) {
	lists, err := m.ShoppingListCollection.Find(ctx, bson.D{{"id", id}})
	if err != nil {
		log.Fatal(err)
	}
	var ls []model.ShoppingList
	err = lists.All(ctx, &ls)
	if err != nil {
		log.Fatal(err)
	}
	if len(ls) == 0 {
		return nil, errors.New("NOT FOUND")
	}
	var sl model.UpdateShoppingListRequest
	sl.UpdatedAt = time.Now().UTC()
	err = json.Unmarshal(body, &sl)
	if err != nil {
		return nil, errors.New("ERROR TO UNMARSHALL")
	}
	update := bson.D{
		{"$set", sl},
	}
	res, err := m.ShoppingListCollection.UpdateOne(ctx, bson.D{{"id", id}}, update)
	if err != nil {
		log.Fatal(err)
	}
	return res, nil
}

func (m *MongoDb) UpdateItem(ctx context.Context, id string, body []byte) (*mongo.UpdateResult, error) {
	items, err := m.ItemCollection.Find(ctx, bson.D{{"id", id}})
	if err != nil {
		log.Fatal(err)
	}
	var ls []model.ShoppingList
	err = items.All(ctx, &ls)
	if err != nil {
		log.Fatal(err)
	}
	if len(ls) == 0 {
		return nil, errors.New("NOT FOUND")
	}
	var item model.UpdateItemRequest
	item.UpdatedAt = time.Now().UTC()
	err = json.Unmarshal(body, &item)
	if err != nil {
		return nil, errors.New("ERROR TO UNMARSHALL")
	}
	update := bson.D{
		{"$set", item},
	}
	res, err := m.ItemCollection.UpdateOne(ctx, bson.D{{"id", id}}, update)
	if err != nil {
		log.Fatal(err)
	}
	return res, nil
}

func (m *MongoDb) Registration(ctx context.Context, name, password string) *model.User {
	var user *model.User
	user = model.NewUser(name, password)
	_, err := m.UserCollection.InsertOne(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
	return user
}

func (m *MongoDb) Login(ctx context.Context, user *model.User) (string, error) {
	if user.State != 1 {
		return "", errors.New("USER NOT ACTIVE")
	}
	users, err := m.UserCollection.Find(ctx, bson.D{{"id", user.Id}})
	if err != nil {
		log.Fatal(err)
	}
	var u []model.User
	err = users.All(ctx, &u)
	if err != nil {
		log.Fatal(err)
	}
	if len(u) == 0 {
		return "", errors.New("NOT FOUND")
	} else {
		if u[0].State != 1 {
			return "", errors.New("USER NOT ACTIVE")
		} else if u[0].Password == user.Password && u[0].Name == user.Name {
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
		}
	}
	return "", errors.New("NOT FOUND")

}

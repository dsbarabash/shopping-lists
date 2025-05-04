package repository

import (
	"context"
	"errors"
	"github.com/dsbarabash/shopping-lists/internal/config"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

type MongoDb struct {
	ShoppingListCollection *mongo.Collection
	ItemCollection         *mongo.Collection
	UserCollection         *mongo.Collection
}

func (m *MongoDb) AddItem(ctx context.Context, item *model.Item) {
	_, err := m.ItemCollection.InsertOne(ctx, item)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *MongoDb) AddShoppingList(ctx context.Context, sl *model.ShoppingList) {
	_, err := m.ShoppingListCollection.InsertOne(ctx, sl)
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
	items, err := m.ItemCollection.Find(ctx, bson.D{primitive.E{}})
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
	items, err := m.ItemCollection.Find(ctx, bson.D{primitive.E{Key: "id", Value: id}})
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
	items, err := m.ItemCollection.Find(ctx, bson.D{primitive.E{Key: "id", Value: id}})
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
	res, err := m.ItemCollection.DeleteOne(ctx, bson.D{primitive.E{Key: "id", Value: id}})
	if err != nil {
		log.Fatal(err)
	}
	return res, nil
}

func (m *MongoDb) GetSls(ctx context.Context) []model.ShoppingList {
	lists, err := m.ShoppingListCollection.Find(ctx, bson.D{primitive.E{}})
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
	items, err := m.ShoppingListCollection.Find(ctx, bson.D{primitive.E{Key: "id", Value: id}})
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
	lists, err := m.ShoppingListCollection.Find(ctx, bson.D{primitive.E{Key: "id", Value: id}})
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
	res, err := m.ShoppingListCollection.DeleteOne(ctx, bson.D{primitive.E{Key: "id", Value: id}})
	if err != nil {
		log.Fatal(err)
	}
	return res, nil
}

func (m *MongoDb) UpdateSl(ctx context.Context, id string, sl model.UpdateShoppingListRequest) (*mongo.UpdateResult, error) {
	lists, err := m.ShoppingListCollection.Find(ctx, bson.D{primitive.E{Key: "id", Value: id}})
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
	update := bson.D{
		primitive.E{Key: "$set", Value: sl},
	}
	res, err := m.ShoppingListCollection.UpdateOne(ctx, bson.D{primitive.E{Key: "id", Value: id}}, update)
	if err != nil {
		log.Fatal(err)
	}
	return res, nil
}

func (m *MongoDb) FindItem(ctx context.Context, id string) (*model.Item, error) {
	items, err := m.ItemCollection.Find(ctx, bson.D{primitive.E{Key: "id", Value: id}})
	if err != nil {
		log.Fatal(err)
	}
	var it []model.Item
	err = items.All(ctx, &it)
	if len(it) == 0 {
		return nil, errors.New("NOT FOUND")
	}
	return &it[0], nil
}

func (m *MongoDb) UpdateItem(ctx context.Context, id string, item model.UpdateItemRequest) (*mongo.UpdateResult, error) {
	update := bson.D{
		primitive.E{Key: "$set", Value: item},
	}
	res, err := m.ItemCollection.UpdateOne(ctx, bson.D{primitive.E{Key: "id", Value: id}}, update)
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
	users, err := m.UserCollection.Find(ctx, bson.D{primitive.E{Key: "id", Value: user.Id}})
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

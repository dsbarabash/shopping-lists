package mongo

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
)

type MongoDb struct {
	ShoppingListCollection *mongo.Collection
	ItemCollection         *mongo.Collection
	UserCollection         *mongo.Collection
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
	log.Println("Connected to redis")
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
	log.Println("Connected to mongodb")

	return &MongoDb{
		ShoppingListCollection: shoppingListCollection,
		ItemCollection:         itemCollection,
		UserCollection:         userCollection,
	}, nil
}

func (m *MongoDb) GetSlById(ctx context.Context, id string) (*model.ShoppingList, error) {
	items, err := m.ShoppingListCollection.Find(ctx, bson.D{primitive.E{Key: "id", Value: id}})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var ls []*model.ShoppingList
	err = items.All(ctx, &ls)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(ls) == 0 {
		return nil, errors.New("NOT FOUND")
	}
	log.Println("Get shopping list: ", ls[0])
	return ls[0], nil
}

func (m *MongoDb) AddShoppingList(ctx context.Context, sl *model.ShoppingList) error {
	_, err := m.ShoppingListCollection.InsertOne(ctx, sl)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Inserted shopping list: ", sl.String())
	return nil
}

func (m *MongoDb) GetSls(ctx context.Context) ([]*model.ShoppingList, error) {
	lists, err := m.ShoppingListCollection.Find(ctx, bson.D{primitive.E{}})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var ls []*model.ShoppingList
	err = lists.All(ctx, &ls)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("Get shopping lists: ", ls)
	return ls, nil
}

func (m *MongoDb) UpdateSl(ctx context.Context, id string, sl *model.ShoppingList) error {
	lists, err := m.ShoppingListCollection.Find(ctx, bson.D{primitive.E{Key: "id", Value: id}})
	if err != nil {
		log.Println(err)
		return err
	}
	var ls []model.ShoppingList
	err = lists.All(ctx, &ls)
	if err != nil {
		log.Println(err)
		return err
	}
	if len(ls) == 0 {
		return errors.New("NOT FOUND")
	}
	update := bson.D{
		primitive.E{Key: "$set", Value: sl},
	}
	_, err = m.ShoppingListCollection.UpdateOne(ctx, bson.D{primitive.E{Key: "id", Value: id}}, update)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Update shopping list: ", ls[0])
	return nil
}

func (m *MongoDb) DeleteSlById(ctx context.Context, id string) error {
	lists, err := m.ShoppingListCollection.Find(ctx, bson.D{primitive.E{Key: "id", Value: id}})
	if err != nil {
		log.Println(err)
		return err
	}
	var ls []model.ShoppingList
	err = lists.All(ctx, &ls)
	if err != nil {
		log.Println(err)
		return err
	}
	if len(ls) == 0 {
		return errors.New("NOT FOUND")
	}
	_, err = m.ShoppingListCollection.DeleteOne(ctx, bson.D{primitive.E{Key: "id", Value: id}})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Delete shopping list: ", ls[0])
	return nil
}

func (m *MongoDb) AddItem(ctx context.Context, item *model.Item) error {
	_, err := m.ItemCollection.InsertOne(ctx, item)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Inserted item: ", item.String())
	return nil
}

func (m *MongoDb) GetItemById(ctx context.Context, id string) (*model.Item, error) {
	items, err := m.ItemCollection.Find(ctx, bson.D{primitive.E{Key: "id", Value: id}})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var ls []*model.Item
	err = items.All(ctx, &ls)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(ls) == 0 {
		return nil, errors.New("NOT FOUND")
	}
	log.Println("Get item: ", ls[0])
	return ls[0], nil
}

func (m *MongoDb) GetItems(ctx context.Context) ([]*model.Item, error) {
	items, err := m.ItemCollection.Find(ctx, bson.D{primitive.E{}})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var ls []*model.Item
	err = items.All(ctx, &ls)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("Get items: ", ls)
	return ls, nil
}

func (m *MongoDb) UpdateItem(ctx context.Context, id string, item *model.Item) error {
	items, err := m.ItemCollection.Find(ctx, bson.D{primitive.E{Key: "id", Value: id}})
	if err != nil {
		log.Println(err)
		return err
	}
	var it []model.Item
	err = items.All(ctx, &it)
	if err != nil {
		log.Println(err)
		return err
	}
	if len(it) == 0 {
		return errors.New("NOT FOUND")
	}
	update := bson.D{
		primitive.E{Key: "$set", Value: item},
	}
	_, err = m.ItemCollection.UpdateOne(ctx, bson.D{primitive.E{Key: "id", Value: id}}, update)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Update item: ", it[0])
	return nil
}

func (m *MongoDb) DeleteItemById(ctx context.Context, id string) error {
	items, err := m.ItemCollection.Find(ctx, bson.D{primitive.E{Key: "id", Value: id}})
	if err != nil {
		log.Println(err)
		return err
	}
	var ls []model.Item
	err = items.All(ctx, &ls)
	if err != nil {
		log.Println(err)
		return err
	}
	if len(ls) == 0 {
		return errors.New("NOT FOUND")
	}
	_, err = m.ItemCollection.DeleteOne(ctx, bson.D{primitive.E{Key: "id", Value: id}})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Delete item: ", ls[0])
	return nil
}

func (m *MongoDb) Registration(ctx context.Context, name, password string) (*model.User, error) {
	var user *model.User
	user = model.NewUser(name, password)
	_, err := m.UserCollection.InsertOne(ctx, user)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("Inserted user: ", user.Name)
	return user, nil
}

func (m *MongoDb) Login(ctx context.Context, user *model.User) (string, error) {
	if user.State != 1 {
		return "", errors.New("USER NOT ACTIVE")
	}
	users, err := m.UserCollection.Find(ctx, bson.D{primitive.E{Key: "id", Value: user.Id}})
	if err != nil {
		log.Println(err)
		return "", err
	}
	var u []model.User
	err = users.All(ctx, &u)
	if err != nil {
		log.Println(err)
		return "", err
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
	log.Println("Login user: ", user.Name)
	return "", errors.New("NOT FOUND")
}

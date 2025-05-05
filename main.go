package main

import (
	"context"
	"github.com/dsbarabash/shopping-lists/internal/app"
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"github.com/dsbarabash/shopping-lists/internal/repository/mongo"
	"log"
)

func main() {
	MongoDb, err := mongo.ConnectMongoDb()
	if err != nil {
		log.Fatal(err)
	}
	RedisDB, err := mongo.ConnectRedisDb()
	if err != nil {
		log.Fatal(err)
	}
	logWriter := repository.NewLogWriter(RedisDB)
	log.SetOutput(logWriter)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	newApp, err := app.NewService(ctx, MongoDb)
	if err != nil {
		log.Fatal(err)
	}
	err = newApp.Start()
	if err != nil {
		log.Fatal(err)
	}
}

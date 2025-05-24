package main

import (
	"context"
	"github.com/dsbarabash/shopping-lists/internal/app"
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"github.com/dsbarabash/shopping-lists/internal/repository/postgres"
	"github.com/dsbarabash/shopping-lists/internal/repository/redis"
	"github.com/dsbarabash/shopping-lists/internal/service"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//MongoDb, err := mongo.ConnectMongoDb()
	//if err != nil {
	//	log.Fatal(err)
	//}
	RedisDB, err := redis.ConnectRedisDb()
	if err != nil {
		log.Fatal(err)
	}
	logWriter := repository.NewLogWriter(RedisDB)
	log.SetOutput(logWriter)

	PostgresDB, err := postgres.ConnectPostgresDb()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = PostgresDB.Close(); err != nil {
			log.Fatal("cannot close psql connection", err)
		}
	}()
	if err = PostgresDB.Migrate(ctx, "db/migrations"); err != nil {
		log.Fatal(err)
	}

	Service, err := service.NewService(PostgresDB)
	UserService, err := service.NewUserService(PostgresDB)
	newApp, err := app.NewService(ctx, Service, UserService)
	if err != nil {
		log.Fatal(err)
	}
	err = newApp.Start()
	if err != nil {
		log.Fatal(err)
	}
}

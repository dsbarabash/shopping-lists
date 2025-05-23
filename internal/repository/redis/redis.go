package redis

import (
	"context"
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/config"
	"github.com/go-redis/redis/v8"
	"log"
)

func ConnectRedisDb() (*redis.Client, error) {
	ctx := context.Background()
	cfg := config.NewRedisConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), // Адрес и порт Redis-сервера
		Password: "",                                       // Пароль (если есть)
		DB:       0,                                        // Номер базы данных
	})

	// Проверка соединения
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Println("Failed to connect to redis: ", err)
		return nil, err
	}
	return client, nil
}

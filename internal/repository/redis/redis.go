package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

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

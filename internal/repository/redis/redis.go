package redis

import (
	"context"
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/config"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
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

func ConnectRedisDbWithRetries(maxAttempts int, delay time.Duration) (*redis.Client, error) {
	var client *redis.Client
	var err error
	ctx := context.Background()
	cfg := config.NewRedisConfig()

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), // Адрес и порт Redis-сервера
			Password: "",                                       // Пароль (если есть)
			DB:       0,                                        // Номер базы данных
		})
		_, err = client.Ping(ctx).Result()
		if err != nil {
			log.Printf("Attempt %d: connection error: %v", attempt, err)
			time.Sleep(delay)
			continue
		} else {
			log.Printf("Successfully connected after %d attempts", attempt)
			return client, nil
		}

	}
	log.Printf("Ping failed: %v", err)
	return nil, err
}

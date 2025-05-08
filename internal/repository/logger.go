package repository

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"io"
	"time"
)

func NewLogWriter(redisDB *redis.Client) io.Writer {
	return &logWriter{
		redisDB,
	}
}

type logWriter struct {
	redisDB *redis.Client
}

func (l *logWriter) Write(p []byte) (n int, err error) {
	err = l.redisDB.Set(l.redisDB.Context(), time.Now().String(), p, 30*time.Minute).Err()
	if err != nil {
		fmt.Println("Ошибка при кешировании данных:", err)
		return 0, err
	}
	return len(p), nil
}

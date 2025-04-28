package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func NewRedisConfig() *RedisConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env.example: %v", err)
	}

	redisDBStr := os.Getenv("REDIS_DB")
	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		log.Fatalf("Ошибка преобразования REDIS_DB в число: %v", err)
	}

	return &RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB,
	}
}

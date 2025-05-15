// config/config.go
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HostConfig *HostConfig
	DB         *DBConfig
	Redis      *RedisConfig
	JWT        *JWTConfig
	SMTP       *SMTPEmailConfig
	Telegram   *TelegramConfig
	S3         *S3Config
	Kafka      *KafkaConfig
}

type SMTPEmailConfig struct {
	Host       string
	Port       string
	Email      string
	Password   string
	MaxRetries int
	RetryDelay time.Duration
}

type JWTConfig struct {
	Secret []byte
	Exp    time.Duration
}

type DBConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	Name         string
	SSLMode      string
	PoolCapacity int
}

type HostConfig struct {
	Port string
}

type TelegramConfig struct {
	Token     string
	IsEnabled bool
}

type KafkaConfig struct {
	URL string
}

// GetConfig загружает конфигурацию из файла .env и переменных окружения
func GetConfig() *Config {
	// Загружаем переменные из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}

	// Читаем переменные окружения и преобразуем необходимые
	dbPortStr := os.Getenv("DB_PORT")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Ошибка преобразования DB_PORT в число: %v", err)
	}

	dbPoolCapacityStr := os.Getenv("DB_POOL_CAPACITY")
	dbPoolCapacity, err := strconv.Atoi(dbPoolCapacityStr)
	if err != nil {
		log.Fatalf("Ошибка преобразования DB_POOL_CAPACITY в число: %v", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	jwtExpStr := os.Getenv("JWT_EXP") // Время жизни токена в минутах
	jwtExp, err := strconv.Atoi(jwtExpStr)
	if err != nil {
		log.Fatalf("Ошибка преобразования JWT_EXP в число: %v", err)
	}
	jwtExpDuration := time.Duration(jwtExp * int(time.Minute))

	redisConfig := NewRedisConfig()

	// Получение макс. числа отправок на почту 1 письма
	smtpMaxRetriesStr := os.Getenv("SMTP_MAX_RETRIES")
	smtpMaxRetries, err := strconv.Atoi(smtpMaxRetriesStr)
	if err != nil {
		log.Fatalf("Ошибка преобразования SMTP_MAX_RETRIES в число: %v", err)
	}

	// Задержка меду отправками
	smtpRetryDelayStr := os.Getenv("SMTP_RETRY_DELAY")
	smtpRetryDelay, err := strconv.Atoi(smtpRetryDelayStr)
	if err != nil {
		log.Fatalf("Ошибка преобразования SMTP_RETRY_DELAY в число: %v", err)
	}
	// Преобразуем в time.Duration
	retryDelayDuration := time.Duration(smtpRetryDelay) * time.Second

	tgEnabled := os.Getenv("TG_ENABLED") == "true"
	tgToken := os.Getenv("TG_TOKEN")
	if tgToken == "" && tgEnabled {
		log.Fatalf("Не задан токен Telegram бота")
	}

	s3Config, err := GetS3Config()
	if err != nil {
		log.Fatal(err)
	}
	kafka := &KafkaConfig{
		URL: os.Getenv("KAFKA_URL"),
	}
	return &Config{
		HostConfig: &HostConfig{
			Port: os.Getenv("PORT"),
		},
		DB: &DBConfig{
			Host:         os.Getenv("DB_HOST"),
			Port:         dbPort,
			User:         os.Getenv("DB_USER"),
			Password:     os.Getenv("DB_PASSWORD"),
			Name:         os.Getenv("DB_NAME"),
			SSLMode:      os.Getenv("DB_SSLMODE"),
			PoolCapacity: dbPoolCapacity,
		},
		Redis: redisConfig,
		JWT: &JWTConfig{
			Secret: []byte(jwtSecret),
			Exp:    jwtExpDuration,
		},
		SMTP: &SMTPEmailConfig{
			Host:       os.Getenv("SMTP_HOST"),
			Port:       os.Getenv("SMTP_PORT"),
			Email:      os.Getenv("SMTP_EMAIL"),
			Password:   os.Getenv("SMTP_PASSWORD"),
			MaxRetries: smtpMaxRetries,
			RetryDelay: retryDelayDuration,
		},
		Telegram: &TelegramConfig{
			Token:     tgToken,
			IsEnabled: tgEnabled,
		},
		S3:    s3Config,
		Kafka: kafka,
	}
}

// GetDSN возвращает строку подключения к базе данных
func (db *DBConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.Name, db.SSLMode)
}

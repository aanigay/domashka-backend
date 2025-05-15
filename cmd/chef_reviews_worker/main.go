package main

import (
	"context"
	"domashka-backend/config"
	"domashka-backend/pkg/postgres"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	topic     = "chef_reviews"
	batchSize = 1
)

func getKafkaReader(kafkaURL, topic string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Topic:       topic,
		Brokers:     brokers,
		MaxWait:     1 * time.Second,
		GroupID:     "chef_reviews_worker",
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.LastOffset,
	})
}

type chefReview struct {
	ChefID int64 `json:"chef_id"`
	Rating int16 `json:"rating"`
}

func main() {
	// get kafka reader using environment variables.
	fmt.Println("KAFKA_URL 1", os.Getenv("KAFKA_URL"))
	err := godotenv.Load()

	fmt.Println("KAFKA_URL 2", os.Getenv("KAFKA_URL"))
	reader := getKafkaReader(os.Getenv("KAFKA_URL"), topic)
	defer reader.Close()
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
	dbConfig := config.DBConfig{
		Host:         os.Getenv("DB_HOST"),
		Port:         dbPort,
		User:         os.Getenv("DB_USER"),
		Password:     os.Getenv("DB_PASSWORD"),
		Name:         os.Getenv("DB_NAME"),
		SSLMode:      os.Getenv("DB_SSLMODE"),
		PoolCapacity: dbPoolCapacity,
	}

	pg, err := postgres.New(dbConfig.GetDSN(), postgres.MaxPoolSize(dbConfig.PoolCapacity))
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer pg.Close()

	fmt.Println("start consuming ... !!")
	reviews := make([]chefReview, 0)
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}

		review := chefReview{}
		err = json.Unmarshal(m.Value, &review)
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		reviews = append(reviews, review)
		if len(reviews) >= batchSize {
			err = UpdateChefRating(pg, reviews...)
			if err != nil {
				log.Println(err)
			}
			log.Printf("Review processed: %s", string(m.Key))
			reviews = make([]chefReview, 0)
		}
	}
}

func UpdateChefRating(pg *postgres.Postgres, reviews ...chefReview) error {
	ctx := context.Background()
	for _, review := range reviews {
		var currentRating float32
		var currentCount int32
		err := pg.Pool.QueryRow(ctx, `SELECT rating, reviews_count FROM chef_ratings WHERE chef_id = $1`, review.ChefID).Scan(&currentRating, &currentCount)
		if err != nil {
			return err
		}
		newRating := (float32(review.Rating) + currentRating*float32(currentCount)) / float32(currentCount+1)
		_, err = pg.Pool.Exec(ctx, `UPDATE chef_ratings SET rating = $1, reviews_count = reviews_count + 1 WHERE chef_id = $2`, newRating, review.ChefID)
		if err != nil {
			return err
		}
	}
	return nil
}

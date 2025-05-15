package integration

import (
	"context"
	"domashka-backend/config"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"strings"
	"testing"
	"time"
)

func TestIntegration_KafkaInteraction(t *testing.T) {
	cfg := config.GetConfig()
	dishReviewsWriter := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Kafka.URL),
		Topic:        "dish_reviews",
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    1,
		BatchTimeout: 10 * time.Millisecond,
	}
	defer func() {
		err := dishReviewsWriter.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	chefReviewsWriter := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Kafka.URL),
		Topic:        "chef_reviews",
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    1,
		BatchTimeout: 10 * time.Millisecond,
	}
	defer func() {
		err := chefReviewsWriter.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	brokers := strings.Split(cfg.Kafka.URL, ",")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Topic:       "chef_reviews",
		Brokers:     brokers,
		MaxWait:     1 * time.Second,
		GroupID:     "chef_reviews_worker",
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.LastOffset,
	})
	defer func() {
		err := reader.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	go func() {
		if err := chefReviewsWriter.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte("test"),
			Value: []byte("test"),
		}); err != nil {
			log.Println(err)
			t.Fatalf(err.Error())
		}
	}()

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		m, err := reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				t.Errorf("timeout waiting for message")
				cancel()
				return
			}
			log.Fatalln(err)
		}
		cancel()

		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}

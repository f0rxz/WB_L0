package main

import (
	"context"
	"encoding/json"
	"log"
	"orderservice/config"
	"orderservice/pkg/generator"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	order := generator.RandomOrder()

	data, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("failed to marshal order: %v", err)
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.KafkaBrokers),
		Topic:        cfg.KafkaOrderTopic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}

	defer func() {
		if err := writer.Close(); err != nil {
			log.Printf("kafka writer close error: %v", err)
		}
	}()

	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("order_key_1"),
			Value: data,
			Time:  time.Now(),
		},
	)

	if err != nil {
		log.Fatalf("failed to write message: %v", err)
	}

	log.Println("Order successfully sent to Kafka")
}

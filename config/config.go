package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PostgresDSN     string        `envconfig:"POSTGRES_DSN" default:"postgres://postgres:postgres@localhost:5432/orderservice?sslmode=disable"`
	KafkaBrokers    string        `envconfig:"KAFKA_BROKERS" default:"localhost:9092"`
	KafkaOrderTopic string        `envconfig:"KAFKA_ORDER_TOPIC" default:"orders"`
	KafkaRetryTopic string        `envconfig:"KAFKA_RETRY_TOPIC" default:"orders_retry"`
	KafkaDLQTopic   string        `envconfig:"KAFKA_DLQ_TOPIC" default:"orders_dlq"`
	KafkaGroupID    string        `envconfig:"KAFKA_GROUP_ID" default:"order-service"`
	HTTPPort        string        `envconfig:"HTTP_PORT" default:":8080"`
	CacheTTL        time.Duration `envconfig:"CACHE_TTL" default:"24h"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &cfg, nil
}

func PrintUsage() {
	var cfg Config
	if err := envconfig.Usage("", &cfg); err != nil {
		fmt.Printf("fail to print envconfig usage: %s\n", err)
	}
}

package config

import "os"

type Config struct {
	PostgresDSN  string
	KafkaBrokers string
	KafkaTopic   string
	KafkaGroupID string
	HTTPPort     string
}

func LoadConfig() (*Config, error) {
	return &Config{
		PostgresDSN:  getEnv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/orderservice?sslmode=disable"),
		KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost:9092"),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "orders"),
		KafkaGroupID: getEnv("KAFKA_GROUP_ID", "order-service"),
		HTTPPort:     getEnv("HTTP_PORT", ":8080"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

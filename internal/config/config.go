package config

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	PostgresDSN  string
	KafkaBrokers string
	KafkaTopic   string
	KafkaGroupID string
	HTTPPort     string
	CacheTTL     time.Duration
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		PostgresDSN:  mustEnv("POSTGRES_DSN"),
		KafkaBrokers: mustEnv("KAFKA_BROKERS"),
		KafkaTopic:   mustEnv("KAFKA_TOPIC"),
		KafkaGroupID: mustEnv("KAFKA_GROUP_ID"),
		HTTPPort:     getEnvOrWarn("HTTP_PORT", "8080"),
	}

	ttlStr := getEnvOrWarn("CACHE_TTL", "24h")
	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		return nil, fmt.Errorf("invalid CACHE_TTL value %q: %w", ttlStr, err)
	}
	cfg.CacheTTL = ttl

	return cfg, nil
}

func mustEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic(fmt.Sprintf("missing required environment variable: %s", key))
	}
	return value
}

func getEnvOrWarn(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("[WARN] env %q not set, using default: %q", key, defaultValue)
		return defaultValue
	}
	return value
}

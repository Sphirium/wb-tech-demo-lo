package config

import "os"

type Config struct {
	HTTPPort      string
	KafkaBroker   string
	KafkaTopic    string
	PostgresURL   string
	RedisAddr     string
	RedisPassword string
	LogLevel      string
}

func Load() *Config {
	return &Config{
		HTTPPort:      getEnv("HTTP_PORT", "8081"),
		KafkaBroker:   getEnv("KAFKA_BROKER", "localhost:9092"),
		KafkaTopic:    getEnv("KAFKA_TOPIC", "orders"),
		PostgresURL:   getEnv("POSTGRES_URL", "host=127.0.0.1 user=wbuser password=wbpass dbname=wb_orders port=5432 sslmode=disable"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

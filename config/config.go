package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	API        APIConfig        `yaml:"api"`
	RabbitMQ   RabbitMQConfig   `yaml:"rabbitmq"`
	Redis      RedisConfig      `yaml:"redis"`
	Processing ProcessingConfig `yaml:"processing"`
}

type APIConfig struct {
	Port string `yaml:"port"`
}

type RabbitMQConfig struct {
	URL         string `yaml:"url"`
	MainQueue   string `yaml:"main_queue"`
	RetryQueue  string `yaml:"retry_queue"`
	DLQ         string `yaml:"dlq"`
	DLXExchange string `yaml:"dlx_exchange"`
}

type RedisConfig struct {
	Addr string `yaml:"addr"`
}

type ProcessingConfig struct {
	MaxRetryCount            int `yaml:"max_retry_count"`
	RetryDelayMilliseconds   int `yaml:"retry_delay_milliseconds"`
	ConsumerTimeoutSeconds   int `yaml:"consumer_timeout_seconds"`
	IdempotencyProcessingTTL int `yaml:"idempotency_processing_ttl_seconds"`
	IdempotencyCompletedTTL  int `yaml:"idempotency_completed_ttl_hours"`
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}

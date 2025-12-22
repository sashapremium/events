package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env      string         `yaml:"env"`
	Postgres DatabaseConfig `yaml:"postgres"`
	Kafka    KafkaConfig    `yaml:"kafka"`
	HTTP     HTTPConfig     `yaml:"http"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}

type HTTPConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func LoadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	if len(cfg.Kafka.Brokers) == 0 {
		return nil, fmt.Errorf("kafka.brokers is required")
	}
	if cfg.Kafka.Topic == "" {
		return nil, fmt.Errorf("kafka.topic is required")
	}

	return &cfg, nil
}

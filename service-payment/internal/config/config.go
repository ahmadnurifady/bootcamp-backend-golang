package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

type KafkaConfig struct {
	KafkaBrokers   string
	KafkaGroupId   string
	OrchestraTopic string
}

type Config struct {
	KafkaConfig
}

func (c *Config) readConfig() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	c.KafkaConfig = KafkaConfig{
		KafkaBrokers:   os.Getenv("KAFKA_BROKERS"),
		KafkaGroupId:   os.Getenv("KAFKA_GROUP_ID"),
		OrchestraTopic: os.Getenv("ORCHESTRA_TOPIC"),
	}

	if c.KafkaConfig.KafkaBrokers == "" || c.KafkaConfig.KafkaGroupId == "" || c.KafkaConfig.OrchestraTopic == "" {
		return errors.New("missing required environment variables")
	}

	return nil
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := cfg.readConfig(); err != nil {
		return nil, err
	}

	return cfg, nil
}

package kafka

import "github.com/kelseyhightower/envconfig"

type Config struct {
	KafkaBrokers []string `envconfig:"KAFKA_BROKERS" required:"true"`
	KafkaGroupID string   `envconfig:"KAFKA_GROUP_ID" default:"my-consumer-group"`
	KafkaTopic   string   `envconfig:"KAFKA_TOPIC" required:"true"`
	KafkaVersion string   `envconfig:"KAFKA_VERSION" default:"2.8.0"`
}

func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

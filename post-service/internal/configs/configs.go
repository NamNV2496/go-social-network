package configs

import (
	"fmt"

	"github.com/namnv2496/post-service/configs"
	"gopkg.in/yaml.v2"
)

type Config struct {
	GRPC     GRPC     `yaml:"grpc"`
	Kafka    Kafka    `yaml:"kafka"`
	Database Database `yaml:"database"`
	Logger   Logger   `yaml:"logger"`
}

func NewConfig() (Config, error) {
	var (
		configBytes []byte = configs.DefaultConfigBytes
		config      Config
		err         error
	)
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshal configuration file: %w", err)
	}

	return config, nil
}

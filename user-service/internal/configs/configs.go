package configs

import (
	"fmt"

	"github.com/namnv2496/user-service/configs"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Auth  Auth  `yaml:"auth"`
	Redis Redis `yaml:"redis"`
	GRPC  GRPC  `yaml:"grpc"`
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

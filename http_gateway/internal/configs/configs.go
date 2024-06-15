package configs

import (
	"fmt"

	"github.com/namnv2496/http_gateway/configs"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Auth    Auth    `yaml:"auth"`
	Gateway Gateway `yaml:"gateway"`
	GRPC    GRPC    `yaml:"grpc"`
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

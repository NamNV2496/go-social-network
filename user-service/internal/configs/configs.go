package configs

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Auth          Auth
	Redis         Redis
	GRPC          GRPC
	Database      Database
	ElasticSearch ElasticSearch
	Email         Email
	SMS           SMS
	Location      Location
}

func NewConfig() (*Config, error) {
	var dbConfig Config
	if err := env.Parse(&dbConfig); err != nil {
		fmt.Printf("Failed to parse environment variables: %v", err)
		return nil, err
	}
	return &dbConfig, nil
}

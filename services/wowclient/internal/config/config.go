package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type General struct {
	ServerAddress string `envconfig:"SERVER_ADDRESS" default:"localhost:8000"`
	NumOfClients  int    `envconfig:"NUM_OF_CLIENTS" default:"50"`
}

type Configuration struct {
	General General
}

func LoadConfig() (*Configuration, error) {
	cfg := &Configuration{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("error process config: %w", err)
	}
	return cfg, nil
}

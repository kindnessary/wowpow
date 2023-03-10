package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type SQLDatabase struct {
	Host               string        `envconfig:"POSTGRES_HOST" default:"localhost"`
	Port               int           `envconfig:"POSTGRES_PORT" default:"5432"`
	User               string        `envconfig:"POSTGRES_USER" default:"postgres"`
	Password           string        `envconfig:"POSTGRES_PASS" default:"postgres"`
	Name               string        `envconfig:"POSTGRES_DB" default:"wowserver"`
	MaxOpenConns       int           `envconfig:"POSTGRES_MAX_OPEN_CONNS" default:"10"`
	ConnMaxLifetime    time.Duration `envconfig:"POSTGRES_CONN_MAX_LIFETIME" default:"10m"`
	MaxIdleConns       int           `envconfig:"POSTGRES_MAX_IDLE_CONNS" default:"10"`
	StorageDriver      string        `envconfig:"POSTGRES_STORAGE_DRIVER" default:"pgx"`
	SSLMode            string        `envconfig:"POSTGRES_SSL_MODE" default:"disable"`
	ConnectionAttempts int           `envconfig:"POSTGRES_CONNECTION_ATTEMPTS" default:"5"`
	ConnectionInterval time.Duration `envconfig:"POSTGRES_CONNECTION_INTERVAL" default:"3s"`
}

type General struct {
	Address            string        `envconfig:"ADDRESS" default:"0.0.0.0:8000"`
	NumOfQuotes        int           `envconfig:"NUM_OF_QUOTES" default:"35"`
	Difficulty         uint8         `envconfig:"DIFFICULTY" default:"15"`
	ConnectionLifetime time.Duration `envconfig:"CONNECTION_LIFETIME" default:"30s"`
}

type Configuration struct {
	SQLDatabase SQLDatabase
	General     General
}

func LoadConfig() (*Configuration, error) {
	cfg := &Configuration{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("error process config: %w", err)
	}
	return cfg, nil
}

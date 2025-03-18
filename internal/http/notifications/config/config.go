package config

import (
	"context"
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port        string        `envconfig:"PORT" default:"3001"`
	Environment string        `envconfig:"ENVIRONMENT" default:"dev"`
	ServiceName string        `envconfig:"SERVICE_NAME" default:"casino_loyalty_reward_system"`
	DatabaseURI string        `envconfig:"DB_URI" default:"postgres://tester:testing@localhost:5432/main"`
	RedisURI    string        `envconfig:"REDIS_URI" default:"redis://localhost:6379"`
	JWTKey      string        `envconfig:"JWT_KEY" default:"true"`
	JWTDuration time.Duration `envconfig:"JWT_DURATION" default:"24h"`
}

func newConfig(ctx context.Context) (*Config, error) {
	var config Config

	err := envconfig.Process("", &config)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch environment variables: %w", err)
	}

	return &config, nil
}
